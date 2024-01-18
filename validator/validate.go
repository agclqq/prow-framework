package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"github.com/agclqq/prow/dynamicstruct"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := TransInit(v, "zh"); err != nil {
			return
		}
		RegisterTagName(v, "alias")
	}
}

func GetValidationError(err error) error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {
		TransErrs := errs.Translate(trans)
		firstErr := ""
		for _, v := range TransErrs {
			firstErr = v
			break
		}
		return errors.New(firstErr)
	} else {
		return err
	}
}
func GetValidationErrors(err error) []error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {
		TransErrs := errs.Translate(trans)
		resErrs := make([]error, 0)
		for _, v := range TransErrs {
			resErrs = append(resErrs, errors.New(v))
		}
		return resErrs
	} else {
		return []error{err}
	}
}

type SetsValidator struct {
	rsMetaField map[string][]reflect.StructField
}

func getMetaField(rsMetaField map[string][]reflect.StructField, job any) map[string][]reflect.StructField {
	vJob := reflect.ValueOf(job)
	if vJob.Kind() == reflect.Ptr {
		return getMetaField(rsMetaField, vJob.Elem().Interface())
	}
	tJob := reflect.TypeOf(job)
	for i := 0; i < tJob.NumField(); i++ {
		if tJob.Field(i).Type.Kind() == reflect.Struct {
			getMetaField(rsMetaField, vJob.Field(i).Interface())
			continue
		}
		if tJob.Field(i).Tag.Get("header") != "" {
			rsMetaField["header"] = append(rsMetaField["header"], tJob.Field(i))
			continue
		}
		if tJob.Field(i).Tag.Get("query") != "" {
			rsMetaField["query"] = append(rsMetaField["query"], tJob.Field(i))
			continue
		}
		if tJob.Field(i).Tag.Get("uri") != "" {
			rsMetaField["uri"] = append(rsMetaField["uri"], tJob.Field(i))
			continue
		}
		//如果有其他的集合验证规则，可以继续补充
		rsMetaField["form"] = append(rsMetaField["form"], tJob.Field(i))
	}
	return rsMetaField
}

// Verify  集合验证
// 由于gin的验证器，每个分类都需要分开调，可能一次请求中需要验证多次，代码繁琐
// 封装一个方法，做成集合验证
func (s *SetsValidator) Verify(ctx *gin.Context, job any) error {
	//解析
	s.rsMetaField = make(map[string][]reflect.StructField)
	getMetaField(s.rsMetaField, job)

	type Struct struct {
		strct reflect.Type
		index map[string]int
	}
	//封装，并验证
	var err error
	for k, v := range s.rsMetaField {
		tmpStruct := make([]reflect.StructField, 0)
		for _, vv := range v {
			//t := reflect.StructField{
			//	Name: vv.Name,
			//	Type: vv.Type,
			//	Tag:  vv.Tag,
			//}
			tmpStruct = append(tmpStruct, vv)
		}
		newStruct := dynamicstruct.NewInstance(tmpStruct)
		fmt.Printf("%v\n", newStruct)

		switch strings.ToLower(k) {
		case "header":
			err = ctx.ShouldBindHeader(newStruct)
		case "uri":
			err = ctx.ShouldBindUri(newStruct)
		case "query":
			err = ctx.ShouldBindQuery(newStruct)
			//如果有其他的集合验证规则，可以继续补充
		default:
			err = ctx.ShouldBind(newStruct)
		}
		fmt.Printf("%v\n", newStruct)
		if err != nil {
			return err
		}
		//ts := reflect.TypeOf(newStruct)
		//vs := reflect.ValueOf(newStruct)
		//for i := 0; i < ts.NumField(); i++ {
		//	fmt.Println(ts.Field(i).Name)
		//	fmt.Println(vs.Field(i))
		//}
	}
	return nil
}

func ValidElementName(name string) bool {
	r := regexp.MustCompile(`[a-z][_a-z0-9]*`)
	return r.FindString(name) == name
}

func ValidUsername(email string, username *string) error {
	emailUsername := strings.Split(email, "@")[0]
	if *username == "" {
		*username = emailUsername
		return nil
	}

	if emailUsername != *username {
		*username = emailUsername
		return errors.New("用户名必须和邮箱用户名一致")
	}
	return nil
}
