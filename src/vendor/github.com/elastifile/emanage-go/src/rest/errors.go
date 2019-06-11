package rest

import (
	"net/http"
	"strings"
)

type restError struct {
	Description string
	Response    *http.Response
	Body        string
}

func NewRestError(description string, response *http.Response, resBody []byte) *restError {
	err := &restError{
		Description: description,
		Response:    response,
		Body:        string(resBody),
	}
	return err
}

func (e *restError) Error() string {
	return strings.Join(
		[]string{
			e.Description,
			"Status: " + e.Response.Status,
			"Body: " + e.removePrefix(e.Body),
		},
		"\n",
	)
}

func (e *restError) removePrefix(body string) string {
	const prefix = ")]}',\n"
	if strings.HasPrefix(body, prefix) {
		body = body[len(prefix):]
	}
	return body
}
