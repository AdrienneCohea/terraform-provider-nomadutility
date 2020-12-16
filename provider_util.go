package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
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
