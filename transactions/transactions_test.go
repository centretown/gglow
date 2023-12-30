package transactions

import "testing"

var list = []string{
	"test_transactions.yml",
}

// const path = "/home/dave/gglow/transactions/"

func TestTransactions(t *testing.T) {
	for _, item := range list {
		// name := path + item
		t.Logf("testing '%s'", item)

		tr, err := ReadTransaction(item)
		if err != nil {
			t.Fatal(err)
		}

		err = tr.Process()
		if err != nil {
			t.Fatal(err)
		}

		tr.ShowLogs()

		if tr.HasErrors() {
			t.Fatal()
		}
	}
}
