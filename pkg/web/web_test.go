package web

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHttpClient is a mock implementation of the HttpClient.
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Post(url string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Head(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Options(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Put(url string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Patch(url string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Delete(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) CloseIdleConnections() {
	m.Called()
}

func (m *MockHttpClient) Transport() http.RoundTripper {
	args := m.Called()
	return args.Get(0).(http.RoundTripper)
}

func (m *MockHttpClient) Jar() http.CookieJar {
	args := m.Called()
	return args.Get(0).(http.CookieJar)
}
