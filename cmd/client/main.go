package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meddion/pkg/core"
)

const (
	_dbFile                = "_test_db_file_"
	_testAddr              = ""
	_testPort              = "2022"
	_isAliveInterval       = time.Minute * 2
	_peerDiscoveryInterval = time.Minute * 5
)

func main() {
	log := log.Default()

	db, err := core.NewBlockRepo(_dbFile)
	if err != nil {
		log.Fatalf("on creating a block repo %s", err)
	}

	blkchain, err := core.NewBlockchain(db, log)
	if err != nil {
		log.Fatalf("on creating the Blockchain instance: %s", err)
	}

	peerPool := core.NewPeerPool(log, _peerDiscoveryInterval, _isAliveInterval)
	rcv := core.NewReceiverRPC(blkchain, peerPool, log)

	serv, err := core.NewServer(rcv)
	if err != nil {
		log.Fatalf("on creating the Server: %s", err)
	}

	servDone := make(chan struct{}, 1)
	go func() {
		defer func() {
			close(servDone)
		}()

		log.Printf("Starting listening for incoming connections on %s:%s", _testAddr, _testPort)

		if err := serv.Start(_testAddr, _testPort); err != http.ErrServerClosed {
			log.Printf("on starting the Server: %s", err)
		}
		log.Print("The Server has been closed.")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Print("Got an OS signal. Preparing to close...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := serv.Close(ctx); err != nil {
		log.Printf("on calling close on the Server: %s", err)
	}

	<-servDone
}
