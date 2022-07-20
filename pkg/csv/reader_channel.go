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

func (p *ReaderChannel[T]) Read(fileReader io.Reader, outputs chan<- T, count chan<- int) {
	csvReader := csv.NewReader(fileReader)
	header, _ := csvutil.Header(new(T), "csv")
	decoder, _ := csvutil.NewDecoder(csvReader, header...)

	counter := 0
	for {
		var record T
		err := decoder.Decode(&record)
		if err == nil {
			outputs <- record
			counter++
		} else if err == io.EOF {
			break
		}
	}
	count <- counter

	close(count)
	close(outputs)
}
