package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Provider defines the schema and resource map
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"nomadutility_acl_bootstrap": aclBootstrap(),
		},
	}
}
