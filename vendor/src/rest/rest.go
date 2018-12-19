package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"gopkg.in/inconshreveable/log15.v2"

	"ejson"
	"helputils"
	"logging"
	"retry"
	"runtimeutil"
)

var logger = logging.NewLogger("rest")

type HttpMethod string

const (
	MethodPost   HttpMethod = "POST"
	MethodPut    HttpMethod = "PUT"
	MethodGet    HttpMethod = "GET"
	MethodDelete HttpMethod = "DELETE"
)

const sessionsUri = "api/sessions"

var securityPrefix = []byte(")]}',\n")

var defaultTimeout time.Duration = 2 * time.Hour
var DumpHTTP bool
var DumpHTTPOnError = true

var AfterShutdown func()
var BeforeStart func()
var BeforeForceReset func()

type Session struct {
	baseURL     *url.URL // Base URL of server to connect to, e.g. http://func11-cm/
	credentials credentials
	cookies     []*http.Cookie
	Client      http.Client
	xsrf        string
	SessionsUri string
	timeout     time.Duration
}

func NewSession(baseURL *url.URL, timeout time.Duration) *Session {
	result := &Session{
		baseURL: baseURL,
		timeout: timeout,
	}
	result.init()
	logger.Debug("new session", "url", baseURL)
	return result
}

func NewDefaultSession(baseURL *url.URL) *Session {
	return NewSession(baseURL, defaultTimeout)
}

func CheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (rs *Session) init() {
	rs.Client = http.Client{
		CheckRedirect: CheckRedirect,
		Timeout:       rs.timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	rs.SessionsUri = sessionsUri
}

type credentials struct {
	User     string `json:"login"`
	Password string `json:"password"`
}

func (rs *Session) Url() *url.URL {
	return rs.baseURL
}

func (rs *Session) Login(user string, password string) error {
	creds := credentials{user, password}
	params := struct {
		User credentials `json:"user"`
	}{creds}
	rs.credentials = creds

	jsonBody, stdErr := json.Marshal(params)
	if stdErr != nil {
		return errors.Wrap(stdErr, 0)
	}

	resp, _, err := rs.requestHttp(MethodPost, rs.SessionsUri, jsonBody)
	if err != nil {
		return errors.Wrap(err, 0)
	} else {
		if resp.StatusCode == http.StatusFound {
			if strings.Contains(resp.Header.Get("Location"), "https://") {
				rs.baseURL.Scheme = "https"
				logger.Debug("Received redirect status code, Changing scheme http -> https and retrying login", "status code", resp.StatusCode)
			}
			resp, _, err = rs.requestHttp(MethodPost, rs.SessionsUri, jsonBody)
			if err != nil {
				logger.Error("Failed to login", "method", MethodPost, "url", rs.baseURL, "uri", rs.SessionsUri, "err", err)
				return err
			}
			logger.Debug("logged-in", "method", MethodPost, "url", rs.baseURL, "uri", rs.SessionsUri)
		}
	}

	rs.cookies = resp.Cookies()
	xsrf, err := func() (string, error) {
		for _, cookie := range rs.cookies {
			if cookie.Name == "XSRF-TOKEN" {
				xsrf, e := url.QueryUnescape(cookie.Value)
				if e != nil {
					return "", e
				}

				return xsrf, nil
			}
		}
		return "", NewRestError("XSRF cookie not found", resp, []byte(""))
	}()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	rs.xsrf = xsrf

	return nil
}

func (rs *Session) RetriedLogin(user string, password string) error {
	return rs.RetriedLoginTimeout(user, password, 30*time.Second)
}

// Returns the inner error if the argument is an Error, otherwise
// return it as-is.
func InnerError(e error) error {
	switch e := e.(type) {
	case *errors.Error:
		return e.Err
	default:
		return e
	}
}

func (rs *Session) RetriedLoginTimeout(user string, password string, timeout time.Duration) error {
	interval := 1 * time.Second
	err := retry.Basic{
		Timeout: interval,
		Retries: int(timeout / interval),
	}.Do(func() error {
		e := rs.Login(user, password)
		if ue, ok := InnerError(e).(*url.Error); ok {
			if _, ok := ue.Err.(*net.OpError); ok {
				return &retry.TemporaryError{Err: e}
			}
		}
		return e
	})
	return err
}

func (rs *Session) Logout() error {
	if err := rs.Request(MethodDelete, rs.SessionsUri, nil, nil); err != nil {
		return err
	}

	rs.cookies = nil
	rs.xsrf = ""
	return nil
}

func (rs *Session) requestHttp(method HttpMethod, relURL string, body []byte) (resp *http.Response, resBody []byte, resErr error) {
	// FIXME: Change this to be (but there seems to be an issue with an extra leading or trailing slash):
	// fullURL := path.Join(baseURL, relURL)
	fullURL := fmt.Sprintf("%s/%s", rs.baseURL, relURL)

	// logger.Debug("Request", "method", method, "relURL", relURL, "body", string(body))
	req, err := http.NewRequest(string(method), fullURL, bytes.NewReader(body))
	if err != nil {
		resErr = errors.Wrap(err, 0)
		return
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	if rs.cookies != nil {
		for _, cookie := range rs.cookies {
			req.AddCookie(cookie)
		}
		req.Header.Add("X-XSRF-TOKEN", rs.xsrf)
	}

	if DumpHTTP {
		e := dumpRequest(req)
		if e != nil {
			resErr = e
			return
		}
	}

	resp, err = rs.Client.Do(req)
	if err != nil {
		fmt.Println("ERROR: HTTP Request failed!")
		if DumpHTTPOnError {
			_ = dumpRequest(req)
		}
		if resp != nil && DumpHTTPOnError {
			_ = dumpResponse(resp)
		}
		resErr = errors.Wrap(err, 0)
		return
	}

	if DumpHTTP {
		e := dumpResponse(resp)
		if e != nil {
			resErr = e
			return
		}
	}

	shouldClose := true
	defer func() {
		if !shouldClose {
			return
		}
		e := resp.Body.Close()
		if e != nil && resErr == nil {
			resErr = errors.Wrap(e, 0)
		}
	}()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			// Connection may be closed by server on logout; ignore.
			if method == MethodDelete && relURL == rs.SessionsUri {
				shouldClose = false
				return
			}
		}
		logger.Error("Error while reading response body", "resp.Body", resp.Body)
		resErr = err
		return
	}

	// logger.Debug("Response", "method", method, "uri", relURL, "resBody", string(resBody))
	if resp.StatusCode >= http.StatusBadRequest {
		resErr = NewRestError("HTTP request failed", resp, resBody)
		return
	}

	return
}

