package http

import (
	"io/ioutil"
	"llm_hub/pkg/json"
	"log"
	"net/http"
	"strings"
)

func PostWithHeader(url string, bodyMap map[string]interface{}, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(json.SafeDump(bodyMap)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("post with header, url = %v, request = %v, response = %v", url, json.SafeDump(bodyMap), string(body))

	if err != nil {
		return body, err
	}
	return body, nil
}
