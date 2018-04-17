package main

import (
	"net"
	"github.com/aspcartman/mysocks/socks"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	lst, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	socks.Serve(lst, &socks.Config{
		Auth: authcheck(),
		HandleError: func(stage string, err error) {
			logrus.WithError(err).Error(stage)
		},
		Dial: net.Dial,
	})
}

func authcheck() func(user, password string) bool {
	envvar := os.Getenv("PROXY_AUTH")
	if len(envvar) == 0 {
		return nil
	}

	// Expects user:pass;user;user:pass in PROXY_AUTH
	usrs := map[string]string{}
	for _, authstr := range strings.Split(envvar, ";") {
		r := strings.Split(authstr, ":")
		if len(r) == 1 {
			usrs[r[0]] = ""
		} else {
			usrs[r[0]] = r[1]
		}
	}

	return func(user, password string) bool {
		if p, ok := usrs[user]; ok && len(p) == 0 {
			return true // no password specified
		} else if ok {
			return p == password
		}
		return false
	}
}
