package raven

// Scheduling bits

import (
  "bytes"
  "log"
  "time"
  "os/exec"
  "syscall"
  "math/rand"
)

// the basic entry for scheduling a check against a host
type StatusEntry struct {
  Check     *CheckEntry
  Host      *HostEntry
  ExitCode  int
  Queued    bool
  Next      time.Time
  Last      time.Time
}

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
      t.Check = GetCheckEntry(cn)
      t.Host  = GetHostEntry(ch)
      t.ExitCode  = 3
      t.Last  = time.Unix(0, 0)
      t.Next  = time.Now().Add(time.Duration(rand.Intn(60)) * time.Second)
      status = append(status,t)
    }
  }
}

// Runs the checks
func runner(id int, rec, done chan *StatusEntry) {
  for {
    job := <-rec
    log.Printf( "worker %d, got Job %s(%s)", id, job.Host.DisplayName,job.Check.DisplayName)
		cmd := exec.Command("/usr/bin/ping", "-c", "5", job.Host.Hostname)
		var out bytes.Buffer
		cmd.Stdout = &out
    if err := cmd.Start(); err != nil {
      log.Fatalf("cmd.Start: %v")
    }
    job.ExitCode = 0
    if err := cmd.Wait(); err != nil {
      if exiterr, ok := err.(*exec.ExitError); ok {
      // The program has exited with an exit code != 0
      // This works on both Unix and Windows. Although package
      // syscall is generally platform dependent, WaitStatus is
      // defined for both Unix and Windows and in both cases has
      // an ExitStatus() method with the same signature.
        if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
         job.ExitCode = status.ExitStatus()
         log.Printf("Exit Status: %d", status.ExitStatus())
        }
      } else {
        log.Fatalf("cmd.Wait: %v", err)
      }
    }
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
        log.Printf( "Disbatching %s(%s)",
          this.Host.DisplayName,this.Check.CheckN)
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
    if job.ExitCode < 0 || job.ExitCode > 3 {
      job.ExitCode = 3
    }
    job.Last = time.Now()
    job.Next = job.Last.Add( job.Check.Interval[job.ExitCode]).
      Add(time.Duration(rand.Intn(10)-5) * time.Second)
    log.Printf( "Rescheduling %s(%s) in %s Exit: %d",
      job.Host.DisplayName,job.Check.CheckN,
      time.Until(job.Next).Round(time.Second), job.ExitCode)
    job.Queued = false
  }
}

// Starts up the scheduler, and workers
func StartSchedule(work int) {
  var disbatchQ = make( chan *StatusEntry)
  var returnQ = make( chan *StatusEntry)

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
      status[i].Last.Truncate(0).Local(), status[i].ExitCode,
      status[i].Next.Truncate(0).Local())
  }
}
