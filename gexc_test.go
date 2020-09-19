package gexc

import (
	"github.com/fufuceng/gexc/internal/openex"
	"github.com/fufuceng/gexc/response"
	time2 "github.com/fufuceng/gexc/time"
	"github.com/fufuceng/gexc/types"
	"reflect"
	"testing"
	"time"
)

type testClient struct{}

func (t testClient) Latest(params openex.LatestParams) (*response.SingleDate, error) {
	return &response.SingleDate{
		Rates: types.RateItem{
			"EUR": 8.0,
			"GBP": 9.0,
			"USD": 7.0,
		},
		Base: "TRY",
		Date: time2.Gexc{
			Time: time.Date(2020, 12, 29, 12, 30, 0, 0, time.UTC),
		},
	}, nil
}

func (t testClient) SingleDate(params openex.SingleDateParams) (*response.SingleDate, error) {
	return &response.SingleDate{
		Rates: types.RateItem{
			"EUR": 8.0,
			"GBP": 9.0,
			"USD": 7.0,
		},
		Base: "TRY",
		Date: time2.Gexc{
			Time: time.Date(2020, 11, 29, 12, 30, 0, 0, time.UTC),
		},
	}, nil
}

func (t testClient) History(params openex.HistoryParams) (*response.History, error) {
	return &response.History{
		Rates: types.TimeRateItem{
			"2020-12-29": types.RateItem{
				"EUR": 8.0,
				"GBP": 9.0,
				"USD": 7.0,
			},

			"2020-12-28": types.RateItem{
				"EUR": 7.0,
				"GBP": 8.0,
				"USD": 6.0,
			},

			"2020-12-27": types.RateItem{
				"EUR": 6.0,
				"GBP": 7.0,
				"USD": 5.0,
			},
		},

		Base:    "TRY",
		StartAt: time2.NewGexc(time.Date(2020, 12, 27, 0, 0, 0, 0, time.UTC)),
		EndAt:   time2.NewGexc(time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC)),
	}, nil
}

func TestFx_HistoryOfFullChain(t *testing.T) {
	type fields struct {
		openexClient openex.Client
	}

	type args struct {
		baseCurrency string
		against      []string
		From         time.Time
		Until        time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    response.History
		wantErr bool
	}{
		{
			name: "should return correct history for TRY with given valid date range",
			fields: fields{
				openexClient: testClient{},
			},

			args: args{
				baseCurrency: "TRY",
				against:      []string{"eur", "usd", "gbp"},
				From:         time.Date(2020, 12, 27, 0, 0, 0, 0, time.UTC),
				Until:        time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
			},

			want: response.History{
				Base: "TRY",
				StartAt: time2.Gexc{
					Time: time.Date(2020, 12, 27, 0, 0, 0, 0, time.UTC),
				},
				EndAt: time2.Gexc{
					Time: time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
				},
				Rates: types.TimeRateItem{

					"2020-12-29": types.RateItem{
						"EUR": 8.0,
						"GBP": 9.0,
						"USD": 7.0,
					},

					"2020-12-28": types.RateItem{
						"EUR": 7.0,
						"GBP": 8.0,
						"USD": 6.0,
					},

					"2020-12-27": types.RateItem{
						"EUR": 6.0,
						"GBP": 7.0,
						"USD": 5.0,
					},
				},
			},
		},

		{
			name: "should raise an error if from field is empty",
			fields: fields{
				openexClient: testClient{},
			},

			args: args{
				baseCurrency: "TRY",
				against:      []string{"eur", "usd", "gbp"},
				From:         time.Time{},
				Until:        time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
			want:    response.History{},
		},

		{
			name: "should raise an error if to field is empty",
			fields: fields{
				openexClient: testClient{},
			},

			args: args{
				baseCurrency: "TRY",
				against:      []string{"eur", "usd", "gbp"},
				From:         time.Time{},
				Until:        time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
			want:    response.History{},
		},

		{
			name: "should raise an error if from value equal to until",
			fields: fields{
				openexClient: testClient{},
			},

			args: args{
				baseCurrency: "TRY",
				against:      []string{"eur", "usd", "gbp"},
				From:         time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
				Until:        time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
			want:    response.History{},
		},

		{
			name: "should raise an error if from value bigger than until",
			fields: fields{
				openexClient: testClient{},
			},

			args: args{
				baseCurrency: "TRY",
				against:      []string{"eur", "usd", "gbp"},
				From:         time.Date(2020, 12, 30, 0, 0, 0, 0, time.UTC),
				Until:        time.Date(2020, 12, 29, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
			want:    response.History{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fx{
				openexClient: tt.fields.openexClient,
			}
			got, err := f.
				HistoryOf(tt.args.baseCurrency).
				Against(tt.args.against...).
				From(tt.args.From).
				Until(tt.args.Until)

			if (err != nil) != tt.wantErr {
				t.Errorf("TestFx_HistoryOfFullChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestFx_HistoryOfFullChain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFx_AmountChain(t *testing.T) {
	type fields struct {
		openexClient openex.Client
	}

	type args struct {
		amount float64
		from   string
		to     string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "should convert 5 TRY to EUR successfully",
			fields: fields{
				openexClient: testClient{},
			},
			args: args{
				amount: 5,
				from:   "TRY",
				to:     "EUR",
			},
			want:    40.0,
			wantErr: false,
		},
		{
			name: "should raise an error if unknown currency exist",
			fields: fields{
				openexClient: testClient{},
			},
			args: args{
				amount: 5,
				from:   "UNKNOWN",
				to:     "EUR",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fx{
				openexClient: tt.fields.openexClient,
			}
			got, err := f.
				Amount(tt.args.amount).
				From(tt.args.from).
				To(tt.args.to)

			if (err != nil) != tt.wantErr {
				t.Errorf("TestFx_AmountChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TestFx_AmountChain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFx_Convert(t *testing.T) {
	type fields struct {
		openexClient openex.Client
	}

	type args struct {
		amount float64
		from   string
		to     string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "should convert 5 TRY to EUR successfully",
			fields: fields{
				openexClient: testClient{},
			},
			args: args{
				amount: 5,
				from:   "TRY",
				to:     "EUR",
			},
			want:    40.0,
			wantErr: false,
		},
		{
			name: "should raise an error if unknown currency exist",
			fields: fields{
				openexClient: testClient{},
			},
			args: args{
				amount: 5,
				from:   "UNKNOWN",
				to:     "EUR",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fx{
				openexClient: tt.fields.openexClient,
			}
			got, err := f.Convert(tt.args.amount, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
