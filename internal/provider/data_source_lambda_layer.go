package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLambdaLayer() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Bref PHP Lambda layer for published runtime version.",

		ReadContext: readerContextFuncProvider("bref_lambda_layer"),

		Schema: map[string]*schema.Schema{
			"layer_name": {
				// This description is used by the documentation generator and the language server.
				Description: "The Bref PHP runtime lambda layer name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "The Bref PHP runtime lambda layer version.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"arn": {
				Description: "The Bref PHP runtime lambda layer ARN.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"layer_arn": {
				Description: "The Bref PHP runtime lambda layer ARN.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
