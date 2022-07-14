package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func NewGithubClient(timeout int, token string) *GithubClient {
	return &GithubClient{
		client:  http.Client{},
		timeout: timeout,
		token:   token,
	}
}

type Filters map[string]string

type GithubClient struct {
	client  http.Client
	timeout int
	token   string
}


func (gc GithubClient) call(baseUrl string, output interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(gc.timeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Accept-Charset", "UTF-8")
	if len(gc.token) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("token %s",gc.token))
	}

	res, err := gc.client.Do(req)
	if err != nil {
		return err
	}

	if res == nil {
		return fmt.Errorf("response is nil")
	}
	if (res.StatusCode / 100) != 2 {
		return fmt.Errorf(res.Status)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(resBody, output)
}
