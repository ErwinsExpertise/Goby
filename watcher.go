package main

import (
	"context"
	"time"

	"github.com/ErwinsExpertise/Goby/snapshot"
)

// Scheduler performs snapshots at given interval
func Scheduler(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(SnapshotFreq) * time.Minute)
	create := true
	for {
		if create == false {
			continue
		}
		create = false
		go func() {
			for {
				select {
				case <-ticker.C:
					snapshot.SnapshotDroplet(DOToken, DropletID)
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()
	}

}

// Cleaner removes backups that are older than given frequency
func Cleaner(ctx context.Context) {

}
