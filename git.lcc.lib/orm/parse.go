package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type OrmParse struct {
	db    *sql.DB
	tx    *sql.Tx
	marks []interface{}
}

func NewOrmParse(driver, dsn string, maxopenconns, maxidleconns int) *OrmParse {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxopenconns)
	db.SetMaxIdleConns(maxidleconns)
	return &OrmParse{db: db}
}

func (this *OrmParse) ParseSql(args map[string]interface{}) string {
	this.marks = make([]interface{}, 0)
	mode, _ := args["mode"].(string)
	tables, _ := args["table"].([]tableSt)
	query := ""
	switch mode {
	case "select":
		field, _ := args["field"].(string)
		if field == "" {
			field = "*"
		}
		query = "SELECT " + field + " FROM " + this.parseSqlTable(tables, true)
		if _, isok := args["where"]; isok {
			wheres, _ := args["where"].([]whereSt)
			if len(wheres) > 0 {
				query += " WHERE " + this.parseSqlWhere(wheres)
			}
		}
		if _, isok := args["group"]; isok {
			groups, _ := args["group"].([]string)
			if len(groups) > 0 {
				query += " GROUP BY " + strings.Join(groups, ",")
			}
			if _, isok := args["having"]; isok {
				having, _ := args["having"].(string)
				query += " HAVING " + having
			}
		}
		if _, isok := args["order"]; isok {
			orders, _ := args["order"].([]string)
			if len(orders) > 0 {
				query += "  ORDER BY " + strings.Join(orders, ",")
			}
		}
	case "update":
		values, _ := args["value"].([]valueSt)
		query = "UPDATE " + this.parseSqlTable(tables, false) + " SET " + this.parseSqlUpdateValue(values)
		if _, isok := args["where"]; isok {
			wheres, _ := args["where"].([]whereSt)
			if len(wheres) > 0 {
				query += " WHERE " + this.parseSqlWhere(wheres)
			}
		}
	case "insert":
		values, _ := args["value"].([]valueSt)
		query = "INSERT INTO " + this.parseSqlTable(tables, false) + this.parseSqlInsertValue(values)
	case "delete":
		query = "DELETE FROM " + this.parseSqlTable(tables, false)
		if _, isok := args["where"]; isok {
			wheres, _ := args["where"].([]whereSt)
			if len(wheres) > 0 {
				query += " WHERE " + this.parseSqlWhere(wheres)
			}
		}
	}
	return query
}

func (this *OrmParse) parseSqlWhere(wheres []whereSt) string {
	ocwhere := make([]string, 0)
	logical := false
	for _, where := range wheres {
		upername := strings.ToUpper(where.name)
		if (logical && where.name != ")" && upername != "OR") || where.logical == OP_NOT {
			ocwhere = append(ocwhere, fmt.Sprintf("%s ", where.logical))
		}
		if where.name == "(" || where.name == ")" {
			ocwhere = append(ocwhere, where.name)
			logical = where.name != "("
		} else if upername == "OR" {
			ocwhere = append(ocwhere, where.name)
			logical = false
		} else {
			switch where.opt {
			case OP_BETWEEN:
				if val, isok := where.value.([2]int); isok {
					ocwhere = append(ocwhere, fmt.Sprintf("%s %s ? AND ?", where.name, where.opt))
					this.marks = append(this.marks, val[0], val[1])
				}
			case OP_IN, OP_NOTIN:
				if val, isok := where.value.(string); isok {
					ocwhere = append(ocwhere, fmt.Sprintf("%s %s (%s)", where.name, where.opt, val))
				}
			case OP_SQL:
				ocwhere = append(ocwhere, fmt.Sprintf("%s", where.name))
			case OP_ISNULL, OP_ISNOTNULL:
				ocwhere = append(ocwhere, fmt.Sprintf("%s %s", where.name, where.opt))
			default:
				if where.ftype == DT_SQL {
					ocwhere = append(ocwhere, fmt.Sprintf("%s %s %s", where.name, where.opt, where.value))
				} else {
					ocwhere = append(ocwhere, fmt.Sprintf("%s %s ?", where.name, where.opt))
					this.marks = append(this.marks, where.value)
				}
			}
		}
	}
	return strings.Join(ocwhere, " ")
}

