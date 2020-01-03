package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stobias123/gosolar"
)

// this is the global mutex for use with this plugin
var orionMutexKV = mutexkv.NewMutexKV()

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORION_HOST", nil),
				Description: "Hostname for orion server",
			},
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORION_USER", nil),
				Description: "API User",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORION_PASSWORD", nil),
				Description: "API Password",
			},
			"ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORION_SSL", nil),
				Description: "Use https?",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"orion_ip_address": resourceIPAddress(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"orion_subnet": dataSourceSubnet(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	user := d.Get("user").(string)
	pass := d.Get("password").(string)
	ssl := d.Get("ssl").(bool)

	client := gosolar.NewClient(host, user, pass, ssl, true)

	return client, nil
}
