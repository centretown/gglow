package main

import (
	"flag"
	"fmt"
	"gglow/transactions"
	"os"
)

var transactionFile string

const (
	transactionDefault = "transaction.yaml"
	transactionUsage   = "transaction file"
)

func init() {
	flag.StringVar(&transactionFile, "t", transactionDefault, transactionUsage+" (short form)")
	flag.StringVar(&transactionFile, "transaction", transactionDefault, transactionUsage)
}

func main() {
	flag.Parse()

	if transactionFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	transaction, err := transactions.ReadTransaction(transactionFile)
	if err != nil {
		flag.Usage()
		fmt.Println(transactionFile, err)
		os.Exit(1)
	}

	transaction.Process()
	transaction.ShowLogs()
}
