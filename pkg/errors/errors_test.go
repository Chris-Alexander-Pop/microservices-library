package errors_test

import (
	"errors"
	"net/http"
	"testing"

	appErrors "github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
)

type ErrorsSuite struct {
	*test.Suite
}

func TestErrorsSuite(t *testing.T) {
	test.Run(t, &ErrorsSuite{Suite: test.NewSuite()})
}

func (s *ErrorsSuite) TestAppError() {
	originalErr := errors.New("database connection failed")

	// Test New Wrapper
	e := appErrors.New(appErrors.CodeInternal, "Something went wrong", originalErr)

	s.Equal(appErrors.CodeInternal, e.Code)
	s.Equal("Something went wrong", e.Message)
	s.Equal(originalErr, e.Err)
	// Update expected error string format: [CODE] Message: Err
	s.Equal("[INTERNAL] Something went wrong: database connection failed", e.Error())

	// Test Unwrap
	s.Equal(originalErr, errors.Unwrap(e))
}

func (s *ErrorsSuite) TestHelpers() {
	err := errors.New("oops")

	notFound := appErrors.NotFound("Not Found", err)
	s.Equal(appErrors.CodeNotFound, notFound.Code)
	s.Equal(http.StatusNotFound, appErrors.HTTPStatus(notFound))

	badReq := appErrors.InvalidArgument("Bad Request", err)
	s.Equal(appErrors.CodeInvalidArgument, badReq.Code)
	s.Equal(http.StatusBadRequest, appErrors.HTTPStatus(badReq))
}
