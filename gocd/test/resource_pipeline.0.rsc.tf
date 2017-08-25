resource "gocd_pipeline" "test-pipeline" {
  name = "pipeline0-terraform"
  template = "${gocd_pipeline_template.test-pipeline.id}"
  materials = [
    {
      type = "git"
      attributes {
        name = "gocd-src"
        url = "git@github.com:gocd/gocd"
        branch = "feature/my-addition"
        destination = "gocd-dir"
        auto_update = true
        filter {
          ignore = [
            "one",
            "two"
          ]
        }
      }
    }
  ]
}


resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"
  stages = [
    <<STAGE
{
  "name": "test-stage",
  "fetch_materials": false,
  "clean_working_directory": false,
  "never_cleanup_artifacts": false,
  "approval": {
    "type": "success"
  },
  "jobs": [
    {
      "name": "job1",
      "tasks": [
        {
          "type": "exec",
          "attributes": {
            "run_if": [
              "passed"
            ],
            "command": "terraform"
          }
        }
      ]
    }
  ]
}
STAGE
  ]
}
