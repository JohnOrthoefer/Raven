/*
   Raven Network Discovery and Monitoring
   Copyright (C) 2019 John{at}Orthoefer{dot}org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package ravenLog

import (
	"container/ring"
	"fmt"
	"sort"
	"time"
)

type Message struct {
	Who  string
	Msg  string
	When string
	time time.Time
	lvl  int
}

type MsgChan chan *Message
type MsgMap map[string]*Message
type strArray []string
type MessageArray []*Message

var logChan MsgChan
var errChan MsgChan
var msgBuff *ring.Ring
var errBuff *ring.Ring
var lastMsg MsgMap

func init() {
	logChan = make(MsgChan, 10)
	errChan = make(MsgChan, 10)
	lastMsg = make(MsgMap)
	msgBuff = ring.New(200)
	errBuff = ring.New(200)
	go rcvMessage()
	go rcvErrors()
}

func SendMessage(lvl int, w, m string) {
	msg := new(Message)
	msg.time = time.Now()
	msg.When = msg.time.Format("2006/01/02 15:04:05")
	msg.Who = w
	msg.Msg = m
	msg.lvl = lvl

	logChan <- msg
}

func SendError(lvl int, w, m string) {
	msg := new(Message)
	msg.time = time.Now()
	msg.When = msg.time.Format("2006/01/02 15:04:05")
	msg.Who = w
	msg.Msg = m
	msg.lvl = lvl
	errChan <- msg
}

func rcvErrors() {
	for {
		m := <-errChan
		errBuff.Value = m
		errBuff = errBuff.Next()
		fmt.Printf("%s (%s) %s\n", m.When, m.Who, m.Msg)
	}
}

func GetErrors() MessageArray {
	var rtn MessageArray
	errBuff.Do(func(p interface{}) {
		if p != nil {
			rtn = append(rtn, p.(*Message))
		}
	})
	return rtn
}

func GetLastMessage() MessageArray {
	var rtn MessageArray
	var keys sort.StringSlice

	for k, _ := range lastMsg {
		keys = append(keys, k)
	}
	keys.Sort()
	for _, k := range keys {
		rtn = append(rtn, lastMsg[k])
	}
	return rtn
}

func GetLog() MessageArray {
	var rtn MessageArray
	msgBuff.Do(func(p interface{}) {
		if p != nil {
			rtn = append(rtn, p.(*Message))
		}
	})
	return rtn
}

func rcvMessage() {
	for {
		m := <-logChan
		msgBuff.Value = m
		msgBuff = msgBuff.Next()
		lastMsg[m.Who] = m
		fmt.Printf("%s (%s) %s\n", m.When, m.Who, m.Msg)
	}
}
