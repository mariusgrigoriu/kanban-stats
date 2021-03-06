package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mariusgrigoriu/kanban-stats/trello"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

type flags struct {
	trelloKey, trelloToken, trelloBoardID, influxURL, influxDB, influxUser, influxPassword string
	verbose, dryRun, version                                                               bool
}

func getCommandLineFlags() (flags flags) {
	flag.StringVar(&flags.trelloKey, "trellokey", "", "Trello application key")
	flag.StringVar(&flags.trelloToken, "trellotoken", "", "Trello access token")
	flag.StringVar(&flags.trelloBoardID, "boardid", "", "Trello board ID")
	flag.StringVar(&flags.influxURL, "influxurl", "http://localhost:8086", "http://host:port")
	flag.StringVar(&flags.influxDB, "influxdb", "kanban", "Influx datbase name")
	flag.StringVar(&flags.influxUser, "influxuser", "root", "Influx username")
	flag.StringVar(&flags.influxPassword, "influxpass", "root", "Influx password")
	flag.BoolVar(&flags.verbose, "v", false, "Print verbose information")
	flag.BoolVar(&flags.dryRun, "d", false, "Dry run does not output to database")
	flag.BoolVar(&flags.version, "version", false, "Print version and exit")
	flag.Parse()
	return
}

func main() {
	flags := getCommandLineFlags()

	if flags.version {
		fmt.Println(GetVersion())
		return
	}

	if flags.trelloKey == "" {
		flags.trelloKey = os.Getenv("TRELLO_KEY")
	}

	if flags.trelloToken == "" {
		flags.trelloToken = os.Getenv("TRELLO_TOKEN")
	}

	if flags.influxPassword == "" {
		flags.influxPassword = os.Getenv("INFLUX_PASS")
	}

	log.Print(ApplicationName, ": start")
	defer log.Print(ApplicationName, ": end")

	trello := trello.NetworkClient{
		Key:   flags.trelloKey,
		Token: flags.trelloToken,
	}

	board := GetBoardFromTrello(trello, flags.trelloBoardID)

	if flags.verbose {
		printInfo(board)
	}
	if !flags.dryRun {
		influxdb, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
			Addr:     flags.influxURL,
			Username: flags.influxUser,
			Password: flags.influxPassword,
		})
		if err != nil {
			log.Fatal(ApplicationName, ": influxdb.NewClient: ", err)
		}

		err = writePointsToDatabase(influxdb, BuildPointBatch(flags.influxDB, board))
		if err != nil {
			log.Fatal(ApplicationName, ": writeStatsToDatabase: ", err)
		}
		//writeListsToDatabase(influxdbClient, board.Columns, flags.trelloBoardID)
	}
}

func printInfo(board Board) {
	fmt.Printf("Board ID: %v\n\nLists\n", board.GetID())
	labelsFound := make(map[string]string, 1000)
	for _, column := range board.GetColumns() {
		fmt.Printf("%s(%s): %d\n", column.GetName(), column.GetID(), len(column.GetCards()))
		for _, card := range column.GetCards() {
			for _, label := range card.Labels {
				labelsFound[label.Name] = label.ID
			}
		}
	}

	fmt.Printf("\nLabels In Use\n")
	for name, id := range labelsFound {
		fmt.Printf("%s: %s\n", name, id)
	}
}
