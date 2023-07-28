package comment

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type basicComment struct {
	delimiterKey   string
	prefix         string
	controllerName string
	timestamp      time.Time
	endpointHash   string
}

func BasicCommentNew() *basicComment {
	comment := basicComment{}
	return &comment
}

func (ctx *basicComment) containsDelimiterKey(input string) bool {
	return strings.Contains(input, ctx.delimiterKey)
}

func BasicCommentNewFromData(delimiterKey string, prefix string, controllerName string, timestamp time.Time, endpointHash string) (basicComment, error) {
	if strings.ContainsAny(prefix, delimiterKey) ||
		strings.ContainsAny(controllerName, delimiterKey) ||
		strings.ContainsAny(endpointHash, delimiterKey) {
		return basicComment{}, errors.New("comment parameter cannot have delimiter key")
	}
	comment := basicComment{
		delimiterKey:   delimiterKey,
		prefix:         prefix,
		controllerName: controllerName,
		timestamp:      timestamp,
		endpointHash:   endpointHash,
	}
	return comment, nil
}

func (ctx *basicComment) FromString(comment string) error {
	commentParts := strings.Split(comment, ctx.delimiterKey)
	if len(commentParts) != 4 {
		// Do not report here errors, because not all rules can have valid comment structure
		return nil
	}
	timestamp, err := strconv.ParseInt(commentParts[2], 10, 64)
	if err != nil {
		return err
	}
	ctx.prefix = commentParts[0]
	ctx.controllerName = commentParts[1]
	ctx.timestamp = time.Unix(timestamp, 0)
	ctx.endpointHash = commentParts[3]

	return nil
}

func BasicCommentNewFromString(comment string, delimiter string) (basicComment, error) {
	obj := basicComment{
		delimiterKey: delimiter,
	}
	err := obj.FromString(comment)
	return obj, err
}

func (ctx basicComment) ToString() string {
	comment := ctx.prefix + ctx.delimiterKey
	comment += ctx.controllerName + ctx.delimiterKey
	comment += fmt.Sprintf("%d", ctx.timestamp.Unix()) + ctx.delimiterKey
	comment += ctx.endpointHash
	return comment
}

func (ctx basicComment) IsSameFamily(other string) bool {
	otherObj, err := BasicCommentNewFromString(other, ctx.delimiterKey)
	if err != nil {
		return false
	}
	if otherObj.GetPrefix() != ctx.GetPrefix() {
		return false
	}
	if otherObj.GetControllerName() != ctx.GetControllerName() {
		return false
	}
	if otherObj.GetEndpointHash() != ctx.GetEndpointHash() {
		return false
	}
	return true
}

func (ctx *basicComment) SetPrefix(prefix string) {
	ctx.prefix = prefix
}

func (ctx basicComment) GetPrefix() string {
	return ctx.prefix
}

func (ctx *basicComment) SetControllerName(controller string) {
	ctx.controllerName = controller
}

func (ctx basicComment) GetControllerName() string {
	return ctx.controllerName
}

func (ctx *basicComment) SetTimestamp(timestamp time.Time) {
	ctx.timestamp = timestamp
}

func (ctx basicComment) GetTimestamp() time.Time {
	return ctx.timestamp
}

func (ctx *basicComment) SetEndpointHash(endpointHash string) {
	ctx.endpointHash = endpointHash
}

func (ctx basicComment) GetEndpointHash() string {
	return ctx.endpointHash
}

func (ctx basicComment) GetDelimiter() string {
	return ctx.delimiterKey
}

func (ctx *basicComment) SetDelimiter(delimiter string) {
	ctx.delimiterKey = delimiter
}
