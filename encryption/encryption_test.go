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
