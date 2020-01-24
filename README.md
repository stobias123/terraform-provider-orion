# terraform-provider-orion

Provides two useful things `data_source_subnet` and `resource_ip_address`


`resource_ip_address` - Gets the first available ip address from a subnet.

`data_source_subnet` - Provides info about the chosen VLAN. 

Usage:


```
resource "orion_ip_address" "test" {
    vlan = 148
}

data "orion_subnet" "test" {
    vlan = 148
}

output "orion_ip" {
    value = orion_ip_address.test.address
}

output "orion_subnet" {
    value = data.orion_subnet.test.address
}
```

## Additional info
This all pairs well with `solarcmd` found [here](https://github.com/stobias123/gosolar/tree/master/command)
