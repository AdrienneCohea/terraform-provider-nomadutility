package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform/helper/schema"
)

func aclBootstrap() *schema.Resource {
	return &schema.Resource{
		Create: bootstrapACLs,
		Read:   noop,
		Update: noop,
		Delete: forget,
		Importer: &schema.ResourceImporter{
			State: importState,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_ADDR", "http://127.0.0.1:4646"),
				Description: "URL of the root of the target Nomad agent.",
			},
			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_CACERT", ""),
				Description: "A path to a PEM-encoded certificate authority used to verify the remote agent's certificate.",
			},
			"ca_pem": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "PEM-encoded certificate authority used to verify the remote agent's certificate.",
			},
			"cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_CLIENT_CERT", ""),
				Description: "A path to a PEM-encoded certificate provided to the remote agent; requires use of key_file.",
			},
			"key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_CLIENT_KEY", ""),
				Description: "A path to a PEM-encoded private key, required if cert_file is specified.",
			},
			"tls_server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_CLIENT_KEY", ""),
				Description: "Specifies an optional string used to set the SNI host when connecting to Vault via TLS.",
			},
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
	c := meta.(*Config)
	client := getClient(d)

	return backoff.Retry(func() error {
		resp, _, err := client.ACLTokens().Bootstrap(nil)
		if err != nil {
			return maybeRetry(err)
		}

		log.Printf("[DEBUG] Created ACL token %q", resp.AccessorID)

		d.SetId(resp.AccessorID)

		err = multiError(
			d.Set("accessor_id", resp.AccessorID),
			d.Set("secret_id", resp.SecretID),
			d.Set("name", resp.Name),
			d.Set("type", resp.Type),
			d.Set("policies", resp.Policies),
			d.Set("global", resp.Global),
			d.Set("create_time", resp.CreateTime.UTC().String()))

		log.Printf("[DEBUG] Saved ACL token in state %q", resp.AccessorID)

		return err
	}, c.retryBackoff)
}

func importState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// id will be set by the user in the format accessor_id:secret_id
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id, expected accessor_id:secret_id")
	}
	d.Set("accessor_id", parts[0])
	d.Set("secret_id", parts[1])
	return []*schema.ResourceData{d}, nil
}
