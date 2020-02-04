// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package untappd_test

import (
	"net/http"
)

type mockHTTPClient struct {
	GetCalled chan bool
	GetInput  struct {
		Url chan string
	}
	GetOutput struct {
		R   chan *http.Response
		Err chan error
	}
	DoCalled chan bool
	DoInput  struct {
		Req chan *http.Request
	}
	DoOutput struct {
		R   chan *http.Response
		Err chan error
	}
}

func newMockHTTPClient() *mockHTTPClient {
	m := &mockHTTPClient{}
	m.GetCalled = make(chan bool, 100)
	m.GetInput.Url = make(chan string, 100)
	m.GetOutput.R = make(chan *http.Response, 100)
	m.GetOutput.Err = make(chan error, 100)
	m.DoCalled = make(chan bool, 100)
	m.DoInput.Req = make(chan *http.Request, 100)
	m.DoOutput.R = make(chan *http.Response, 100)
	m.DoOutput.Err = make(chan error, 100)
	return m
}
func (m *mockHTTPClient) Get(url string) (r *http.Response, err error) {
	m.GetCalled <- true
	m.GetInput.Url <- url
	return <-m.GetOutput.R, <-m.GetOutput.Err
}
func (m *mockHTTPClient) Do(req *http.Request) (r *http.Response, err error) {
	m.DoCalled <- true
	m.DoInput.Req <- req
	return <-m.DoOutput.R, <-m.DoOutput.Err
}