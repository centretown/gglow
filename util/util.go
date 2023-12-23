package main

import (
	"flag"
	"fmt"
	"gglow/resources"
	"gglow/transactions"
	"os"

	"fyne.io/fyne/v2/app"
)

var transactionFile string

const (
	transactionDefault = "transactions.yaml"
	transactionUsage   = "transaction file"
)

func init() {
	flag.StringVar(&transactionFile, "t", transactionDefault, transactionUsage+" (short form)")
	flag.StringVar(&transactionFile, "transaction", transactionDefault, transactionUsage)
}

func main() {
	app.NewWithID(resources.AppID)
	flag.Parse()

	transaction, err := transactions.ReadTransaction(transactionFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	transaction.Process()
	transaction.ShowLogs()
}
