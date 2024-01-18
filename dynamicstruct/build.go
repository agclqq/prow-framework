package dynamicstruct

import "reflect"

type Builder struct {
	Fields []reflect.StructField
}

func (b *Builder) Build() Struct {
	typ := reflect.StructOf(b.Fields)

	index := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		index[typ.Field(i).Name] = i
	}

	return Struct{typ, index}
}

type Struct struct {
	strct reflect.Type
	index map[string]int
}

func (s *Struct) New() *Instance {
	instance := reflect.New(s.strct).Elem()
	return &Instance{instance, s.index}
}
func NewInstance(fields []reflect.StructField) *Instance {
	strc := (&Builder{fields}).Build()
	return strc.New()
}

type Instance struct {
	internal reflect.Value
	index    map[string]int
}

func (i *Instance) Field(name string) reflect.Value {
	return i.internal.Field(i.index[name])
}

func (i *Instance) SetString(name, value string) {
	i.Field(name).SetString(value)
}

func (i *Instance) SetBool(name string, value bool) {
	i.Field(name).SetBool(value)
}

func (i *Instance) SetInt64(name string, value int64) {
	i.Field(name).SetInt(value)
}

func (i *Instance) SetFloat64(name string, value float64) {
	i.Field(name).SetFloat(value)
}

func (i *Instance) Interface() interface{} {
	return i.internal.Interface()
}

func (i *Instance) Addr() interface{} {
	return i.internal.Addr().Interface()
}
