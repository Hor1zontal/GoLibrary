/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import (
	"aliens/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

func HttpGet(paramUrl string) []byte {
	resp, err := http.Get(paramUrl)
	if err != nil {
		log.Error("%v", err)
		return []byte{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("%v", err)
	}
	return body
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
