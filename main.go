package main

// libraries
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// setup json structure
type Response struct {
	Event           interface{}                       `json:"event"`
	EventType       interface{}                       `json:"event_type"`
	AppID           interface{}                       `json:"app_id"`
	UserID          interface{}                       `json:"user_id"`
	MessageID       interface{}                       `json:"message_id"`
	PageTitle       interface{}                       `json:"page_title"`
	PageURL         interface{}                       `json:"page_url"`
	BrowserLanguage interface{}                       `json:"browser_language"`
	ScreenSize      interface{}                       `json:"screen_size"`
	Attributes      map[string]map[string]interface{} `json:"attributes"`
	Traits          map[string]map[string]interface{} `json:"traits"`
}

// Initiate the processing and adapt the JSON according to the specified criteria.
func worker(requests <-chan map[string]interface{}) {

	// establish the request json within for loop
	for req := range requests {

		var response Response
		attributes := make(map[string]map[string]interface{})
		traits := make(map[string]map[string]interface{})

		// Extract the key and value separately and verify if the specified key exists. If the key is present, modify it to the desired output key.
		for k, v := range req {

			if k == "ev" {
				response.Event = v

			} else if k == "et" {
				response.EventType = v

			} else if k == "id" {
				response.AppID = v

			} else if k == "uid" {
				response.UserID = v

			} else if k == "mid" {
				response.MessageID = v

			} else if k == "t" {
				response.PageTitle = v

			} else if k == "p" {
				response.PageURL = v

			} else if k == "l" {
				response.BrowserLanguage = v

			} else if k == "sc" {
				response.ScreenSize = v

			} else {
				if strings.HasPrefix(k, "atrk") {
					attrk := req[k].(string)

					atrv_value := "atrv" + strings.TrimPrefix(k, "atrk")
					atrt_value := "atrt" + strings.TrimPrefix(k, "atrk")

					value := req[atrv_value]
					types := req[atrt_value]

					attributes[attrk] = map[string]interface{}{
						"value": value,
						"type":  types,
					}

					response.Attributes = attributes
				} else if strings.HasPrefix(k, "uatrk") {
					traitk := req[k].(string)
					traitValue := req["uatrv"+strings.TrimPrefix(k, "uatrk")]
					traitType := req["uatrt"+strings.TrimPrefix(k, "uatrk")]

					traits[traitk] = map[string]interface{}{
						"value": traitValue,
						"type":  traitType,
					}

					response.Traits = traits
				}
			}

		}

		//Marshal to response to convert struct to bytes
		responseBytes, err := json.Marshal(response)

		if err != nil {

			return
		}

		// Utilize the API_KEY and webhook URL to upload the JSON data to the webhook website
		API_KEY := "142261c0-1a7b-41c1-be3e-c2ac30d43f3c"
		url := "https://webhook.site/" + API_KEY

		res, err := http.Post(url, "application/json", bytes.NewBuffer(responseBytes))

		if err != nil {
			fmt.Println(err)
			return

		}

		defer res.Body.Close()

		fmt.Println("Webhook response status:", res.Status)

	}

}

func main() {
	// Created a HTTP server
	http.HandleFunc("/post_json", uploadJsonRequest)
	http.ListenAndServe(":8080", nil)

}

func uploadJsonRequest(w http.ResponseWriter, r *http.Request) {

	requestsChannel := make(chan map[string]interface{})

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Created a Golang channel to send this request to a golang worker
	go worker(requestsChannel)
	requestsChannel <- data

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
