package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"url_shortener_go/src/url"
)

var (
	port    int
	baseUrl string
	stats   chan string
)

type Headers map[string]string

type Redirect struct {
	stats chan string
}

func init() {
	port = 8888
	baseUrl = fmt.Sprintf("http://localhost:%d", port)
}

func main() {
	stats = make(chan string)
	defer close(stats)
	go registerStatistics(stats)

	http.HandleFunc("/api/shortener", Shortener)
	http.HandleFunc("/api/stats", ShowStats)
	http.Handle("/r/", &Redirect{stats})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func Shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		answerWith(w, http.StatusMethodNotAllowed, Headers{
			"Allow": "POST",
		})
		return
	}

	url, newUrl, err := url.LookOrCreateNewUrl(extractUrl(r))

	if err != nil {
		answerWith(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if newUrl {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	urlShortened := fmt.Sprintf("%s/r/%s", baseUrl, url.Id)

	answerWith(w, status, Headers{"Location": urlShortened, "Link": fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", baseUrl, url.Id)})
}

func (red *Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if url := url.LookUp(id); url != nil {
		http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
		red.stats <- id
	} else {
		http.NotFound(w, r)
	}
}

func ShowStats(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if url := url.LookUp(id); url != nil {
		json, err := json.Marshal(url.Stats())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		answerWithJSON(w, string(json))
	} else {
		http.NotFound(w, r)
	}
}

func answerWithJSON(
	w http.ResponseWriter,
	answer string,
) {
	answerWith(w, http.StatusOK, Headers{
		"Content-Type": "application/json",
	})
	fmt.Fprint(w, answer)
}

func registerStatistics(ids <-chan string) {
	for id := range ids {
		url.RegisterClick(id)
		fmt.Printf("Click registered successfully for %s", id)
	}
}

func answerWith(
	w http.ResponseWriter,
	status int,
	headers Headers,
) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}

func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}
