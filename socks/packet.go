package socks

import (
	"encoding/binary"
	"net"
	"strconv"
)

type address struct {
	addrType byte
	host     []byte
	port     []byte
}

type request struct {
	version byte
	command byte
	*address
}

func (r *request) ToPacket() []byte {
	packet := []byte{
		r.version,
		r.command,
		0x00,
		r.addrType,
	}
	packet = append(packet, r.host...)
	packet = append(packet, r.port...)
	return packet
}

func (r *request) Address() string {
	var host string
	switch r.addrType {
	case socks5AddressTypeIPv4, socks5AddressTypeIPv6:
		host = net.IP(r.host).String()
	case socks5AddressTypeFQDN:
		host = string(r.host)
	}
	port := strconv.Itoa(int(binary.BigEndian.Uint16(r.port)))
	return net.JoinHostPort(host, port)
}
