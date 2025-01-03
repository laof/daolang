package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type Rule struct {
	Prefix   string
	Protocol string
}

func getUrl(str string) string {

	ma := []Rule{
		{Prefix: "/0/", Protocol: ""},
		{Prefix: "/1/", Protocol: "http://"},
		{Prefix: "/2/", Protocol: "https://"},
	}

	for _, v := range ma {
		if strings.HasPrefix(str, v.Prefix) {
			return strings.Replace(str, v.Prefix, v.Protocol, 1)
		}
	}

	return ""
}

func forwardRequest(target string, w http.ResponseWriter, r *http.Request) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()

	forwardedRequest, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	for key, value := range r.Header {
		forwardedRequest.Header[key] = value
	}

	response, err := transport.RoundTrip(forwardedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		w.Header()[key] = value
	}
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RequestURI()

	if len(url) < 5 {
		fmt.Fprintf(w, "%s", "not found: "+url)
		return
	}

	url = getUrl(r.URL.RequestURI())

	if url == "" {
		fmt.Fprintf(w, "404 not found")
	} else {
		forwardRequest(url, w, r)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
