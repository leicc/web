package orm

import (
	"git.lcc.lib/core"
	"strconv"
)

const (
	DT_SQL         = "sql"
	DT_AUTO        = "auto"
	DT_BINARY      = "binary"
	OP_AS          = "AS"
	OP_MAX         = "MAX"
	OP_MIN         = "MIN"
	OP_SUM         = "SUM"
	OP_AVG         = "AVG"
	OP_COUNT       = "COUNT"
	OP_EQ          = "="
	OP_NE          = "<>"
	OP_GT          = ">"
	OP_LT          = "<"
	OP_GE          = ">="
	OP_LE          = "<="
	OP_BETWEEN     = "BETWEEN"
	OP_LIKE        = "LIKE"
	OP_NOTLIKE     = "NOT LIKE"
	OP_REGEXP      = "REGEXP"
	OP_ISNULL      = "IS NULL"
	OP_ISNOTNULL   = "IS NOT NULL"
	OP_IN          = "IN"
	OP_NOTIN       = "NOT IN"
	OP_AND         = "AND"
	OP_OR          = "OR"
	OP_NOT         = "NOT"
	OP_SQL         = "SQL"
	ASC            = "ASC"
	DESC           = "DESC"
	SQLMODE_QUERY  = 1
	SQLMODE_OPDATA = 0
)

var (
	dbIniFile string
)

func init() {
	dbIniFile = "./config/db.ini"
}

func SetDBIni(inifile string) {
	dbIniFile = inifile
}

func DB(dbini string) *OrmQuery {
	var err error
	max_open_conns, max_idle_conns := 0, 0
	ini := core.NewIni(dbIniFile)
	dsn := ini.GetItem(dbini, "Dsn")
	driver := ini.GetItem(dbini, "Driver")
	open_conns := ini.GetItem(dbini, "MaxOpenConns")
	idle_conns := ini.GetItem(dbini, "MaxIdleConns")
	if max_open_conns, err = strconv.Atoi(open_conns); err != nil || max_open_conns < 1 {
		max_open_conns = 128
	}
	if max_idle_conns, err = strconv.Atoi(idle_conns); err != nil || max_idle_conns < 1 {
		max_idle_conns = 64
	}
	db := NewOrmQuery(driver, dsn, max_open_conns, max_idle_conns)
	return db
}

type Model struct {
	db     *OrmQuery
	table  string
	pri_id int64
	fields map[string]string
}

func (this *Model) Init(dbini, table string, field map[string]string) {
	this.table = table
	this.field = field
	this.db = DB(dbini)
}

func (this *Model) NewOne(fields map[string]string) int64 {
	this.db.Table(this.table).Insert(fields)
}

func (this *Model) GetOne(id int64) map[string]string {
	this.db.Table(this.table).Where(this.pri_id, id).GetRow()
}
