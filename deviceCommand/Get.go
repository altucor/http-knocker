package deviceCommand

import (
	"net/http"

	"github.com/altucor/http-knocker/deviceCommon"
)

type Get struct {
	cmdType deviceCommon.DeviceCommandType
}

func GetNew() Get {
	cmd := Get{
		cmdType: deviceCommon.DeviceCommandGet,
	}
	return cmd
}

func (ctx Get) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Get) Rest() (string, string, string, error) {
	method := http.MethodGet
	url := "/ip/firewall/filter"
	body := ""

	return method, url, body, nil
}

func (ctx Get) IpTables() (string, error) {
	return "iptables -S INPUT", nil
}
