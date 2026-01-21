package tests

import (
	"testing"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/auth/session"
	"github.com/chris-alexander-pop/system-design-library/pkg/auth/session/adapters/memory"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
)

type SessionTestSuite struct {
	test.Suite
	manager session.Manager
}

func (s *SessionTestSuite) SetupTest() {
	s.Suite.SetupTest()
	s.manager = memory.New(session.Config{TTL: time.Hour})
}

func (s *SessionTestSuite) TestCreateGetDelete() {
	userID := "user-123"

	sess, err := s.manager.Create(s.Ctx, userID, nil)
	s.NoError(err)
	s.NotEmpty(sess.ID)

	got, err := s.manager.Get(s.Ctx, sess.ID)
	s.NoError(err)
	s.Equal(userID, got.UserID)

	err = s.manager.Delete(s.Ctx, sess.ID)
	s.NoError(err)

	_, err = s.manager.Get(s.Ctx, sess.ID)
	s.Error(err)
}

func TestSessionSuite(t *testing.T) {
	test.Run(t, new(SessionTestSuite))
}
