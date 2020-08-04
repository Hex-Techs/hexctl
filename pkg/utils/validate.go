package utils

import (
	"fmt"
	"net"
	"net/url"
)

func ValidataIP(ip ...string) error {
	for _, i := range ip {
		if err := net.ParseIP(i); err == nil {
			return fmt.Errorf("%s is not a legitimate IP address", i)
		}
	}
	return nil
}

func ValidataURL(urlSlice ...string) error {
	for _, i := range urlSlice {
		u, err := url.Parse(i)
		if err != nil {
			return err
		}
		if u.Scheme != "https" {
			return fmt.Errorf("scheme must https")
		}
		if u.Host == "" {
			return fmt.Errorf("not found host in url %s", i)
		}
	}
	return nil
}
