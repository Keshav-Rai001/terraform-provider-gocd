package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceJobDefinition(t *testing.T) {
	testSteps := []resource.TestStep{}
	for i := 0; i <= 1; i++ {
		testSteps = append(
			testSteps,
			testStepComparisonCheck(TestStepJsonComparison{
				Id:           "data.gocd_job_definition.test",
				Config:       testFile(fmt.Sprintf("data_source_job_template.%d.rsc.tf", i)),
				ExpectedJSON: testFile(fmt.Sprintf("data_source_job_template.%d.rsp.json", i)),
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGocdProviders,
		Steps:     testSteps,
	})
}
