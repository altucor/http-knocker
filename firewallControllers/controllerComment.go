package firewallControllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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