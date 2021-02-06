package forge_test

import (
	"net/http"
	"testing"

	"github.com/aaronellington/forge"
)

func Test_Security_Placeholder(t *testing.T) {
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
