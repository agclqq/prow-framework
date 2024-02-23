package validator

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func Test_trans(t *testing.T) {
	user := User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            45,
		Email:          "@a.com",
		FavouriteColor: "blue",
		Addresses: []*Address{
			&Address{
				Street: "Eavesdown Docks",
				City:   "Osiris",
				Planet: "Universe",
				Phone:  "none",
			},
		},
	}

	vld := validator.New()
	tr, err := New(vld, WithLocal("zh"), WithAliasTag())
	if err != nil {
		t.Error(err)
		return
	}
	err = vld.Struct(user)
	if err != nil {
		fmt.Println(tr.GetError(err))
		return
	} else {
		t.Error("want and error but got nil")
	}
}