func dumpRequest(req *http.Request) error {
	fmt.Println("## HTTP DUMP REQUEST ##")
	if req == nil {
		fmt.Printf("Nothing to dump (empty request object)\n\n")
		return nil
	}

	data, e := httputil.DumpRequest(req, true)
	if e != nil {
		return errors.Wrap(e, 0)
	}
	fmt.Printf("Header: %+v\n", req.Header)
	fmt.Printf("%v\n", string(data))
	return nil
}

func dumpResponse(resp *http.Response) error {
	fmt.Println("## HTTP DUMP RESPONSE ##")
	if resp == nil {
		fmt.Printf("Nothing to dump (empty response object)\n\n")
		return nil
	}

	data, e := httputil.DumpResponse(resp, true)
	if e != nil {
		return errors.Wrap(e, 0)
	}
	fmt.Printf("Header: %+v\n", resp.Header)
	fmt.Printf("%v\n", string(data))
	return nil
}

//go:generate stringer -type=ControlTaskStatus

type controlTask struct {
	Id          int               `json:"id"`
	Status      ControlTaskStatus `json:"status"`
	LastError   string            `json:"last_error"`
	Name        string            `json:"name"`
	CurrentStep string            `json:"current_step"`
}

type ControlTaskStatus int

const (
	ControlTaskStatusSuccess    ControlTaskStatus = 0
	ControlTaskStatusError      ControlTaskStatus = 1
	ControlTaskStatusCanceled   ControlTaskStatus = 2
	ControlTaskStatusIncomplete ControlTaskStatus = 3
	ControlTaskStatus_last                        = ControlTaskStatusIncomplete
)

var taskStatuses = map[string]ControlTaskStatus{
	"success":     ControlTaskStatusSuccess,
	"error":       ControlTaskStatusError,
	"canceled":    ControlTaskStatusCanceled,
	"in_progress": ControlTaskStatusIncomplete,
}

func (cts *ControlTaskStatus) UnmarshalJSON(data []byte) (err error) {
	var (
		intVal int
		strVal string
	)
	if err = json.Unmarshal(data, &intVal); err == nil {
		*cts = ControlTaskStatus(intVal)
		return
	}
	if err = json.Unmarshal(data, &strVal); err == nil {
		*cts = taskStatuses[strVal]
		return
	}
	return
}

