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
	"bufio"
	cryptorand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/stvp/go-toml-config"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
)

// IRC configuration parameters
var (
	ircServer    = config.String("irc.server", "irc.twitch.tv:6667")
	ircSSL       = config.Bool("irc.ssl", false)
	ircNick      = config.String("irc.nick", "tasred")
	ircName      = config.String("irc.name", "tasred")
	ircUser      = config.String("irc.user", "tasred")
	ircPass      = config.String("irc.password", "")
	ircChannel   = config.String("irc.channel", "#agdq")
	badWordsFile = config.String("data.badwords", "bad-words.txt")
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

type Filter struct {
	badWords []string
}

// Okay determines if a phrase contains any words from a blacklist
func (f *Filter) Okay(phrase string) bool {
	for _, word := range f.badWords {
		if strings.Contains(phrase, word) {
			return false
		}
	}
	return true
}

// readWordList reads words, one per line, from a file, or aborts if it cannot
// be read.
func readWordList(filename string) (*Filter, error) {
	var (
		f   Filter
		r   io.ReadCloser
		err error
	)
	if r, err = os.Open(filename); err != nil {
		return nil, err
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := scanner.Text()
		f.badWords = append(f.badWords, word)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &f, nil
}

func main() {
	configFile := flag.String("config", "config.toml", "Config file")
	flag.Parse()

	config.Parse(*configFile)
	if err := seed(); err != nil {
		log.Fatalf("Unable to seed the PRNG: %s\n", err)
	}

	var filter *Filter
	var err error
	if filter, err = readWordList(*badWordsFile); err != nil {
		log.Fatalf("Unable to read bad words: %s\n", err)
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
		//channel := line.Args[0]
		msg := line.Args[1]
		prettyLine := fmt.Sprintf("%s: %s", line.Nick, msg)
		if filter.Okay(prettyLine) {
			fmt.Printf("NICE    %s\n", prettyLine)
		} else {
			fmt.Printf("NAUGHTY %s\n", prettyLine)
		}
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
			ic.Quit("TASBOT be with you!")
		case <-quit:
			log.Printf("Disconnected")
			break main
		}
	}
}
