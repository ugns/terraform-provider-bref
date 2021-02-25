---
page_title: "bref_lambda_layer Data Source - terraform-provider-bref"
subcategory: ""
description: |-
  Bref PHP Lambda layer for published runtime version.
---

# Data Source `bref_lambda_layer`

Bref PHP Lambda layer for published runtime version.

## Example Usage

```terraform
data "bref_lambda_layer" "console" {
  layer_name = "console"
}
```

## Schema

### Required

- **layer_name** (String) The Bref PHP runtime lambda layer name.

### Optional

- **id** (String) The ID of this resource.

### Read-only

- **arn** (String) The Bref PHP runtime lambda layer ARN.
- **layer_arn** (String) The Bref PHP runtime lambda layer ARN.
- **version** (Number) The Bref PHP runtime lambda layer version.


