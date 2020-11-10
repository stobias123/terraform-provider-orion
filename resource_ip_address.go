package main

import (
	"fmt"

	guuid "github.com/google/uuid"
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
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"subnet_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func genUUID() string {
	id := guuid.New()
	return id.String()
}

func resourceIPAddressCreate(d *schema.ResourceData, m interface{}) error {
	log.Info("StartResource")
	randID := genUUID()
	d.SetId(randID)
	mutexKey := fmt.Sprintf("ip-address")
	orionMutexKV.Lock(mutexKey)
	defer orionMutexKV.Unlock(mutexKey)

	client := m.(*gosolar.Client)
	var subnet gosolar.Subnet
	subnetName := d.Get("subnet_name").(string)
	vlanName := d.Get("vlan").(string)
	if len(subnetName) > 1 {
		log.Infof("Subnet")
		log.Infof(subnetName)
		subnet = client.GetSubnet(subnetName)
	} else if len(vlanName) > 1 {
		log.Infof("vlan")
		log.Infof(vlanName)
		subnet = client.GetSubnetByVLAN(vlanName)
	} else {
		return fmt.Errorf("Provide either subnet_name or vlan")
	}

	log.Info(subnet)
	suggestedIP := client.GetFirstAvailableIP(subnet.Address, fmt.Sprintf("%d", subnet.CIDR))
	log.Info(suggestedIP)
	reservedIP := client.ReserveIP(suggestedIP.Address)

	// This is really kind of sloppy, I should be returning err from gosolar.
	if len(d.Get("hostname").(string)) > 0 {
		ipHostname := client.AddHostnametoIPNode(suggestedIP.Address, d.Get("hostname").(string))

		log.Info("SetHostname")
		log.Info(ipHostname)
	}

	log.Info("Reserved IP")
	log.Info(reservedIP)

	log.Info("Suggested IP")
	log.Info(suggestedIP)

	d.Set("vlan", subnet.VLAN)
	//d.Set("subnet_name", subnet.DisplayName)
	d.Set("address", suggestedIP.Address)
	d.Set("status", reservedIP.Status)
	//d.Set("status_string", suggestedIP.StatusString)
	log.Info(d)

	return resourceIPAddressRead(d, m)
}

// If the status is "Available" then we can use it.
func resourceIPAddressRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gosolar.Client)
	log.Info("READ IP")
	log.Info(d)
	requestedAddress := d.Get("address").(string)
	ip := client.GetIP(requestedAddress)

	// If the ip is available... reset ID so that terraform will re-create.
	if ip.Status == 2 {
		d.SetId("")
		return nil
	}

	d.Set("address", ip.Address)
	d.Set("status", ip.Status)
	//d.Set("status_string", ip.StatusString)
	log.Info(d)
	return nil
}

func resourceIPAddressUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIPAddressRead(d, m)
}

func resourceIPAddressDelete(d *schema.ResourceData, m interface{}) error {
	log.Info("StartResource")
	client := m.(*gosolar.Client)
	client.ReleaseIP(d.Get("address").(string))
	d.SetId("")
	return nil
}
