package emanage

import (
	"fmt"
	"path"

	"github.com/go-errors/errors"

	"rest"
)

const sysEventsUri = "api/events"

type events struct {
	conn *rest.Session
}

type Events struct {
	Id           int    `json:"id"`
	EventTypeId  int    `json:"event_type_id"`
	Message      string `json:"message"`
	Timestamp    string `json:"timestamp,omitempty"`
	Severity     string `json:"severity,omitempty"`
	Acknowledged struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"acknowledged,omitempty"`
}

type EventsRecentOpts struct {
	Since    int    `json:"since,omitempty"`
	Severity string `json:"severity,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Unacked  bool   `json:"unacked,omitempty"`
}

type EventsHistoryOpts struct {
	Limit    int    `json:"since,omitempty"`
	Severity string `json:"severity,omitempty"`
	Bookmark int    `json:"bookmark,omitempty"`
}

func (cr *events) GetAll() (result []Events, err error) {
	if err = cr.conn.Request(rest.MethodGet, sysEventsUri, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (cr *events) GetAllRecent(opts *EventsRecentOpts) (result []Events, err error) {
	if err = cr.conn.Request(rest.MethodGet, path.Join(sysEventsUri, "recent"), opts, &result); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	return result, nil
}

func (cr *events) GetHistory(opts *EventsHistoryOpts) (result []Events, err error) {
	if err = cr.conn.Request(rest.MethodGet, path.Join(sysEventsUri, "history"), opts, &result); err != nil {
		return nil, errors.Wrap(err, 0)
	}
	return result, nil
}

func (cr *events) Ack(Id int) (err error) {
	ackURI := path.Join(sysEventsUri, fmt.Sprintf("%d/ack", Id))
	if err = cr.conn.Request(rest.MethodPut, ackURI, nil, nil); err != nil {
		logger.Error("Ack Error", "err", err, "req", ackURI)
		return err
	}
	return nil
}

func (cr *events) UnAck(Id int) (err error) {
	unackURI := path.Join(sysEventsUri, fmt.Sprintf("%d/unack", Id))
	if err = cr.conn.Request(rest.MethodPut, unackURI, nil, nil); err != nil {
		logger.Error("UnAck Error", "err", err, "req", unackURI)
		return err
	}
	return nil
}

func (cr *events) AckAll(Ids []int) (err error) {
	if len(Ids) == 0 {
		return errors.New("AckAll : Event Id list provided is empty")
	}
	for _, Id := range Ids {
		err = cr.Ack(Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cr *events) UnAckAll(Ids []int) (err error) {
	if len(Ids) == 0 {
		return errors.New("NackAll : Event Id list provided is empty")
	}
	for _, Id := range Ids {
		err = cr.UnAck(Id)
		if err != nil {
			return err
		}
	}
	return nil
}
