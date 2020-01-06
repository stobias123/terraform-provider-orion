provider "orion" {
    host = "localhost"
    user = "test"
    password = "pass"
    ssl = false
}

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