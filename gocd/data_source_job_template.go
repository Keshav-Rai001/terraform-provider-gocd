package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"encoding/json"
	"strconv"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/drewsonne/go-gocd/gocd"
)

func dataSourceGocdJobTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGocdJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"environment_variables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resources": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"elastic_profile_id"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"elastic_profile_id": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"resources"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tab": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"artifact": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"property": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"xpath": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdJobTemplateRead(d *schema.ResourceData, meta interface{}) error {

	tasks := []gocd.Task{}
	for _, task_string := range d.Get("tasks").([]interface{}) {
		task := gocd.Task{}
		err := json.Unmarshal([]byte(task_string.(string)), &task)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	doc := gocd.Job{
		Name:  d.Get("name").(string),
		Tasks: tasks,
	}
	jsonDoc, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(d.Get("name").(string))))

	return nil
}
