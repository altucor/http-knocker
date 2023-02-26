package deviceCommand

import (
	"encoding/json"
	"fmt"
	"httpKnocker/deviceCommon"
	"httpKnocker/logging"
	"net/http"
)

type Move struct {
	cmdType deviceCommon.DeviceCommandType
	from    uint64
	to      uint64
}

func MoveNew(from uint64, to uint64) Move {
	cmd := Move{
		cmdType: deviceCommon.DeviceCommandMove,
		from:    from,
		to:      to,
	}

	return cmd
}

func (ctx Move) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Move) Rest() (string, string, string, error) {
	// ip/firewall/filter/move numbers=4 destination=0
	method := http.MethodPost
	url := "/ip/firewall/filter/move"

	bodyMap := make(map[string]string)
	bodyMap["numbers"] = fmt.Sprintf("*%X", ctx.from)
	bodyMap["destination"] = fmt.Sprintf("*%X", ctx.to)

	body, err := json.Marshal(bodyMap)
	if err != nil {
		logging.CommonLog().Errorf("[DeviceCommand Move] error marshaling data to REST: %s\n", err)
		return "", "", "", err
	}

	return method, url, string(body), nil
}

func (ctx Move) IpTables() (string, error) {
	return "", nil
}
