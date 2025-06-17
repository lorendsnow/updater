// Package main implements a simple updater service that periodically downloads CSV files from a
// website, and then updates a mysql database with those values. The service uses a "blue/green"
// strategy using alternating tables to update the database since there are no unique identifiers
// to use as keys in the data. Once a table has been updated, the service emits an event notifying
// any subscribers that the active table has changed.
package main

import (
	"os"

	"github.com/lorendsnow/updater/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
