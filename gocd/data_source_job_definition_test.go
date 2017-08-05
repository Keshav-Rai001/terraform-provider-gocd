package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceJobDefinition(t *testing.T) {

	for i := 0; i <= 1; i++ {
		t.Run(
			fmt.Sprintf("gocd_job_definition.%d", i),
			DataSourceJobDefinition(t, i,
				fmt.Sprintf("data_source_job_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_job_definition.%d.rsp.json", i),
			),
		)
	}
}

func DataSourceJobDefinition(t *testing.T, index int, configPath string, expectedPath string) func(t *testing.T) {
	return func(t *testing.T) {
		config := testFile(configPath)
		expected := testFile(expectedPath)
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testGocdProviders,
			Steps: []resource.TestStep{testStepComparisonCheck(&TestStepJSONComparison{
				Index:        index,
				ID:           "data.gocd_job_definition.test",
				Config:       config,
				ExpectedJSON: expected,
			})},
		})
	}
}
