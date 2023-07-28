package firewallControllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/altucor/http-knocker/comment"
	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/endpoint"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type RemoveAnythingMatching struct {
	OlderThen        common.DurationSeconds `yaml:"older-then"`
	SameFirewallName bool                   `yaml:"same-firewall-name"`
	SamePrefix       bool                   `yaml:"same-prefix"`
}

type ControllerCfg struct {
	DropRuleComment        string                 `yaml:"drop-rule-comment"`
	RemoveAnythingMatching RemoveAnythingMatching `yaml:"remove-anything-matching"`
	UpdateRulesEvery       common.DurationSeconds `yaml:"update-rules-every"`
}

type controllerBasic struct {
	watchdogRunning  bool
	url              string
	comment          comment.IRuleComment
	deviceController *devices.DeviceController
	endpoint         *endpoint.Endpoint
	controllerCfg    ControllerCfg
}

// additional arguments - dev devices.IDevice, cfg endpoint.Endpoint
func ControllerBasicNew(controllerCfg ControllerCfg) *controllerBasic {
	ctx := controllerBasic{
		watchdogRunning:  false,
		comment:          comment.BasicCommentNew(),
		deviceController: nil,
		endpoint:         nil,
		controllerCfg:    controllerCfg,
	}
	ctx.comment.SetDelimiter("-")
	ctx.comment.SetPrefix("httpKnocker")
	ctx.comment.SetControllerName("basicController")
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

func (ctx *controllerBasic) SetUrl(url string) {
	ctx.url = url
}

func (ctx *controllerBasic) SetDevice(dev *devices.DeviceController) {
	ctx.deviceController = dev
}

func (ctx *controllerBasic) SetEndpoint(endpoint *endpoint.Endpoint) {
	ctx.endpoint = endpoint
	ctx.comment.SetEndpointHash(ctx.endpoint.GetHash(ctx.url))
}

func (ctx *controllerBasic) GetHash() string {
	return ctx.endpoint.GetHash(ctx.url)
}

func (ctx *controllerBasic) Start() error {
	logging.CommonLog().Info("[ControllerBasic] Starting...")
	// go UpdateRulesEveryThread(ctx)
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

func (ctx *controllerBasic) HttpCallbackAddClient(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[knock] accessing knock endpoint:", ctx.url)

	if clientAddr, err := ctx.endpoint.IpAddrSource.GetFromRequest(r); err != nil {
		logging.CommonLog().Error("[knock] Error getting client address:", err)
	} else {
		// Perform adding client in another thread
		// To be able response to HTTP client faster
		// And prevent timing attacks
		go ctx.AddClient(clientAddr)
		if ctx.endpoint.ResponseCodeOnSuccess != 0 {
			w.WriteHeader(int(ctx.endpoint.ResponseCodeOnSuccess))
			fmt.Fprintf(w, "%d\n", ctx.endpoint.ResponseCodeOnSuccess)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404\n")
		}
		//http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (ctx *controllerBasic) GetHttpCallback() (string, func(w http.ResponseWriter, r *http.Request)) {
	// ctx.endpoint.Middlewares["test"].Middleware.Register(ctx.HttpCallbackAddClient)
	return "/" + ctx.url, ctx.endpoint.RegisterMiddlewares(ctx.HttpCallbackAddClient)
}

func (ctx *controllerBasic) DeduplicationCleanup(ip_addr firewallField.Address) error {
	frwRules, err := ctx.deviceController.GetRulesFiltered(ctx.comment, false)
	if err != nil {
		logging.CommonLog().Error("[ControllerBasic] DeduplicationCleanup cannot check is client exist: ", err)
		return err
	}
	for _, element := range frwRules {
		// TODO: Fix deduplication for SSH
		if element.SrcAddress == ip_addr {
			_, err := ctx.deviceController.RunCommandWithReply(command.RemoveNew(element.Id.GetValue()))
			if err != nil {
				logging.CommonLog().Error("[ControllerBasic] DeduplicationCleanup error removing client with duplicated src-address: ", err)
				return err
			}
		}
	}
	return nil
}

func (ctx *controllerBasic) AddClient(ip_addr firewallField.Address) error {
	// First of all check is client with src-address already present
	// to prevent duplication rules for one ip addr

	err := ctx.DeduplicationCleanup(ip_addr)
	if err != nil {
		logging.CommonLog().Error("[ControllerBasic] AddClient cannot check is client exist: ", err)
		return err
	}

	// If drop rule comment is set
	// Then find it's id on device firewall
	// Otherwise ignore searching of it
	var dropRuleId uint64 = 0
	if ctx.controllerCfg.DropRuleComment != "" {
		dropRuleId, err = ctx.deviceController.FindRuleIdByComment(ctx.controllerCfg.DropRuleComment, false)
		if err != nil {
			logging.CommonLog().Error("[ControllerBasic] AddClient cannot find drop rule id")
			return err
		}
	}

	ctx.comment.SetTimestamp(time.Now())
	frwCmdAdd := command.AddNew(
		ip_addr.GetValue(),
		ctx.endpoint.Port,
		ctx.endpoint.Protocol.GetValue(),
		ctx.comment.ToString(),
		dropRuleId,
	)
	_, err = ctx.deviceController.RunCommandWithReply(frwCmdAdd)
	if err != nil {
		logging.CommonLog().Error("[ControllerBasic] Add command execution error:", err)
		return err
	}
	return err
}

func (ctx *controllerBasic) CheckIsRuleExpired(rule firewallCommon.FirewallRule) bool {
	commentObj, err := comment.BasicCommentNewFromString(rule.Comment.GetValue(), ctx.comment.GetDelimiter())
	if err != nil {
		logging.CommonLog().Error("Error parsing comment:", rule.Comment.GetValue())
		return false
	}
	curTime := time.Now()
	timeRemaining := ctx.endpoint.DurationSeconds.GetValue() - curTime.Sub(commentObj.GetTimestamp())
	logging.CommonLog().Debugf("[ControllerBasic] CheckIsRuleExpired Rule %d Diff timestamp: %f Max duration: %d Time remaining: %f sec",
		rule.Id.GetValue(),
		curTime.Sub(commentObj.GetTimestamp()).Seconds(),
		ctx.endpoint.DurationSeconds.GetSeconds(),
		timeRemaining.Seconds(),
	)
	return curTime.Sub(commentObj.GetTimestamp()) > ctx.endpoint.DurationSeconds.GetValue()
}

func (ctx *controllerBasic) CleanupExpiredClients() error {
	frwRules, err := ctx.deviceController.GetRulesFiltered(ctx.comment, true)
	if err != nil {
		return err
	}

	var pendingForRemove []uint64
	for _, rule := range frwRules {
		commentData, err := comment.BasicCommentNewFromString(rule.Comment.GetValue(), ctx.comment.GetDelimiter())
		if err != nil {
			logging.CommonLog().Error("[basicController] CleanupExpiredClients err:", err)
			continue
		}
		logging.CommonLog().Infof("Rule diff timestamp: %d Max duration: %d", time.Since(commentData.GetTimestamp()), ctx.endpoint.DurationSeconds.GetValue())
		if time.Since(commentData.GetTimestamp()) > ctx.endpoint.DurationSeconds.GetValue() {
			pendingForRemove = append(pendingForRemove, rule.Id.GetValue())
		}
	}
	return ctx.deviceController.RemoveBatch(pendingForRemove)
}

func ClientsWatchdog(firewall *controllerBasic) {
	for {
		if !firewall.watchdogRunning {
			return
		}
		time.Sleep(time.Second)
		// logging.CommonLog().Debugf("[controllerBasic] ClientsWatchdog worked %d", uint64(time.Now().Unix()))
		frwRules, err := firewall.deviceController.GetRulesFiltered(firewall.comment, false)
		if err != nil {
			logging.CommonLog().Error("[ControllerBasic] ClientsWatchdog error getting rules: ", err)
			continue
		}

		if len(frwRules) == 0 {
			continue
		}
		for _, element := range frwRules {
			if firewall.CheckIsRuleExpired(element) {
				err := firewall.CleanupExpiredClients()
				if err != nil {
					logging.CommonLog().Error("[ControllerBasic] Error cleaning expired clients:", err)
				}
			}
		}
	}
}
