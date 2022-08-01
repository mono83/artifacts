package cmd

import (
	"github.com/mono83/artifacts/data"
	"github.com/mono83/artifacts/db"
	"github.com/mono83/artifacts/influx"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
	"math/rand"
	"sync"
	"time"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start", "daemon"},
	Short:   "Starts daemon that will perform regular artifact data read",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Establishing connection to MySQL database
		cnf, conn, artifacts, err := configure()
		if err != nil {
			return err
		}

		xray.BOOT.Info("Scheduling :count artifacts for work", args.Count(len(artifacts)))
		wg := sync.WaitGroup{}
		wg.Add(len(artifacts))
		for _, a := range artifacts {
			initialWait := time.Duration(rand.Intn(len(artifacts)*1000)) * time.Millisecond
			go func(wait time.Duration, a data.Artifact) {
				delay := a.IntervalOrDefault(cnf.ScheduleOrDefault())
				xray.BOOT.Info("Using initial wait :elapsed for :name to avoid startup load", args.Elapsed(wait), args.Name(a.Metric))
				time.Sleep(wait)
				for {
					ray := xray.ROOT.Fork().With(args.Name(a.Metric))
					ray.Debug("Performing work on :name")
					results, err := db.Read(conn, a)
					if err != nil {
						ray.Error("Error obtaining data for :name - :err", args.Error{Err: err})
					} else {
						var toSend []data.ResultsTable
						toSend = append(toSend, *results)
						toSend = append(toSend, results.Recombine()...)
						err = influx.Send(cnf.InfluxDBAddr, toSend)
						if err != nil {
							ray.Error("Error sending data to InfluxDB - :err", args.Error{Err: err})
						}
					}

					ray.Debug("Sleeping :elapsed for :name", args.Elapsed(delay))
					time.Sleep(delay)
				}
				wg.Done()
			}(initialWait, a)
		}
		wg.Wait()
		return nil
	},
}
