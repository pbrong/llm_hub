package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetWithHeader(url string, queryMap map[string]interface{}, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	url += "?"
	for key, value := range queryMap {
		url += fmt.Sprintf("%v=%v&", key, value)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("get with header, url = %v, response = %v", url, string(body))
	if err != nil {
		return nil, err
	}
	return body, nil
}
