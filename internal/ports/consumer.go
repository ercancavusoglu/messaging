package ports

type Consumer interface {
	Start() error
	Stop()
}
