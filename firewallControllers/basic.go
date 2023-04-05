package firewallControllers

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/endpoint"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type controllerBasicComment struct {
	delimiterKey string
	prefix       string
	firewallName string
	timestamp    time.Time
	endpointHash string
}

func FirewallCommentNew(delimiterKey string, prefix string, firewallName string, timestamp time.Time, endpointHash string) (controllerBasicComment, error) {
	if strings.ContainsAny(prefix, delimiterKey) ||
		strings.ContainsAny(firewallName, delimiterKey) ||
		strings.ContainsAny(endpointHash, delimiterKey) {
		return controllerBasicComment{}, errors.New("comment parameter cannot have delimiter key")
	}
	comment := controllerBasicComment{
		delimiterKey: delimiterKey,
		prefix:       prefix,
		firewallName: firewallName,
		timestamp:    timestamp,
		endpointHash: endpointHash,
	}
	return comment, nil
}

func FirewallCommentNewFromString(comment string, delimiterKey string) (controllerBasicComment, error) {
	commentParts := strings.Split(comment, delimiterKey)
	if len(commentParts) != 4 {
		// Do not report here errors, because not all rules acan have valid comment structure
		return controllerBasicComment{}, nil
	}
	timestamp, err := strconv.ParseInt(commentParts[2], 10, 64)
	if err != nil {
		return controllerBasicComment{}, err
	}
	commentObj := controllerBasicComment{
		delimiterKey: delimiterKey,
		prefix:       commentParts[0],
		firewallName: commentParts[1],
		timestamp:    time.Unix(timestamp, 0),
		endpointHash: commentParts[3],
	}

	return commentObj, nil
}

func (ctx controllerBasicComment) build() string {
	comment := ctx.prefix + ctx.delimiterKey
	comment += ctx.firewallName + ctx.delimiterKey
	comment += fmt.Sprintf("%d", ctx.timestamp.Unix()) + ctx.delimiterKey
	comment += ctx.endpointHash
	return comment
}

func (ctx controllerBasicComment) getPrefix() string {
	return ctx.prefix
}

func (ctx controllerBasicComment) getFirewallName() string {
	return ctx.firewallName
}

func (ctx controllerBasicComment) getTimestamp() time.Time {
	return ctx.timestamp
}

func (ctx controllerBasicComment) getEndpointHash() string {
	return ctx.endpointHash
}

type ClientAdded struct {
	Id    uint64
	Added time.Time
}

type SafeAddedClientsStorage struct {
	mu      sync.Mutex
	clients []ClientAdded
}

type ControllerCfg struct {
	DropRuleComment string `yaml:"drop-rule-comment"`
	Device          string `yaml:"device"`
}

type controllerBasic struct {
	watchdogRunning       bool
	prefix                string
	name                  string
	device                devices.IDevice
	endpoint              *endpoint.Endpoint
	controllerCfg         ControllerCfg
	addedClients          SafeAddedClientsStorage
	hash                  string
	delimiterKey          string
	needUpdateClientsList bool
}

// additional arguments - dev devices.IDevice, cfg endpoint.Endpoint
func ControllerBasicNew(controllerCfg ControllerCfg) *controllerBasic {
	ctx := controllerBasic{
		watchdogRunning:       false,
		name:                  "basicfirewall",
		prefix:                "httpKnocker",
		device:                nil,
		endpoint:              nil,
		controllerCfg:         controllerCfg,
		delimiterKey:          "-",
		needUpdateClientsList: true,
	}
	return &ctx
}

func ControllerBasicNewFromYaml(value *yaml.Node) (*controllerBasic, error) {
	var temp struct {
		Config ControllerCfg `yaml:"config"`
	}
	if err := value.Decode(&temp); err != nil {
		return nil, err
	}
	return ControllerBasicNew(temp.Config), nil
}

func (ctx *controllerBasic) GetDeviceName() string {
	return ctx.controllerCfg.Device
}

func (ctx *controllerBasic) SetDevice(dev devices.IDevice) {
	ctx.device = dev
}

