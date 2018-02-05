package addr

import "net"

func New(Protocol, Address string) net.Addr{

	return net.Addr( &Config{Protocol, Address} )
}

type Config struct{

	Protocol, Address string
}

func (c *Config)Network() string{
	return c.Protocol
}

func (c *Config)String() string{
	return c.Address
}
