provider "orion" {
    host = "localhost"
    user = "test"
    password = "pass"
    ssl = false
}

resource "orion_ip_address" "test" {
    vlan = 148
}