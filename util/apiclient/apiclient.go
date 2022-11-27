package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Host string
	*http.Client
}

func (c *Client) Do(method, path string, payload, dst interface{}) error {
	encodedPayload, err := json.Marshal(payload)

	if err != nil {
		return errors.New("cannot marshal request")
	}

	request, err := http.NewRequest(method, c.Host+path, bytes.NewBuffer(encodedPayload))

	if err != nil {
		return errors.New("cannot create request")
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(request)
	if err != nil {
		return errors.New("cannot process request: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("cannot execute request: code %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	return json.Unmarshal(data, &dst)
}
