package openex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_client_toUrl(t *testing.T) {
	type fields struct {
		config Config
	}
	type args struct {
		path string
		qp   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   url.URL
		want1  string
	}{
		{
			name:   "should build correct url object without query parameter",
			fields: fields{config: defaultConfig},
			args:   args{path: "path"},
			want: url.URL{
				Scheme: defaultConfig.Protocol,
				Host:   defaultConfig.BaseUrl,
				Path:   "path",
			},
			want1: fmt.Sprintf("%s://%s/%s", defaultConfig.Protocol, defaultConfig.BaseUrl, "path"),
		},
		{
			name:   "should build correct url object with query parameter",
			fields: fields{config: defaultConfig},
			args:   args{path: "path", qp: "qp1=true&qp2=false"},
			want: url.URL{
				Scheme:   defaultConfig.Protocol,
				Host:     defaultConfig.BaseUrl,
				Path:     "path",
				RawQuery: "qp1=true&qp2=false",
			},
			want1: fmt.Sprintf("%s://%s/%s?%s", defaultConfig.Protocol, defaultConfig.BaseUrl, "path", "qp1=true&qp2=false"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client{
				config: tt.fields.config,
			}

			got := c.toUrl(tt.args.path, tt.args.qp)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toUrl() = %v, want %v", got, tt.want)
			}

			if got.String() != tt.want1 {
				t.Errorf("toUrl() - String = %v, want %v", got.String(), tt.want1)
			}
		})
	}
}

func Test_client_doRequest(t *testing.T) {
	type fields struct {
		config     Config
		httpGetter httpGetter
	}
	type args struct {
		method string
		url    url.URL
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *response
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "should return correct response and code if method is GET",
			fields: fields{
				config: defaultConfig,
				httpGetter: func(url string) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(bytes.NewBufferString("Hello World")),
						StatusCode: 200,
					}, nil
				},
			},
			args: args{method: "GET", url: url.URL{}},
			want: &response{
				Body: []byte("Hello World"),
				Code: 200,
			},
		},
		{
			name:      "should raise an panic message if method is not GET",
			fields:    fields{config: defaultConfig},
			args:      args{method: "POST", url: url.URL{}},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantPanic {
						t.Errorf("doRequest() panic = %v, wantErr %v", r, tt.wantPanic)
					}
				}
			}()

			c := client{
				config:     tt.fields.config,
				httpGetter: tt.fields.httpGetter,
			}
			got, err := c.doRequest(tt.args.method, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
