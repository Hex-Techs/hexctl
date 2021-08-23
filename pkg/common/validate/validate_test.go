package validate

import (
	"testing"
)

func TestValidateIP(t *testing.T) {
	var errs []error
	ipArr := []string{"192.168.1.2", "aa.bb.cc.dd", "192.168.1.256", "192.168.$.a"}
	for _, ip := range ipArr {
		errs = append(errs, ValidateIP(ip))
	}
	for i, e := range errs {
		if i == 0 {
			if e != nil {
				t.Error(e)
				return
			}
		} else {
			if e == nil {
				t.Errorf("%d is nil", i)
				return
			}
		}
	}
}

func TestValidateURL(t *testing.T) {
	url1 := "aabbcc"
	url2 := "http://sdfdf"
	url3 := "https://"
	url4 := "https://ab.com"
	var err error
	if err = ValidateURL(url1); err == nil {
		t.Errorf("%s is not a url", url1)
		return
	}
	if err = ValidateURL(url2); err == nil {
		t.Errorf("%s is not a url", url1)
		return
	}
	if err = ValidateURL(url3); err == nil {
		t.Errorf("%s is not a url", url1)
		return
	}
	if err = ValidateURL(url4); err != nil {
		t.Errorf("%s is a url", url1)
		return
	}
}
