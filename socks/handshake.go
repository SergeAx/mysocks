package socks

import "bytes"

func (c *connection) handshake() error {
	method, err := c.buf.ReadByte()
	if err != nil {
		return err
	}

	methods := make([]byte, method)
	if _, err = c.buf.Read(methods); err != nil {
		return err
	}

	if c.conf.Auth == nil {
		c.sendAuthReply(socks5AuthMethodNoRequired)
		return nil
	}

	if err = c.authBasedPassword(methods); err != nil {
		return err
	}

	return nil
}

func (c *connection) authBasedPassword(methods []byte) error {
	method := socks5AuthMethodPassword

	if !bytes.Contains(methods, []byte{method}) {
		return ErrAuthMethodNotSupported
	}

	c.sendAuthReply(method)

	if version, err := c.buf.ReadByte(); err != nil {
		return err
	} else if version != 0x01 {
		return ErrVersionError
	}

	usernameLength, err := c.buf.ReadByte()
	if err != nil {
		return err
	}
	username := make([]byte, usernameLength)
	if _, err = c.buf.Read(username); err != nil {
		return err
	}
	passwordLength, err := c.buf.ReadByte()
	if err != nil {
		return err
	}
	password := make([]byte, passwordLength)
	if _, err = c.buf.Read(password); err != nil {
		return err
	}

	if c.conf.Auth(string(username), string(password)) {
		c.conn.Write([]byte{0x01, 0x00})
	} else {
		c.conn.Write([]byte{0x01, 0x01})
		return ErrAuthenticationFailure
	}

	return nil
}

