package tests

import (
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/security/secrets"
	"github.com/chris-alexander-pop/system-design-library/pkg/security/secrets/adapters/memory"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
)

type SecretsTestSuite struct {
	test.Suite
	manager secrets.SecretManager
}

func (s *SecretsTestSuite) SetupTest() {
	s.Suite.SetupTest()
	s.manager = memory.New()
}

func (s *SecretsTestSuite) TestSetGet() {
	name := "db-password"
	value := "super-secret"

	err := s.manager.Set(s.Ctx, name, value)
	s.NoError(err)

	retrieved, err := s.manager.Get(s.Ctx, name)
	s.NoError(err)
	s.Equal(value, retrieved)
}

func (s *SecretsTestSuite) TestGet_NotFound() {
	_, err := s.manager.Get(s.Ctx, "unknown")
	s.Error(err)
}

func TestSecretsSuite(t *testing.T) {
	test.Run(t, new(SecretsTestSuite))
}
