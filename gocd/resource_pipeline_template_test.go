package gocd

import (
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/drewsonne/go-gocd/gocd"
)

func TestResourcePipelineTemplate(t *testing.T) {
	t.Run("Basic", testResourcePipelineTemplateBasic)
	t.Run("ImportBasic", testResourcePipelineTemplateImportBasic)
	t.Run("Exists", testResourcePipelineTemplateExists)
	t.Run("PipelineReadHelper", testResourcePipelineTemplateReadHelper)
}

func testResourcePipelineTemplateReadHelper(t *testing.T) {

	t.Run("MissingName", testResourcePipelineTemplateReadHelperMissingName)
	t.Run("JSONFail", testResourcePipelineTemplateReadHelperJSONFail)
}

func testResourcePipelineTemplateReadHelperJSONFail(t *testing.T) {
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{
		Attributes: map[string]string{"name": "mock-name",},
	})

	p := gocd.PipelineTemplate{
		Name: "mock-name",
		Stages: []*gocd.Stage{
			{Name: ""},
		},
	}

	err := readPipelineTemplate(rd, &p, nil)

	assert.EqualError(t, err, "`gocd.Stage.Name` is empty")

}

func testResourcePipelineTemplateReadHelperMissingName(t *testing.T) {

	rd := (&schema.Resource{Schema: map[string]*schema.Schema{}}).Data(&terraform.InstanceState{})
	e := errors.New("mock-error")
	err := readPipelineTemplate(rd, nil, e)

	assert.EqualError(t, err, "mock-error")
}

func testResourcePipelineTemplateExists(t *testing.T) {
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{})

	exists, err := resourcePipelineTemplateExists(rd, nil)
	assert.False(t, exists)
	assert.EqualError(t, err, "`name` can not be empty")
}

func testResourcePipelineTemplateBasic(t *testing.T) {

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineTemplateDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_template.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline"),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template0-terraform"),
					testCheckPipelineTemplate1StageCount("gocd_pipeline_template.test-pipeline"),
				),
			},
			{
				Config: testFile("resource_pipeline_template.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline"),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template0-terraform"),
					testCheckPipelineTemplate2StageCount("gocd_pipeline_template.test-pipeline"),
				),
			},
		},
	})

}

func testCheckPipelineTemplate1StageCount(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource].Primary; rs.Attributes["stages.#"] != "1" {
			return fmt.Errorf("Expected 1 stage. Found '%s'", rs.Attributes["stages.#"])
		}
		return nil
	}
}

func testCheckPipelineTemplate2StageCount(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource].Primary; rs.Attributes["stages.#"] != "2" {
			return fmt.Errorf("Expected 2 stages. Found '%s'", rs.Attributes["stages.#"])
		}
		return nil
	}
}

func testCheckPipelineTemplateName(resource string, id string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource]; rs.Primary.ID != id {
			return fmt.Errorf("Expected id 'template1-terraform', got '%s", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckPipelineTemplateExists(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pipeline template name is set")
		}

		return nil
	}
}
