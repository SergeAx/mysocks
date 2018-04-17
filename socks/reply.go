package socks

import (
	"net"
	"syscall"
)

func (c *connection) sendReply(status byte) {
	reply := []byte{socks5version, status, 0x00}
	hostName, hostPort, _ := splitHostPort(c.conn.LocalAddr().String())
	ip := net.ParseIP(string(hostName))
	addrType := socks5AddressTypeIPv4
	if ip.To16() != nil {
		addrType = socks5AddressTypeIPv6
	}
	reply = append(reply, addrType)
	reply = append(reply, ip...)
	reply = append(reply, hostPort...)
	c.conn.Write(reply)
}

func (c *connection) sendAuthReply(status byte) {
	c.conn.Write([]byte{socks5version, status})
}

func (c *connection) sendUDPReply(request *request) {
	reply := []byte{0, 0}
	fragment := byte(0)
	reply = append(reply, fragment)
	reply = append(reply, request.addrType)
	reply = append(reply, request.host...)
	reply = append(reply, request.port...)
	c.conn.Write(reply)
}

func (c *connection) sendReplyWithError(err error) {
	switch err := err.(type) {
	case net.Error:
		if err.Timeout() {
			c.sendReply(socks5StatusHostUnreachable)
		}
	case *net.OpError:
		switch err.Op {
		case "dial":
			c.sendReply(socks5StatusHostUnreachable)
		case "read":
			c.sendReply(socks5StatusConnectionRefused)
		}
	case syscall.Errno:
		switch err {
		case syscall.ECONNREFUSED:
			c.sendReply(socks5StatusConnectionRefused)
		}
	default:
		switch err {
		case ErrAddressTypeNotSupported:
			c.sendReply(socks5StatusAddressTypeNotSupported)
		case ErrCommandNotSupported:
			c.sendReply(socks5StatusCommandNotSupported)
		case ErrAuthMethodNotSupported:
			c.sendAuthReply(socks5AuthMethodNoAcceptable)
		case ErrAuthenticationFailure:
			// do nothing, since the handshake() has already responded to the client
		default:
			c.sendReply(socks5StatusGeneral)
		}
	}
}