func (this *OrmParse) parseSqlTable(tables []tableSt, isalias bool) string {
	query := ""
	for _, table := range tables {
		if table.alise != "" && isalias {
			table.name = fmt.Sprintf("%s %s %s", table.name, OP_AS, table.alise)
		}
		if table.on != "" {
			query += fmt.Sprintf(" LEFT JOIN %s ON (%s)", table.name, table.on)
		} else if query == "" {
			query += table.name
		} else {
			query += fmt.Sprintf(" ,%s", table.name)
		}
	}
	return query
}

func (this *OrmParse) parseSqlUpdateValue(values []valueSt) string {
	ovals := make([]string, len(values))
	for idx, val := range values {
		if val.ftype == DT_SQL {
			ovals[idx] = fmt.Sprintf("`%s`=%s", val.name, val.value)
		} else {
			ovals[idx] = fmt.Sprintf("`%s`=?", val.name)
			this.marks = append(this.marks, val.value)
		}
	}
	return strings.Join(ovals, ",")
}

func (this *OrmParse) parseSqlInsertValue(values []valueSt) string {
	ovals := make([]string, len(values))
	fields := make([]string, len(values))
	for idx, val := range values {
		fields[idx] = fmt.Sprintf("`%s`", val.name)
		if val.ftype == DT_SQL {
			szstr, _ := val.value.(string)
			fields[idx] = szstr
		} else {
			ovals[idx] = "?"
			this.marks = append(this.marks, val.value)
		}
	}
	return fmt.Sprintf("(%s)VALUES(%s)", strings.Join(fields, ","), strings.Join(ovals, ","))
}

func (this *OrmParse) GetAll(query string, offset, limit int) []map[string]string {
	if limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}
	res := this.Execute(query, SQLMODE_QUERY)
	rows, _ := res.(*sql.Rows)
	datalist := this.fetch(rows, false)
	return datalist
}

func (this *OrmParse) GetFirst(query string) map[string]string {
	query += " limit 1"
	fmt.Println(query)
	res := this.Execute(query, SQLMODE_QUERY)
	rows, _ := res.(*sql.Rows)
	datalist := this.fetch(rows, true)
	fmt.Println(datalist)
	if len(datalist) == 1 {
		return datalist[0]
	}
	return nil
}

func (this *OrmParse) fetch(res *sql.Rows, isonlyfirst bool) []map[string]string {
	columns, _ := res.Columns()
	nlen := len(columns)
	datalist := make([]map[string]string, 0)
	defer res.Close()
	valslice := make([]interface{}, nlen)
	for i := 0; i < nlen; i++ {
		valslice[i] = new(string)
	}
	for res.Next() {
		res.Scan(valslice...)
		valmap := make(map[string]string)
		for i := 0; i < nlen; i++ {
			str, _ := valslice[i].(*string)
			valmap[columns[i]] = *str
		}
		datalist = append(datalist, valmap)
		if isonlyfirst {
			break
		}
	}
	return datalist[0:len(datalist)]
}

func (this *OrmParse) Execute(query string, sqlmode int8) interface{} {
	var res interface{}
	if err := this.db.Ping(); err != nil {
		panic(err)
	}
	stmt, err := this.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	nmark := len(this.marks)
	if sqlmode == SQLMODE_QUERY {
		if nmark > 0 {
			res, err = stmt.Query(this.marks...)
		} else {
			res, err = stmt.Query()
		}
	} else {
		if this.tx == nil && nmark > 0 {
			res, err = stmt.Exec(this.marks...)
		} else if this.tx == nil {
			res, err = stmt.Exec()
		} else if nmark > 0 && this.tx != nil {
			res, err = this.tx.Stmt(stmt).Exec(this.marks...)
		} else {
			res, err = this.tx.Stmt(stmt).Exec()
		}
	}
	if err != nil {
		panic(err)
	}
	return res
}

func (this *OrmParse) Begin() bool {
	var err error
	this.tx, err = this.db.Begin()
	if err != nil {
		return false
	}
	return true
}

func (this *OrmParse) Rollback() {
	if this.tx != nil {
		this.tx.Rollback()
		this.tx = nil
	}
}

func (this *OrmParse) Commit() {
	if this.tx != nil {
		this.tx.Commit()
		this.tx = nil
	}
}

func (this *OrmParse) Close() {
	this.db.Close()
	if this.tx != nil {
		this.tx = nil
	}
}
