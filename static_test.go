package forge_test

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/forge"
)

func Test_Static_Success(t *testing.T) {
	static := &forge.Static{
		FileSystem: http.Dir("./test_files"),
	}

	request, _ := http.NewRequest(http.MethodGet, "/success.txt", nil)

	handlerTest(t, handlerTestCase{
		Handler:          static,
		Request:          request,
		TargetStatusCode: http.StatusOK,
		TargetBody:       "Success!\n",
	})
}

func Test_Static_DirectoryWithIndex(t *testing.T) {
	static := &forge.Static{
		FileSystem: http.Dir("./test_files"),
	}

	request, _ := http.NewRequest(http.MethodGet, "/directory_with_index", nil)

	handlerTest(t, handlerTestCase{
		Handler:          static,
		Request:          request,
		TargetStatusCode: http.StatusOK,
		TargetBody:       "Directory With Index\n",
	})
}

func Test_Static_DirectoryNoIndex(t *testing.T) {
	static := &forge.Static{
		FileSystem: http.Dir("./test_files"),
	}

	request, _ := http.NewRequest(http.MethodGet, "/directory_no_index", nil)

	handlerTest(t, handlerTestCase{
		Handler:          static,
		Request:          request,
		TargetStatusCode: http.StatusNotFound,
		TargetBody:       forge.ResponseTextNotFound,
	})
}

func Test_Static_NotFound(t *testing.T) {
	static := &forge.Static{
		FileSystem: http.Dir("./test_files"),
	}

	request, _ := http.NewRequest(http.MethodGet, "/not-found", nil)

	handlerTest(t, handlerTestCase{
		Handler:          static,
		Request:          request,
		TargetStatusCode: http.StatusNotFound,
		TargetBody:       forge.ResponseTextNotFound,
	})
}

func Test_Static_CustomNotFound(t *testing.T) {
	customNotFoundResponse := "Custom Not Found"
	static := &forge.Static{
		FileSystem: http.Dir("./test_files"),
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(customNotFoundResponse))
		}),
	}

	request, _ := http.NewRequest(http.MethodGet, "/not-found", nil)

	handlerTest(t, handlerTestCase{
		Handler:          static,
		Request:          request,
		TargetStatusCode: http.StatusNotFound,
		TargetBody:       customNotFoundResponse,
	})
}
