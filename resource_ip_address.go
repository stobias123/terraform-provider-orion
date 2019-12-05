package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	"github.com/stobias123/gosolar"
)

func resourceIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPAddressCreate,
		Read:   resourceIPAddressRead,
		Update: resourceIPAddressUpdate,
		Delete: resourceIPAddressDelete,

		Schema: map[string]*schema.Schema{
			"vlan": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"subnet_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_string": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIPAddressCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gosolar.Client)
	var subnet gosolar.Subnet
	subnetName := d.Get("subnet_name").(string)
	vlanName := d.Get("vlan").(string)

	if len(subnetName) < 1 {
		subnet = client.GetSubnet(subnetName)
	} else if len(vlanName) < 1 {
		subnet = client.GetSubnetByVLAN(vlanName)
	} else {
		log.Errorf("Provide either subnet_name or vlan")
	}

	suggestedIP := client.GetFirstAvailableIP(subnet.Address, fmt.Sprintf("%d", subnet.CIDR))

	d.Set("vlan", subnet.VLAN)
	d.Set("subnet_name", subnet.DisplayName)
	d.Set("address", suggestedIP.Address)

	reservedIP := client.ReserveIP(suggestedIP.Address)

	d.Set("status", reservedIP.Status)
	d.Set("status_string", suggestedIP.StatusString)

	return resourceIPAddressRead(d, m)
}

// If the status is "Available" then we can use it.
func resourceIPAddressRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gosolar.Client)

	requested_address := d.Get("address").(string)
	ip := client.GetIP(requested_address)

	// If the ip is available... reset ID so that terraform will re-create.
	if ip.Status == 2 {
		d.SetId("")
		return nil
	}

	d.Set("address", ip.Address)
	d.Set("status", ip.Status)
	d.Set("status_string", ip.StatusString)

	return nil
}

func resourceIPAddressUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIPAddressRead(d, m)
}

func resourceIPAddressDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
