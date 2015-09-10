package gogomailer

import (
	"appengine/urlfetch"
	"fmt"
	mg "github.com/mailgun/mailgun-go"
	"strings"
)

// sendWithMailgun handles the main send and logging logic for any single email deploy
func (c *Client) sendWithMailgun() (string, error) {
	var err error

	c.Body, err = c.prepareTmpl()
	if err != nil {
		return "", err
	}

	gun := mg.NewMailgun(Conf["mailgun"]["domain"], Conf["mailgun"]["secret"], Conf["mailgun"]["public"])

	// override the http client
	client := urlfetch.Client(c.context)
	gun.SetClient(client)

	// generate mailgun message
	message := mg.NewMessage(fmt.Sprintf("%s <%s>", Conf["default"]["fromname"], Conf["default"]["fromemail"]), c.Subject, c.Body, c.Recipient[0].Render)
	message.SetHtml(strings.Replace(c.Body, "\\", "", -1))

	// add additional recipients
	for k, v := range c.Recipient {
		if k > 0 {
			err := message.AddRecipient(v.Render)
			if err != nil {
				c.context.Errorf("Could not append [%s] as Mailgun recipient: %v", v.Render, err)
				return "", err
			}
		}
	}

	// send the email
	_, id, err := gun.Send(message)
	if err != nil {
		c.context.Errorf("Error: %v", err)
		return "", err
	}

	if Conf["default"]["logmessages"] == "true" {
		// if the id is not empty then add to the message log
		if id != "" {
			_, logErr := c.addMessageLog()
			if logErr != nil {
				c.context.Errorf("Failed to add message to log: %v", logErr)
			}
		}
	}

	return id, err
}
