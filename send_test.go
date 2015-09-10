package gogomailer

import (
	"appengine/aetest"
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

var (
	basic_email_test = `{
	    "template": "test",
	    "params": {
	        "test": "one",
	        "test2": "2"
	    },
	    "recipient": [
	        {
	            "name": "Mark",
	            "email": "hi+one@markhayden.me"
	        },
	        {
	            "name": "Hayden",
	            "email": "hi+two@markhayden.me"
	        }
	    ],
	    "subject": "test",
	    "body": "hup"
	}`
)

func TestNewClient(t *testing.T) {
	// test setup
	b := []byte(basic_email_test)

	ctx, _ := aetest.NewContext(nil)
	defer ctx.Close()

	// test creating client
	c, err := NewClient(b, ctx)

	assert.Equal(t, len(c.Recipient) == 2, true)
	assert.Equal(t, c.Subject == "test", true)
	assert.Equal(t, c.Params["test"] == "one", true)
	assert.Equal(t, err == nil, true)

	// test template rendering
	tmpl, err := c.prepareTmpl()
	assert.Equal(t, tmpl == "<html><head></head><body>Here I am. Var one.</body></html>", true)
	assert.Equal(t, err == nil, true)

	id, err := c.sendWithMailgun()
	fmt.Println(id)
	fmt.Println(err)
}
