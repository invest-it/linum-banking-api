package nordigen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type headerMap map[string]string

func requestWithHeaders(method string, url string, headers headerMap, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func requestWithAuthorization(method string, url string, token *tokenInfo, body io.Reader, headers headerMap) (*http.Response, error) {
	if headers == nil {
		headers = make(headerMap)
	}
	headers["Authorization"] = "Bearer " + token.Access
	return requestWithHeaders(method, url, headers, body)
}

func getWithAuthorization(url string, token *tokenInfo) (*http.Response, error) {
	return requestWithAuthorization(http.MethodGet, url, token, nil, nil)
}

func mapResponseToStruct[T any](resp *http.Response) (*T, error) {
	var res T
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func getAndMapWithAuthorization[T any](url string, token *tokenInfo) (*T, error) {
	resp, err := getWithAuthorization(url, token)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return mapResponseToStruct[T](resp)
}

func postWithAuthorization(url string, token *tokenInfo, body io.Reader, contentType string) (*http.Response, error) {
	headers := make(headerMap)
	headers["Content-Type"] = contentType
	return requestWithAuthorization(http.MethodPost, url, token, body, headers)
}

func postAndMapWithAuthorization[T any](url string, token *tokenInfo, body io.Reader, contentType string) (*T, error) {
	resp, err := http.Post(url, contentType, body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		// TODO: Log server response
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	return mapResponseToStruct[T](resp)
}
