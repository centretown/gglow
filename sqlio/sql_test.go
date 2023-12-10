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

	frame, err := sqlh.ReadEffect("Scan Complementary")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(frame)

	frame.Interval += 10
	err = sqlh.WriteEffect("Scan Complementary2", frame)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(frame)

	frame.Interval += 10
	err = sqlh.WriteEffect("Scan Complementary4", frame)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(frame)

	frame2, err := sqlh.ReadEffect("Scan Complementary4")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(frame2)

	if frame.Interval != frame2.Interval {
		t.Fatalf("read and write differ: %d != %d", frame2.Interval, frame.Interval)
	}

	frame2.Interval += 10

	err = sqlh.CreateNewEffect("Scan Complementary9", frame2)
	if err != nil {
		t.Log(err)
	}

	err = sqlh.CreateNewEffect("Scan Complementary5", frame2)
	if err != nil {
		t.Log(err)
	}
	t.Log(frame2)
}
