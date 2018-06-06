package main

import (
	"log"
	"net/http"
	"path/filepath"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	"flag"
	"strconv"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		port= flag.Int("p", 80, "host port default is 80")
	)
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if dir == "/" {
		dir += "www"
	} else {
		dir += "/www"
	}
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(dir+"/chat"))))
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Coolpy7 Web Tool Host On Port",strconv.Itoa(*port))

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
