package gocd

import (
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
	"os"
	"runtime"
)

func Provider() terraform.ResourceProvider {
	return SchemaProvider()
}

// SchemaProvider describing the required configs to interact with GoCD server. Environment variables can also be set:
//   baseurl        - GOCD_URL
//   username       - GOCD_USERNAME
//   password       - GOCD_PASSWORD
//   skip_ssl_check - GOCD_SKIP_SSL_CHECK
func SchemaProvider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			//"gocd_stage_definition": dataSourceGocdStageTemplate(),
			"gocd_job_definition":  dataSourceGocdJobTemplate(),
			"gocd_task_definition": dataSourceGocdTaskDefinition(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gocd_environment":             resourceEnvironment(),
			"gocd_environment_association": resourceEnvironmentAssociation(),
			"gocd_pipeline_template":       resourcePipelineTemplate(),
			"gocd_pipeline":                resourcePipeline(),
			"gocd_pipeline_stage":          resourcePipelineStage(),
		},
		Schema: map[string]*schema.Schema{
			"baseurl": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["gocd_baseurl"],
				DefaultFunc: envDefault("GOCD_URL"),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["username"],
				DefaultFunc: envDefault("GOCD_USERNAME"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["password"],
				DefaultFunc: envDefault("GOCD_PASSWORD"),
			},
			"skip_ssl_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["skip_ssl_check"],
				DefaultFunc: envDefault("GOCD_SKIP_SSL_CHECK"),
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"baseurl":  "URL for the GoCD Server",
		"username": "User to interact with the GoCD API with.",
		"password": "Password for User for GoCD API interaction.",
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var url, u, p string
	var rUrl, rU, rP, rB interface{}
	var ok, nossl, b bool
	var cfg *gocd.Configuration

	if rUrl, ok = d.GetOk("baseurl"); ok {
		if url, ok = rUrl.(string); !ok || url == "" {
			url = os.Getenv("GOCD_URL")
		}
	}

	if rU, ok = d.GetOk("username"); ok {
		if u, ok = rU.(string); !ok || u == "" {
			u = os.Getenv("GOCD_USERNAME")
		}
	}

	if rP, ok = d.GetOk("password"); ok {
		if p, ok = rP.(string); !ok || p == "" {
			p = os.Getenv("GOCD_PASSWORD")
		}
	}

	if rB, ok = d.GetOk("skip_ssl_check"); ok {
		if b, ok = rB.(bool); !ok {
			nossl = false
		} else {
			nossl = b
		}
	}

	cfg = &gocd.Configuration{
		Server:       url,
		Username:     u,
		Password:     p,
		SkipSslCheck: nossl,
	}

	hClient := &http.Client{
		Transport: http.DefaultTransport,
	}

	// Add API logging
	hClient.Transport = logging.NewTransport("GoCD", hClient.Transport)
	gc := gocd.NewClient(cfg, hClient)

	versionString := terraform.VersionString()
	gc.UserAgent = fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, versionString)

	return gc, nil

}

func envDefault(e string) schema.SchemaDefaultFunc {
	return schema.MultiEnvDefaultFunc([]string{
		e,
	}, nil)
}
