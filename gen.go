package gs

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type Literal struct {
	V interface{}
}

func (l Literal) Code() string {
	v := reflect.ValueOf(l.V)
	t := v.Type()
	switch t.Kind() {
	case reflect.Map:
		return l.mapCode(v)
	}
	return fmt.Sprintf("%#v", l.V)
}

func (l Literal) mapCode(v reflect.Value) string {
	var w bytes.Buffer
	w.WriteString(v.Type().String())
	w.WriteString("{\n")
	var keys stringValues = v.MapKeys()
	sort.Sort(keys)
	for i := 0; i < v.Len(); i++ {
		fmt.Fprintf(&w, "\t%#v: %#v,\n", keys[i].Interface(), v.MapIndex(keys[i]).Interface())
	}
	w.WriteString("}")
	return w.String()
}

type stringValues []reflect.Value

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return fmt.Sprintf("%#v", sv[i].Interface()) }
