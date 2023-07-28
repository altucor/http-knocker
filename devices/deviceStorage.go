package devices

import (
	"sync"

	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type DeviceStorage struct {
	devices map[string]interface{}
}

func (ctx *DeviceStorage) Init() {
	ctx.devices = make(map[string]interface{})
	ctx.devices["rest"] = DeviceRestNewFromYaml
	ctx.devices["ssh"] = DeviceSshNewFromYaml
	ctx.devices["puller"] = DevicePullerNewFromYaml
	ctx.devices["router-os"] = DeviceRouterOsNewFromYaml
}

func (ctx *DeviceStorage) GetDevice(name string, node *yaml.Node) (IDevice, error) {
	if _, ok := ctx.devices[name]; !ok {
		logging.CommonLog().Fatalf("[DeviceStorage] Cannot find device under name: '%s'", name)
	}
	return ctx.devices[name].(func(*yaml.Node) (IDevice, error))(node)
}

var lock = &sync.Mutex{}
var deviceStorage *DeviceStorage

func GetDeviceStorage() *DeviceStorage {
	lock.Lock()
	defer lock.Unlock()
	if deviceStorage == nil {
		deviceStorage = &DeviceStorage{}
		deviceStorage.Init()
	}

	return deviceStorage
}
