package main

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type fakeRes struct {
	b []byte
}

func (f fakeRes) Header() http.Header {
	return map[string][]string{}
}
func (f *fakeRes) Write(data []byte) (int, error) {
	b := make([]byte, len(data))
	copy(b, data)
	f.b = b
	return len(data), nil
}
func (f fakeRes) WriteHeader(statusCode int) {
}

func Test_memoryTotal(t *testing.T) {
	/*
		fake req res
	*/
	req := http.Request{
		URL: &url.URL{
			RawQuery: "name=impress",
		},
	}
	var f fakeRes
	cacheRate(&f, &req)
	if strings.Count(string(f.b), "\n") != 9 { // 6 nodes + 3 lines
		t.Fatal("cacheRate error")
	}
	slowlog(&f, &req)
	if strings.Count(string(f.b), "\n") != 9 { // 6 nodes + 3 lines
		t.Fatal("slowlog error")
	}
	memoryTotal(&f, &req)
	if strings.Count(string(f.b), "\n") != 9 { // 6 nodes + 3 lines
		t.Fatal("memoryTotal error")
	}
	opsTotal(&f, &req)
	if strings.Count(string(f.b), "\n") != 9 { // 6 nodes + 3 lines
		t.Fatal("memoryTotal error")
	}
}
