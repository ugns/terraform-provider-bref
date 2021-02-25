---
page_title: "bref Provider"
subcategory: ""
description: |-
  
---

# bref Provider



## Example Usage

```terraform
provider "bref" {
  # example configuration here
}
```

## Schema

### Optional

- **bref_version** (String) The Bref PHP runtime version to work with. Can be specified with the `BREF_VERSION` environment variable.
- **region** (String) AWS Region of Bref PHP runtime layers. Can be specified with the `AWS_REGION` or `AWS_DEFAULT_REGION` environment variable.
