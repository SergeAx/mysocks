package socks

import (
	"strconv"

	"encoding/binary"
	"net"
)

func toPort(port string)  []byte {
	buffer := make([]byte, 2)
	if port, err := strconv.ParseUint(port, 10, 16); err == nil {
		binary.BigEndian.PutUint16(buffer, uint16(port))
		return buffer
	}
	return nil
}

func splitHostPort(addr string) (host, port []byte, err error) {
	hostName, hostPort, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}
	return []byte(hostName), toPort(hostPort), err
}
