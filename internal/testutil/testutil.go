package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// ActualResponse return recorded response
func ActualResponse(t *testing.T, router *gin.Engine,
	method, url string, body io.Reader) *http.Response {
	httpRecorder := httptest.NewRecorder()

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	router.ServeHTTP(httpRecorder, req)

	return httpRecorder.Result()
}

// ReqBodyFromInterface return request body that contain json payload.
func ReqBodyFromInterface(t *testing.T, body interface{}) *bytes.Buffer {
	jsonBytes, err := json.Marshal(body)
	require.NoError(t, err)

	return bytes.NewBuffer(jsonBytes)
}

// JSONStringFromInterface return json string by interface.
func JSONStringFromInterface(t *testing.T, res interface{}) string {
	bytes, err := json.Marshal(res)
	require.NoError(t, err)

	return string(bytes)
}

// JSONStringFromResBody return json string by http response body
func JSONStringFromResBody(t *testing.T, body io.Reader) string {
	bytes, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	return string(bytes)
}
