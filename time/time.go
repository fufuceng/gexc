package time

import (
	"errors"
	"time"
)

const GexcLayout = "2006-01-02"

//Gexc represents the time value that will be used in gexc library
//to communicate with the api
type Gexc struct {
	time.Time
}

func (et *Gexc) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var err error
	et.Time, err = time.Parse(`"`+GexcLayout+`"`, string(data))
	return err
}

func (et *Gexc) MarshalJSON() ([]byte, error) {
	if y := et.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Gexc.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(GexcLayout))
	b = et.AppendFormat(b, GexcLayout)

	return b, nil
}

func (et Gexc) String() string {
	return et.Time.Format(GexcLayout)
}

func NewGexc(t time.Time) Gexc {
	return Gexc{t}
}
