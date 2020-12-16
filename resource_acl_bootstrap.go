package main

import (
  "github.com/cenkalti/backoff/v4"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
)

func aclBootstrap() *schema.Resource {
	return &schema.Resource{
		Create: bootstrapACLs,
		Read:   noop,
		Delete: forget,

		Schema: map[string]*schema.Schema{
			"accessor_id": {
				Description: "Nomad-generated ID for this token.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"secret_id": {
				Description: "The value that grants access to Nomad.",
				Computed:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},

			"name": {
				Description: "Human-readable name for this token.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"type": {
				Description: "The type of token to create, 'client' or 'management'.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"policies": {
				Description: "The ACL policies to associate with the token, if it's a 'client' type.",
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"global": {
				Description: "Whether the token should be replicated to all regions or not.",
				Computed:    true,
				Type:        schema.TypeBool,
			},

			"create_time": {
				Description: "The timestamp the token was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func bootstrapACLs(d *schema.ResourceData, meta interface{}) error {
  c := meta.(Config)

	return backoff.Retry(func() error {
		resp, _, err := c.client.ACLTokens().Bootstrap(nil)
		if err != nil {
			return maybeRetry(err)
		}

    log.Printf("[DEBUG] Created ACL token %q", resp.AccessorID)

		d.SetId(resp.AccessorID)

		_ = d.Set("accessor_id", resp.AccessorID)
		_ = d.Set("secret_id", resp.SecretID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("type", resp.Type)
		_ = d.Set("policies", resp.Policies)
		_ = d.Set("global", resp.Global)
		_ = d.Set("create_time", resp.CreateTime.UTC().String())

    log.Printf("[DEBUG] Created ACL token %q", resp.AccessorID)

		return nil
	}, c.retryBackoff)
}
