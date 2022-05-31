package rick

type Config struct {
	//Addr  地址
	Addr string `envconfig:"RICK_ADDR"`
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	return &Client{cfg}
}
