package socks

import (
	"io"
	"sync"
)

func (c *connection) handleConnect(request *request) (err error) {
	remoteConn, err := c.conf.Dial("tcp", request.Address())
	if err != nil {
		return err
	}

	defer remoteConn.Close()

	c.sendReply(socks5StatusSucceeded)

	return proxy(c.conn, remoteConn)
}

func (c *connection) handleUDPAssociate(request *request) error {
	remoteConn, err := c.conf.Dial("udp", request.Address())
	if err != nil {
		return err
	}

	defer remoteConn.Close()

	c.sendUDPReply(request)

	return proxy(c.conn, remoteConn)
}

func proxy(one, two io.ReadWriter) error {
	var err1, err2 error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		_, err1 = io.Copy(one, two)
		wg.Done()
	}()
	go func() {
		_, err2 = io.Copy(two, one)
		wg.Done()
	}()

	wg.Wait()

	switch {
	case err1 != nil:
		return err1
	case err2 != nil:
		return err2
	}

	return nil
}
