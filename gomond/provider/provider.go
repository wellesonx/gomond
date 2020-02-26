package provider

type Provider interface {
	Start() error
	Close() error
	Follow(chan []byte) error
}
