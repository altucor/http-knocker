# http-knocker


---
## Project description
TLDR: http-knocker - project which allows you to open ports on your router by visiting some secret URL instead of old-fashioned port knocking via ping.

Usually port knocking done via ICMP protocol by sending ping packets. There is some common approaches to configure and check correctness of knock:
1. Sending `ping` requests to several listening ports.
   Example: ping `12345`, `12346`, `12347` to open `22 SSH` port.
2. Sending `ping` packets of specific length.\
   Example: ping `12345` port `3` times with payload size `733`, `823`, `341`

## Problem description
The main problem is security. Any service port, publicly exposed to WAN can be potential target for attacker to start bruteforcing it and to gain access to local network and infrastructure.
Some of the services should be publicly available only for small group of members. Exposing ports of these services in to WAN for all, allows scanning utils to find this exposed port and start bruteforcing it to gain access. As prevention of this case, usual solution is to add allowed client IP addresses in to white list, but these clients can have dynamic IP addresses or should have access in critical situations from mobile devices which can have completely different IP address which is not in whitelist.\
VPN usualy resolves all mentioned above problems. But it also can be problem for end users which devices now can have completely different IP addressation/default gateway/dns servers and not accesible local network devices. For example it can be painful for some local network development. VPN produces another problems:
1. Often tunnel drops which can require full re-authorization with 2FA and other stuff
2. Often force logout/login to be able access `local-network` devices or `remote-network` devices


## The main targets of this project:
1. Make usage of `private` services more secure for allowed group of users via `public` networks
2. Simplify getting access to `private` services for allowed end users
3. Potentially replace VPN services in some cases


## How it works
System administrator activities:
1. Deploy `http-knocker` in local network
2. Configure HTTP/HTTPS(80/443) port forwarding for `http-knocker` which allows it to listen incoming requests from users.
3. Configuring border gateway. Setup port forwarding to interesting internal services like SSH(22) RDP(3389). And also block connections from all external ip addresses to it.
4. Describing public enpoints which should be secured in `.yml` config file. Writing secret URL which will open access for specific port. Adding in config routers with with username, password, ip, public key, etc. So `http-knocker` will be able to connect to router and control it.

End user activities:
1. Ask system administrator for secret URL for example for SSH(22) service.
2. Just visit this URL from you device and you will be granted to access this port for limited amount of time, or unlimited, depends on config from system administrator.

## Examples

Here simple example of `http-knocker` config file which will open port `22` for `1 hour 30 minutes` for IP address from which user will open `HTTP/HTTPS` link `https://<http-knocker-service-ip-address>:8000/giurhft45iedeohulgtirhguesolkvgdtruheixlggirdnth`. It will open this port on `MikroTik RouterOS` device with ip address `192.168.0.1` via `REST` protocol. If user wrote URL correctly he will receive in response `222` HTTP code, it is example of customization. Like to stay more hidden, system administrator can even write `404` code on success so it will be harder to bruteforce. Also if user will enter invalid URL it will receive `404` HTTP code as default response from server. Under the hood if user visit correct URL, `http-knocker` will connect to MikroTik device with provided credentials via REST protocol, it will find in firewall rule with comment `http-knocker-drop-all-rule-22-tcp`, where it should be `DROP` rule for all incoming connections to `22` port `TCP` protocol. `http-knocker` will find this rule and will place `ACCEPT` rule with IP address of client visited secret URL, exactly before `DROP`.

```
server: # Common configuration of http-knocker listener
  host: 0.0.0.0 # Listen on all interfaces
  port: 8000 # Listen on port 8000
  default-response-code: 404 # If someone tries to access non-existent URL response with HTTP 404

devices: # Block with description of devices(routers)
  deviceRouterOsRest: # Unique device name
    type: rest # Device type, can be for example: rest, ssh, router-os. Which specifies type of communication with router.
    protocol: rest-router-os # Device protocol. Some device types can support several different protocols. Like via ssh type we potentially can communicate with IP tables or with anoother RouterOS device.
    connection: # Router connection info
      username: httpKnockerUser # Username
      password: mySuperPassword # Password
      endpoint: https://192.168.0.1/rest # Router endpoint

endpoints: # Description of endpoints
  endpoint-first: # Endpoint unique name
    duration: 1h30s # Timeout, how much port will be opened for client
    port: 22 # Which port will be opened
    protocol: tcp # On which protocol port will be opened
    response-code-on-success: 222 # Which HTTP code user will receive on visiting this endpoint

controllers: # Controllers section, combines devices and endpoints, allows re-using endpoints with other devices and deduplication of config.
  firewall-rest-first: # Uniqie name of controller
    type: basic # Type of controller
    url: giurhft45iedeohulgtirhguesolkvgdtruheixlggirdnth # Secret URL which should be accessed to open endpoint port
    device: deviceRouterOsRest # Device which should open port
    endpoint: endpoint-first # Endpoint which describes which port and how long should be opeend
    config: # Additional configuration of controller
      drop-rule-comment: http-knocker-drop-all-rule-22-tcp # Comment from DROP rule from device, before which ACCEPT clients should be placed

```

## Features
Suppoerted Devices:
1. RouterOS via REST
2. iptables via SSH

Endpoint:
1. Can open port by just visiting secret URL
2. Can open port by visiting secret URL and passing Basic HTTP Auth
3. Can extract client IP address from different like HTTP Headers section

## Roadmap
Add support for devices:
1. Cisco SSH, UART console port
2. UFW via SSH
3. pfSense via HTTP or REST
4. Netgear, D-Link, TP-Link, etc
5. Small SOHO like ASUS
6. Puller device. For cases where routers cannot be accessed locally from http-knocker. Solution where instead of directly connecting to routers, http-knocker will work as listener to which other routers can connect, push their current status and pull new information about new clients which should be added and which should be kicked.


Knocking validation:
1. Third party auth providers like Authelia, DUO
2. Reducing 


General:
1. Prepare easy to deploy containerized solutions like Docker etc
2. Writing better documnetation
3. Writing more unit tests
4. Adding more devices and protocols
5. Refactoring for easier extendability



