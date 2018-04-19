package socks

import (
	"net"
	"bufio"
	"github.com/aspcartman/mysocks/monitoring"
)

var metricRequests = monitoring.NewEventMetric("requests", nil, "command", "addr", "in", "out", "err")

type connection struct {
	conn net.Conn
	conf *Config

	buf *bufio.Reader
}

func (c connection) Handle() {
	c.buf = bufio.NewReader(c.conn)
	defer c.conn.Close()

	var (
		err error
		req *request
		st  stat
	)

	t := metricRequests.Start()
	defer func() {
		var addr string
		var cmd int
		var errs string
		if req != nil {
			cmd = int(req.command)
			addr = req.Address()
		}
		if err != nil {
			errs = err.Error()
		}
		metricRequests.Stop(t, cmd, addr, st.in, st.out, errs)
	}()

	if err = c.verifyVersion(); err != nil {
		c.error("version verification", err)
		return
	}

	if err = c.handshake(); err != nil {
		c.error("handshake", err)
		return
	}

	if req, err = c.readRequest(); err != nil {
		c.error("reading request", err)
		return
	} else if req.version != socks5version {
		err = ErrVersionError
		c.error("reading request", err)
		return
	}

	switch req.command {
	case commandConnect:
		st, err = c.handleConnect(req)
	case commandUDPAssociate:
		st, err = c.handleUDPAssociate(req)
	default:
		err = ErrCommandNotSupported
	}

	if err != nil {
		c.error("command handling", err)
	}
}

func (c *connection) verifyVersion() error {
	b, err := c.buf.ReadByte()
	if err != nil {
		return err
	}

	if b != socks5version {
		return ErrSocks4VersionError
	}

	return nil
}

func (c *connection) error(stage string, err error) {
	c.sendReplyWithError(err)
	c.conf.HandleError(stage, err)
}

func (c *connection) readRequest() (*request, error) {
	var err error

	request := &request{}
	if request.version, err = c.buf.ReadByte(); err != nil {
		return nil, err
	}
	if request.command, err = c.buf.ReadByte(); err != nil {
		return nil, err
	}
	if _, err = c.buf.ReadByte(); err != nil {
		return nil, err
	}
	if request.address, err = c.readAddr(); err != nil {
		return nil, err
	}

	return request, nil
}

func (c *connection) readAddr() (*address, error) {
	var addr address
	var err error

	if addr.addrType, err = c.buf.ReadByte(); err != nil {
		return nil, err
	}

	switch addr.addrType {
	case socks5AddressTypeIPv4, socks5AddressTypeIPv6:
		length := net.IPv4len
		if addr.addrType == socks5AddressTypeIPv6 {
			length = net.IPv6len
		}
		addr.host = make([]byte, length)
		if _, err = c.buf.Read(addr.host); err != nil {
			return nil, err
		}
	case socks5AddressTypeFQDN:
		length, err := c.buf.ReadByte()
		if err != nil {
			return nil, err
		}
		addr.host = make([]byte, length)
		if _, err = c.buf.Read(addr.host); err != nil {
			return nil, err
		}
	default:
		return nil, ErrAddressTypeNotSupported
	}
	addr.port = make([]byte, 2)
	if _, err = c.buf.Read(addr.port); err != nil {
		return nil, err
	}

	return &addr, nil
}
