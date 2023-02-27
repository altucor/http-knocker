# Example config for ssh device
<pre>
server:
  host: 0.0.0.0
  port: 8000
  default-response-code: 404

firewalls:
  firewall-second:
    firewallType: firewallBasic
    dropRuleCommnet: http-knocker-drop-all-rule
    protocol: iptables
    deviceSsh:
      username: <ssh device username>
      password: <ssh device password>
      host: <ip address or domain name of firewall>
      port: <port of ssh service>
      knownHosts: "/test/path/known_hosts"

endpoints:
  endpoint-second:
    url: a123
    duration: 1h
    port: 2222
    protocol: udp #tcp/udp/icmp
    response-code-on-success: 200

knocks:
  knock-second:
    firewall: firewall-second
    enpoint: endpoint-second
</pre>
