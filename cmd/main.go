package cmd

import (
	"database/sql"
	"github.com/mono83/artifacts/data"
	"github.com/mono83/artifacts/db"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
)

var (
	flagMySQLDSN   string
	flagMySQLTable string
)

// MainCmd is a main command
var MainCmd = &cobra.Command{
	Use: "artifacts",
}

func init() {
	MainCmd.AddCommand(
		testCmd,
		serveCmd,
	)

	MainCmd.PersistentFlags().StringVarP(
		&flagMySQLDSN,
		"dsn",
		"d",
		"root:root@tcp(127.0.0.1:3308)/",
		"MySQL DSN",
	)
	MainCmd.PersistentFlags().StringVarP(
		&flagMySQLTable,
		"table",
		"t",
		"__artifacts",
		"Database table name with artifacts",
	)
}

func mysql() (*sql.DB, []data.Artifact, error) {
	xray.BOOT.Info("Establishing connection to MySQL server")
	conn, err := db.Connect(flagMySQLDSN)
	if err != nil {
		return nil, nil, err
	}

	xray.BOOT.Info("Reading artifacts configuration from table :name", args.Name(flagMySQLTable))
	artifacts, err := db.ReadFromConfigTable(conn, flagMySQLTable)
	if err != nil {
		return nil, nil, err
	}

	return conn, artifacts, nil
}
