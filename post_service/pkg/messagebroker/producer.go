package messagebroker

type Producer interface {
	Start() error
	Stop() error
	Produce(key, body []byte, logBody string) error
}