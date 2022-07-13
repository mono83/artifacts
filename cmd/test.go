package cmd

import (
	"fmt"
	"github.com/mono83/artifacts/db"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs all artifact queries sequentially",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Establishing connection to MySQL database
		conn, artifacts, err := mysql()
		if err != nil {
			return err
		}

		for i, x := range artifacts {
			fmt.Printf(
				"%03d. Working on %s\n     %s\n",
				i+1,
				x.Metric,
				x.Query,
			)

			fetchedData, err := db.Read(conn, x)
			if err != nil {
				fmt.Println("     ERROR: ", err)
			} else {
				for _, datum := range fetchedData {
					if len(datum.Group) > 0 {
						fmt.Println("    ", datum.Value, datum.Group)
					} else {
						fmt.Println("    ", datum.Value)
					}
				}
			}
		}
		return nil
	},
}

func init() {

}
