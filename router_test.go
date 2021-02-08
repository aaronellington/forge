package forge_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/fuzzingbits/forge"
)

type handlerTestCase struct {
	Handler               http.Handler
	Request               *http.Request
	TargetStatusCode      int
	TargetBody            string
	CustomResponseChecker func(t *testing.T, response *http.Response)
}

func Test_Router_Success(t *testing.T) {
	testResponse := "Test Response"

	router := &forge.Router{}
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testResponse))
	}))

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          router,
		Request:          request,
		TargetStatusCode: http.StatusOK,
		TargetBody:       testResponse,
	})
}

func Test_Router_NotFound(t *testing.T) {
	router := &forge.Router{}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          router,
		Request:          request,
		TargetStatusCode: http.StatusNotFound,
		TargetBody:       forge.ResponseTextNotFound,
	})
}
func Test_Router_CustomNotFound(t *testing.T) {
	testResponse := "Custom Not Found"

	router := &forge.Router{
		NotFoundHander: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(testResponse))
		}),
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          router,
		Request:          request,
		TargetStatusCode: http.StatusOK,
		TargetBody:       testResponse,
	})
}

func Test_Router_StripTrailingSlash(t *testing.T) {
	testResponse := "Test Response"
	testPath := "/hello"

	router := &forge.Router{}
	router.Handle(testPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testResponse))
	}))

	request, _ := http.NewRequest(http.MethodGet, testPath+"/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          router,
		Request:          request,
		TargetStatusCode: http.StatusOK,
		TargetBody:       testResponse,
	})
}

func Test_RespondHTML(t *testing.T) {
	type customType struct {
		Success bool
	}

	targetBody := "<b>Bold!</b>"
	router := &forge.Router{}
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forge.RespondHTML(w, http.StatusInternalServerError, []byte(targetBody))
	}))

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          router,
		Request:          request,
		TargetStatusCode: http.StatusInternalServerError,
		TargetBody:       targetBody,
	})
}

func Test_RespondJSON(t *testing.T) {
	type customType struct {
		Success bool
	}

	router := &forge.Router{}
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forge.RespondJSON(w, http.StatusInternalServerError, customType{Success: true})
	}))

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          router,
		Request:          request,
		TargetStatusCode: http.StatusInternalServerError,
		TargetBody:       "{\"Success\":true}\n",
	})
}

func handlerTest(t *testing.T, testCase handlerTestCase) {
	server := httptest.NewServer(testCase.Handler)
	defer server.Close()

	serverURL, _ := url.Parse(server.URL)
	testCase.Request.URL.Host = serverURL.Host
	testCase.Request.URL.Scheme = serverURL.Scheme

	response, err := http.DefaultClient.Do(testCase.Request)
	if err != nil {
		t.Fatalf(
			"Request Failed: %s",
			err,
		)
	}
	defer response.Body.Close()

	if testCase.CustomResponseChecker != nil {
		testCase.CustomResponseChecker(t, response)
		return
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf(
			"Body Read Failed: %s",
			err,
		)
	}

	if !reflect.DeepEqual(string(responseBytes), testCase.TargetBody) {
		t.Fatalf(
			"response body: %s expected: %s",
			responseBytes,
			testCase.TargetBody,
		)
	}
}
