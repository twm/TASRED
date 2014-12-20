// Copyright 2013, 2014 Tom Most
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/stvp/go-toml-config"
	"log"
	"math/rand"
	"os"
	"os/signal"
)

// IRC configuration parameters
var (
	ircServer  = config.String("irc.server", "irc.freenode.net:7000")
	ircSSL     = config.Bool("irc.ssl", true)
	ircNick    = config.String("irc.nick", "carrotfacts")
	ircName    = config.String("irc.name", "")
	ircUser    = config.String("irc.user", "carrotfacts")
	ircPass    = config.String("irc.password", "")
	ircChannel = config.String("irc.channel", "#carrotfacts-test")
)

// seed feeds the Go PRNG a cryptographically random number so we don't always
// choose facts in the same order.
func seed() (err error) {
	var seed int64
	err = binary.Read(cryptorand.Reader, binary.LittleEndian, &seed)
	if err == nil {
		rand.Seed(seed)
	}
	return
}

func main() {
	configFile := flag.String("config", "config.toml", "Config file")
	flag.Parse()

	config.Parse(*configFile)
	if err := seed(); err != nil {
		log.Fatalf("Unable to seed the PRNG: %s\n", err)
	}

	c := irc.NewConfig(*ircNick)
	c.SSL = *ircSSL
	c.Server = *ircServer
	c.Pass = *ircPass
	c.Me.Ident = *ircNick
	c.Me.Name = *ircName
	ic := irc.Client(c)

	ic.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			conn.Join(*ircChannel)
		})

	ic.HandleFunc("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		msg := line.Args[1]
		fmt.Printf("%s %s %s\n", channel, line.Nick, msg)
	})

	quit := make(chan bool)
	ic.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	if err := ic.ConnectTo(*ircServer); err != nil {
		log.Fatalf("Unable to connect: %s\n", err)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

main:
	for {
		select {
		case <-interrupt:
			ic.Quit("Carrot be with you!")
		case <-quit:
			log.Printf("Disconnected")
			break main
		}
	}
}
