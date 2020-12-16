package main

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform/helper/schema"
)

//Provider defines the schema and resource map
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_ADDR", "http://127.0.0.1:4646"),
				Description: "URL of the root of the target Nomad agent.",
			},
			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NOMAD_CACERT", ""),
				Description: "A path to a PEM-encoded certificate authority used to verify the remote agent's certificate.",
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
	conf := api.DefaultConfig()
	conf.Address = d.Get("address").(string)
	conf.TLSConfig.CACert = d.Get("ca_file").(string)
	conf.TLSConfig.ClientCert = d.Get("cert_file").(string)
	conf.TLSConfig.ClientKey = d.Get("key_file").(string)
	conf.TLSConfig.TLSServerName = d.Get("tls_server_name").(string)

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = MustDuration(d.Get("initial_backoff_interval").(string))
	b.Multiplier = d.Get("backoff_multiplier").(float64)
	b.MaxInterval = MustDuration(d.Get("max_backoff_interval").(string))
	b.MaxElapsedTime = MustDuration(d.Get("timeout").(string))

	return &Config{
		retryBackoff: b,
		client:       client,
	}, nil
}

type Config struct {
	retryBackoff *backoff.ExponentialBackOff
	client       *api.Client
}
