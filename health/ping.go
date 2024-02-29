package healthchecks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type Ping struct {
	URL     string
	Method  string
	Timeout time.Duration
	client  *fasthttp.Client
	Body    interface{}
	Headers map[string]string
}

// NewPingChecker : time - millisecond
func NewPingChecker(URL, method string, timeout time.Duration, body interface{}, headers map[string]string) *Ping {
	if method == "" {
		method = "GET"
	}

	if timeout == 0 {
		timeout = 500
	}

	pingChecker := Ping{
		URL:     URL,
		Method:  method,
		Timeout: timeout,
		Body:    body,
		Headers: headers,
	}
	pingChecker.client = &fasthttp.Client{}

	return &pingChecker
}

func (p Ping) Check(name string) Integration {
	var (
		start        = time.Now()
		status       = true
		errorMessage = ""
	)

	jsonBody, err := json.Marshal(p.Body)
	if err != nil {
		status = false
		errorMessage = fmt.Sprintf("request failed: %s -> %s with error: %s", p.Method, p.URL, err)
		return Integration{
			Name:         name,
			Error:        errorMessage,
			Status:       status,
			ResponseTime: time.Since(start).Milliseconds(),
		}
	}

	byteBody := []byte(jsonBody)
	bodyReader := bytes.NewReader(byteBody)

	req, err := http.NewRequest(p.Method, p.URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range p.Headers {
		req.Header.Set(key, value)
	}

	client := http.Client{
		Timeout: p.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 500 {
		status = false
		errorMessage = fmt.Sprintf("request failed: %s -> %s with error: %s", p.Method, p.URL, err)
		return Integration{
			Name:         name,
			Error:        errorMessage,
			Status:       status,
			ResponseTime: time.Since(start).Milliseconds(),
		}
	}

	return Integration{
		Name:         name,
		Status:       status,
		ResponseTime: time.Since(start).Milliseconds(),
		Error:        errorMessage,
	}

}
