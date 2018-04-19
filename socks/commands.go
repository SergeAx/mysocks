package socks

import (
	"io"
	"sync"
)

type stat struct {
	in, out int64
}

func (c *connection) handleConnect(request *request) (stat, error) {
	remoteConn, err := c.conf.Dial("tcp", request.Address())
	if err != nil {
		return stat{}, err
	}

	defer remoteConn.Close()

	c.sendReply(socks5StatusSucceeded)

	return proxy(c.conn, remoteConn)
}

func (c *connection) handleUDPAssociate(request *request) (stat, error) {
	remoteConn, err := c.conf.Dial("udp", request.Address())
	if err != nil {
		return stat{}, err
	}

	defer remoteConn.Close()

	c.sendUDPReply(request)

	return proxy(c.conn, remoteConn)
}

func proxy(one, two io.ReadWriter) (stat, error) {
	var w1, w2 int64
	var err1, err2 error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		w1, err1 = io.Copy(two, one)
		wg.Done()
	}()
	go func() {
		w2, err2 = io.Copy(one, two)
		wg.Done()
	}()

	wg.Wait()

	st := stat{w1, w2}
	switch {
	case err1 != nil:
		return st, err1
	case err2 != nil:
		return st, err2
	}

	return st, nil
}
