package rabbitmq

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	VHost    string
}

func (c *Config) URL() string {
	return "amqp://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.VHost
}
