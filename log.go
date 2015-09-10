package gogomailer

import (
	"appengine/datastore"
	goutil "github.com/markhayden/goutils"
	"time"
)

const (
	MESSAGE_LOG_TABLE = "GoGoMessageLog"
)

type MessageLog struct {
	Id        string      `json:"id"`
	Subject   string      `json:"subject"`
	Body      string      `json:"body"`
	Recipient []Recipient `json:"recipient"`
	Sent      time.Time
}

// addMessageLog takes a client struct and creates a db log for the message
// this will allow us to look at old reports and eventually build in some
// safeguards so we don't spam any one recipient
func (c *Client) addMessageLog() (string, error) {
	// make a new unique id
	id, err := goutil.MakeUniqueId()
	if err != nil {
		return "", err
	}

	log := MessageLog{
		Id:        id,
		Subject:   c.Subject,
		Body:      c.Body,
		Recipient: c.Recipient,
		Sent:      time.Now(),
	}

	// establish data key
	key := datastore.NewKey(c.context, MESSAGE_LOG_TABLE, id, 0, nil)

	_, err = datastore.Put(c.context, key, &log)
	if err != nil {
		c.context.Errorf("Failed to add message log: %v", err)
		return "", err
	}

	return id, nil
}
