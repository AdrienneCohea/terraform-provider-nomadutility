package main

import "github.com/hashicorp/terraform/helper/schema"

func aclBootstrap() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Default:  "127.0.0.1",
				Optional: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	d.SetId("token-from-api")
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
