package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExtraDataSourceLambdaLayer(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExtraDataSourceLambdaLayer,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.bref_extra_lambda_layer.foo", "layer_name", regexp.MustCompile("^yaml")),
				),
			},
		},
	})
}

const testAccExtraDataSourceLambdaLayer = `
data "bref_extra_lambda_layer" "foo" {
  layer_name = "yaml-php-82"
}
`
