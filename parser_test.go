package main

import (
	"io"
	"os"
	"testing"

	"io-lib/dbio"
	"io-lib/text_files"

	_ "github.com/mattn/go-sqlite3"
)

func TestParseUSXFiles(t *testing.T) {
	//dir := "/Users/gary/FCBH2024/download/ABIWBT/ABIWBTN_ET-usx"
	dir := "/Users/gary/FCBH2024/download/ENGWEB/ENGWEBN_ET-usx"

	files, err := text_files.ReadDir(dir, ".usx")
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("no .usx files found in %s", dir)
	}
	t.Logf("found %d .usx files", len(files))

	db, dbPath, err := dbio.OpenDB()
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}
	defer dbio.CloseDB(db, dbPath)

	log := dbio.NewLogger(db, "usx-parser-test")

	parser := USXParser{db: db, log: log}
	if err = parser.parseUSXFiles(files); err != nil {
		t.Fatalf("parseUSXFiles failed: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM scripts").Scan(&count)
	if err != nil {
		t.Fatalf("counting scripts: %v", err)
	}
	t.Logf("inserted %d script records", count)
	if count == 0 {
		t.Error("expected script records but got 0")
	}

	// Copy the database to the usx-parser directory
	if err = copyFile(dbPath, "usx-parser.db"); err != nil {
		t.Fatalf("copying database: %v", err)
	}
	t.Logf("wrote usx-parser.db")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
