package main

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform/helper/schema"
)

//Provider defines the schema and resource map
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"initial_backoff_interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "500ms",
				Description: "Specifies an initial backoff interval for retries.",
			},
			"backoff_multiplier": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     1.5,
				Description: "Specifies an multiplier for each wait time between retries.",
			},
			"max_backoff_interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "30s",
				Description: "Specifies an maximum backoff interval for retries.",
			},
			"timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10m",
				Description: "Specifies an upper bound on the time spent waiting for the ACLs to bootstrap before failing.",
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"nomadutility_acl_bootstrap": aclBootstrap(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = MustDuration(d.Get("initial_backoff_interval").(string))
	b.Multiplier = d.Get("backoff_multiplier").(float64)
	b.MaxInterval = MustDuration(d.Get("max_backoff_interval").(string))
	b.MaxElapsedTime = MustDuration(d.Get("timeout").(string))

	return &Config{
		retryBackoff: b,
	}, nil
}

type Config struct {
	retryBackoff *backoff.ExponentialBackOff
}
