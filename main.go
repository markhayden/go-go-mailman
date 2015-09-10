package gogomailer

import (
	"appengine"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/pat"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Id        string                 `json:"id"`
	context   appengine.Context      `json:"context"`
	Token     string                 `json:"token"`
	Template  string                 `json:"template"`
	Params    map[string]interface{} `json:"params"`
	Recipient []Recipient            `json:"recipient"`
	Subject   string                 `json:"subject"`
	Body      string                 `json:"body"`
}

type Recipient struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Render string `json:"render"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func init() {
	var err error

	// load the configuration
	Conf, err = loadConfig()
	if err != nil {
		panic(err)
	}

	m := pat.New()

	m.Post("/deploy/email", http.HandlerFunc(DeployHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m.ServeHTTP(w, r)
	})
}

// DeployHandler handles email deployment. Need to break this out into a taskmanager
// maintained method so that we don't get bogged down processing calls as we scale.
func DeployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create appengine context
	ctx := appengine.NewContext(r)

	ctx.Infof("conf: %s", Conf)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.Errorf("Unable to read the body of the request: %v", err)
		http.Error(w, "Unable to read the body of the request", 400)
		return
	}

	// create a new send client
	c, err := NewClient(body, ctx)

	// verify token
	if c.Token != Conf["default"]["token"] {
		ctx.Errorf("Invalid token: %s", c.Token)
		http.Error(w, "Acess Denied", 400)
		return
	}

	// render the template
	tmpl, err := c.prepareTmpl()
	if err != nil {
		ctx.Errorf("Unable to render template: %v", err)
		http.Error(w, "Unable to render template", 400)
		return
	}

	c.Body = tmpl

	var sentId string
	switch Conf["default"]["smtp"] {
	case "mailgun":
		sentId, err = c.sendWithMailgun()
		if err != nil {
			ctx.Errorf("Unable to deploy email: %v", err)
			http.Error(w, "Unable to deploy email", 400)
			return
		}
	default:
		ctx.Errorf("No mail client found, unable to send.")
		http.Error(w, "No mail client found, unable to send.", 400)
	}

	// everything worked, give success
	w.WriteHeader(200)

	// render message
	out, _ := json.Marshal(Response{
		Status:  200,
		Message: "Success",
		Data: map[string]string{
			"id": sentId,
		},
	})

	fmt.Fprint(w, string(out))
}
