package crud

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pkg.oars.vip/go-pkg/server/base"
)

const (
	prefix    = "-- name: "
	comment   = "--"
	newline   = "\n"
	delimiter = ";"
)

// Statement represents a statment in the sql file.
type Statement struct {
	Name   string
	Value  string
	Driver string
}

// SqlFileParser parses the sql file.
type SqlFileParser struct {
	prefix string
}

// New returns a new parser.
func NewSqlFileParser() *SqlFileParser {
	return NewSqlFileParserPrefix(prefix)
}

// NewPrefix returns a new parser with the given prefix.
func NewSqlFileParserPrefix(prefix string) *SqlFileParser {
	return &SqlFileParser{prefix: prefix}
}

// ParseFile parses the sql file.
func (p *SqlFileParser) ParseFile(filepath string) ([]*Statement, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return p.Parse(f)
}

// Parse parses the sql file and returns a list of statements.
func (p *SqlFileParser) Parse(r io.Reader) ([]*Statement, error) {
	var (
		stmts []*Statement
		stmt  *Statement
	)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, prefix) {
			stmt = new(Statement)
			stmt.Name, stmt.Driver = parsePrefix(line, p.prefix)
			stmts = append(stmts, stmt)
		}
		if strings.HasPrefix(line, comment) {
			continue
		}
		if stmt != nil {
			stmt.Value += line + newline
		}
	}
	for _, stmt := range stmts {
		stmt.Value = strings.TrimSpace(stmt.Value)
	}
	return stmts, nil
}

func parsePrefix(line, prefix string) (name string, driver string) {
	line = strings.TrimPrefix(line, prefix)
	line = strings.TrimSpace(line)
	fmt.Sscanln(line, &name, &driver)
	return
}

var QuerySqlFiles embed.FS

func InitQuerySql(sqlFiles embed.FS) error {
	QuerySqlFiles = sqlFiles
	return readQuerySql()
}

var querySqls sync.Map

func readQuerySql() error {
	des, err := QuerySqlFiles.ReadDir("query")
	if err != nil {
		return err
	}
	p := NewSqlFileParser()
	for _, d := range des {
		fd, err := QuerySqlFiles.ReadFile("query/" + d.Name())
		if err != nil {
			return err
		}
		stmts, err := p.Parse(bytes.NewBuffer(fd))
		if err != nil {
			return err
		}
		for _, stmt := range stmts {
			querySqls.Store(strings.TrimSuffix(d.Name(), ".sql")+"."+stmt.Name, stmt.Value)
		}

	}
	return nil
}

func BuildSubQuery(db *gorm.DB, name string, args map[string]any) (*gorm.DB, error) {

	sqlstr := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		qdb, err := BuildQuery(tx, name, args)
		if err != nil {
			return nil
		}
		var m []string
		return qdb.Find(&m)
	})
	tx := db.Table("(" + sqlstr + ") as m")
	tx.Statement.Table = "m"
	return tx, nil
}

func BuildQuery(db *gorm.DB, name string, args map[string]any) (*gorm.DB, error) {
	sqlstr, ok := querySqls.Load(name)
	if !ok {
		return db, errors.New("name:" + name + " not exist")
	}
	// var ps []any
	// for k, v := range args {
	// 	if v == "" {
	// 		continue
	// 	}
	// 	if strings.HasPrefix(k, "__") {
	// 		continue
	// 	}
	// 	ps = append(ps, sql.Named(k, v))
	// }
	t, err := template.New("sql").Parse(sqlstr.(string))
	if err != nil {
		return db, err
	}
	b := bytes.NewBuffer(nil)
	err = t.Execute(b, args)
	if err != nil {
		return db, err
	}

	return db.Raw(b.String(), args), nil
}

func BuildLike(s string) string {
	if s == "" {
		return s
	}
	return "%" + s + "%"
}
func FindWithPage(db *gorm.DB, g *gin.Context, res any) (any, error) {
	var page base.Page
	g.ShouldBindQuery(&page)
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return nil, err
	}
	if page.PageNum > 0 && page.PageSize > 0 {
		db = db.Offset((page.PageNum - 1) * page.PageSize).Limit(page.PageSize)
	}
	err = db.Find(&res).Error
	if err != nil {
		return nil, err
	}
	if page.PageNum > 0 && page.PageSize > 0 {
		return res, nil
	}
	resp := base.PageResp{
		List: res,
		Page: page,
	}
	return resp, err
}

func init() {
	readQuerySql()
}
