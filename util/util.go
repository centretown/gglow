package main

import (
	"flag"
	"fmt"
	"glow-gui/resources"
	"glow-gui/transactions"
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

	// err = logs.Process()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// f := func(dataIn effects.IoHandler, dataOut effects.IoHandler) error {
	// 	err := dataOut.CreateNewDatabase()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = WriteDatabase(dataIn, dataOut)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	dataOut.OnExit()
	// 	return nil
	// }

	// err := f(dataIn, sqlio.NewSqlLiteHandler())
	// if err != nil {
	// 	return
	// }
	// f(dataIn, sqlio.NewMySqlHandler())
	// if err != nil {
	// 	return
	// }
	// fmt.Println("Complete!")
}
