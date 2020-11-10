// File : resource_fake_object.go
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	"github.com/stobias123/gosolar"
)

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSubnetRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet FriendlyName you'd like to search",
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_mask": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"reserved": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"reserved_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"transient_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetRead(d *schema.ResourceData, meta interface{}) error {
	// return findIPAddress()
	client := meta.(*gosolar.Client)

	var s gosolar.Subnet

	vlan := d.Get("vlan").(string)
	address := d.Get("address").(string)
	if len(vlan) > 1 {
		s = client.GetSubnetByVLAN(vlan)
	} else if len(address) > 1 {
		s = client.GetSubnetByAddress(address)
	} else {
		log.Errorf("Provide vlan")
	}

	log.Infof("Subnet found: %s", s)
	log.Infof("CIDR found: %s", s.CIDR)
	d.SetId(s.VLAN)
	d.Set("address", s.Address)
	d.Set("cidr", s.CIDR)
	d.Set("vlan", s.VLAN)
	d.Set("address_mask", s.AddressMask)
	d.Set("display_name", s.DisplayName)
	d.Set("reserved", s.FriendlyName)
	d.Set("total_count", s.TotalCount)
	d.Set("used_count", s.UsedCount)
	d.Set("available_count", s.AvailableCount)
	d.Set("reserved_count", s.ReservedCount)
	log.Printf("Subnet found: %s", d)

	return nil
}
