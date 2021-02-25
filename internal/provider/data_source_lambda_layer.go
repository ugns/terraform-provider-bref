package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLambdaLayer() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Bref PHP Lambda layer for published runtime version.",

		ReadContext: dataSourceLambdaLayerRead,

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

func dataSourceLambdaLayerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	ch := cleanhttp.DefaultPooledClient()

	layerName := d.Get("layer_name").(string)

	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("https://raw.githubusercontent.com/brefphp/bref/%s/layers.json", client.Version), nil)
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
	version, err := strconv.Atoi(regions[client.Region].(string))
	if err != nil {
		return diag.Errorf("Unable to locate a Bref v%s lambda layer version for %s in %s region", client.Version, layerName, client.Region)
	}
	arn := fmt.Sprintf("arn:aws:lambda:%s:%s:layer:%s:%d", client.Region, client.AccountId, layerName, version)

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