func (ctx *controllerBasic) GetEndpointName() string {
	return ctx.controllerCfg.Device
}

func (ctx *controllerBasic) SetEndpoint(endpoint *endpoint.Endpoint) {
	ctx.endpoint = endpoint
}

func (ctx *controllerBasic) Start() error {
	logging.CommonLog().Info("[ControllerBasic] Starting...")
	go ClientsWatchdog(ctx)
	ctx.watchdogRunning = true
	logging.CommonLog().Info("[ControllerBasic] Starting... DONE")
	return nil
}

func (ctx *controllerBasic) Stop() error {
	logging.CommonLog().Info("[ControllerBasic] Stopping...")
	ctx.watchdogRunning = false
	logging.CommonLog().Info("[ControllerBasic] Stopping... DONE")
	return nil
}

func (ctx *controllerBasic) GetDevice() devices.IDevice {
	return ctx.device
}

func (ctx *controllerBasic) GetEndpoint() endpoint.Endpoint {
	return *ctx.endpoint
}

func (ctx *controllerBasic) AddClient(ip_addr firewallField.Address) error {
	// First of all check is client with src-address already present
	// to prevent duplication rules for one ip addr
	frwRules, err := ctx.GetRules()
	if err != nil {
		logging.CommonLog().Error("[ControllerBasic] AddClient cannot check is client exist: %s", err)
		return err
	}
	for _, element := range frwRules {
		// TODO: Fix deduplication for SSH
		if element.SrcAddress == ip_addr {
			_, err := ctx.device.RunCommandWithReply(command.RemoveNew(element.Id.GetValue()))
			if err != nil {
				logging.CommonLog().Error("[ControllerBasic] AddClient error removing client with duplicated src-address: %s", err)
				return err
			}
		}
	}

	dropRuleId, err := ctx.GetDropRuleId()
	if err != nil {
		logging.CommonLog().Error("[ControllerBasic] AddClient cannot find drop rule id %d", dropRuleId)
		return err
	}

	comment, err := FirewallCommentNew(
		ctx.delimiterKey,
		ctx.prefix,
		ctx.name,
		time.Now(),
		ctx.endpoint.GetHash(),
	)
	if err != nil {
		return err
	}
	frwCmdAdd := command.AddNew(
		ip_addr.GetValue(),
		ctx.endpoint.Port,
		ctx.endpoint.Protocol.GetValue(),
		comment.build(),
		dropRuleId,
	)
	_, err = ctx.device.RunCommandWithReply(frwCmdAdd)
	if err != nil {
		logging.CommonLog().Error("[ControllerBasic] Add command execution error: %s", err)
		return err
	}
	ctx.needUpdateClientsList = true
	return err
}

func (ctx *controllerBasic) GetRules() ([]firewallCommon.FirewallRule, error) {
	getResponse, err := ctx.device.RunCommandWithReply(command.GetNew())
	if err != nil {
		return nil, err
	}
	return getResponse.(response.Get).GetRules(), nil
}

func (ctx *controllerBasic) FindRuleIdByComment(comment string) (uint64, error) {
	frwRules, err := ctx.GetRules()
	if err != nil {
		return 0, err
	}
	for _, element := range frwRules {
		if element.Comment.GetValue() == comment {
			return element.Id.GetValue(), nil
		}
	}

	logging.CommonLog().Error("[ControllerBasic] FindRuleIdByComment Cannot find target rule")
	return 0, errors.New("cannot find target rule")
}

func (ctx *controllerBasic) GetDropRuleId() (uint64, error) {
	return ctx.FindRuleIdByComment(ctx.controllerCfg.DropRuleComment)
}

func (ctx *controllerBasic) IsClientWithAddrExist(ip_addr netip.Addr) (bool, error) {
	frwRules, err := ctx.GetRules()
	if err != nil {
		return false, err
	}
	for _, element := range frwRules {
		if element.SrcAddress.GetValue() == ip_addr {
			return true, nil
		}
	}
	return false, nil
}

