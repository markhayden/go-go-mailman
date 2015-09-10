package gogomailer

import (
	"appengine"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2"
	goutil "github.com/markhayden/goutils"
	"io/ioutil"
)

// NewClient creates the primary send client for interface with mailgun and
// all template rendering
func NewClient(payload []byte, context appengine.Context) (Client, error) {
	var c Client
	err := json.Unmarshal(payload, &c)
	if err != nil {
		return c, err
	}

	// prepare the recipient payload for email client
	for k, v := range c.Recipient {
		c.Recipient[k].Render = fmt.Sprintf("%s <%s>", v.Name, v.Email)
	}

	c.context = context

	return c, err
}

// loadTmpl takes a file name and loads the template from directory
func loadTmpl(tmp string) (string, error) {
	content, err := ioutil.ReadFile(tmp)
	if err != nil {
		return "", err
	}
	bodyPretty := goutil.PrepHtmlForJson(string(content), false)
	return bodyPretty, nil
}

// prepareTmpl takes the loaded template and merges all params available
// on the client into the template
func (c *Client) prepareTmpl() (string, error) {
	// load local template file into memory
	body, err := loadTmpl(fmt.Sprintf("%s%s.html", Conf["template"]["path"], c.Template))
	if err != nil {
		c.context.Errorf("Unable to load template: %v", err)
		return "", err
	}

	tmpl, err := pongo2.FromString(body)
	if err != nil {
		c.context.Errorf("Unable to prepare body of template for pongo: %v", err)
		return "", err
	}

	output, err := tmpl.Execute(c.Params)
	if err != nil {
		c.context.Errorf("Unable to merge variables with pongo template: %v", err)
		return "", err
	}

	return output, nil
}
