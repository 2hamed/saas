package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockDispatcher is a mock implementation of the Dispatcher interface
// it will return the values supplied to it in response to respective function calls
type mockDispatcher struct {
	enqueueErr error

	result    string
	resultErr error

	statusExists     bool
	statusIsFinished bool
	statusIsPending  bool
	statusErr        error
}

func (md mockDispatcher) Enqueue(url string) error {
	return md.enqueueErr
}
func (md mockDispatcher) FetchResult(url string) (string, error) {
	return md.result, md.resultErr
}
func (md mockDispatcher) FetchStatus(url string) (exists bool, isPending bool, isFinished bool, err error) {
	return md.statusExists, md.statusIsPending, md.statusIsFinished, md.statusErr
}

func TestNewJobHandlerSuccess(t *testing.T) {

	d := mockDispatcher{}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/new", strings.NewReader(url.Values{"urls": []string{"http://google.com;http://stackoverflow.com"}}.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(NewJobHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Result().StatusCode)

}

func TestNewJobHandlerEmptyUrl(t *testing.T) {

	d := mockDispatcher{}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/new", strings.NewReader(url.Values{"urls": []string{""}}.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(NewJobHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Result().StatusCode)

}

func TestNewJobHandlerEnqueueError(t *testing.T) {

	d := mockDispatcher{
		enqueueErr: errors.New("arbitrary error"),
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/new", strings.NewReader(url.Values{"urls": []string{"http://google.com"}}.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(NewJobHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Result().StatusCode)

}

func TestGetResultHandlerSuccess(t *testing.T) {

	d := mockDispatcher{
		result: "",

		statusExists:     true,
		statusIsFinished: true,
		statusErr:        nil,
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/result/aHR0cDovL2dvb2dsZS5jb20gLW4K", nil)

	handler := http.HandlerFunc(GetResultHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Result().StatusCode)

}

func TestGetResultHandlerNotFound(t *testing.T) {

	d := mockDispatcher{
		result: "",

		statusExists:     false,
		statusIsFinished: false,
		statusErr:        nil,
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/result/aHR0cDovL2dvb2dsZS5jb20gLW4K", nil)

	handler := http.HandlerFunc(GetResultHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Result().StatusCode)

}

func TestGetResultHandlerIsPending(t *testing.T) {

	d := mockDispatcher{
		result: "",

		statusExists:     true,
		statusIsFinished: false,
		statusIsPending:  true,
		statusErr:        nil,
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/result/aHR0cDovL2dvb2dsZS5jb20gLW4K", nil)

	handler := http.HandlerFunc(GetResultHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 204, rec.Result().StatusCode)

}

func TestGetResultHandlerJobFailed(t *testing.T) {

	d := mockDispatcher{
		result: "",

		statusExists:     true,
		statusIsFinished: false,
		statusIsPending:  false,
		statusErr:        nil,
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/result/aHR0cDovL2dvb2dsZS5jb20gLW4K", nil)

	handler := http.HandlerFunc(GetResultHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 501, rec.Result().StatusCode)

}

func TestGetResultHandlerStatusErr(t *testing.T) {

	d := mockDispatcher{
		result: "",

		statusExists:     true,
		statusIsFinished: true,
		statusErr:        errors.New("arbitrary error"),
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/result/aHR0cDovL2dvb2dsZS5jb20gLW4K", nil)

	handler := http.HandlerFunc(GetResultHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Result().StatusCode)

}

func TestGetResultHandlerResultErr(t *testing.T) {

	d := mockDispatcher{
		result:    "",
		resultErr: errors.New("arbitrary error"),

		statusExists:     true,
		statusIsFinished: true,
		statusErr:        nil,
	}
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/result/aHR0cDovL2dvb2dsZS5jb20gLW4K", nil)

	handler := http.HandlerFunc(GetResultHandler(d))

	handler.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Result().StatusCode)

}
