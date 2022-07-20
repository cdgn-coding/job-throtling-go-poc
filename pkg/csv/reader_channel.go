package csv

import (
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"io"
)

type ReaderChannel[T any] struct{}

func NewReaderChannel[T any]() *ReaderChannel[T] {
	return &ReaderChannel[T]{}
}

func (p *ReaderChannel[T]) Read(fileReader io.Reader, outputs chan<- T) {
	csvReader := csv.NewReader(fileReader)
	header, _ := csvutil.Header(new(T), "csv")
	decoder, _ := csvutil.NewDecoder(csvReader, header...)

	for {
		var record T
		err := decoder.Decode(&record)
		if err == nil {
			outputs <- record
		} else {
			break
		}
	}

	close(outputs)
}
