# Example config for REST device
<pre>
server: # Here common configuration for server which handles connections to knock endpoints
  host: 0.0.0.0
  port: 8000
  default-response-code: 404

# Here array of firewall(devices with firewall additional info)
# Each of them describe one firewall device like:
1) Linux machine iwth iptables and remote access via SSH
2) MikroTik device with RouterOS and access via REST
# Also here is raw example puller device, which actually should host firewall rules for firewalls
# In future each firewall can connect by its own to puller endpoint and pull new firewall rules for him
firewalls:
  firewall-pull:
    firewallType: firewallPull
    devicePuller:
      username: <username>
      password: <password>
      port: 8001
      endpoint: puller/test
  firewall-rest-test:
    firewallType: firewallBasic
    dropRuleCommnet: http-knocker-drop-all-rule
    protocol: router-os-rest
    deviceRest:
      username: <username on router os>
      password: <password to router os user>
      endpoint: <https link to router os REST endpoint>

# Here is list of endpoints
# Each endpoint it is generally long URL which should be 
# taken in secure place and shared only with people whoc need access
# When client access endpoint(open url in browser or via curl)
# It triggers opening of specific port for specific time and on specific protocol
# The port will be opened for ip address frmo which URL was accessed
# There is also configurations which allows you to chose from where to read client ip address
# Also optioanlly you can use BasicAuth, so doesnt matter if your URL is leaked
# Clients anyway will be asked for password
# For endpoints which doesnt use authentication there is option response-code-on-success
# It alows you to hide your endpoint while attacker can try to bruteforce you
# So by setting response-code-on-success to 404. Attacker will receive 404 code even if he hit valid endpoint
endpoints:
  endpoint-second:
    ip-source:
      type: http-headers
      field-name: X-Http-Knocker-IP
    ip-source-type: http-request-param
    ip-source-field-name: targetAddr
    auth:
      auth-type: basic-auth
      # users-file: /Users/username/Desktop/credentials
      users:
        - "john:$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
    url: a123
    duration: 30s
    port: 2222
    protocol: udp # tcp/udp/icmp # /etc/protocols # See also http://www.iana.org/assignments/protocol-numbers
    response-code-on-success: 200

# Here list of knocks
# Knocks describes connection between firewall and endpoint
# So you can use same endpoint for several firewalls
# Or you can assign several endpoints for one device
# But for now it is not tested so good
knocks:
  knock-second:
    firewall: firewall-rest-test
    enpoint: endpoint-second
</pre>
