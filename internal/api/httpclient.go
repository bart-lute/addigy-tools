package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func doRequest(method string, endPoint string, requestBody any, responseBody any) {

	// Requirements
	baseUrl := viper.GetString("api.url")
	timeout := time.Duration(viper.GetInt("api.timeout")) * time.Second
	apiKey := viper.GetString("api.key")

	// Convert the Request Body to a Reader Object
	rb, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	requestBodyReader := strings.NewReader(string(rb))

	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", baseUrl, endPoint), requestBodyReader)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-api-key", apiKey)

	httpClient := &http.Client{
		Timeout: timeout,
	}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("API error: %s", response.Status))
	}

	if err := json.Unmarshal(body, responseBody); err != nil {
		log.Fatal(err)
	}
}
