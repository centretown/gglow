package sqlio

import "testing"

func TestSql(t *testing.T) {
	sqlh := NewSqlHandler()

	err := sqlh.Ping()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlh.OnExit()

	folder := "Grid 10x10"
	effects, err := sqlh.RefreshKeys(folder)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sqlh.Folder)
	for _, effect := range effects {
		t.Log(effect)
	}

	if folder != sqlh.Folder {
		t.Errorf("folders don't match %s!=%s", folder, sqlh.Folder)
	}

	for _, effect := range sqlh.KeyList() {
		t.Log(effect)
	}

}
