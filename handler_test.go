package publicerror_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	publicerror "github.com/xeoncross/public-error-go"
)

var errStartingItAll = errors.New("real error msg")
var publicErrorMessage = "sorry, something went wrong"

// This test demonstrates how publicerror might be used in an http.Handler to
// provide pre-defined public error messages to clients without loss of the
// original error chain

// maybe each app can implement this so they can add whatever metrics or logging they want
type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		// todo: add r.Context trace id to logs

		// Report all errors and nested errors
		log.Printf("HTTP %d - %s", publicerror.StatusCode(err), err)

		// Show the client a "safe" error or http.StatusInternalServerError
		http.Error(w, publicerror.Message(err), publicerror.StatusCode(err))

		// Extra: verify we can still check the error type and didn't loose the
		// power of .Is() and .Unwrap()!
		if errors.Is(err, errStartingItAll) == false {
			panic("unexpected!")
		}
	}
}

// Lets create a handler that tries to load the user or something
func myhandler(w http.ResponseWriter, r *http.Request) error {
	err := fetchUser()
	if err != nil {
		return fmt.Errorf("myhandler: %s: %w", r.URL.String(), err)
	}
	return nil
}

// the database call fails and we want to show the real error in the log, but
// the "user safe" error to the client
func fetchUser() error {
	return publicerror.Error{
		Err:        errStartingItAll,
		Message:    publicErrorMessage,
		StatusCode: http.StatusBadRequest,
	}
}

func TestHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("/", appHandler(myhandler))
	mux.ServeHTTP(rr, req)

	t.Logf("HTTP %d: %s", rr.Result().StatusCode, rr.Body)

	if rr.Body.String() != publicErrorMessage+"\n" {
		t.Errorf("unexpected response body: %q", rr.Body.String())
	}
}
