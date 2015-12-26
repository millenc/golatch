package golatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LatchAPI struct {
	Proxy             *url.URL
	OnRequestStart    func(request *LatchRequest)
	OnResponseReceive func(request *LatchRequest, response *http.Response, responseBody string)
}

func (l *LatchAPI) DoRequest(request *LatchRequest, responseType LatchResponse) (response *LatchResponse, err error) {
	var client *http.Client
	var resp *http.Response
	var body []byte

	//Initialize the client
	client = &http.Client{}
	if l.Proxy != nil {
		client.Transport = &http.Transport{Proxy: http.ProxyURL(l.Proxy)}
	}

	//Perform the request
	req := request.GetHttpRequest()

	if l.OnRequestStart != nil {
		l.OnRequestStart(request)
	}
	if resp, err = client.Do(req); err != nil {
		return
	}

	//Get the response's body
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if l.OnResponseReceive != nil {
		l.OnResponseReceive(request, resp, string(body))
	}

	//Handle HTTP errors
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("HTTP error [%d] body: %s", resp.StatusCode, body))
		return
	}

	//Check if the response is an error before decoding it
	latch_error_response := &LatchErrorResponse{}
	if err = json.Unmarshal(body, latch_error_response); err != nil {
		return
	} else if (*latch_error_response).Err.Code != 0 {
		err = &latch_error_response.Err
		return
	}

	//Decode response into a typed response (if one has been specified)
	if responseType != nil {
		err = responseType.Unmarshal(string(body))
		response = &responseType
	}

	return response, err
}

//Sets the proxy URL to be used in all requests to the API
func (l *Latch) SetProxy(proxyURL *url.URL) {
	l.Proxy = proxyURL
}

//Gets the complete url for a request
func GetLatchURL(queryString string) *url.URL {
	latch_url, err := (&url.URL{}).Parse(fmt.Sprint(API_URL, API_PATH, "/", API_VERSION, "/", queryString))
	if err != nil {
		latch_url = &url.URL{}
	}

	return latch_url
}
