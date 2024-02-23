package validator

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/fa"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/locales/it"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/lv"
	"github.com/go-playground/locales/nl"
	"github.com/go-playground/locales/pt"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/tr"
	"github.com/go-playground/locales/vi"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	arTranslations "github.com/go-playground/validator/v10/translations/ar"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	esTranslations "github.com/go-playground/validator/v10/translations/es"
	faTandslations "github.com/go-playground/validator/v10/translations/fa"
	frTandslations "github.com/go-playground/validator/v10/translations/fr"
	idTandslations "github.com/go-playground/validator/v10/translations/id"
	itTandslations "github.com/go-playground/validator/v10/translations/it"
	jaTandslations "github.com/go-playground/validator/v10/translations/ja"
	lvTandslations "github.com/go-playground/validator/v10/translations/lv"
	nlTandslations "github.com/go-playground/validator/v10/translations/nl"
	ptTandslations "github.com/go-playground/validator/v10/translations/pt"
	ptbrTandslations "github.com/go-playground/validator/v10/translations/pt_BR"
	ruTandslations "github.com/go-playground/validator/v10/translations/ru"
	trTandslations "github.com/go-playground/validator/v10/translations/tr"
	viTandslations "github.com/go-playground/validator/v10/translations/vi"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	zhtwTranslations "github.com/go-playground/validator/v10/translations/zh_tw"
)

type Trans struct {
	vld     *validator.Validate
	uniTran *ut.UniversalTranslator
	tran    ut.Translator
}

var trans ut.Translator
var std = newStd()

type Option func(trans *Trans) error

func newStd() *Trans {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		fmt.Errorf("failed to init validator")
		return nil
	}
	t, err := New(v, WithLocal("zh"), WithAliasTag())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return t
}

func New(vld *validator.Validate, opts ...Option) (*Trans, error) {
	uni := ut.New(ar.New(), en.New(), es.New(), fa.New(), fr.New(), id.New(), it.New(), ja.New(), lv.New(), nl.New(), pt.New(), ru.New(), tr.New(), vi.New(), zh.New())
	tran := &Trans{vld: vld, uniTran: uni}
	for _, opt := range opts {
		err := opt(tran)
		if err != nil {
			return nil, err
		}
	}
	return tran, nil
}
func WithAliasTag() Option {
	return func(trans *Trans) error {
		RegisterTagName(trans.vld, "alias")
		return nil
	}
}
func WithLocal(local string) Option {
	return func(trans *Trans) error {
		tran, ok := trans.uniTran.GetTranslator(local)
		if !ok {
			return errors.New(fmt.Sprintf("cannot be translated as %s", local))
		}
		trans.tran = tran
		var err error
		switch local {
		case "ar":
			err = arTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "en":
			err = enTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "es":
			err = esTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "fa":
			err = faTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "fr":
			err = frTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "id":
			err = idTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "it":
			err = itTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "ja":
			err = jaTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "lv":
			err = lvTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "nl":
			err = nlTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "pt":
			err = ptTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "pt_BR":
			err = ptbrTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "ru":
			err = ruTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "tr":
			err = trTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "vi":
			err = viTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "zh_tw":
			err = zhtwTranslations.RegisterDefaultTranslations(trans.vld, tran)
		default:
			err = enTranslations.RegisterDefaultTranslations(trans.vld, tran)
		}
		return err
	}
}
func (t *Trans) GetError(err error) error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {
		TransErrs := errs.Translate(t.tran)
		for _, v := range TransErrs {
			return errors.New(v)
		}
	}
	return err
}

func (t *Trans) GetErrors(err error) []error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {
		TransErrs := errs.Translate(trans)
		resErrs := make([]error, 0)
		for _, v := range TransErrs {
			resErrs = append(resErrs, errors.New(v))
		}
		return resErrs
	}
	return []error{err}
}

func GetError(err error) error {
	return std.GetError(err)
}

func GetErrors(err error) []error {
	return std.GetErrors(err)
}

//func TransInit(v *validator.Validate, local string) (err error) {
//	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//	zhT := zh.New()
//	enT := en.New()
//	uni := ut.New(enT, zhT)
//	tran, ok := uni.GetTranslator(local)
//	trans = tran
//	if !ok {
//		return fmt.Errorf("uni.GetTranslator(%s) failed", local)
//	}
//	switch local {
//	case "en":
//		err = enTranslations.RegisterDefaultTranslations(v, trans)
//	case "zh":
//		err = zhTranslations.RegisterDefaultTranslations(v, trans)
//	default:
//		err = enTranslations.RegisterDefaultTranslations(v, trans)
//	}
//	return
//}
//
//func extendZhTrans(v *validator.Validate, trans ut.Translator) (err error) {
//	translations := []struct {
//		tag             string
//		translation     string
//		override        bool
//		customRegisFunc validator.RegisterTranslationsFunc
//		customTransFunc validator.TranslationFunc
//	}{
//		{
//			tag:         "required_if",
//			translation: "{0}为必填字段",
//			override:    false,
//			//customRegisFunc: func(ut ut.Translator) error {
//			//	return ut.Add("required_if", "{0}为必填字段", false)
//			//},
//			//customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
//			//	t, _ := ut.T("required_if", fe.Field())
//			//	return t
//			//},
//		},
//	}
//
//	for _, t := range translations {
//		if t.customTransFunc != nil && t.customRegisFunc != nil {
//			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)
//		} else if t.customTransFunc != nil && t.customRegisFunc == nil {
//			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)
//		} else if t.customTransFunc == nil && t.customRegisFunc != nil {
//			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)
//		} else {
//			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
//		}
//		if err != nil {
//			return
//		}
//	}
//	return
//}
//
//func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
//	return func(ut ut.Translator) (err error) {
//		if err = ut.Add(tag, translation, override); err != nil {
//			return
//		}
//		return
//	}
//}
//
//func translateFunc(ut ut.Translator, fe validator.FieldError) string {
//	t, err := ut.T(fe.Tag(), fe.Field())
//	if err != nil {
//		log.Printf("警告: 翻译字段错误: %#v", fe)
//		return fe.(error).Error()
//	}
//	return t
//}
