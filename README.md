# http-knocker


---
## Project description
TLDR: http-knocker - project which allows you to open ports on your router by visiting some secret URL instead of old-fashioned port knocking via ping.

Usually port knocking done via ICMP protocol by sendin ping packets. There is some common approaches to configure and check correctness of knock:
1. Sending `ping` requests to several listening ports.
   Example: ping `12345`, `12346`, `12347` to open `22 SSH` port.
2. Sending `ping` packets of specific length.\
   Example: ping `12345` port `3` times with payload size `733`, `823`, `341`

## Problem description
The main problem is security. Any service publicly exposed port to WAN can be potential target for attacker to start bruteforcing it and to gain access to local network and infrastructure.
Some of the services should be publicly available only for small group of members. Exposing ports of these services in to WAN for all, allows scanning utils to find this exposed port and start bruteforcing it to gain access. As prevention of this case, usual solution is to add allowed client IP addresses in to white list, but these clients can have dynamic IP addresses or should have access in critical situations from mobile devices which can have completely different IP address which is not in whitelist.\
VPN usualy resolves all mentioned above problems. But it also can be problem for end users which devices now can have completely different IP addressation/default gateway/dns servers and not accesible local network devices. For example it can be painful for some local network development. VPN produces another problems:
1. Often tunnel drops which can require full re-authorization with 2FA and other stuff
2. Often force logout/login to be able access `local-network` devices or `remote-network` devices


## The main targets of this project:
1. Make usage of `private` services more secure for allowed group of users via `public` networks
2. Simplify getting access to `private` services for allowed end users
3. Potentially replace VPN services in some cases


## How it works

## Examples


## Roadmap


