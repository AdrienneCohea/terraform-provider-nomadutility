package main

import "github.com/hashicorp/terraform/helper/schema"

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
