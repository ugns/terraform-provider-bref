package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLambdaLayer(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLambdaLayer,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.bref_lambda_layer.foo", "layer_name", regexp.MustCompile("^con")),
				),
			},
		},
	})
}

const testAccDataSourceLambdaLayer = `
data "bref_lambda_layer" "foo" {
  layer_name = "console"
}
`
