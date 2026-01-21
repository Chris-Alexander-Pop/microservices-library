package tests

import (
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/security/captcha"
	"github.com/chris-alexander-pop/system-design-library/pkg/security/captcha/adapters/memory"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
)

type CaptchaTestSuite struct {
	test.Suite
	verifier captcha.Verifier
}

func (s *CaptchaTestSuite) SetupTest() {
	s.Suite.SetupTest()
	s.verifier = memory.New("valid-token")
}

func (s *CaptchaTestSuite) TestVerify() {
	// Memory adapter usually allows any token or specific ones
	err := s.verifier.Verify(s.Ctx, "valid-token")
	s.NoError(err)
}

func TestCaptchaSuite(t *testing.T) {
	test.Run(t, new(CaptchaTestSuite))
}
