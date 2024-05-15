package utils

import (
	"testing"
)

func TestMatchOneLabel(t *testing.T) {
	matchLabels := map[string]string{"app": "applicationName"}
	appLabels := map[string]string{"app": "applicationName"}

	result := IsMatchSelectorsInclude(appLabels, matchLabels)
	if result != true {
		t.Error("Result was incorrect")
	}
}

func TestMatchOneLabelAmongSeveral(t *testing.T) {
	matchLabels := map[string]string{"app": "applicationName"}
	appLabels := map[string]string{"api": "prod", "app": "applicationName"}

	result := IsMatchSelectorsInclude(appLabels, matchLabels)
	if result != true {
		t.Error("Result was incorrect")
	}
}

func TestIncorrectMatchKey(t *testing.T) {
	matchLabels := map[string]string{"app": "applicationName"}
	appLabels := map[string]string{"name": "applicationName"}

	result := IsMatchSelectorsInclude(appLabels, matchLabels)
	if result != false {
		t.Error("Result was incorrect")
	}
}

func TestIncorrectMatchValue(t *testing.T) {
	matchLabels := map[string]string{"app": "applicationName2"}
	appLabels := map[string]string{"app": "applicationName"}

	result := IsMatchSelectorsInclude(appLabels, matchLabels)
	if result != false {
		t.Error("Result was incorrect")
	}
}
