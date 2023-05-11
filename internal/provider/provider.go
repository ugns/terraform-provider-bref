package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"region": {
					Type:     schema.TypeString,
					Required: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"AWS_REGION",
						"AWS_DEFAULT_REGION",
					}, "us-east-1"),
					Description: "AWS Region of Bref PHP runtime layers. Can be specified with the `AWS_REGION` " +
						"or `AWS_DEFAULT_REGION` environment variable.",
					InputDefault: "us-east-1",
				},
				"bref_version": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BREF_VERSION", "2.0.5"),
					Description: "The Bref PHP runtime version to work with. Can be specified with the " +
						"`BREF_VERSION` environment variable.",
					InputDefault: "2.0.5",
				},
				"bref_extra_version": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BREF_EXTRA_VERSION", "1.1.1"),
					Description: "The Bref Extra PHP runtime version to work with. Can be specified with the " +
						"`BREF_EXTRA_VERSION` environment variable.",
					InputDefault: "1.1.1",
				},
				"bref_aws_account": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BREF_AWS_ACCOUNT", "534081306603"),
					Description: "The Bref AWS account to pull layers from. Can be specified with the " +
						"`BREF_AWS_ACCOUNT` environment variable.",
					InputDefault: "534081306603",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"bref_lambda_layer":       dataSourceLambdaLayer(),
				"bref_extra_lambda_layer": extraDataSourceLambdaLayer(),
			},
			ResourcesMap: map[string]*schema.Resource{},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	Region       string
	Version      string
	ExtraVersion string
	AccountIds   map[string]string
	URLs         map[string]string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiClient := apiClient{
			Region:       d.Get("region").(string),
			Version:      d.Get("bref_version").(string),
			ExtraVersion: d.Get("bref_extra_version").(string),
			AccountIds: map[string]string{
				"bref_lambda_layer":       d.Get("bref_aws_account").(string),
				"bref_extra_lambda_layer": "403367587399",
			},
			URLs: map[string]string{
				"bref_lambda_layer":       fmt.Sprintf("https://raw.githubusercontent.com/brefphp/bref/%s/layers.json", d.Get("bref_version").(string)),
				"bref_extra_lambda_layer": fmt.Sprintf("https://raw.githubusercontent.com/brefphp/extra-php-extensions/%s/layers.json", d.Get("bref_extra_version").(string)),
			},
		}

		return &apiClient, nil
	}
}

func readerContextFuncProvider(source string) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// use the meta value to retrieve your client from the provider configure method
		client := meta.(*apiClient)
		ch := cleanhttp.DefaultPooledClient()

		layerName := d.Get("layer_name").(string)

		var diags diag.Diagnostics
		url := client.URLs[source]
		accountId := client.AccountIds[source]

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return diag.Errorf("Unable to build request for %s version of Bref runtime layers", client.Version)
		}

		r, err := ch.Do(req)
		if err != nil {
			return diag.Errorf("Error retrieving Bref runtime layers: %s", err.Error())
		}
		defer r.Body.Close()

		var layers map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&layers)
		if err != nil {
			return diag.Errorf("Error parsing Bref runtime layers: %s", err.Error())
		}

		regions := layers[layerName].(map[string]interface{})
		var version int

		// We do this type conversion because some versions are float64 and some
		// others are strings. Terraform spec expects int.
		switch t := regions[client.Region].(type) {
		case string:
			version, err = strconv.Atoi(t)
			if err != nil {
				return diag.Errorf("Unable to locate a Bref v%s lambda layer version for %s in %s region", client.Version, layerName, client.Region)
			}
		case float64:
			version = int(t)
		default:
			return diag.Errorf("Unable to locate a Bref v%s lambda layer version for %s in %s region", client.Version, layerName, client.Region)
		}

		arn := fmt.Sprintf("arn:aws:lambda:%s:%s:layer:%s:%d", client.Region, accountId, layerName, version)

		if err := d.Set("version", version); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("arn", arn); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("layer_arn", arn); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

		return diags
	}
}
