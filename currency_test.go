package gexc

import (
	"reflect"
	"testing"
)

func Test_sanitizeCurrencyCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should handle lowercase characters",
			args: args{code: "eur"},
			want: "EUR",
		},
		{
			name: "should handle space characters",
			args: args{code: " eur "},
			want: "EUR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeCurrencyCode(tt.args.code); got != tt.want {
				t.Errorf("sanitizeCurrencyCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sanitizeCurrencyName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should handle upper case characters",
			args: args{name: "turkish LIRA"},
			want: "turkish lira",
		},
		{
			name: "should handle spaces",
			args: args{name: " turkish LIRA "},
			want: "turkish lira",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeCurrencyName(tt.args.name); got != tt.want {
				t.Errorf("sanitizeCurrencyName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyByCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name  string
		args  args
		want  Currency
		want1 bool
	}{
		{
			name:  "should return true and corresponding currency object if code exist in map",
			args:  args{code: "EUR"},
			want:  Currency{Code: "EUR", Name: "euro"},
			want1: true,
		},
		{
			name:  "should return false and empty currency object if code not exist in map",
			args:  args{code: "NOT_EXIST"},
			want:  Currency{},
			want1: false,
		},
		{
			name:  "should sanitize code before searching",
			args:  args{code: " eur "},
			want:  Currency{Code: "EUR", Name: "euro"},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CurrencyByCode(tt.args.code)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CurrencyByCode() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CurrencyByCode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCurrencyByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  Currency
		want1 bool
	}{
		{
			name:  "should return true and corresponding currency object if name exist in map",
			args:  args{name: "euro"},
			want:  Currency{Code: "EUR", Name: "euro"},
			want1: true,
		},
		{
			name:  "should return false and empty currency object if name not exist in map",
			args:  args{name: "NOT_EXIST"},
			want:  Currency{},
			want1: false,
		},
		{
			name:  "should sanitize name before searching",
			args:  args{name: " euro "},
			want:  Currency{Code: "EUR", Name: "euro"},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CurrencyByName(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CurrencyByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CurrencyByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_currencyByX(t *testing.T) {
	type args struct {
		source *map[string]*Currency
		key    string
	}
	tests := []struct {
		name  string
		args  args
		want  Currency
		want1 bool
	}{
		{
			name: "should return currency object and true if given key exist in given resource",
			args: args{
				source: &map[string]*Currency{
					"EUR": {Code: "EUR", Name: "euro"},
				},
				key: "EUR",
			},
			want:  Currency{Code: "EUR", Name: "euro"},
			want1: true,
		},

		{
			name: "should return empty currency object and false if given key not exist in given resource",
			args: args{
				source: &map[string]*Currency{
					"EUR": {Code: "EUR", Name: "euro"},
				},
				key: "TRY",
			},
			want:  Currency{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := currencyByX(tt.args.source, tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("currencyByX() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("currencyByX() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
