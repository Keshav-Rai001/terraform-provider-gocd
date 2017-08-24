package gocd

import (
	"testing"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
)

func testResourcePipeline(t *testing.T) {
	t.Run("Basic", testResourcePipelineBasic)
	t.Run("ImportBasic", testResourcePipelineImportBasic)
	t.Run("ExistsFail", testResourcePipelineExistsFail)
	t.Run("PipelineReadHelper", testResourcePipelineReadHelper)
}

func testResourcePipelineBasic(t *testing.T) {

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

func testResourcePipelineExistsFail(t *testing.T) {
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{})

	exists, err := resourcePipelineExists(rd, nil)
	assert.False(t, exists)
	assert.EqualError(t, err, "`name` can not be empty")
}

func testResourcePipelineReadHelper(t *testing.T) {

}
