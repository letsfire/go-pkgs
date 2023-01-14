package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func IsUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func UrlSuffix(url_ string) string {
	urlInfo, _ := url.Parse(url_)
	return strings.TrimLeft(filepath.Ext(urlInfo.Path), ".")
}

func UrlAdd(url_ string, args map[string]string) string {
	var values = make(url.Values)
	for key, value := range args {
		values.Add(key, value)
	}
	if strings.Contains(url_, "?") {
		return url_ + "&" + values.Encode()
	} else {
		return url_ + "?" + values.Encode()
	}
}

func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error, url=%v, statusCode=%v", url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func HttpPost(url string, args map[string]interface{}) ([]byte, error) {
	bts, _ := JsonUnEscape(args)
	buf := bytes.NewBuffer(bts)
	response, err := http.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error, url=%v, args=%+v, statusCode=%v", url, args, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func DownloadFile(url string, path string, force bool) error {
	if !force && FileExist(path) {
		return nil
	}
	if res, err := http.Get(url); err != nil {
		return err
	} else {
		defer res.Body.Close()
		file, err := os.Create(path)
		if err == nil {
			defer file.Close()
			_, err = io.Copy(file, res.Body)
		}
		return err
	}
}