type TaskID struct {
	Url          string            `json:"url"`
	Status       ControlTaskStatus `json:"status"`
	ErrorMessage string            `json:"last_error"`
	Name         string
	CurrentStep  string
	Error        error
}

func (tid *TaskID) String() string {
	return fmt.Sprintf("%+v", *tid)
}

type AsyncRequest struct {
	Async bool `json:"async,omitempty"`
}

func (rs *Session) WaitAllTasks(tasks []*TaskID) error {
	urls := make([]string, len(tasks))
	for i, tid := range tasks {
		urls[i] = tid.Url
	}
	logger.Info("waiting", "tasks", urls)

	s := retry.Sigmoid{
		Limit:   1 * time.Minute,
		Retries: 100,
	}
	return s.Do(func() error {
		var (
			newTasks []*TaskID
			tempErr  error
		)
		watcher := make(chan *TaskID)

		for _, task := range tasks {
			task := task
			go func() {
				var ct controlTask
				taskUrl, _ := url.Parse(task.Url)
				err := rs.Request(MethodGet, taskUrl.Path, nil, &ct)
				if err != nil {
					task.ErrorMessage = fmt.Sprintf("Couldn't connect to %s", task.Url)
				} else {
					task.ErrorMessage = ct.LastError
				}

				task.Error = err
				task.Status = ct.Status

				if task.Name != ct.Name || task.CurrentStep != ct.CurrentStep {
					task.Name, task.CurrentStep = ct.Name, ct.CurrentStep
					if ct.Name != "nop" {
						logger.Info("waiting for task", "url", task.Url, "name", task.Name, "step", task.CurrentStep)
					}
				}

				watcher <- task
			}()
		}

		for range tasks {
			tid := <-watcher
			if tid.Error != nil {
				logger.Error(tid.Error.Error())
				// Change this to ControlTaskStatusIncomplete if you
				// want to retry on HTTP error rather than to give up.
				tid.Status = ControlTaskStatusError
			}
			switch tid.Status {
			case ControlTaskStatusSuccess:
				logger.Info("task succeeded", "tid", tid)
			case ControlTaskStatusError, ControlTaskStatusCanceled:
				err := fmt.Errorf("Task <%s> failed due to %v. Cause: %s",
					tid.Url,
					tid.Status,
					tid.ErrorMessage,
				)
				logger.Error("task error or cancel", "err", err)
				return err
			default:
				tempErr = &retry.TemporaryError{
					Err: fmt.Errorf(
						"Task <%s> didn't complete yet: %v",
						tid.Url,
						tid.Status,
					),
				}
				logger.Debug("task incomplete", "tempErr", tempErr)
				newTasks = append(newTasks, tid)
			}
		}
		tasks = newTasks
		return tempErr
	})
}

func (rs *Session) AsyncRequest(method HttpMethod, uri string, body interface{}) ([]*TaskID, error) {
	var tIDs []*TaskID
	if body == nil {
		body = AsyncRequest{Async: true}
	}
	err := rs.Request(method, uri, body, &tIDs)
	return tIDs, err
}

