package models

import (
	"fmt"
	"net/url"
	"strings"
)

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func paif(err error) {
	if err != nil {
		panic(err)
	}
}

func Foo(src string, dist string) {
	r := url.QueryEscape(src)
	r = strings.Replace(r, "+", "%20", -1)
	if r != dist {
		fmt.Printf("ensrc:%s\ngo:%s\njs:%s\n\n", src, r, dist)
	}

	r, err := url.QueryUnescape(dist)
	paif(err)
	if r != src {
		fmt.Printf("desrc:%s\ngo:%s\njs:%s\n\n", src, r, dist)
	}
}
