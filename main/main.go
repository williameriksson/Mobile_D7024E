package main

import (
	//"log"
	"time"
    "Mobile_D7024E/api"
)

const running_time time.Duration = 60 * time.Second

// build: go build -o main.exe main.go

func main() {

	//log.Printf("main: starting HTTP server")

    api.StartServer()
/*

    srv := api.StartNewHttpServer()
    log.Printf("main: serving for %v", running_time)

    time.Sleep(running_time)

    log.Printf("main: stopping HTTP server")

    // now close the server gracefully ("shutdown")
    // timeout could be given instead of nil as a https://golang.org/pkg/context/
    if err := srv.Shutdown(nil); err != nil {
        panic(err) // failure/timeout shutting down the server gracefully
    }

    log.Printf("main: done. exiting")*/
}