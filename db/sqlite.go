package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wutzi15/storago/files"
	"github.com/wutzi15/storago/interactive"
)

func OpenSqlite() *sql.DB {
	filename := "storago.sqlite"
	newFile := false

	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		newFile = true
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Panic(err)
	}

	if !newFile {
		return db
	}

	sqlStmt := "CREATE TABLE  IF NOT EXISTS files (id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL,  size INTERGER DEFAULT 0,  scantime INTERGER DEFAULT 0, isDir INTEGER DEFAULT 0)"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Panic(err)
	}

	return db
}

func InsertFilesIntoDB(file *files.File, scantime int64, db *sql.DB) {
	sqlStmt := fmt.Sprintf(`INSERT INTO files (name, size, scantime, isDir) VALUES (?,?,%d,?)`, scantime)
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	recurseThroughFiles(file, stmt)
	tx.Commit()
}

func recurseThroughFiles(file *files.File, stmt *sql.Stmt) {
	isDir := 0
	if file.IsDir {
		isDir = 1
	}
	name := interactive.GetParentName(file)
	// fmt.Printf("Inserting: %s, %d, %d\n", name, file.Size, isDir)
	_, err := stmt.Exec(name, file.Size, isDir)
	if err != nil {
		log.Panic(err)
	}
	for _, subfile := range file.Files {
		recurseThroughFiles(subfile, stmt)
	}

}