func (rs *Session) RetriedRequest(method HttpMethod, uri string, body interface{}, result interface{}, timeout time.Duration) (err error) {
	interval := 5 * time.Second
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(interval) {
		err = rs.Request(method, uri, body, result)
		if err == nil {
			return err
		}
	}

	return err
}
func (rs *Session) Request(method HttpMethod, uri string, body interface{}, result interface{}) error {
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			logger.Error("http request failed during json.Marshal",
				"err", err, "method", method, "uri", uri, "body", body)
			return errors.Wrap(err, 0)
		}
	}

	logger.Debug("request",
		"method", method,
		"url", rs.Url(),
		"uri", uri,
		"body", string(jsonBody),
		"result", result,
	)

	if method != MethodGet {
		rs.logAction(method, uri, body)
	}

	resp, resBody, err := rs.requestHttp(method, uri, jsonBody)
	if err != nil {
		if e, ok := err.(*restError); ok {
			if e.Response.StatusCode == http.StatusUnauthorized {
				logger.Debug("Got eManage Auth Error, maybe session expired, relogin and retrying ...")
				rs.init() // invalidate connection before we retry
				e := rs.Login(rs.credentials.User, rs.credentials.Password)
				if e != nil {
					return errors.WrapPrefix(e, "Failed login", 0)
				}
				resp, resBody, err = rs.requestHttp(method, uri, jsonBody)
				if err != nil {
					logger.Error("http request failed",
						"host", rs.baseURL.Host,
						"method", method,
						"uri", uri,
						"header", resp.Header,
						"body", body,
						"err", err,
					)
					return err
				}
			} else {
				logger.Error("http request failed, will dump response from server",
					"host", rs.baseURL.Host,
					"method", method,
					"uri", uri,
					"header", resp.Header,
					"body", body,
					"err", err,
				)
				dumpResponse(resp)
				return err
			}
		} else {
			dumpResponse(resp)
			return err
		}
	}

	if result == nil {
		if resp.StatusCode != http.StatusNoContent {
			logger.Debug("Request returns content, but no result structure was provided", "method", method, "uri", uri)
		}
		return nil
	}

	// trim security prefix
	if bytes.HasPrefix(resBody, securityPrefix) {
		resBody = bytes.TrimPrefix(resBody, securityPrefix)
	}

	err = json.Unmarshal(resBody, result)
	if err != nil {
		logger.Error("http request failed during json.Unmarshal",
			"err", err, "method", method, "uri", uri, "body", body, "result", result, "resBody", string(resBody))
		return ejson.NewError(err, resBody)
	}

	logger.Debug("response",
		"method", method,
		"url", rs.Url(),
		"uri", uri,
		"body", string(jsonBody),
		"resp", resp,
		"resBody", string(resBody),
	)
	return nil
}

var methodActionDict = map[HttpMethod]string{
	MethodGet:    "Read",
	MethodPut:    "Update",
	MethodPost:   "Create",
	MethodDelete: "Delete",
}

func (rs *Session) logAction(method HttpMethod, uri string, body interface{}) {
	kvIfs := []interface{}{}

	if body != nil {
		if helputils.IsReflectingStruct(body) {
			kvIfs = append(kvIfs, helputils.MustStructToKeyValueInterfaces(body)...)
		} else {
			kvIfs = append(kvIfs, "body", body)
		}
	}

	kvIfs = append(kvIfs, "REST", fmt.Sprintf("%s http://%s/%s", method, rs.baseURL.Host, uri))

	caller := actionCaller()
	logLvl := uriLogLvl(fmt.Sprintf("%s /%s", method, uri))
	if logging.ShowCaller(caller, logLvl) {
		kvIfs = append(kvIfs, "caller", caller)
	}

	logMsg := methodActionDict[method] + " " + uriToMsg(uri)
	//logger.WriteLvl(logLvl, logMsg, kvIfs...)
	logger.Debug(logMsg, kvIfs...)
}

func uriToMsg(uri string) (msg string) {
	if uri[0] == '/' {
		uri = uri[1:]
	}
	uriToks := strings.Split(uri, "/")

	entity := uriToks[1]
	msg = singular(entity)

	for {
		if len(uriToks) > 2 {
			id := uriToks[2]
			if _, err := strconv.ParseInt(id, 10, 64); err == nil {
				id = "#" + id
			}
			msg += " " + id

			if len(uriToks) > 3 {
				opName := ": " + uriToks[3]
				msg += singular(opName)
			}
		}

		if len(uriToks) > 4 {
			uriToks = uriToks[2:]
		} else {
			break
		}
	}

	return strings.Replace(msg, "_", " ", -1)
}

func actionCaller() string {
	callStr := runtimeutil.CallerString(0)
	for i := 1; i < 6 && (strings.Contains(callStr, "infra/rest") || strings.Contains(callStr, "infra/emanage")); i++ {
		callStr = runtimeutil.CallerString(i)
	}
	return callStr
}

var restUriRegexpToLogLvlMap = map[string]log15.Lvl{
	"PUT /api/events/\\d+/ack": log15.LvlDebug,
}

func uriLogLvl(uri string) log15.Lvl {
	for re, lvl := range restUriRegexpToLogLvlMap {
		if matched, err := regexp.MatchString(re, uri); err != nil {
			logger.Warn("regex match failure", "err", err)
		} else if matched {
			return lvl
		}
	}
	return log15.LvlInfo
}

func singular(plural string) string {
	if strings.HasSuffix(plural, "ies") {
		return plural[:len(plural)-3] + "y"
	} else if plural[len(plural)-1] == 's' {
		return plural[:len(plural)-1]
	}
	return plural
}
