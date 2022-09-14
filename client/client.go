package client

import (
  "fmt"
	"io/ioutil"
	"net/http"
	"strings"
  "time"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// get proxer header values
	r.URL.Scheme = r.Header.Get("X-ProxyTo-Scheme")
	r.URL.Host = r.Header.Get("X-ProxyTo-Host")

	// set default schema
	if len(r.URL.Scheme) == 0 {
		r.URL.Scheme = "http"
	}

	// check if host is passed
	if len(r.URL.Host) == 0 {
		http.Error(w, "Target HOST must be defined in `X-ProxyTo-Host` header", 470)
		return
	}

	fmt.Println(fmt.Sprintf("[%s] >> Request: %s %s - To: %s", time.Now(), r.Method, r.URL.Path, r.URL.Scheme + "://" + r.URL.Host))

	// prepare request
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ct := r.Header.Get("Content-Type")

	if len(ct) == 0 {
		ct = "application/json"
	}

	req.Header.Set("Content-Type", ct)

	// set all the headers
	for k, v := range r.Header {
		req.Header.Set(k, strings.Join(v, ","))
	}

	// make request to PROXY
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rsp.Body.Close()

	// read the data
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// set headers
	for k, v := range rsp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	
	// set status
	w.WriteHeader(rsp.StatusCode)
	
	// write the data
	w.Write(b)
}

func New() http.Handler {
	return &Handler{}
}
