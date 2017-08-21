// +build !android,!ios

// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package daylight

import (
	"net"
	"net/http"
	"time"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"
	"github.com/EGaaS/go-egaas-mvp/packages/tcpserver"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"

	"github.com/sirupsen/logrus"
)

func httpListener(ListenHTTPHost string, BrowserHTTPHost string, route http.Handler) {
	host := ListenHTTPHost
	listener, err := net.Listen("tcp4", host)
	if err != nil {
		log.WithFields(logrus.Fields{"host": host, "error": err}).Fatal("Error start listening on host, when trying start serve http")
	}
	go func() {
		srv := &http.Server{Handler: route}
		if err := srv.Serve(listener); err != nil {
			log.WithFields(logrus.Fields{"host": host, "error": err}).Fatal("Error serving http on host")
		}
	}()
}

// For ipv6 on the server
func httpListenerV6(route http.Handler) {
	port := *utils.ListenHTTPPort
	listener, err := net.Listen("tcp6", ":"+port)
	if err != nil {
		log.WithFields(logrus.Fields{"port": port, "error": err}).Fatal("Error listeining ipv6 on port, when trying start serve http")
	}

	go func() {
		srv := &http.Server{Handler: route}
		if err := srv.Serve(listener); err != nil {
			log.WithFields(logrus.Fields{"port": port, "error": err}).Fatal("Error serving http on host")
		}
	}()
}

func tcpListener() {
	go func() {
		// switch on the listing by TCP-server and the processing of incoming requests
		listener, err := net.Listen("tcp4", *utils.TCPHost+":"+consts.TCP_PORT)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, "host": *utils.TCPHost, "port": consts.TCP_PORT}).Fatal("Error starting listen, when starting tcp server")
		}
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.WithFields(logrus.Fields{"error": err, "host": *utils.TCPHost, "port": consts.TCP_PORT}).Fatal("error accepting, when starting tcp server")
					time.Sleep(time.Second)
				} else {
					go func(conn net.Conn) {
						tcpserver.HandleTCPRequest(conn)
						conn.Close()
					}(conn)
				}
			}
		}()
	}()
}
