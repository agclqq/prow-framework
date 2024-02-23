package validator

import (
	"testing"
)

type HParams struct {
	ID int `json:"id" header:"id"`
}
type QParams struct {
	Age uint8 `json:"age" query:"age"`
}
type All struct {
	HParams
	QParams
	Hao uint8 `json:"hao" form:"hao"`
}

func Test_getMetaField(t *testing.T) {
	//var all = All{}
	//var m = make(map[string][]reflect.StructField)
	//rs := getMetaField(m, all)
	//if len(rs) != 3 {
	//	t.Errorf("错误")
	//	return
	//}
}

//func TestSetsValidator_Verify(t *testing.T) {
//	set := &SetsValidator{}
//	var a = All{}
//	set.Verify(&gin.Context{}, a)
//}
