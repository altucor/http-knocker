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

func (ctx basicComment) GetPrefix() string {
	return ctx.prefix
}

func (ctx basicComment) GetControllerName() string {
	return ctx.controllerName
}

func (ctx basicComment) GetTimestamp() time.Time {
	return ctx.timestamp
}

func (ctx basicComment) GetEndpointHash() string {
	return ctx.endpointHash
}
