package forge_test

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/forge"
)

func Test_Logger_Placeholder(t *testing.T) {
	logger := &forge.Logger{
		Handler: &forge.Static{FileSystem: http.Dir("./test_files")},
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          logger,
		Request:          request,
		TargetStatusCode: http.StatusNotFound,
		TargetBody:       forge.ResponseTextNotFound,
	})
}
