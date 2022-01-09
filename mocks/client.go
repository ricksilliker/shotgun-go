package mocks

import "net/http"

type MockShotgunClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockShotgunClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}
