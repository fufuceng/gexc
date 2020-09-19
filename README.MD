# Gexc

Gexc is an unofficial `go` client for free [Exchange Rates Api](https://exchangeratesapi.io/) 
which is published by European Central Bank.    

## Installation

```bash
go get github.com/fufuceng/gexc
```

## Usage

### Conversion - Long Version
```go
converted, err := gexc.New().Amount(100).From("EUR").To("TRY")
if err != nil {
    log.Fatal(err)
}

fmt.Println(converted) // -> 896₺
```

### Conversion - Short Version
```go
converted, err := gexc.New().Convert(100, "EUR", "TRY")
if err != nil {
    log.Fatal(err)
}

fmt.Println(converted) // -> 896₺
```

### Latest

```go
latest, err := gexc.New().HistoryOf("EUR").Against("TRY", "USD").Latest()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Latest values of %v\n", latest.Base)
for curr, val := range latest.Rates {
    fmt.Printf("\t * Currency: %v, Value: %v\n", curr, val)
}

/*
Latest values of EUR
    * Currency: TRY, Value: 8.96
    * Currency: USD, Value: 1.1833
*/
```

### History
```go
from := time.Now().AddDate(0, 0, -1)
until := time.Now()

history, err := gexc.New().HistoryOf("EUR").Against("TRY", "USD").From(from).Until(until)
if err != nil {
    log.Fatalf(err)
}

fmt.Printf("History of %v: from: %v, until: %v\n", history.Base, history.StartAt, history.EndAt)

for date, item := range history.Rates {
	fmt.Printf("Date: %v\n", date)
	for curr, val := range item {
		fmt.Printf("\t * Currency: %v, Value: %v\n", curr, val)
	}
}

/*
History of EUR: from: 2020-09-17, until: 2020-09-18
Date: 2020-09-18
     * Currency: TRY, Value: 8.96
     * Currency: USD, Value: 1.1833
Date: 2020-09-17
     * Currency: TRY, Value: 8.9037
     * Currency: USD, Value: 1.1797
*/

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)