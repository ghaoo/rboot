package gorc

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strconv"
)

func sendGet(url string, address string, start int64, end int64) (len int64, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10))
	req.Header.Set("Connection", "close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	var resp *http.Response
	resp, err = client.Do(req)
	defer resp.Body.Close()
	file, err := createFile(address)
	len, err = io.Copy(file, resp.Body)
	return len, err
}

func sendHead(url string) (l string, support bool, err error) {
	var req *http.Request
	req, err = http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Println("create HEAD failed")
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		log.Println("HEAD response failed")
		return
	}
	defer resp.Body.Close()
	l = resp.Header.Get("Content-Length")
	s := resp.Header.Get("Accept-Ranges")
	if s != "" {
		support = true
	}
	return
}
