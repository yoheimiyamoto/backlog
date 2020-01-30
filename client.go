package backlog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	get(path string, params url.Values) ([]byte, error)
	put(path string, params url.Values) ([]byte, error)
	// post(path string, body []byte) ([]byte, error)
	// put(path string, body []byte) ([]byte, error)
	// delete(path string, body []byte) ([]byte, error)
}

type client struct {
	apiKey       string
	endpointBase *url.URL
	httpClient   *http.Client
}

func newClient(subdomain, apiKey string, httpClient *http.Client) *client {
	c := client{
		apiKey: apiKey,
	}

	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
		// c.httpClient = &http.Client{Timeout: time.Duration(6000)}
	}

	u, _ := url.ParseRequestURI(fmt.Sprintf(APIEndpointBase, subdomain))
	c.endpointBase = u

	return &c
}

func (c *client) get(path string, query url.Values) ([]byte, error) {
	url := c.newURL(path, query)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

func (c *client) put(path string, params url.Values) ([]byte, error) {
	url := c.newURL(path, nil)
	reader := strings.NewReader(params.Encode())

	req, err := http.NewRequest("PATCH", url.String(), reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

func (c *client) do(req *http.Request) ([]byte, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

func (c *client) newURL(path string, query url.Values) *url.URL {
	u := *c.endpointBase
	u.Path = path
	params := url.Values{"apiKey": {c.apiKey}}
	if query != nil {
		for k, v := range query {
			params.Add(k, v[0])
		}
	}
	u.RawQuery = params.Encode()
	return &u
}
