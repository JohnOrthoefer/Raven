package raven

// Scheduling bits

import (
  "log"
  "time"
  "math/rand"
  ."./ravenTypes"
)

// an array that tracks everything
var status []*StatusEntry

// loops though the checks looking for hosts
func BuildSchedule() {
  for _,cn := range ListChecks() {
    log.Printf( "Scheduling %s", cn)
    log.Printf( "Hosts %v", ListCheckHosts(cn))
    for _,ch := range ListCheckHosts(cn) {
      log.Printf( "-- host %s", ch)
      t:=new( StatusEntry)
      t.Check   = GetCheckEntry(cn)
      t.Host    = GetHostEntry(ch)
      t.Last    = time.Unix(0, 0)
      t.Change  = t.Last
      t.Next    = time.Now().Add(time.Duration(rand.Intn(60)) * time.Second)
      t.Return  = new( ExitReturn)
      t.Return.Exit = 3
      t.Return.Text = ""
      t.Return.Perf = ""
      t.Return.Long = ""
      t.OldRtn = t.Return
      status = append(status,t)
    }
  }
}

// Runs the checks
func runner(id int, rec, done chan *StatusEntry) {
  for {
    job := <-rec
    log.Printf( "worker %d, got Job %s(%s)", id, job.Host.DisplayName,job.Check.DisplayName)
    job.Return = job.Check.CheckF( job.Host, job.Check.Options)
    done<-job
  }
}

// single thead to disbatch tasks to the runners
func disbatcher(send chan *StatusEntry) {
  for {
    sentJob := false
    now := time.Now()
    for _,this:=range status {
      if this.Queued {
        continue
      }
      if this.Next.Before(now) {
        sentJob = true
        this.Queued = true
        this.OldRtn = this.Return
        this.Return = nil
        log.Printf( "Disbatching %s(%s)",
          this.Host.DisplayName,this.Check.DisplayName)
        send <- this
      }
    }
    if !sentJob {
      // Max sleep time
      when := time.Now().Add(30*time.Second)
      // see if we should sleep shorter
      for _,this:=range status {
        if this.Queued {
          continue
        }
        if this.Next.Before(when) {
          when = this.Next
        }
      }
      sleepTime := time.Until( when)
      log.Printf( "Sleeping for %s", sleepTime.Round(time.Second))
      time.Sleep( sleepTime)
    }
  }
}

// this does the clean up when the runner is done, and resubmits the job 
// to the runqueue
func receiver(r chan *StatusEntry) {
  for {
    job := <-r
    if job.Return.Exit < 0 || job.Return.Exit > 3 {
      job.Return.Exit = 3
    }
    job.Last = time.Now()
    job.Next = job.Last.Add( job.Check.Interval[job.Return.Exit]).
      Add(time.Duration(rand.Intn(10)-5) * time.Second)
    if job.OldRtn.Exit != job.Return.Exit {
     job.Change = job.Last
    }
    log.Printf( "Rescheduling %s(%s) in %s Exit: %d",
      job.Host.DisplayName,job.Check.DisplayName,
      time.Until(job.Next).Round(time.Second), job.Return.Exit)
    job.Queued = false
  }
}

// Starts up the scheduler, and workers
func StartSchedule(work int) {
  var disbatchQ = make( chan *StatusEntry, work)
  var returnQ = make( chan *StatusEntry, work)

  for i:=0; i < work; i++ {
    log.Printf( "Starting runner %d", i)
    go runner(i, disbatchQ, returnQ)
  }
  go disbatcher(disbatchQ)
  go receiver(returnQ)
}

// prints the schedule not preaty but it's debugging
func DumpSchedule() {
  for i:=range status {
    log.Printf( "%s[%s] - Last:%s(Exit:%d) Next:%s ",
      status[i].Host.DisplayName, status[i].Check.CheckN,
      status[i].Last.Truncate(0).Local(), status[i].Return.Exit,
      status[i].Next.Truncate(0).Local())
  }
}
