package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		utils.LogWarning(fmt.Sprintf("%s [y/n]: ", s))

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func mapToArray(myMap map[string]string) []string {
	v := make([]string, 0, len(myMap))
	for key, value := range myMap {
		v = append(v, key+"="+value)
	}
	slices.Sort(v)
	return v
}

func arrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
