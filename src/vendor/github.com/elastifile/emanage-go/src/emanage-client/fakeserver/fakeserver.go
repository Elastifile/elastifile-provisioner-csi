package fakeserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"emanage-client"
	"logging"
)

var logger = logging.NewLogger("fakeserver")

func New() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/sessions", login)
	mux.HandleFunc("/api/systems/", list)
	mux.HandleFunc("/api/systems/1/setup", setup)
	mux.HandleFunc("/api/systems/1/start", start)
	mux.HandleFunc("/api/systems/1/shutdown", shutdown)

	mux.HandleFunc("/api/policies", policies)
	mux.HandleFunc("/api/data_containers", dataContainers)
	mux.HandleFunc("/api/exports", exports)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("fakeServer: Got unknown request", "request", *r)
		http.NotFound(w, r)
	})

	return httptest.NewServer(mux)
}

func login(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	logger.Debug("fakeServer: login", "request", *r, "body", string(data))
	writeString(w, `{"info":"Logged in","user":{"id":1,"login":"admin","admin":true,"first_name":"Super","surname":"Admin","email":"admin@example.com","created_at":"2015-11-15T09:37:06.000Z","updated_at":"2015-11-24T15:24:44.551Z"}}`)
}

func list(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeServer: list systems", "request", *r)
	writeSystemDetails(w)
}

func setup(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeServer: setup", "request", *r)
	writeSystemDetails(w)
}

func start(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeServer: start", "request", *r)
	time.Sleep(200 * time.Millisecond)
	system.state = emanage.StateInService
	writeSystemDetails(w)
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeServer: shutdown", "request", *r)
	time.Sleep(200 * time.Millisecond)
	system.state = emanage.StateDown
	writeSystemDetails(w)
}

////////////////////////////////////////////////////////////////////////////////

func policies(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeserver: policies", "request", *r)
	p := emanage.Policy{
		Id:   1,
		Name: "mypolicy",
	}
	ps := []emanage.Policy{p}
	data, err := json.Marshal(ps)
	if err != nil {
		panic(err)
	}
	logger.Debug("fakeserver: policies", "response", string(data))
	write(w, data)
}

////////////////////////////////////////////////////////////////////////////////

func dataContainers(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeserver: data containers", "request", *r)
	dc := emanage.DataContainer{
		Id:       1,
		Name:     "mydc",
		PolicyId: 1,
	}
	dcs := []emanage.DataContainer{dc}
	data, err := json.Marshal(dcs)
	if err != nil {
		panic(err)
	}
	logger.Debug("fakeserver: data containers", "response", string(data))
	write(w, data)
}

////////////////////////////////////////////////////////////////////////////////

func exports(w http.ResponseWriter, r *http.Request) {
	logger.Debug("fakeserver: exports", "request", *r)
	exp := emanage.Export{
		Name:            "myexport",
		DataContainerId: 1,
	}
	exps := []emanage.Export{exp}
	data, err := json.Marshal(exps)
	if err != nil {
		panic(err)
	}
	logger.Debug("fakeserver: exports", "response", string(data))
	write(w, data)
}

////////////////////////////////////////////////////////////////////////////////

var securityPrefix = []byte(")]}',\n")

func write(w http.ResponseWriter, data []byte) {
	http.SetCookie(w, &http.Cookie{
		Name:  "XSRF-TOKEN",
		Value: "bogus cookie",
	})
	_, _ = w.Write(securityPrefix)
	_, _ = w.Write(data)
}

func writeString(w http.ResponseWriter, s string) {
	write(w, []byte(s))
}

////////////////////////////////////////////////////////////////////////////////

var system struct {
	state emanage.SystemState
}

func init() {
	system.state = emanage.StateInService
}

func writeSystemDetails(w http.ResponseWriter) {
	sysDetails := &emanage.SystemDetails{
		Id:   1,
		Name: "FakeSystem",
		// UUID:   uuid.UUID("62cadcb6-94eb-48ab-aca0-2d28b9492b36"),
		Status: system.state,
	}

	data, _ := json.Marshal(sysDetails)
	write(w, data)
}
