package ref

import "reflect"

import "errors"

import "fmt"

import "strings"

const Tag = "field"

type Field struct {
	Name    string
	Ref     string
	Ignored bool
	Value   interface{}
}

type Entity struct {
	table        string
	Value        interface{}
	Fields       map[string]Field
	FieldsString string
}

const (
	AND = iota
	OR
)

var _logical = []string{"AND", "OR"}

func NewEntity(table string, value interface{}) (*Entity, error) {
	types := reflect.TypeOf(value)
	if types.Kind() != reflect.Struct {
		return nil, errors.New("Value is not struct")
	}

	values := reflect.ValueOf(value)

	fields := make(map[string]Field)
	fieldsString := ""

	n := types.NumField()
	for i := 0; i < n; i++ {
		field := types.Field(i)
		name := field.Name
		value := values.Field(i).Interface()
		tag := field.Tag.Get(Tag)

		fields[name] = Field{
			Name:    name,
			Ref:     tag,
			Value:   value,
			Ignored: false,
		}

		if tag == "-" {
			f := fields[name]
			f.Ignored = true
			fields[name] = f
			continue
		}

		if i == 0 {
			fieldsString = fmt.Sprintf("`%s`", tag)
		} else {
			fieldsString = fmt.Sprintf("%s, `%s`", fieldsString, tag)
		}
	}

	entity := &Entity{
		table:        table,
		Value:        value,
		Fields:       fields,
		FieldsString: fieldsString,
	}

	return entity, nil
}

func (entity *Entity) Insert() (string, []interface{}) {

	values := []interface{}{}

	for _, field := range entity.Fields {
		if field.Ignored {
			continue
		}
		values = append(values, field.Value)
	}

	qmark := []string{}

	for i := 0; i < len(entity.Fields); i++ {
		qmark = append(qmark, "?")
	}
	qmarks := strings.Join(qmark, ",")
	sql := fmt.Sprintf("INSERT `%s` INTO (%s) VALUES(%s)", entity.table, entity.FieldsString, qmarks)

	return sql, values
}

func (entity *Entity) FindBy(fields []string, logical int) (string, error) {
	if len(fields) == 0 {
		return "", errors.New("条件个数为0")
	}

	sub := ""

	// 检查fieds是否存在
	for i, field := range fields {
		f, ok := entity.Fields[field]
		if !ok {
			return "", fmt.Errorf("字段中不包含%s", f.Ref)
		}
		if i == 0 {
			sub = fmt.Sprintf("`%s`= ?", f.Ref)
		} else {
			sub = fmt.Sprintf("%s %s `%s`= ?", sub, _logical[logical], f.Ref)
		}
	}

	return fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", entity.FieldsString, entity.table, sub), nil
}

func (entity *Entity) DeleteBy(fields []string, logical int) (string, error) {
	if len(fields) == 0 {
		return fmt.Sprintf("DELETE FROM `%s`", entity.table), errors.New("条件个数为0")
	}

	sub := ""

	// 检查fieds是否存在
	for i, field := range fields {
		f, ok := entity.Fields[field]
		if !ok {
			return "", fmt.Errorf("字段中不包含%s", f.Ref)
		}
		if i == 0 {
			sub = fmt.Sprintf("`%s`= ?", f.Ref)
		} else {
			sub = fmt.Sprintf("%s %s `%s`= ?", sub, _logical[logical], f.Ref)
		}
	}

	return fmt.Sprintf("DELETE FROM `%s` WHERE %s", entity.table, sub), nil
}

func (entity *Entity) Update(fields ...string) string {

	qmark := []string{}

	if len(fields) == 0 {
		for _, et := range entity.Fields {
			qmark = append(qmark, fmt.Sprintf("`%s`", et.Ref))
		}
	} else {
		fieldsMap := map[string]struct{}{}
		for _, field := range fields {
			fieldsMap[field] = struct{}{}
		}
		for _, et := range entity.Fields {
			if _, ok := fieldsMap[et.Name]; ok {
				qmark = append(qmark, fmt.Sprintf("`%s`", et.Ref))
			}
		}
	}

	qmarks := strings.Join(qmark, "=?,")
	qmarks = fmt.Sprintf("%s=?", qmarks)
	sql := fmt.Sprintf("UPDATE `%s` SET %s", entity.table, qmarks)

	return sql
}
