package validator

import (
	"errors"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	arTranslations "github.com/go-playground/validator/v10/translations/ar"
	"golang.org/x/exp/slices"
)

type testUser struct {
	FirstName      string         `validate:"required"`
	LastName       string         `validate:"required"`
	Age            uint8          `validate:"gte=0,lte=130"`
	Email          string         `validate:"required,email"`
	FavouriteColor string         `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*testAddress `validate:"required,dive,required"` // a person can have a home and cottage...
}
type testAddress struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

type testAlias struct {
	FirstName string `validate:"required" alias:"姓"`
	LastName  string `binding:"required" alias:"名"` // for gin
}

var tu = testUser{
	FirstName:      "Badger",
	LastName:       "Smith",
	Age:            45,
	Email:          "@a.com",
	FavouriteColor: "blue",
	Addresses: []*testAddress{
		{
			Street: "Eavesdown Docks",
			City:   "Osiris",
			Planet: "Universe",
			Phone:  "none",
		},
	},
}
var ta = testAlias{
	FirstName: "",
}

func Test_stdErrs(t *testing.T) {
	vld := validator.New()
	err := vld.Struct(tu)
	if err == nil {
		t.Error("want and error but got nil")
	}

	getErrs := GetErrors(err)
	if len(getErrs) != 2 {
		t.Error("want 2, but got zero")
		return
	}

	//set std to nil
	std = nil
	getErrs = GetErrors(err)
	if len(getErrs) != 1 {
		t.Errorf("want 1, but %d", len(getErrs))
	}
	if !errors.Is(getErrs[0], ErrValidatorNotInit) {
		t.Errorf("want %s, got %s", ErrValidatorNotInit.Error(), getErrs[0].Error())
	}
}

func Test_stdErr(t *testing.T) {
	vld := validator.New()
	err := vld.Struct(tu)
	if err == nil {
		t.Error("want and error but got nil")
	}

	getErr := GetError(err)
	if getErr == nil {
		t.Error("want error, got nil")
		return
	}

	//set std to nil
	std = nil
	getErr = GetError(err)
	if !errors.Is(getErr, ErrValidatorNotInit) {
		t.Errorf("want %s, got %s", ErrValidatorNotInit.Error(), getErr.Error())
	}
}

func Test_ChangeGinVldLang(t *testing.T) {
	SwitchGinVldLang("zh_hant_tw")
	//err := validator.New().Struct(ta)
	err := std.Vld.Struct(ta)
	vldErr := GetError(err)
	if vldErr == nil {
		t.Error("want error, got nil")
		return
	}
	if vldErr.Error() != "名為必填欄位" {
		t.Errorf("want 名為必填欄位, got %s", vldErr.Error())
	}
}

func Test_Alias(t *testing.T) {
	vld := validator.New()
	tr, err := New(vld, WithLocal("zh"), WithAliasTag())
	if err != nil {
		t.Error(err)
		return
	}
	vliErr := tr.Vld.Struct(ta)
	if vliErr == nil {
		t.Error("want and error but got nil")
	}
	getErr := tr.GetError(vliErr)
	if getErr == nil {
		t.Error("want error, got nil")
		return
	}
	if getErr.Error() != "姓为必填字段" {
		t.Errorf("want 姓为必填字段, got %s", getErr.Error())
	}
}

func Test_UnsupportedLocal(t *testing.T) {
	//force fewer translation language packs to test unsupported languages
	oldTt := targetTrans
	targetTrans = map[Lang]func(*validator.Validate, ut.Translator) error{
		Ar: arTranslations.RegisterDefaultTranslations,
	}
	defer func() {
		targetTrans = oldTt
	}()
	_, err := New(validator.New(), WithLocal("zh"))
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestMultilingual(t *testing.T) {
	type args struct {
		vld  *validator.Validate
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Example 1", args: args{validator.New(), []Option{WithLocal("ar")}}, want: []string{"يجب أن يكون Email عنوان بريد إلكتروني صالح", "يجب أن يكون FavouriteColor لون صالح"}},
		{name: "Example 2", args: args{validator.New(), []Option{WithLocal("en")}}, want: []string{"Email must be a valid email address", "FavouriteColor must be a valid color"}},
		{name: "Example 3", args: args{validator.New(), []Option{WithLocal("es")}}, want: []string{"Email debe ser una dirección de correo electrónico válida", "FavouriteColor debe ser un color válido"}},
		{name: "Example 4", args: args{validator.New(), []Option{WithLocal("fa")}}, want: []string{"Email باید یک ایمیل معتبر باشد", "FavouriteColor باید یک رنگ معتبر باشد"}},
		{name: "Example 5", args: args{validator.New(), []Option{WithLocal("fr")}}, want: []string{"Email doit être une adresse email valide", "FavouriteColor doit être une couleur valide"}},
		{name: "Example 6", args: args{validator.New(), []Option{WithLocal("id")}}, want: []string{"Email harus berupa alamat email yang valid", "FavouriteColor harus berupa warna yang valid"}},
		{name: "Example 7", args: args{validator.New(), []Option{WithLocal("it")}}, want: []string{"Email deve essere un indirizzo email valido", "FavouriteColor deve essere un colore valido"}},
		{name: "Example 8", args: args{validator.New(), []Option{WithLocal("ja")}}, want: []string{"Emailは正しいメールアドレスでなければなりません", "FavouriteColorは正しい色でなければなりません"}},
		{name: "Example 9", args: args{validator.New(), []Option{WithLocal("lv")}}, want: []string{"Email jābūt derīgai e-pasta adresei", "FavouriteColor jābūt derīgai krāsai"}},
		{name: "Example 10", args: args{validator.New(), []Option{WithLocal("nl")}}, want: []string{"Email moet een geldig email adres zijn", "FavouriteColor moet een geldige kleur zijn"}},
		{name: "Example 11", args: args{validator.New(), []Option{WithLocal("pt")}}, want: []string{"Email deve ser um endereço de e-mail válido", "FavouriteColor deve ser uma cor válida"}},
		{name: "Example 12", args: args{validator.New(), []Option{WithLocal("pt_BR")}}, want: []string{"Email deve ser um endereço de e-mail válido", "FavouriteColor deve ser uma cor válida"}},
		{name: "Example 13", args: args{validator.New(), []Option{WithLocal("ru")}}, want: []string{"Email должен быть email адресом", "FavouriteColor должен быть цветом"}},
		{name: "Example 14", args: args{validator.New(), []Option{WithLocal("tr")}}, want: []string{"Email geçerli bir e-posta adresi olmalıdır", "FavouriteColor geçerli bir renk olmalıdır"}},
		{name: "Example 15", args: args{validator.New(), []Option{WithLocal("vi")}}, want: []string{"Email phải là giá trị email address", "FavouriteColor phải là màu sắc hợp lệ"}},
		{name: "Example 16", args: args{validator.New(), []Option{WithLocal("zh")}}, want: []string{"Email必须是一个有效的邮箱", "FavouriteColor必须是一个有效的颜色"}},
		{name: "Example 17", args: args{validator.New(), []Option{WithLocal("zh_hant_tw")}}, want: []string{"Email必須是一個有效的信箱", "FavouriteColor必須是一個有效的顏色"}},
		{name: "Example 18", args: args{validator.New(), []Option{WithLocal("notExist")}}, want: []string{"cannot be translated as notExist"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := New(tt.args.vld, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, nil)
				return
			}
			if err != nil {
				return
			}
			getErr := tr.GetError(tt.args.vld.Struct(tu))
			if getErr == nil {
				t.Error("want error, got nil")
			} else {
				if !slices.Contains(tt.want, getErr.Error()) {
					t.Errorf("want %s, got %s", tt.want, getErr.Error())
				}
			}
		})
	}
}
