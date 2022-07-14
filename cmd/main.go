package cmd

import (
	"database/sql"
	"github.com/mono83/artifacts/config"
	"github.com/mono83/artifacts/data"
	"github.com/mono83/artifacts/db"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
)

var (
	flagConfigurationFile string
)

// MainCmd is a main command
var MainCmd = &cobra.Command{
	Use:     "artifacts",
	Version: "1.0.1",
}

func init() {
	MainCmd.AddCommand(
		testCmd,
		serveCmd,
	)

	MainCmd.PersistentFlags().StringVarP(
		&flagConfigurationFile,
		"config",
		"c",
		"config.yaml",
		"Configuration file location",
	)
}

func configure() (*config.Configuration, *sql.DB, []data.Artifact, error) {
	xray.BOOT.Info("Reading configuration file :name", args.Name(flagConfigurationFile))
	cnf, err := config.Read(flagConfigurationFile)
	if err != nil {
		return nil, nil, nil, err
	}

	xray.BOOT.Info("Establishing connection to MySQL server")
	conn, err := db.Connect(cnf.MySQLDSN)
	if err != nil {
		return nil, nil, nil, err
	}

	xray.BOOT.Info("Reading artifacts configuration from table :name", args.Name(cnf.MySQLTable))
	artifacts, err := db.ReadFromConfigTable(conn, cnf.MySQLTable)
	if err != nil {
		return nil, nil, nil, err
	}

	return cnf, conn, artifacts, nil
}
