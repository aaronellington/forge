package forge_test

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/forge"
)

func Test_Router_Placeholder(t *testing.T) {
	security := &forge.Security{
		Handler: &forge.Static{FileSystem: http.Dir("./test_files")},
	}

	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	handlerTest(t, handlerTestCase{
		Handler:          security,
		Request:          request,
		TargetStatusCode: http.StatusNotFound,
		TargetBody:       forge.ResponseTextNotFound,
	})
}
