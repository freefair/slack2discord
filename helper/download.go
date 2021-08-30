package helper

import (
	"io/ioutil"
	"net/http"
)

func HTTPDownload(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
	}
	return d, err
}

func HTTPDownloadWithToken(uri string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + token)
	client := http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
	}
	return d, err
}
