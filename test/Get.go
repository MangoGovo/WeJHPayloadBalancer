package test

import (
	"io/ioutil"
	"net/http"
)

// Get 封装的 GET 请求函数
func Get(url string) (string, int) {
	response, err := http.Get(url)
	if err != nil {
		return "", http.StatusInternalServerError
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", http.StatusInternalServerError
	}

	return string(body), response.StatusCode
}
