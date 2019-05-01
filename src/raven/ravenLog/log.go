package ravenLog

import (
  "fmt"
  "time"
  "sort"
  "container/ring"
)

type Message struct {
  Who string
  Msg string
  When string
  time time.Time
  lvl int
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
  logChan = make( MsgChan, 10)
  errChan = make( MsgChan, 10)
  lastMsg = make( MsgMap)
  msgBuff = ring.New( 200)
  errBuff = ring.New( 200)
  go rcvMessage()
  go rcvErrors()
}

func SendMessage( lvl int, w, m string) {
  msg := new( Message)
  msg.time = time.Now()
  msg.When = msg.time.Format("2006/01/02 15:04:05")
  msg.Who = w
  msg.Msg = m
  msg.lvl = lvl


  logChan <- msg
}

func SendError( lvl int, w, m string) {
  msg := new( Message)
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
    fmt.Printf( "%s (%s) %s\n", m.When, m.Who, m.Msg)
  }
}

func GetErrors() MessageArray {
  var rtn MessageArray
  errBuff.Do( func(p interface{}) {
    if p != nil {
      rtn = append( rtn, p.(*Message))
    }
  })
  return rtn
}

func GetLastMessage() MessageArray {
  var rtn MessageArray
  var keys sort.StringSlice

  for k,_ := range lastMsg {
    keys = append(keys, k)
  }
  keys.Sort()
  for _,k := range keys {
    rtn = append( rtn, lastMsg[k])
  }
  return rtn
}

func GetLog() MessageArray {
  var rtn MessageArray
  msgBuff.Do( func(p interface{}) {
    if p != nil {
      rtn = append( rtn, p.(*Message))
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
    fmt.Printf( "%s (%s) %s\n", m.When, m.Who, m.Msg)
  }
}
