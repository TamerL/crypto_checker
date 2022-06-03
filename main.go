package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/", Handler)
	_ = http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	endpoint := strings.ToUpper(strings.Split(r.URL.String(), "/")[1])
	if endpoint == "USD" || endpoint == "EUR" || endpoint == "GBP" || endpoint == "JPY" {
		response, err := http.Get(fmt.Sprintf("https://api.coinbase.com/v2/prices/spot?currency=%s", endpoint))
		if err != nil {
			fmt.Println("The backend API failed with error $s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		}
	} else if endpoint == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Please append endpoint: USD or EUR or GBP or JPY"))
	} else if endpoint == "HEALTH" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Application is OK"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}
