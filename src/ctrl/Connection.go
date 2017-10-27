package ctrl

import "net"

type connectionConf struct{

	network, address string
}

func NewConnectionConf(network, address string) net.Addr{

	return net.Addr(&connectionConf{network, address})
}

func (c *connectionConf)Network() string{
	return c.network
}

func (c *connectionConf)String() string{
	return c.address
}
