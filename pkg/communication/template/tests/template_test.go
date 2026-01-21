package tests

import (
	"testing"

	"github.com/chris-alexander-pop/system-design-library/pkg/communication/template"
	"github.com/chris-alexander-pop/system-design-library/pkg/communication/template/adapters/memory"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
)

type TemplateTestSuite struct {
	test.Suite
	engine template.Engine
}

func (s *TemplateTestSuite) SetupTest() {
	s.Suite.SetupTest()
	mem := memory.New()
	mem.AddTemplate("welcome", "Hello {{.name}}")
	s.engine = mem
}

func (s *TemplateTestSuite) TestRender() {
	// Memory adapter usually mocks rendering or returns simple string
	input := map[string]string{"name": "World"}
	result, err := s.engine.Render(s.Ctx, "welcome", input)
	s.NoError(err)
	s.NotEmpty(result)
}

func TestTemplateSuite(t *testing.T) {
	test.Run(t, new(TemplateTestSuite))
}
