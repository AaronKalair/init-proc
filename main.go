package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func sigHandler(pid *int, signalChannel chan os.Signal) {
	log.Println("Sig handler registered")
	var sigToSend syscall.Signal = syscall.SIGHUP
	for {
		sig := <-signalChannel
		switch sig {
		// Sent went the controlling terminal is closed, typically used
		// by daemonised processes to reload config
		case syscall.SIGHUP:
			sigToSend = syscall.SIGHUP
		// Like pressing CTRL+C
		case syscall.SIGINT:
			sigToSend = syscall.SIGINT
		// Core Dump
		case syscall.SIGQUIT:
			sigToSend = syscall.SIGQUIT
		// Invalid instruction
		case syscall.SIGILL:
			sigToSend = syscall.SIGILL
		// A debugger wants to know something has happened
		case syscall.SIGTRAP:
			sigToSend = syscall.SIGTRAP
		// Usually when a process itself calls abort()
		case syscall.SIGABRT:
			sigToSend = syscall.SIGABRT
		// Unimplemented instruction
		//case syscall.SIGEMT:
		//	sigToSend = syscall.SIGEMT
		// Die immediately, can not be ignored
		case syscall.SIGKILL:
			sigToSend = syscall.SIGKILL
		// Memory access error
		case syscall.SIGBUS:
			sigToSend = syscall.SIGBUS
		// Invalid memory reference, probably access memory you dont own
		case syscall.SIGSEGV:
			sigToSend = syscall.SIGSEGV
		// Bad argument to a syscall or violated a seccomp rule
		case syscall.SIGSYS:
			sigToSend = syscall.SIGSYS
		// Attempted to write to pipe with a process on the other end
		case syscall.SIGPIPE:
			sigToSend = syscall.SIGPIPE
		// timer (real / clock time) set earlier has elapsed
		case syscall.SIGALRM:
			sigToSend = syscall.SIGALRM
		// Request termination similar to SIGINT
		case syscall.SIGTERM:
			sigToSend = syscall.SIGTERM
		// Socket has urgent / out of bound data to read
		case syscall.SIGURG:
			sigToSend = syscall.SIGURG
		// Stop a process for resuming later CTRL+Z
		case syscall.SIGSTOP:
			sigToSend = syscall.SIGSTOP
		//  Similar to SIGSTOP but cant be ignored
		case syscall.SIGTSTP:
			sigToSend = syscall.SIGTSTP
		// Resume after stopping
		case syscall.SIGCONT:
			sigToSend = syscall.SIGCONT
		// Child process has terminated, interupted or resumed
		case syscall.SIGCHLD:
			var status syscall.WaitStatus
			var rusage syscall.Rusage
			log.Println("Waiting on Children")
			for {
				retValue, err := syscall.Wait4(-1, &status, syscall.WNOHANG, &rusage)
				if err != nil {
					panic(err)
				}
				if retValue <= 0 {
					break
				}
			}
			log.Println("Done waiting on Children")
			sigToSend = syscall.SIGCHLD
		// Alias for SIGCHLD
		//case syscall.SIGCLD:
		//	sigToSend = syscall.SIGCLD
		// Attempted to read from TTY whilst in background
		case syscall.SIGTTIN:
			sigToSend = syscall.SIGTTIN
		// Attempted to write to TTY whilst in background
		//case syscall.SIGTTOU:
		//	sigToSend = syscall.SIGTOU
		// File descriptor is ready for IO
		case syscall.SIGIO:
			sigToSend = syscall.SIGIO
		// Arithmetic error such as divide by 0
		case syscall.SIGFPE:
			sigToSend = syscall.SIGFPE
		// Usually when a process itself calls abort() alias for ABRT
		//case syscall.SIGIOT:
		//	sigToSend = syscall.SIGIOT
		// Alias for SIGIO
		//case syscall.SIGPOLL:
		//	sigToSend = syscall.SIGPOLL
		// Timer that measures CPU time used by current process and CPU
		// time used on behalf of process, used for profiling
		case syscall.SIGPROF:
			sigToSend = syscall.SIGPROF
		// System experiences power failure
		case syscall.SIGPWR:
			sigToSend = syscall.SIGPWR
		// Coprocessor experienced stack fault
		case syscall.SIGSTKFLT:
			sigToSend = syscall.SIGSTKFLT
		// System call with unused sys call number is made, same as
		// SIGSYS normally
		//case syscall.SIGUNUSED:
		//	sigToSend = syscall.SIGUNUSED
		// User defined
		case syscall.SIGUSR1:
			sigToSend = syscall.SIGUSR1
		// User defined
		case syscall.SIGUSR2:
			sigToSend = syscall.SIGUSR2
		// CPU time used by timer elapsed
		case syscall.SIGVTALRM:
			sigToSend = syscall.SIGVTALRM
		// Terminal window size changed
		case syscall.SIGWINCH:
			sigToSend = syscall.SIGWINCH
		// Used up CPU duration previously set
		case syscall.SIGXCPU:
			sigToSend = syscall.SIGXCPU
		// File has grown too large
		case syscall.SIGXFSZ:
			sigToSend = syscall.SIGXFSZ
		}
		log.Println("About to send", sigToSend, "to PID", *pid)
		syscall.Kill(*pid, sigToSend)
	}
}

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel)
	pid := -1

	go sigHandler(&pid, signalChannel)

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	pid = cmd.Process.Pid

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
