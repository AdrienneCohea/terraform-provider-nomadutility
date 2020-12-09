package main

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}

func nomadClient() (*api.Client, error) {
	cfg := &api.Config{}
	return api.NewClient(cfg)
}

func bootstrap(client *api.Client) {
	token, _, err := client.ACLTokens().Bootstrap(nil)
	if err != nil {
		fmt.Printf("Got this: %v", err)
		return
	}

	fmt.Printf("Got this: %v", token)
}
