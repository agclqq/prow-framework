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
	Vld     *validator.Validate
	uniTran *ut.UniversalTranslator
	tran    ut.Translator
}

var (
	ErrValidatorNotInit    = errors.New("validator is not init")
	ErrValidatorFailedInit = errors.New("failed to init validator")
)

var trans ut.Translator
var std = newGinZh()

type Option func(trans *Trans) error

func newGinZh() *Trans {
	t, _ := newGinVld("zh")
	return t
}

func SwitchGinVldLang(local string) {
	t, _ := newGinVld(local)
	std = t
}

func newGinVld(local string) (*Trans, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, ErrValidatorFailedInit
	}
	t, err := New(v, WithLocal(local), WithAliasTag())
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return t, err
}

func New(vld *validator.Validate, opts ...Option) (*Trans, error) {
	uni := ut.New(en.New(), ar.New(), en.New(), es.New(), fa.New(), fr.New(), id.New(), it.New(), ja.New(), lv.New(), nl.New(), pt.New(), pt_BR.New(), ru.New(), tr.New(), vi.New(), zh.New(), zh_Hant_TW.New())
	tran := &Trans{Vld: vld, uniTran: uni}
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
		RegisterTagName(trans.Vld, "alias")
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
			err = arTranslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "en":
			err = enTranslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "es": //西班牙语
			err = esTranslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "fa": //波斯语
			err = faTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "fr": //法语
			err = frTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "id": //印尼语
			err = idTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "it": //意大利语
			err = itTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "ja": //日语
			err = jaTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "lv": //拉脱维亚语
			err = lvTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "nl": //荷兰语
			err = nlTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "pt": //葡萄牙语
			err = ptTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "pt_BR": //巴西葡萄牙语
			err = ptbrTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "ru": //俄语
			err = ruTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "tr": //土耳其语
			err = trTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "vi": //越南语
			err = viTandslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "zh": //中文
			err = zhTranslations.RegisterDefaultTranslations(trans.Vld, tran)
		case "zh_hant_tw": //繁体中文台湾
			err = zhtwTranslations.RegisterDefaultTranslations(trans.Vld, tran)
		default:
			err = enTranslations.RegisterDefaultTranslations(trans.Vld, tran)
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
	if std == nil {
		return ErrValidatorNotInit
	}
	return std.GetError(err)
}

func GetErrors(err error) []error {
	if std == nil {
		return []error{ErrValidatorNotInit}
	}
	return std.GetErrors(err)
}
