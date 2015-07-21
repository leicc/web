package orm

import (
	"database/sql"
)

type OrmQuery struct {
	data   map[string]interface{}
	parser *OrmParse
}

type tableSt struct {
	name  string
	alise string
	on    string
}

type valueSt struct {
	name  string
	value interface{}
	ftype string
}

type whereSt struct {
	name    string
	value   interface{}
	opt     string
	ftype   string
	logical string
}

func NewOrmQuery(driver, dsn string, maxopenconns, maxidleconns int) *OrmQuery {
	parser := NewOrmParse(driver, dsn, maxopenconns, maxidleconns)
	return &OrmQuery{data: make(map[string]interface{}), parser: parser}
}

func (this *OrmQuery) GetParse() *OrmParse {
	return this.parser
}

func (this *OrmQuery) AsSQL(stype string) string {
	this.data["mode"] = stype
	query := this.parser.ParseSql(this.data)
	return query
}

func (this *OrmQuery) Delete() int64 {
	query := this.AsSQL("delete")
	res := this.parser.Execute(query, SQLMODE_OPDATA)
	if ores, isok := res.(sql.Result); isok {
		nrow, _ := ores.RowsAffected()
		return nrow
	}
	return 0
}

func (this *OrmQuery) Insert(fields map[string]interface{}) int64 {
	for field, value := range fields {
		this.Value(field, value, DT_AUTO)
	}
	query := this.AsSQL("insert")
	res := this.parser.Execute(query, SQLMODE_OPDATA)
	if ores, isok := res.(sql.Result); isok {
		nid, _ := ores.LastInsertId()
		return nid
	}
	return 0
}

func (this *OrmQuery) Update(fields map[string]interface{}) int64 {
	for field, value := range fields {
		this.Value(field, value, DT_AUTO)
	}
	query := this.AsSQL("update")
	res := this.parser.Execute(query, SQLMODE_OPDATA)
	if ores, isok := res.(sql.Result); isok {
		nrow, _ := ores.RowsAffected()
		return nrow
	}
	return 0
}

func (this *OrmQuery) GetRow() map[string]string {
	query := this.AsSQL("select")
	return this.parser.GetFirst(query)
}

func (this *OrmQuery) GetList(offset, limit int64) []map[string]string {
	query := this.AsSQL("select")
	return this.parser.GetAll(query, offset, limit)
}

func (this *OrmQuery) GetColumn(offset, limit int64) []string {
	query := this.AsSQL("select")
	data := this.parser.GetAll(query, offset, limit)
	column := make([]string, len(data))
	for idx, val := range data {
		for _, item := range val {
			column[idx] = item
			break
		}
	}
	return column
}

func (this *OrmQuery) GetValue() string {
	query := this.AsSQL("select")
	data := this.parser.GetFirst(query)
	for _, val := range data {
		return val
	}
	return ""
}

func (this *OrmQuery) Where(field string, value interface{}, args ...string) *OrmQuery {
	var wheres []whereSt
	if _, isok := this.data["where"]; !isok {
		wheres = make([]whereSt, 0)
	} else {
		wheres, _ = this.data["where"].([]whereSt)
	}
	opt, ftype, logical := OP_EQ, DT_AUTO, OP_AND
	nlen := len(args)
	if nlen == 1 {
		if args[0] == OP_AND || args[0] == OP_OR || args[0] == OP_NOT {
			logical = args[0]
		} else {
			opt = args[0]
		}
	} else if nlen == 2 {
		opt = args[0]
		if args[1] == OP_AND || args[1] == OP_OR || args[1] == OP_NOT {
			logical = args[1]
		} else {
			ftype = args[1]
		}
	} else if nlen == 3 {
		opt = args[0]
		ftype = args[1]
		logical = args[2]
	}
	wheres = append(wheres, whereSt{name: field, value: value, opt: opt, ftype: ftype, logical: logical})
	this.data["where"] = wheres
	return this
}

func (this *OrmQuery) Clear(parts ...string) *OrmQuery {
	if len(parts) > 0 {
		for _, part := range parts {
			if _, isok := this.data[part]; isok {
				delete(this.data, part)
			}
		}
	} else {
		for idx, _ := range this.data {
			delete(this.data, idx)
		}
		this.data = make(map[string]interface{})
		this.data["mode"] = "select"
	}
	return this
}

func (this *OrmQuery) Value(field string, value interface{}, ftype string) *OrmQuery {
	var values []valueSt
	if _, isok := this.data["value"]; !isok {
		values = make([]valueSt, 0)
	} else {
		values, _ = this.data["value"].([]valueSt)
	}
	values = append(values, valueSt{name: field, value: value, ftype: ftype})
	this.data["value"] = values
	return this
}

func (this *OrmQuery) Field(field string) *OrmQuery {
	this.data["field"] = field
	return this
}

func (this *OrmQuery) Table(table ...string) *OrmQuery {
	var tables []tableSt
	if _, isok := this.data["table"]; !isok {
		tables = make([]tableSt, 0)
	} else {
		tables, _ = this.data["table"].([]tableSt)
	}
	nlen := len(table)
	alise, on := "", ""
	if nlen == 3 {
		alise = table[1]
		on = table[2]
	} else if nlen == 2 {
		alise = table[1]
	}
	tables = append(tables, tableSt{name: table[0], alise: alise, on: on})
	this.data["table"] = tables
	return this
}

func (this *OrmQuery) GroupBy(field ...string) *OrmQuery {
	var groups []string
	if _, isok := this.data["group"]; !isok {
		groups = make([]string, 0)
	} else {
		groups, _ = this.data["group"].([]string)
	}
	groups = append(groups, field...)
	this.data["group"] = groups
	return this
}

func (this *OrmQuery) OrderBy(field, orderby string) *OrmQuery {
	var orders []string
	if _, isok := this.data["order"]; !isok {
		orders = make([]string, 0)
	} else {
		orders, _ = this.data["order"].([]string)
	}
	orders = append(orders, field+" "+orderby)
	this.data["order"] = orders
	return this
}

func (this *OrmQuery) Having(having string) *OrmQuery {
	this.data["having"] = having
	return this
}
