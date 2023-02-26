package firewalls

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"httpKnocker/common"
	"httpKnocker/deviceCommand"
	"httpKnocker/deviceCommon"
	"httpKnocker/devices"
	"httpKnocker/firewallCommon/firewallField"
	"httpKnocker/logging"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"time"
)

type firewallBasicComment struct {
	delimiterKey string
	prefix       string
	firewallName string
	timestamp    uint64
	endpointHash string
}

func FirewallCommentNew(delimiterKey string, prefix string, firewallName string, timestamp uint64, endpointHash string) firewallBasicComment {
	comment := firewallBasicComment{
		delimiterKey: delimiterKey,
		prefix:       prefix,
		firewallName: firewallName,
		timestamp:    timestamp,
		endpointHash: endpointHash,
	}
	return comment
}

func FirewallCommentNewFromString(comment string, delimiterKey string) firewallBasicComment {
	commentParts := strings.Split(comment, delimiterKey)
	if len(commentParts) != 4 {
		return firewallBasicComment{}
	}
	timestamp, _ := strconv.ParseInt(commentParts[2], 10, 64)
	commentObj := firewallBasicComment{
		delimiterKey: delimiterKey,
		prefix:       commentParts[0],
		firewallName: commentParts[1],
		timestamp:    uint64(timestamp),
		endpointHash: commentParts[3],
	}

	return commentObj
}

func (ctx firewallBasicComment) build() string {
	comment := ctx.prefix + ctx.delimiterKey
	comment += ctx.firewallName + ctx.delimiterKey
	comment += fmt.Sprintf("%d", ctx.timestamp) + ctx.delimiterKey
	comment += ctx.endpointHash
	return comment
}

func (ctx firewallBasicComment) getPrefix() string {
	return ctx.prefix
}

func (ctx firewallBasicComment) getFirewallName() string {
	return ctx.firewallName
}

func (ctx firewallBasicComment) getTimestamp() uint64 {
	return ctx.timestamp
}

func (ctx firewallBasicComment) getEndpointHash() string {
	return ctx.endpointHash
}

type ClientAdded struct {
	Id    uint64
	Added uint64
}

type SafeAddedClientsStorage struct {
	mu      sync.Mutex
	clients []ClientAdded
}

type firewallBasic struct {
	prefix       string
	name         string
	device       devices.IDevice
	endpoint     common.EndpointCfg
	firewallCfg  common.FirewallCfg
	addedClients SafeAddedClientsStorage
	hash         string
	// comment               string
	delimiterKey          string
	needUpdateClientsList bool
}

func FirewallBasicNew(dev devices.IDevice, cfg common.EndpointCfg, firewallCfg common.FirewallCfg) *firewallBasic {
	ctx := firewallBasic{
		name:                  "basicfirewall",
		prefix:                "httpknocker",
		device:                dev,
		endpoint:              cfg,
		firewallCfg:           firewallCfg,
		delimiterKey:          "-",
		needUpdateClientsList: true,
	}
	h := sha1.New()
	h.Write([]byte(ctx.endpoint.Url + fmt.Sprintf("%d", ctx.endpoint.DurationSeconds) + fmt.Sprint(ctx.endpoint.Port) + ctx.endpoint.Protocol.GetString()))
	ctx.hash = fmt.Sprintf("%x", h.Sum(nil))
	// ctx.comment = "http-knocker-rule-" + ctx.hash
	// Syntax of basic firewall rule:
	// http-knocker-rule-1436773875-438hcor4gho34thc3ch2843h8t2
	// logging.CommonLog().Info("[firewallBasic] hashed addrlist:", ctx.comment)
	// ctx.firewallCfg.Protocol
	go ClientsWatchdog(&ctx)
	return &ctx
}

func (ctx *firewallBasic) GetDevice() devices.IDevice {
	return ctx.device
}

func (ctx *firewallBasic) GetEndpoint() common.EndpointCfg {
	return ctx.endpoint
}

func (ctx *firewallBasic) AddClient(ip_addr firewallField.Address) error {
	// First of all check is client with src-address already present
	// to prevent duplication rules for one ip addr
	frwRules, err := ctx.GetRules()
	if err != nil {
		logging.CommonLog().Error("[FirewallBasic] AddClient cannot check is client exist: %s\n", err)
		return err
	}
	for _, element := range frwRules.GetRules().GetList() {
		// TODO: Fix deduplication for SSH
		if element.SrcAddress == ip_addr {
			_, err := ctx.device.RunCommandWithReply(deviceCommand.RemoveNew(element.Id.GetValue()), ctx.firewallCfg.Protocol)
			if err != nil {
				logging.CommonLog().Error("[FirewallBasic] AddClient error removing client with duplicated src-address: %s\n", err)
				return err
			}
		}
	}

	dropRuleId, err := ctx.GetDropRuleId()
	if err != nil {
		logging.CommonLog().Error("[FirewallBasic] AddClient cannot find drop rule id %d\n", dropRuleId)
		return err
	}

	addedTimestamp := uint64(time.Now().Unix())
	comment := FirewallCommentNew(
		ctx.delimiterKey,
		ctx.prefix,
		ctx.name,
		addedTimestamp,
		ctx.hash,
	)
	frwCmdAdd := deviceCommand.AddNew(
		ip_addr.GetValue(),
		ctx.endpoint.Port,
		ctx.endpoint.Protocol.GetValue(),
		ctx.endpoint.DurationSeconds.GetValue(),
		comment.build(),
		dropRuleId,
	)
	_, err = ctx.device.RunCommandWithReply(frwCmdAdd, ctx.firewallCfg.Protocol)
	if err != nil {
		logging.CommonLog().Error("[FirewallBasic] Add command execution error: %s\n", err)
		return err
	}
	ctx.needUpdateClientsList = true
	return err
}

