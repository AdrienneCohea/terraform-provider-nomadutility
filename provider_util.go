package main

import (
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func forget(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	_ = m
	return nil
}

func noop(d *schema.ResourceData, m interface{}) error {
	_ = d
	_ = m
	return nil
}

func MustDuration(dur string) time.Duration {
	actual, err := time.ParseDuration(dur)
	if err != nil {
		panic(err)
	}
	return actual
}

func getClient(d *schema.ResourceData) *api.Client {
	conf := api.DefaultConfig()
	conf.Address = d.Get("address").(string)
	conf.TLSConfig.CACert = d.Get("ca_file").(string)
	conf.TLSConfig.CACertPEM = []byte(d.Get("ca_pem").(string))
	conf.TLSConfig.ClientCert = d.Get("cert_file").(string)
	conf.TLSConfig.ClientKey = d.Get("key_file").(string)
	conf.TLSConfig.TLSServerName = d.Get("tls_server_name").(string)

	client, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	return client
}
