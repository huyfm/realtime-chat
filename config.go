package rtc

import (
	"os"
)

type Config struct {
	SrvHost      string
	SrvPort      string
	TLSCert      string
	TLSPrivKey   string
	DSN          string
	ClientID     string
	ClientSecret string
	HashKey      string
}

func DefaultConfig() (*Config, error) {
	c := Config{
		SrvPort:      os.Getenv("SRV_PORT"),
		TLSCert:      os.Getenv("TLS_CERT"),
		TLSPrivKey:   os.Getenv("TLS_PRIVKEY"),
		DSN:          os.Getenv("POSTGRES_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		HashKey:      os.Getenv("HASH_KEY"),
	}
	if c.TLSCert != "" {
		c.TLSCert = "/app/tls/" + c.TLSCert
	}
	if c.TLSPrivKey != "" {
		c.TLSPrivKey = "/app/tls/" + c.TLSPrivKey
	}
	return &c, nil
}
