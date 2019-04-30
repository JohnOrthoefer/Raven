package raven

// Scheduling bits

import (
  "fmt"
  "time"
  "math/rand"
  "./ravenTypes"
  "./ravenLog"
)

// an array that tracks everything
var status []*ravenTypes.StatusEntry

// loops though the checks looking for hosts
func BuildSchedule() {
  for _,cn := range ListChecks() {
    ravenLog.SendError( 10, "BuildSchedule", fmt.Sprintf( "Scheduling %s", cn))
    ravenLog.SendError( 10, "BuildSchedule", fmt.Sprintf( "Hosts %v", ListCheckHosts(cn)))
    for _,ch := range ListCheckHosts(cn) {
      ravenLog.SendError( 10, "BuildSchedule", fmt.Sprintf( "-- host %s", ch))
      t:=new( ravenTypes.StatusEntry)
      t.Check   = GetCheckEntry(cn)
      t.Host    = GetHostEntry(ch)
      t.Last    = time.Unix(0, 0)
      t.Change  = t.Last
      t.Next    = time.Now().Add(time.Duration(rand.Intn(60)) * time.Second)
      t.Return  = new( ravenTypes.ExitReturn)
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
func runner(id int, rec, done chan *ravenTypes.StatusEntry) {
  for {
    job := <-rec
    ravenLog.SendMessage( 10, fmt.Sprintf( "runner-%d", id),
      fmt.Sprintf( "got Job %s(%s)", job.Host.DisplayName,job.Check.DisplayName))
    job.Return = job.Check.CheckF( job.Host, job.Check.Options)
    done<-job
  }
}

// single thead to disbatch tasks to the runners
func disbatcher(send chan *ravenTypes.StatusEntry) {
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
        ravenLog.SendMessage( 10, "disbatch", fmt.Sprintf( "Disbatching %s(%s)",
          this.Host.DisplayName,this.Check.DisplayName))
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
      ravenLog.SendMessage( 10, "disbatch", fmt.Sprintf( "Sleeping for %s", sleepTime.Round(time.Second)))
      time.Sleep( sleepTime)
    }
  }
}

// this does the clean up when the runner is done, and resubmits the job 
// to the runqueue
func receiver(r chan *ravenTypes.StatusEntry) {
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
    ravenLog.SendMessage( 10, "receiver", fmt.Sprintf( "Rescheduling %s(%s) in %s Exit: %d",
      job.Host.DisplayName,job.Check.DisplayName,
      time.Until(job.Next).Round(time.Second), job.Return.Exit))
    job.Queued = false
  }
}

// Starts up the scheduler, and workers
func StartSchedule(work int) {
  var disbatchQ = make( chan *ravenTypes.StatusEntry, work)
  var returnQ = make( chan *ravenTypes.StatusEntry, work)

  for i:=0; i < work; i++ {
    ravenLog.SendError( 10, "StartSchedule", fmt.Sprintf( "Starting runner %d", i))
    go runner(i, disbatchQ, returnQ)
  }
  go disbatcher(disbatchQ)
  go receiver(returnQ)
}

// prints the schedule not preaty but it's debugging
func DumpSchedule() {
  for i:=range status {
    ravenLog.SendMessage( 10, "DumpSchedule", fmt.Sprintf( "%s[%s] - Last:%s(Exit:%d) Next:%s ",
      status[i].Host.DisplayName, status[i].Check.CheckN,
      status[i].Last.Truncate(0).Local(), status[i].Return.Exit,
      status[i].Next.Truncate(0).Local()))
  }
}
