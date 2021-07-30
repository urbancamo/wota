package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var err error

func main() {

	var b []byte
	b, err = ioutil.ReadFile("sota-upload-form.html")
	if err != nil {
		fmt.Println(err)
	}
	s := string(b)
	s = strings.ReplaceAll(s, "$USER", "username")

	fmt.Println(s)
}
