package validate

import (
	"fmt"
	"net"
	"net/url"
)

// ValidateIP valedate the IP
func ValidateIP(ip ...string) error {
	for _, i := range ip {
		if err := net.ParseIP(i); err == nil {
			return fmt.Errorf("%s is not a legitimate IP address", i)
		}
	}
	return nil
}

// ValidateURL valedate the url scheme
func ValidateURL(urlSlice ...string) error {
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

// ValidateArgs 验证参数长度，获取指定index的值
func ValidateArgs(args []string, idx int) (string, error) {
	if idx == -1 {
		if len(args) != 0 {
			return "", fmt.Errorf("args must be empty")
		}
		return "", nil
	}
	if len(args) < idx+1 {
		return "", fmt.Errorf("args length is %d, but need %d", len(args), idx+1)
	}
	return args[idx], nil
}
