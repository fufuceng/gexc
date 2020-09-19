package time

import (
	"reflect"
	"testing"
	"time"
)

func TestExTime_MarshalJSON(t *testing.T) {
	type fields struct {
		Time time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "should marshal time object according to exTime layout",
			fields:  fields{Time: time.Date(2020, 12, 29, 12, 10, 10, 0, time.UTC)},
			want:    []byte("2020-12-29"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			et := &Gexc{
				Time: tt.fields.Time,
			}

			got, err := et.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestExTime_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Gexc
		wantErr bool
	}{
		{
			name:    "should fill time object correctly",
			fields:  fields{Time: time.Time{}},
			args:    args{data: []byte(`"2020-12-29"`)},
			want:    Gexc{time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name:    "should fill as null if time is null",
			fields:  fields{Time: time.Time{}},
			args:    args{data: []byte("null")},
			want:    Gexc{time.Time{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			et := &Gexc{
				Time: tt.fields.Time,
			}

			if err := et.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !et.Time.Equal(tt.want.Time) {
				t.Errorf("MarshalJSON() got = %v, want %v", et.Time, tt.want.Time)
			}
		})
	}
}

func TestNewExTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want Gexc
	}{
		{
			name: "should create exTime object with correct time value",
			args: args{t: time.Date(2020, 12, 29, 10, 10, 10, 0, time.UTC)},
			want: Gexc{Time: time.Date(2020, 12, 29, 10, 10, 10, 0, time.UTC)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGexc(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGexc() = %v, want %v", got, tt.want)
			}
		})
	}
}
