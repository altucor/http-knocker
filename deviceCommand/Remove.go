package deviceCommand

import (
	"fmt"
	"net/http"

	"github.com/altucor/http-knocker/deviceCommon"
)

type Remove struct {
	cmdType deviceCommon.DeviceCommandType
	ruleId  uint64
}

func RemoveNew(id uint64) Remove {
	frw := Remove{
		cmdType: deviceCommon.DeviceCommandRemove,
		ruleId:  id,
	}
	return frw
}

func (ctx Remove) ToMap() map[string]interface{} {
	cmd := make(map[string]interface{})
	cmd["type"] = string(ctx.cmdType)
	cmd["id"] = ctx.ruleId
	return cmd
}

func (ctx Remove) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Remove) GetId() uint64 {
	return ctx.ruleId
}

func (ctx Remove) Rest() (string, string, string, error) {
	method := http.MethodDelete
	url := "/ip/firewall/filter/*" + fmt.Sprintf("%X", ctx.ruleId)
	body := ""
	return method, url, body, nil
}

func (ctx Remove) IpTables() (string, error) {
	// iptables --delete INPUT 3
	return fmt.Sprintf("iptables --delete INPUT %d", ctx.ruleId), nil
}