func (ctx *controllerBasic) GetAddedClientIdsWithTimings() ([]ClientAdded, error) {
	var clientIds []ClientAdded
	frwRules, err := ctx.GetRules()
	if err != nil {
		return clientIds, err
	}

	for _, element := range frwRules {
		if element.Comment.GetValue() != "" {
			comment, err := FirewallCommentNewFromString(element.Comment.GetValue(), ctx.delimiterKey)
			if err != nil {
				logging.CommonLog().Errorf("Error parsing comment %s", element.Comment.GetValue())
			}
			if comment.getPrefix() == ctx.prefix &&
				comment.getFirewallName() == ctx.name &&
				comment.getEndpointHash() == ctx.endpoint.GetHash() {
				clientIds = append(clientIds, ClientAdded{
					Id:    element.Id.GetValue(),
					Added: comment.getTimestamp(),
				})
			}
		}
	}
	return clientIds, nil
}

func (ctx *controllerBasic) CleanupExpiredClients() error {
	clients, err := ctx.GetAddedClientIdsWithTimings()
	if err != nil {
		return err
	}

	for _, element := range clients {
		logging.CommonLog().Info("Rule diff timestamp: %d Max duration: %d", time.Since(element.Added), ctx.endpoint.DurationSeconds.GetValue())
		if time.Since(element.Added) > ctx.endpoint.DurationSeconds.GetValue() {
			_, err := ctx.device.RunCommandWithReply(command.RemoveNew(element.Id))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ctx *controllerBasic) InitListOfAddedClients() {
	clientsAdded, err := ctx.GetAddedClientIdsWithTimings()
	if err != nil {
		logging.CommonLog().Error("[controllerBasic] InitListOfAddedClients error getting list of clients")
	} else {
		ctx.addedClients.mu.Lock()
		ctx.addedClients.clients = clientsAdded
		ctx.addedClients.mu.Unlock()
	}
}

func ClientsWatchdog(firewall *controllerBasic) {
	for {
		if !firewall.watchdogRunning {
			return
		}
		if firewall.needUpdateClientsList {
			firewall.needUpdateClientsList = false
			firewall.InitListOfAddedClients()
		}
		time.Sleep(time.Second)
		logging.CommonLog().Debugf("[controllerBasic] ClientsWatchdog worked %d", uint64(time.Now().Unix()))
		firewall.addedClients.mu.Lock()
		clientsLength := len(firewall.addedClients.clients)
		firewall.addedClients.mu.Unlock()
		if clientsLength == 0 {
			continue
		}
		logging.CommonLog().Debugf("Clients: %v", firewall.addedClients.clients)
		firewall.addedClients.mu.Lock()
		for index, element := range firewall.addedClients.clients {
			curTime := time.Now()
			logging.CommonLog().Debugf("Rule diff timestamp: %f Max duration: %d",
				curTime.Sub(element.Added).Seconds(),
				firewall.endpoint.DurationSeconds.GetSeconds(),
			)
			timeRemaining := firewall.endpoint.DurationSeconds.GetValue() - curTime.Sub(element.Added)
			logging.CommonLog().Infof("Rule %d time remaining: %f sec", element.Id, timeRemaining.Seconds())
			if curTime.Sub(element.Added) > firewall.endpoint.DurationSeconds.GetValue() {
				_, err := firewall.device.RunCommandWithReply(command.RemoveNew(element.Id))
				if err != nil {
					logging.CommonLog().Errorf("[controllerBasic] ClientsWatchdog error removing client %s", err)
				} else {
					// In case of success regenerate local list of added clients
					logging.CommonLog().Infof("[controllerBasic] ClientsWatchdog Removed client from pending list: %v",
						firewall.addedClients.clients[index],
					)
					// If we removed some client than better to break cycle and try again next time
					// In other case we can hit out of bounds after removing and iterating via modified clients array
					// And also list of indexes of added clients on our side will be invalid in comparison to state on remote firewall
					firewall.needUpdateClientsList = true
					break
				}
			}
		}
		firewall.addedClients.mu.Unlock()
	}
}
