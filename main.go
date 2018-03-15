package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	port := os.Getenv("PORT")
	fmt.Println(port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query()["url"]) != 1 {
			http.Error(w, "url param missing.", http.StatusBadRequest)
			return
		}

		url := r.URL.Query()["url"][0]
		log.Println(url)
		log.Println(r.Method)

		if r.Method == http.MethodOptions {
			headers := w.Header()
			headers.Add("Access-Control-Allow-Origin", "*")
			headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, Authorization")
			headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			return
		}

		newRequest, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to make request.", http.StatusInternalServerError)
			return
		}

		// copy the headers of the request
		for header, index := range r.Header {
			for _, val := range index {
				// fmt.Println(header, val)
				newRequest.Header.Add(header, val)
			}
		}

		fmt.Println(r.Cookies())
		// copy cookies of the request
		for _, cookie := range r.Cookies() {
			fmt.Println("cookie - ", cookie)
			newRequest.AddCookie(cookie)
		}

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

		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Content-Type", res.Header.Get("Content-Type"))
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		// fmt.Println(r.Header.Get("Origin"))

		w.Write(body)
		// log.Println(string(body))

		// copy cookies to response
		for _, cookie := range res.Cookies() {
			http.SetCookie(w, cookie)
		}
	})

	log.Println(http.ListenAndServe(":"+port, nil))
}
