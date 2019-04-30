package ravenLog

import (
  "log"
  "fmt"
  "time"
  "container/ring"
)

type message struct {
  who string
  msg string
  when time.Time
  lvl int
}

type msgChan chan *message
type msgMap map[string]*message
type strArray []string

var logChan msgChan
var msgBuff *ring.Ring
var errBuff *ring.Ring
var lastMsg msgMap

func init() {
  logChan = make( msgChan, 10)
  lastMsg = make( msgMap)
  msgBuff = ring.New( 200)
  errBuff = ring.New( 200)
  go rcvMessage()
}

func SendMessage( lvl int, w, m string) {
  msg := new( message)
  msg.when = time.Now()
  msg.who = w
  msg.msg = m
  msg.lvl = lvl

  lastMsg[w] = msg

  logChan <- msg
}

func SendError( lvl int, w, m string) {
    errBuff.Value = m
    errBuff = errBuff.Next()
}

func GetLastMessage() strArray {
  var rtn strArray
  for _,k := range lastMsg {
    rtn = append( rtn, fmt.Sprintf( "%s - {%s} %s",
        k.when.Format("2006/01/02 15:04:05"),
        k.who,
        k.msg))
  }
  return rtn
}

func GetLog() strArray {
  var rtn strArray
  msgBuff.Do( func(p interface{}) {
    if p != nil {
      rtn = append( rtn, fmt.Sprintf( "%s - {%s} %s",
        p.(*message).when.Format("2006/01/02 15:04:05"),
        p.(*message).who,
        p.(*message).msg))
    }
  })
  return rtn
}

func rcvMessage() {
  for {
    m := <-logChan
    msgBuff.Value = m
    msgBuff = msgBuff.Next()
    log.Printf( "(%s) %s", m.who, m.msg)
  }
}
