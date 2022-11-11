package migrate

import (
	"bytes"
	"database/sql"
	"embed"
	"io/ioutil"
)

func MigrateFS(db *sql.DB, ddlfs embed.FS, dirName string) error {
	fs, err := ddlfs.ReadDir(dirName)
	if err != nil {
		return err
	}
	ddl := ""
	for _, f := range fs {
		data, err := ioutil.ReadFile(dirName + "/" + f.Name())
		if err != nil {
			return err
		}
		ddl += string(data) + "\n"
	}
	return Migrate(db, ddl)
}

func Migrate(db *sql.DB, ddl string) error {
	parse := New()
	statements, err := parse.Parse(bytes.NewBufferString(ddl))
	if err != nil {
		return err
	}
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, statement := range statements {
		if _, ok := completed[statement.Name]; ok {

			continue
		}
		if _, err := db.Exec(statement.Value); err != nil {
			return err
		}
		if err := insertMigration(db, statement.Name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`
