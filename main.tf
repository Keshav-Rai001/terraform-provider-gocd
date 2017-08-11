provider "gocd" {
  baseurl = "http://localhost:8153/go/"
  skip_ssl_check = true
}

data "gocd_task_definition" "task1" {
  type = "exec"
  arguments = [
    "arg1",
    "arg2"]
  run_if = [
    "passed"]
  command = "/usr/loca/bin/terraform"
}

data "gocd_job_definition" "job1" {
  name = "job1"
  tasks = [
    "${data.gocd_task_definition.task1.json}"]
}

data "gocd_stage_definition" "manual-approval" {
  name = "test-stage-manual"
  jobs = [
    "${data.gocd_job_definition.job1.json}"]
  manual_approval = true
  authorization_roles = [
    "one",
    "two"]
}

data "gocd_stage_definition" "success-approval" {
  name = "test-stage-success"
  jobs = [
    "${data.gocd_job_definition.job1.json}"]
  success_approval = true
}

resource "gocd_pipeline_template" "my-server" {
  name = "my-test-template"
  stages = [
    "${data.gocd_stage_definition.manual-approval.json}",
    "${data.gocd_stage_definition.success-approval.json}"]

}


output "manual-approval" {
  value = "${data.gocd_stage_definition.manual-approval.json}"
}

output "success-approval" {
  value = "${data.gocd_stage_definition.success-approval.json}"
}