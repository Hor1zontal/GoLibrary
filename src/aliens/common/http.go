package util

//http封装类
import (
	"aliens/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func HttpGet(paramUrl string) string {
	resp, err := http.Get(paramUrl)
	if err != nil {
		log.Error("%v", err)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("%v", err)
		return ""
	}

	return string(body)
}

func HttpPost(url string, param url.Values) string {
	resp, err := http.PostForm(url, param)
	if err != nil {
		log.Error("%v", err)
		return ""
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}

}

func HttpBodyPost(url string, body io.Reader) string {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Error("%v", err)
		return ""
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}

}

func HttpGetBody(request *http.Request) []byte {
	body := request.Body
	if body == nil {
		return []byte{}
	}
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}
	}
	return data
}
