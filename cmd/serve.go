package cmd

import (
	"errors"
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

var (
	flagInfluxDB            string
	flagDefaultScheduleTime time.Duration
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start", "daemon"},
	Short:   "Starts daemon that will perform regular artifact data read",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if len(flagInfluxDB) == 0 {
			return errors.New("empty InfluxDB listener address")
		}
		xray.BOOT.Info("Using InfluxDB :addr", args.Addr(flagInfluxDB))

		// Establishing connection to MySQL database
		conn, artifacts, err := mysql()
		if err != nil {
			return err
		}

		xray.BOOT.Info("Scheduling :count artifacts for work", args.Count(len(artifacts)))
		wg := sync.WaitGroup{}
		wg.Add(len(artifacts))
		for _, a := range artifacts {
			initialWait := time.Duration(rand.Intn(len(artifacts)*1000)) * time.Millisecond
			go func(wait time.Duration, a data.Artifact) {
				delay := a.IntervalOrDefault(flagDefaultScheduleTime)
				xray.BOOT.Info("Using initial wait :elapsed for :name to avoid startup load", args.Elapsed(wait), args.Name(a.Metric))
				time.Sleep(wait)
				for {
					ray := xray.ROOT.Fork().With(args.Name(a.Metric))
					ray.Debug("Performing work on :name")
					results, err := db.Read(conn, a)
					if err != nil {
						ray.Error("Error obtaining data for :name - :err", args.Error{Err: err})
					} else {
						err = influx.Send(flagInfluxDB, results)
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

func init() {
	serveCmd.Flags().StringVarP(
		&flagInfluxDB,
		"influxdb",
		"i",
		"",
		"InfluxDB UDP listener address",
	)
	serveCmd.Flags().DurationVarP(
		&flagDefaultScheduleTime,
		"schedule",
		"s",
		2*time.Minute,
		"Default schedule interval for metrics",
	)
}
