package main

import (
	"flag"

	"github.com/altucor/http-knocker/knocker"
	"github.com/altucor/http-knocker/logging"
)

/*

Devices:
- SSH
- RouterOS REST
- RouterOS API

Firewalls:
- RouterOS REST
- RouterOS API
- IpTables
- UFW
- PfSense
- Firewalld

Ideas:
- 1) Puller device with additional http endpoint
+ 2) Correct logging with different levels
+ 3) Correct error checking
+ 4) Allow to set in endpoint configuration from which parameter get user client like from headers or from GET parameter
- 5) Optional Basic auth for enpoint url. User should provide file with user:pass pairs generated with htpasswd
- 6) Add authentication option for endpoint through authelia

*/

func main() {
	logging.CommonLog().Info("app starting")
	configPath := flag.String("config-path", "", "path to YAML config file")
	flag.Parse()

	knocker, err := knocker.KnockerNewFromConfig(*configPath)
	if err != nil {
		logging.CommonLog().Error(err)
		return
	}
	knocker.Start()
	defer knocker.Stop()
	knocker.Wait()
	logging.CommonLog().Info("end of app")
}
