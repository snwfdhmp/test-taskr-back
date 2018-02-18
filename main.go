package main

import (
	"bytes"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/rjz/githubhook.v0"
)

var (
	log = logrus.New()
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/webhook", func(rw http.ResponseWriter, req *http.Request) {
		hook, err := githubhook.Parse(appSecret, req)
		if err != nil {
			log.Println("fatal:", err)
			return
		}
		// var payload github.WebHookPayload
		// if err := json.Unmarshal(hook.Payload, &payload); err != nil {
		// 	log.Println("fatal:", err)
		// 	return
		// }
		// log.Println("event on", payload.GetRepo().GetFullName(), "by", payload.GetSender().GetLogin())
		// log.Println(payload)
		payload, err := github.ParseWebHook(github.WebHookType(req), hook.Payload)
		if err != nil {
			log.Println("fatal:", err)
			return
		}

		if _, ok := payload.(*github.Issue); ok {
			log.Println("is issue")
		} else if _, ok := payload.(*github.WebHookPayload); ok {
			log.Println("is WebHookPayload")
		} else if data, ok := payload.(*github.IssueCommentEvent); ok {
			name := data.GetSender().GetLogin()
			repo := data.GetRepo().GetFullName()
			issue := data.GetIssue().GetTitle()
			body := data.GetComment().GetBody()
			log.Printf("New comment '%s' by %s on issue '%s' on repo %s", body, name, issue, repo)
		} else if data, ok := payload.(*github.IssuesEvent); ok {
			name := data.GetSender().GetLogin()
			repo := data.GetRepo().GetFullName()
			issue := data.GetIssue().GetTitle()
			body := data.GetIssue().GetBody()
			log.Printf("New issue '%s' by %s on issue '%s' on repo %s", body, name, issue, repo)
		} else if _, ok := payload.(*github.IssueEvent); ok {
			log.Println("is IssueEvent")
		} else if _, ok := payload.(*github.CreateEvent); ok {
			log.Println("is CreateEvent")
		} else if _, ok := payload.(*github.User); ok {
			log.Println("is User")
		} else {
			log.Errorln("Cannot recognize payload. Printing")
			log.Println(payload)
		}
	})

	r.HandleFunc("/callback", func(rw http.ResponseWriter, req *http.Request) {
		log.Infoln("Trigger /callback")
		print(req)
	})

	r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Infoln("Trigger /")
		print(req)
	})

	panic(http.ListenAndServe(":9876", r))
}

func print(req *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	log.Println(buf.String())
	log.Println(req.Header.Get("x-hub-signature"))
	log.Println(req.Header.Get("x-github-event"))
	log.Println(req.Header.Get("x-github-delivery"))
}

var (
	appSecret = []byte(`c6a5dd0b083f14f547d18335cff34d84d3052a5c`)
)
