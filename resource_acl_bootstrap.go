package main

import (
	"fmt"
	"log"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func aclBootstrap() *schema.Resource {
	return &schema.Resource{
		Create: bootstrapACLs,
		Read:   doNothing,
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
	return backoff.Retry(func() error {
		client := meta.(*api.Client)

		resp, _, err := client.ACLTokens().Bootstrap(nil)
		if err != nil {
			return fmt.Errorf("unable to bootstrap ACL system: %s", err.Error())
		}
		log.Printf("[DEBUG] Created ACL token %q", resp.AccessorID)
		d.SetId(resp.AccessorID)

		d.Set("accessor_id", resp.AccessorID)
		d.Set("secret_id", resp.SecretID)
		d.Set("name", resp.Name)
		d.Set("type", resp.Type)
		d.Set("policies", resp.Policies)
		d.Set("global", resp.Global)
		d.Set("create_time", resp.CreateTime.UTC().String())

		return nil
	}, backoff.NewExponentialBackOff())
}

func forget(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func doNothing(d *schema.ResourceData, m interface{}) error {
	return nil
}
