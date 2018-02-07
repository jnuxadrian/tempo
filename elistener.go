package tempo

import (
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"
)

func init() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {

		eventType := req.Header.Get("X-GitHub-Event")

		if eventType == "" {
			http.Error(res, "X-GitHub-Event is missing", http.StatusBadRequest)
			return
		}

		switch eventType {
		case "push":
			handleEvent(res, req, handlePush)
		case "pull_request":
			handleEvent(res, req, handlePullRequest)
		default:
			http.Error(res, "X-GitHub-Event unknown", http.StatusBadRequest)
			return
		}

	})
}

func handleEvent(
	res http.ResponseWriter,
	req *http.Request,
	handler func(context.Context, string, http.ResponseWriter)) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body", http.StatusInternalServerError)
		return
	}

	handler(appengine.NewContext(req), string(body), res)
}

func handlePush(ctx context.Context, body string, res http.ResponseWriter) {
	log.Debugf(ctx, "got push event %v", body)
}

func handlePullRequest(ctx context.Context, body string, res http.ResponseWriter) {
	log.Debugf(ctx, "got pull_request event %v", body)
}
