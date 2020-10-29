package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	// DOToken is the DigitalOcean API Token
	DOToken string
	// DropletID is the DigitalOcean Droplet ID
	DropletID string
	// SnapshotFreq is the frequency that snapshots should occur
	SnapshotFreq int64 // Value in minutes
	// KeepTime is the amount of time that snapshots should be kept
	KeepTime int64 // Value in days
)

func init() {
	DOToken = os.Getenv("DO_API_TOKEN")
	DropletID = os.Getenv("DO_DROPLET_ID")

	flag.Int64Var(&SnapshotFreq, "freq", 1440, "How often to perform snapshots in minutes") // Default 24 hours
	flag.Int64Var(&KeepTime, "keep", 7, "Amount of days to keep snapshots")                 // Default 7 days
	flag.Parse()
}

func main() {
	// Exit if either environment variable has not been set
	if DOToken == "" || DropletID == "" {
		fmt.Printf("Please ensure all environment variables are set.\nDO_API_TOKEN\nDO_DROPLET_ID\n")
		os.Exit(4)
	}
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	// Create channel of size 1
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			sig := <-c
			log.Printf("Recieved %s signal. Shutting down server...\n", sig)
			cancel()
			os.Exit(0)
		case <-ctx.Done():
			// consume
		}

	}()

	log.Println("Goby has started for droplet " + DropletID)
	go Scheduler(ctx)
	go Cleaner(ctx)

	<-ctx.Done()
}
