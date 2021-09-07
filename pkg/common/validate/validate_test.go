package validate

import "testing"

func TestValidateIP(t *testing.T) {
	type args struct {
		ip []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "", args: args{ip: []string{"192.168.1.2"}}, wantErr: false},
		{name: "", args: args{ip: []string{"a.b.c.d"}}, wantErr: true},
		{name: "", args: args{ip: []string{"192.168.1.256"}}, wantErr: true},
		{name: "", args: args{ip: []string{"192.168.$.a"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateIP(tt.args.ip...); (err != nil) != tt.wantErr {
				t.Errorf("ValidateIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	type args struct {
		urlSlice []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "", args: args{urlSlice: []string{"aabbcc"}}, wantErr: true},
		{name: "", args: args{urlSlice: []string{"http://abcde"}}, wantErr: true},
		{name: "", args: args{urlSlice: []string{"https://"}}, wantErr: true},
		{name: "", args: args{urlSlice: []string{"https://aa.com"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateURL(tt.args.urlSlice...); (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
