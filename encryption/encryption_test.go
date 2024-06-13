package encryption

import "testing"

func TestEasyCrypt_Encrypt(t *testing.T) {
	type fields struct {
		Type string
		Key  string
		Iv   string
	}
	type args struct {
		plaintext string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "test1", fields: fields{
			Type: "aes/ECB/PKCS7/Hex",
			Key:  "ploknht78guqwefh",
			Iv:   "09ji-9uygh4es6ga",
		},
			args:    args{plaintext: "test1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EasyCrypt{
				Type: tt.fields.Type,
				Key:  tt.fields.Key,
				Iv:   tt.fields.Iv,
			}
			got, err := e.Encrypt(tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			orig, err := e.Decrypt(got)
			if err != nil {
				t.Error(err)
				return
			}
			if orig != tt.args.plaintext {
				t.Errorf("Decrypt() got = %s, want %s", orig, tt.args.plaintext)
				return
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		model      string
		plaintext  string
		cipherText string
		key        string
		vi         string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test1-1", args: args{model: "aes/ECB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test1-2", args: args{model: "aes/CBC/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test1-3", args: args{model: "aes/CFB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test1-4", args: args{model: "aes/OFB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test1-5", args: args{model: "aes/CTR/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},

		{name: "test2-1", args: args{model: "des/ECB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test2-2", args: args{model: "des/CBC/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test2-3", args: args{model: "des/CFB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test2-4", args: args{model: "des/OFB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test2-5", args: args{model: "des/CTR/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},

		{name: "test3-1", args: args{model: "3des/ECB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test3-2", args: args{model: "3des/CBC/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test3-3", args: args{model: "3des/CFB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test3-4", args: args{model: "3des/OFB/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
		{name: "test3-5", args: args{model: "3des/CTR/PKCS7/Hex", cipherText: "q 我是一个imya文89￥%……&[a{}()", key: "3q9wcubuc-ju4kho", vi: "yhr1abkil26sik97"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipher := Encrypt(tt.args.model, tt.args.plaintext, tt.args.key, tt.args.vi)
			got := Decrypt(tt.args.model, cipher, tt.args.key, tt.args.vi)
			if got != tt.args.plaintext {
				t.Errorf("Decrypt() got = %s, want %s", got, tt.args.plaintext)
			}
		})
	}
}
