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
	"github.com/go-playground/locales/pt_BR"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/tr"
	"github.com/go-playground/locales/vi"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
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
		fmt.Println("failed to init validator")
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
	uni := ut.New(en.New(), ar.New(), en.New(), es.New(), fa.New(), fr.New(), id.New(), it.New(), ja.New(), lv.New(), nl.New(), pt.New(), pt_BR.New(), ru.New(), tr.New(), vi.New(), zh.New(), zh_Hant_TW.New())
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
		case "ar": //阿拉伯语
			err = arTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "en":
			err = enTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "es": //西班牙语
			err = esTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "fa": //波斯语
			err = faTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "fr": //法语
			err = frTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "id": //印尼语
			err = idTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "it": //意大利语
			err = itTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "ja": //日语
			err = jaTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "lv": //拉脱维亚语
			err = lvTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "nl": //荷兰语
			err = nlTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "pt": //葡萄牙语
			err = ptTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "pt_BR": //巴西葡萄牙语
			err = ptbrTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "ru": //俄语
			err = ruTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "tr": //土耳其语
			err = trTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "vi": //越南语
			err = viTandslations.RegisterDefaultTranslations(trans.vld, tran)
		case "zh": //中文
			err = zhTranslations.RegisterDefaultTranslations(trans.vld, tran)
		case "zh_hant_tw": //繁体中文台湾
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
		resErrs := make([]error, 0, len(TransErrs))
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
