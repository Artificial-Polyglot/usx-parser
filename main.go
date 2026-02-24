package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Artificial-Polyglot/io-lib/dbio"
	"github.com/Artificial-Polyglot/io-lib/text_files"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dir := flag.String("dir", "", "directory containing USX files to parse")
	flag.Parse()

	if *dir == "" {
		fmt.Fprintln(os.Stderr, "error: -dir is required")
		flag.Usage()
		os.Exit(1)
	}

	db, dbPath, err := dbio.OpenDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening database: %v\n", err)
		os.Exit(1)
	}

	log := dbio.NewLogger(db, "usx-parser")
	log.Info("starting with dir:", *dir)

	files, err := text_files.ReadDir(*dir, ".usx")
	if err != nil {
		log.Error(err)
		dbio.CloseDB(db, dbPath)
		os.Exit(1)
	}

	parser := USXParser{db: db, log: log}
	if err = parser.parseUSXFiles(files); err != nil {
		log.Error(err)
		dbio.CloseDB(db, dbPath)
		os.Exit(1)
	}

	log.Info("completed successfully")

	if err = dbio.OutputDB(db, dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "error writing database: %v\n", err)
		os.Exit(1)
	}
}
