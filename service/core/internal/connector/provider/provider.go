package provider

type Provider interface {
	GetType() string
	GetHeight() (uint64, error)
}
