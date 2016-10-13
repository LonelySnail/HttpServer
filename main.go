package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"./config"
	"./server"
)

/*
func init() {
	err := os.Mkdir("../upload", 0777)
	fmt.Println(err)
}
*/

func writePIDFile(filename string) error {
	data := fmt.Sprintf("%d\n", os.Getpid())
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func processKill(enableProfile bool, logger *log.Logger) {
	killCh := make(chan os.Signal, 1)
	signal.Notify(killCh, os.Interrupt, os.Kill)
	sig := <-killCh

	if enableProfile {
		pprof.StopCPUProfile()
	}
	logger.Fatal("Stop Server ", sig)
}

func main() {
	coresPtr := flag.Int("C", 0, "Number Cores can be used, if ignored use all Cores")
	pidPtr := flag.String("p", "", "PID file path, if ignored it will not be created")
	configPtr := flag.String("c", "./config/config.json", "Config file path, if ignored will be load from  /etc/HttpServer/config.json")
	logPtr := flag.String("l", "./log/Develop.log", "Log file path, if ignored log will output to stdout")
	blockPtr := flag.Bool("P", false, "block server")
	flag.Parse()

	var logger *log.Logger

	if *logPtr != "" {
		logFile, err := os.OpenFile(*logPtr, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if *coresPtr != 0 {
		logger.Printf("Use %d Cores\n", *coresPtr)
		runtime.GOMAXPROCS(*coresPtr)
	}

	if *pidPtr != "" {
		logger.Printf("Write PID to: %s\n", *pidPtr)
		err := writePIDFile(*pidPtr)
		if err != nil {
			log.Println(err)
		}
	}

	err := config.LoadConfigFile(*configPtr, logger)
	if err != nil {
		log.Fatal(err)
		return
	}

	go server.StartServer(logger)

	processKill(*blockPtr, logger)
}
