package validator

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
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

type Lang string // support language
const (
	Ar         Lang = "ar"
	En         Lang = "en"
	Es         Lang = "es"
	Fa         Lang = "fa"
	Fr         Lang = "fr"
	Id         Lang = "id"
	It         Lang = "it"
	Ja         Lang = "ja"
	Lv         Lang = "lv"
	Nl         Lang = "nl"
	Pt         Lang = "pt"
	Pt_BR      Lang = "pt_BR"
	Ru         Lang = "ru"
	Tr         Lang = "tr"
	Vi         Lang = "vi"
	Zh         Lang = "zh"
	Zh_Hant_TW Lang = "zh_hant_tw"
)

var sourceLang = map[Lang]locales.Translator{
	Ar:         ar.New(),
	En:         en.New(),
	Es:         es.New(),
	Fa:         fa.New(),
	Fr:         fr.New(),
	Id:         id.New(),
	It:         it.New(),
	Ja:         ja.New(),
	Lv:         lv.New(),
	Nl:         nl.New(),
	Pt:         pt.New(),
	Pt_BR:      pt_BR.New(),
	Ru:         ru.New(),
	Tr:         tr.New(),
	Vi:         vi.New(),
	Zh:         zh.New(),
	Zh_Hant_TW: zh_Hant_TW.New(),
}
var targetTrans = map[Lang]func(*validator.Validate, ut.Translator) error{
	Ar:         arTranslations.RegisterDefaultTranslations,
	En:         enTranslations.RegisterDefaultTranslations,
	Es:         esTranslations.RegisterDefaultTranslations,
	Fa:         faTandslations.RegisterDefaultTranslations,
	Fr:         frTandslations.RegisterDefaultTranslations,
	Id:         idTandslations.RegisterDefaultTranslations,
	It:         itTandslations.RegisterDefaultTranslations,
	Ja:         jaTandslations.RegisterDefaultTranslations,
	Lv:         lvTandslations.RegisterDefaultTranslations,
	Nl:         nlTandslations.RegisterDefaultTranslations,
	Pt:         ptTandslations.RegisterDefaultTranslations,
	Pt_BR:      ptbrTandslations.RegisterDefaultTranslations,
	Ru:         ruTandslations.RegisterDefaultTranslations,
	Tr:         trTandslations.RegisterDefaultTranslations,
	Vi:         viTandslations.RegisterDefaultTranslations,
	Zh:         zhTranslations.RegisterDefaultTranslations,
	Zh_Hant_TW: zhtwTranslations.RegisterDefaultTranslations,
}
var trans ut.Translator
var std = newGinZh()

type Option func(trans *Trans) error

func newGinZh() *Trans {
	t, _ := newGinVld("zh")
	return t
}

func SwitchGinVldLang(local Lang) {
	t, _ := newGinVld(local)
	std = t
}

func newGinVld(local Lang) (*Trans, error) {
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
	uni := ut.New(en.New())
	for _, v := range sourceLang {
		err := uni.AddTranslator(v, true)
		if err != nil {
			return nil, err
		}
	}
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

func WithLocal(local Lang) Option {
	return func(trans *Trans) error {
		tran, ok := trans.uniTran.GetTranslator(string(local))
		if !ok {
			return errors.New(fmt.Sprintf("cannot be translated as %s", local))
		}
		trans.tran = tran
		if f, ok := targetTrans[local]; ok {
			return f(trans.Vld, tran)
		}
		return enTranslations.RegisterDefaultTranslations(trans.Vld, tran)
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
