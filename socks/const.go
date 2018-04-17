package socks

import "errors"

const (
	socks5version byte = 5

	commandConnect      byte = 1
	commandUDPAssociate byte = 3

	socks5AddressTypeIPv4 byte = 1
	socks5AddressTypeFQDN byte = 3
	socks5AddressTypeIPv6 byte = 4

	socks5StatusSucceeded               byte = 0
	socks5StatusGeneral                 byte = 1
	socks5StatusHostUnreachable         byte = 4
	socks5StatusConnectionRefused       byte = 5
	socks5StatusCommandNotSupported     byte = 7
	socks5StatusAddressTypeNotSupported byte = 8

	socks5AuthMethodNoRequired    byte = 0x00
	socks5AuthMethodPassword      byte = 0x02
	socks5AuthMethodTLSNoRequired byte = 0x80
	socks5AuthMethodTLSPassword   byte = 0x82
	socks5AuthMethodNoAcceptable  byte = 0xFF
)

var (
	ErrVersionError            = errors.New("version error")
	ErrSocks4VersionError      = errors.New("socks4 not supported")
	ErrCommandNotSupported     = errors.New("command not supported")
	ErrAddressTypeNotSupported = errors.New("address type not supported")
	ErrAuthMethodNotSupported  = errors.New("authentication method not supported")
	ErrAuthenticationFailure   = errors.New("authentication failure")
)
