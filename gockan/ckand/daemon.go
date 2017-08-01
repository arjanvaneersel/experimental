package main

import (
	"log"
	"os"
	"runtime"
	"syscall"

	"github.com/vuleetu/goconfig/config"
)

// daemonise
func daemon() {
	args := []string{program, cfgfile}
	attr := syscall.ProcAttr{Dir: cwd, Env: os.Environ()}
	pid, err := syscall.ForkExec(program, args, &attr)
	if err != nil {
		log.Fatal("vfork(): ", err)
	}
	log.Printf("gockan (pid %d) daemonised", pid)
	os.Exit(0)
}

// Handle certain signals gracefully.
//
//     SIGINT and SIGTERM terminate the program
//     SIGHUP causes the program to restart
//     SIGINFO causes the program to print status information in the log
//     SIGUSR1 causes the program to dump the contents of its working set
//
func SignalHandler(cfg *config.Config) {
	sigs := make(chan os.Signal, 1)
	for {
		sig := <-sigs
		switch {
		case sig == syscall.SIGINT || sig == syscall.SIGTERM:
			log.Printf("received %s - exiting", sig)
			os.Exit(0)
		case sig == syscall.SIGHUP:
			log.Printf("received %s - restarting", sig)
			_, err := config.ReadDefault(cfgfile)
			if err == nil {
				socket.Close()
				args := []string{program, cfgfile}
				attr := syscall.ProcAttr{Dir: cwd, Env: os.Environ()}
				pid, err := syscall.ForkExec(program, args, &attr)
				if err != nil {
					log.Fatal("vfork(): ", err)
				}
				log.Printf("gockan (pid %d) restarted", pid)
				os.Exit(0)
			} else {
				log.Printf("error in config file: %s", err)
			}
		case sig == syscall.SIGUSR1:
			log.Printf("received %s - dumping", sig)
			err := Dump(cfg)
			if err != nil {
				log.Printf("error during dump: %s", err)
			} else {
				log.Printf("dump completed successfully")
			}
		case sig == syscall.SIGINFO:
			m := new(runtime.MemStats)
			runtime.ReadMemStats(m)
			log.Printf("memory: %dMB alloc %dMB total",
				m.Alloc/(1<<20),
				m.TotalAlloc/(1<<20))
		default:
			log.Print(sig)
		}
	}
}
