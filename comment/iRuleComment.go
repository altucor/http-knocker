package comment

import "time"

type IRuleComment interface {
	FromString(comment string) error
	ToString() string
	IsSameFamily(other string) bool
	GetDelimiter() string
	SetDelimiter(delimiter string)
	GetPrefix() string
	SetPrefix(prefix string)
	GetControllerName() string
	SetControllerName(controller string)
	GetEndpointHash() string
	SetEndpointHash(endpointHash string)
	GetTimestamp() time.Time
	SetTimestamp(timestamp time.Time)
}
