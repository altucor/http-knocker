package comment

type IRuleComment interface {
	FromString(comment string) error
	ToString() string
}