func (ctx *firewallBasic) GetRules() (deviceCommon.IDeviceResponse, error) {
	return ctx.device.RunCommandWithReply(deviceCommand.GetNew(), ctx.firewallCfg.Protocol)
}

func (ctx *firewallBasic) FindRuleIdByComment(comment string) (uint64, error) {
	frwRules, err := ctx.GetRules()
	if err != nil {
		return 0, err
	}
	for _, element := range frwRules.GetRules().GetList() {
		if element.Comment.GetValue() == comment {
			return element.Id.GetValue(), nil
		}
	}

	logging.CommonLog().Error("[FirewallBasic] FindRuleIdByComment Cannot find target rule")
	return 0, errors.New("Cannot find target rule")
}

func (ctx *firewallBasic) GetDropRuleId() (uint64, error) {
	return ctx.FindRuleIdByComment(ctx.firewallCfg.DropRuleComment)
}

func (ctx *firewallBasic) IsClientWithAddrExist(ip_addr netip.Addr) (bool, error) {
	frwRules, err := ctx.GetRules()
	if err != nil {
		return false, err
	}
	for _, element := range frwRules.GetRules().GetList() {
		if element.SrcAddress.GetValue() == ip_addr {
			return true, nil
		}
	}
	return false, nil
}

func (ctx *firewallBasic) GetAddedClientIdsWithTimings() ([]ClientAdded, error) {
	var clientIds []ClientAdded
	frwRules, err := ctx.GetRules()
	if err != nil {
		return clientIds, err
	}

	for _, element := range frwRules.GetRules().GetList() {
		comment := element.Comment
		if comment.GetValue() != "" {
			comment := FirewallCommentNewFromString(comment.GetValue(), ctx.delimiterKey)
			if comment.getPrefix() == ctx.prefix && comment.getFirewallName() == ctx.name {
				clientIds = append(clientIds, ClientAdded{
					Id:    element.Id.GetValue(),
					Added: comment.getTimestamp(),
				})
			}
		}
	}
	return clientIds, nil
}

func (ctx *firewallBasic) CleanupExpiredClients() error {
	clients, err := ctx.GetAddedClientIdsWithTimings()
	if err != nil {
		return err
	}

	for _, element := range clients {
		logging.CommonLog().Info("Rule diff timestamp: %d Max duration: %d\n", uint64(time.Now().Unix())-element.Added, ctx.endpoint.DurationSeconds.GetValue())
		if uint64(time.Now().Unix())-element.Added > ctx.endpoint.DurationSeconds.GetValue() {
			_, err := ctx.device.RunCommandWithReply(deviceCommand.RemoveNew(element.Id), ctx.firewallCfg.Protocol)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ctx *firewallBasic) InitListOfAddedClients() {
	clientsAdded, err := ctx.GetAddedClientIdsWithTimings()
	if err != nil {
		logging.CommonLog().Error("[firewallBasic] InitListOfAddedClients error getting list of clients")
	} else {
		ctx.addedClients.mu.Lock()
		ctx.addedClients.clients = clientsAdded
		ctx.addedClients.mu.Unlock()
	}
}

func ClientsWatchdog(firewall *firewallBasic) {
	for true {
		if firewall.needUpdateClientsList {
			firewall.needUpdateClientsList = false
			firewall.InitListOfAddedClients()
		}
		time.Sleep(time.Second)
		logging.CommonLog().Infof("[firewallBasic] ClientsWatchdog worked %d\n", uint64(time.Now().Unix()))
		firewall.addedClients.mu.Lock()
		clientsLength := len(firewall.addedClients.clients)
		firewall.addedClients.mu.Unlock()
		if clientsLength == 0 {
			continue
		}
		logging.CommonLog().Debugf("Clients: %v\n", firewall.addedClients.clients)
		firewall.addedClients.mu.Lock()
		for index, element := range firewall.addedClients.clients {
			logging.CommonLog().Debugf("Rule diff timestamp: %d Max duration: %d\n",
				uint64(time.Now().Unix())-element.Added,
				firewall.endpoint.DurationSeconds.GetValue(),
			)
			timeRemaining := firewall.endpoint.DurationSeconds.GetValue() - (uint64(time.Now().Unix()) - element.Added)
			logging.CommonLog().Infof("Rule %d time remaining: %d sec\n", element.Id, timeRemaining)
			if uint64(time.Now().Unix())-element.Added > firewall.endpoint.DurationSeconds.GetValue() {
				_, err := firewall.device.RunCommandWithReply(deviceCommand.RemoveNew(element.Id), firewall.firewallCfg.Protocol)
				if err != nil {
					logging.CommonLog().Errorf("[firewallBasic] ClientsWatchdog error removing client %s\n", err)
				} else {
					// In case of success regenerate local list of added clients
					logging.CommonLog().Infof("[firewallBasic] ClientsWatchdog Removed client from pending list: %v\n",
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
