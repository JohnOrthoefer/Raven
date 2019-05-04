package ravenChecks

// utllites 
// Run External
// Return ops

import (
  "bytes"
  "log"
  "os/exec"
  "syscall"
)


// move to it's own package?
// I recycled this code into the config generator.. so it makes
// sense but I'm also thinking about making the config generator 
// a detector check.   so maybe it will all come out in the wash
func runExternal( prog string, args ...string) (int, string) {

  var out bytes.Buffer

  cmd := exec.Command(prog, args...)

  cmd.Stdout = &out
  cmd.Stderr = &out
  if err := cmd.Start(); err != nil {
    log.Fatalf("cmd.Start: %v")
  }

  rtnExit:=0
  if err := cmd.Wait(); err != nil {
    if exiterr, ok := err.(*exec.ExitError); ok {
    // The program has exited with an exit code != 0
    // This works on both Unix and Windows. Although package
    // syscall is generally platform dependent, WaitStatus is
    // defined for both Unix and Windows and in both cases has
    // an ExitStatus() method with the same signature.
      if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
        rtnExit = status.ExitStatus()
        log.Printf("Exit Status: %d", rtnExit)
      }
    } else {
      log.Fatalf("cmd.Wait: %v", err)
    }
  }
  return rtnExit, out.String()
}

