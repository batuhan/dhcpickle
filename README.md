# dhcpickle

A simple DHCP server that assigns IPs from http responses.

## env variables

    DHCP_SUBNET_MASK
    DHCP_ROUTER
    DHCP_DNS
    DHCP_IP_ADDRESS_LEASE_TIME
    DHCP_SERVER_IDENTIFIER
    DHCP_ENDPOINT
    DHCP_AUTH_HEADER
    DHCP_AUTH_TOKEN

all except the last 3 env variables have default dynamic values

minimal configuration is needed

### example

    DHCP_SUBNET_MASK=255.255.255.0
    DHCP_ROUTER=192.168.0.1
    DHCP_DNS=8.8.8.8,8.8.4.4
    DHCP_IP_ADDRESS_LEASE_TIME=24h
    DHCP_SERVER_IDENTIFIER=192.168.0.100
    DHCP_ENDPOINT=http://localhost
    DHCP_AUTH_HEADER=X-Auth-Header
    DHCP_AUTH_TOKEN=SecretKey

## request & response

dhcpickle sends a simple http `POST` request with `text/plain` type.

The body of the request contains the mac address of the request.

Response must be a `text/plain` IP address.

### example

request:

    00:50:56:b6:87:6c


response:

    192.168.0.150

