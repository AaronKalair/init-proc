package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func sigHandler(pid *int, signalChannel chan os.Signal) {
	fmt.Println("Sig handler registered")
	var sigToSend syscall.Signal = syscall.SIGHUP
	for {
		sig := <-signalChannel
		switch sig {
		// #1 - Sent went the controlling terminal is closed, typically used
		// by daemonised processes to reload config
		case syscall.SIGHUP:
			sigToSend = syscall.SIGHUP
		// #2 - Like pressing CTRL+C
		case syscall.SIGINT:
			sigToSend = syscall.SIGINT
		// #3 - Core Dump
		case syscall.SIGQUIT:
			sigToSend = syscall.SIGQUIT
		// #4 - Invalid instruction
		case syscall.SIGILL:
			sigToSend = syscall.SIGILL
		// #5 - A debugger wants to know something has happened
		case syscall.SIGTRAP:
			sigToSend = syscall.SIGTRAP
		// #6 - Usually when a process itself calls abort()
		case syscall.SIGABRT:
			sigToSend = syscall.SIGABRT
		// #7 - Unimplemented instruction
		//case syscall.SIGEMT:
		//	sigToSend = syscall.SIGEMT
		// #9 - Die immediately, can not be ignored
		case syscall.SIGKILL:
			sigToSend = syscall.SIGKILL
		// #10 - memory access error
		case syscall.SIGBUS:
			sigToSend = syscall.SIGBUS
		// #11 - Invalid memory reference, probably access memory you dont own
		case syscall.SIGSEGV:
			sigToSend = syscall.SIGSEGV
		// #12 - Bad argument to a syscall or violated a seccomp rule
		case syscall.SIGSYS:
			sigToSend = syscall.SIGSYS
		// #13 - Attempted to write to pipe with a process on the other end
		case syscall.SIGPIPE:
			sigToSend = syscall.SIGPIPE
		// #14 - timer (real / clock time) set earlier has elapsed
		case syscall.SIGALRM:
			sigToSend = syscall.SIGALRM
		// #15 - Request termination similar to SIGINT
		case syscall.SIGTERM:
			sigToSend = syscall.SIGTERM
		// #16 - Socket has urgent / out of bound data to read
		case syscall.SIGURG:
			sigToSend = syscall.SIGURG
		// #17 - Stop a process for resuming later CTRL+Z
		case syscall.SIGSTOP:
			sigToSend = syscall.SIGSTOP
		// #18 - Similar to SIGSTOP but cant be ignored
		case syscall.SIGTSTP:
			sigToSend = syscall.SIGTSTP
		// #19 - Resume after stopping
		case syscall.SIGCONT:
			sigToSend = syscall.SIGCONT
		// #20 - Child process has terminated, interupted or resumed
		case syscall.SIGCHLD:
			var status syscall.WaitStatus
			var rusage syscall.Rusage
			fmt.Println("Waiting on Child")
			for {
				retValue, err := syscall.Wait4(-1, &status, syscall.WNOHANG, &rusage)
				fmt.Println("RetValue", retValue)
				if err != nil {
					panic(err)
				}
				if retValue <= 0 {
					break
				}
			}
			fmt.Println("Done waiting on child")
			sigToSend = syscall.SIGCHLD
		// #21 - Alias for SIGCHLD
		//case syscall.SIGCLD:
		//	sigToSend = syscall.SIGCLD
		// #22 - Attempted to read from TTY whilst in background
		case syscall.SIGTTIN:
			sigToSend = syscall.SIGTTIN
		// #23 - Attempted to write to TTY whilst in background
		//case syscall.SIGTTOU:
		//	sigToSend = syscall.SIGTOU
		// #24 - File descriptor is ready for IO
		case syscall.SIGIO:
			sigToSend = syscall.SIGIO
		// #25 - Arithmetic error such as divide by 0
		case syscall.SIGFPE:
			sigToSend = syscall.SIGFPE
		// #26 - Usually when a process itself calls abort() alias for ABRT
		//case syscall.SIGIOT:
		//	sigToSend = syscall.SIGIOT
		// #27 - Alias for SIGIO
		//case syscall.SIGPOLL:
		//	sigToSend = syscall.SIGPOLL
		// #28 - Timer that measures CPU time used by current process and CPU
		// time used on behalf of process, used for profiling
		case syscall.SIGPROF:
			sigToSend = syscall.SIGPROF
		// #29 - System experiences power failure
		case syscall.SIGPWR:
			sigToSend = syscall.SIGPWR
		// #30 - Coprocessor experienced stack fault
		case syscall.SIGSTKFLT:
			sigToSend = syscall.SIGSTKFLT
		// #31 - System call with unused sys call number is made, same as
		// SIGSYS normally
		//case syscall.SIGUNUSED:
		//	sigToSend = syscall.SIGUNUSED
		// #32 - User defined
		case syscall.SIGUSR1:
			sigToSend = syscall.SIGUSR1
		// #33 - User defined
		case syscall.SIGUSR2:
			sigToSend = syscall.SIGUSR2
		// #34 - CPU time used by timer elapsed
		case syscall.SIGVTALRM:
			sigToSend = syscall.SIGVTALRM
		// #35 - Terminal window size changed
		case syscall.SIGWINCH:
			sigToSend = syscall.SIGWINCH
		// #36 - Used up CPU duration previously set
		case syscall.SIGXCPU:
			sigToSend = syscall.SIGXCPU
		// #37 - File has grown too large
		case syscall.SIGXFSZ:
			sigToSend = syscall.SIGXFSZ
		}
		fmt.Println("About to send", sigToSend, "to PID", *pid)
		syscall.Kill(*pid, sigToSend)
	}
}

func main() {
	fmt.Println("Starting up")
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel)
	pid := -1

	go sigHandler(&pid, signalChannel)

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	err := cmd.Start()
	pid = cmd.Process.Pid

	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
