package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	// valid endpoints
	// valid currencies
	if endpoint == "USD" || endpoint == "EUR" || endpoint == "GBP" || endpoint == "JPY" {
		response, err := http.Get(fmt.Sprintf("https://api.coinbase.com/v2/prices/spot?currency=%s", endpoint))

		//coinbase response error handling
		if err != nil {
			fmt.Println("coinbase api failed with error $s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			if !isJSON(string(data)) {
				w.WriteHeader(http.StatusInternalServerError)
				err := errors.New("can not be parse coinbase response")
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(data))
		}

		// creating 'health' endpoint
	} else if endpoint == "HEALTH" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"alive": true}`)
		// return 'not found' for empty calls without currency
	} else if endpoint == "" {
		w.WriteHeader(http.StatusNotFound)
		err := errors.New("please append endpoint: USD or EUR or GBP or JPY")
		http.Error(w, err.Error(), 404)
		// return 'method not allowed' for wrong endpoints
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := errors.New("method not allowed")
		http.Error(w, err.Error(), 405)
		// w.Write([]byte("Method Not Allowed"))
	}
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
