// Copyright 2021 cyops.se. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyops-se/dd-proxy/db"
	"github.com/cyops-se/dd-proxy/listeners"
	"github.com/cyops-se/dd-proxy/routes"
	"github.com/cyops-se/dd-proxy/types"
	"golang.org/x/sys/windows/svc"
)

var ctx types.Context
var GitVersion string
var GitCommit string

func main() {
	defer handlePanic()
	svcName := "dd-proxy"

	flag.StringVar(&ctx.Cmd, "cmd", "debug", "Windows service command (try 'usage' for more info)")
	flag.StringVar(&ctx.Wdir, "workdir", ".", "Sets the working directory for the process")
	flag.BoolVar(&ctx.Trace, "trace", false, "Prints traces of OCP data to the console")
	flag.BoolVar(&ctx.Version, "v", false, "Prints the commit hash and exists")
	flag.Parse()

	routes.SysInfo.GitVersion = GitVersion
	routes.SysInfo.GitCommit = GitCommit

	if ctx.Version {
		fmt.Printf("dd-proxy version %s, commit: %s\n", routes.SysInfo.GitVersion, routes.SysInfo.GitCommit)
		return
	}

	if _, err := os.Stat(ctx.Wdir); os.IsNotExist(err) {
		if err = os.MkdirAll(ctx.Wdir, os.ModePerm); err != nil {
			fmt.Printf("dd-proxy, failed to create working directory: %s, error: %s\n", ctx.Wdir, err.Error())
			return
		}
	}

	if ctx.Cmd == "install" {
		if err := installService(svcName, "dd-proxy from cyops-se"); err != nil {
			log.Fatalf("failed to %s %s: %v", ctx.Cmd, svcName, err)
		}
		return
	} else if ctx.Cmd == "remove" {
		if err := removeService(svcName); err != nil {
			log.Fatalf("failed to %s %s: %v", ctx.Cmd, svcName, err)
		}
		return
	}

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}
	if inService {
		runService(svcName, false)
		return
	}

	runService(svcName, true)
}

func runEngine() {
	defer handlePanic()

	db.ConnectDatabase(ctx)
	db.InitContent()

	listeners.RegisterType(&listeners.UDPDataListener{})
	listeners.RegisterType(&listeners.UDPMetaListener{})
	listeners.RegisterType(&listeners.UDPFileListener{})
	listeners.Init()

	go RunWeb()

	// Sleep until interrupted
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Exiting (waiting 1 sec) ...")
	time.Sleep(time.Second * 1)
}
