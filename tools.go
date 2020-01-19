package rboot

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
)

func getRbootRemoteIPv4() (net.IP, error) {
	res, err := http.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(body))
	if ip == nil {
		return nil, errors.New("invalid ip address")
	}
	return ip.To4(), nil
}
