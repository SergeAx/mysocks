package socks

import "net"

type Config struct {
	Auth        func(username, password string) bool
	Dial        func(network, address string) (net.Conn, error)
	HandleError func(string, error)
}

func Serve(listener net.Listener, conf *Config) {
	if conf.HandleError == nil {
		conf.HandleError = func(_ string, _ error) {}
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			conf.HandleError("listen", err)
			continue
		}

		go connection{conn: conn, conf: conf}.Handle()
	}
}
