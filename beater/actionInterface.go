package beater

// Action interface
type Action interface {
	Name() string
	Do() error
	Stop()
}
