package orm

import (
	"git.lcc.lib/core"
	"strconv"
)

const (
	DT_SQL         = "sql"
	DT_AUTO        = "auto"
	DT_BINARY      = "binary"
	DT_BIGINT      = "int64"
	DT_INT         = "int32"
	DT_SMALLINT    = "int16"
	DT_TINYINT     = "int8"
	DT_FLOAT       = "float"
	DT_STRING      = "string"
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
	pri_id string
	dbini  string
	fields map[string]string
}

func (this *Model) Init(dbini, table, priid string, fields map[string]string) {
	this.dbini = dbini
	this.table = table
	this.pri_id = priid
	this.fields = fields
	this.db = DB(dbini)
}

func (this *Model) Query() *OrmQuery {
	return DB(this.dbini)
}

func (this *Model) NewOne(fields map[string]interface{}) int64 {
	return this.db.Clear().Table(this.table).Insert(fields)
}

func (this *Model) GetOne(id int64) map[string]string {
	return this.db.Clear().Table(this.table).Where(this.pri_id, id).GetRow()
}

func (this *Model) Save(fields map[string]interface{}) int64 {
	return this.db.Clear().Table(this.table).Update(fields)
}

func (this *Model) Delete(id int64) int64 {
	return this.db.Clear().Table(this.table).Where(this.pri_id, id).Delete()
}

func (this *Model) parseArgs(args ...interface{}) (fields, sort, dir string) {
	nlen := len(args)
	fields, sort, dir = "*", "", DESC
	fisok, sisok, disok := false, false, false
	if nlen == 0 {
		return
	} else if nlen == 1 {
		fields, fisok = args[0].(string)
		if fisok {
			return
		}
	} else if nlen == 2 {
		fields, fisok = args[0].(string)
		sort, sisok = args[1].(string)
		if sisok && fisok && sisok {
			return
		}
	} else if nlen == 3 {
		fields, fisok = args[0].(string)
		sort, sisok = args[1].(string)
		dir, disok = args[2].(string)
		if fisok && sisok && disok {
			return
		}
	}
	panic("parse Args Argument Error!")
}

func (this *Model) Where(field string, value interface{}, args ...string) *Model {
	this.db.Where(field, value, args...)
	return this
}

func (this *Model) Clear(args ...string) *Model {
	this.db.Clear(args...)
	return this
}

func (this *Model) ListOnly(offset, limit int64, args ...interface{}) []map[string]string {
	fields, sort, dir := this.parseArgs(args...)
	if sort != "" && dir != "" {
		this.db.OrderBy(sort, dir)
	}
	list := this.db.Table(this.table).Field(fields).GetList(offset, limit)
	return list
}

func (this *Model) GetItem(args ...interface{}) map[string]string {
	fields, sort, dir := this.parseArgs(args...)
	if sort != "" && dir != "" {
		this.db.OrderBy(sort, dir)
	}
	row := this.db.Table(this.table).Field(fields).GetRow()
	return row
}

func (this *Model) GetList(offset, limit int64, args ...interface{}) map[string]interface{} {
	fields, sort, dir := this.parseArgs(args...)
	if sort != "" && dir != "" {
		this.db.OrderBy(sort, dir)
	}
	sval := this.db.Table(this.table).Field("count(1) as num").GetValue()
	total, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		total = 0
	}
	var list []map[string]string
	if total > 0 {
		list = this.db.Clear("field").Field(fields).GetList(offset, limit)
	}
	return map[string]interface{}{"total": total, "list": list}
}
