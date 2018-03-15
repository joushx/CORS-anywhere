package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()["url"][0]
		log.Println(url)
		log.Println(r.Method)

		newRequest, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to make request.", http.StatusInternalServerError)
			return
		}
		copyToRequest(r.Header, r.Cookies(), newRequest)

		res, err := client.Do(newRequest)
		if err != nil {
			log.Println(err)
			http.Error(w, "Request failed", http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to read data from response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", res.Header.Get("Content-Type"))

		w.Write(body)

		// copy cookies to response
		for _, cookie := range res.Cookies() {
			http.SetCookie(w, cookie)
		}
	})

	log.Println(http.ListenAndServe(":8080", nil))
}

func copyToRequest(h http.Header, c []*http.Cookie, r *http.Request) {
	// copy the headers of the request
	for header, index := range h {
		for _, val := range index {
			r.Header.Add(header, val)
		}
	}

	// copy cookies to request
	for _, cookie := range c {
		r.AddCookie(cookie)
	}
}
