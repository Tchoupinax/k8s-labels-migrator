package utils

import (
	"fmt"
	"regexp"
)

func CheckLabelKey(key string) error {
	r, _ := regexp.Compile("^[a-z]{1}[a-z./-]*[a-z]{1}$")
	if r.MatchString(key) {
		return nil
	} else {
		return fmt.Errorf("the label key does not respect this regex: ^[a-z]{1}[a-z./-]*[a-z]{1}$")
	}
}

func CheckLabelValue(key string) error {
	r, _ := regexp.Compile("^[a-z-]*$")
	if r.MatchString(key) {
		return nil
	} else {
		return fmt.Errorf("the label value does not respect this regex: ^[a-z-]*$")
	}
}
