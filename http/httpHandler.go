package httpHandler

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func HttpGet(request string) ([]byte, error) {
	resp, err := http.Get(request)
	if err != nil {
		return nil, errors.New("http.get failed")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read http response failed")
	}

	return body, nil
}
