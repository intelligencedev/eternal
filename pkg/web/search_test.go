package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRequest(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := r.ParseForm()
		require.NoError(t, err)

		assert.Equal(t, "test query", r.Form.Get("q"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response"))
	}))
	defer ts.Close()

	// Test the postRequest function
	response, err := postRequest(ts.URL, "test query")
	require.NoError(t, err)
	assert.Equal(t, "Test response", response)
}

func TestExtractURLs(t *testing.T) {
	htmlContent := `
		<html>
			<body>
				<a href="http://example.com">Example</a>
				<a href="https://test.com">Test</a>
				<a href="/relative">Relative</a>
			</body>
		</html>
	`

	urls, err := extractURLs(htmlContent)
	require.NoError(t, err)
	assert.Equal(t, []string{"http://example.com", "https://test.com"}, urls)
}

func TestGetSearXNGResults(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
            <html>
                <body>
                    <a href="http://example.com">Example</a>
                    <a href="https://test.com">Test</a>
                    <a href="/relative">Relative</a>
                </body>
            </html>
        `))
	}))
	defer ts.Close()

	results := GetSearXNGResults(ts.URL, "test query")

	// test.com is an 'unwanted' URL, so it should not be included in the results
	assert.Equal(t, []string{"http://example.com"}, results)
}

func TestGetSearXNGResultsError(t *testing.T) {
	// Create a test server that returns an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	results := GetSearXNGResults(ts.URL, "test query")
	assert.Nil(t, results)
}
