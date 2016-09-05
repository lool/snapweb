/*
 * Copyright (C) 2014-2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package avahi

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/presotto/go-mdns-sd"
)

var logger *log.Logger

var _mdns *mdns.MDNS

var initOnce sync.Once

const hostnameLocalhost = "localhost"
const hostnameDefault   = "webdm"

const timeoutDuration = 3 * time.Second

var netInterfaceAddrs = net.InterfaceAddrs

func ipAddrs() (addrs []net.Addr, err error) {
	ifaces, err := netInterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs = append(addrs, iface)
	}

	return addrs, nil
}

func Init(l *log.Logger) {
	logger = l

	var err error
	hostname := getHostname()
	logger.Println("Registering hostname:", hostname)
	_mdns, err = mdns.NewMDNS(hostname, "", "", false, 1)
	if err != nil {
		logger.Println("Cannot create MDNS instance:", err)
		return
	}
	initOnce.Do(timeoutLoop)
}

func timeoutLoop() {
	timer := time.NewTimer(timeoutDuration)

	for {
		if _mdns != nil {
			_mdns.ScanInterfaces()
		}
		loop()
		timer.Reset(timeoutDuration)
		<-timer.C
	}
}

func getHostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Println("Cannot obtain hostname, using default:", err)
		return hostnameDefault
	}
	hostname = strings.Split(hostname, ".")[0]
	if hostname == hostnameLocalhost {
		hostname = hostnameDefault
	}
	return hostname
}

