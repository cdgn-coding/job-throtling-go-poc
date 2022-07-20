package csv

import (
	"github.com/stretchr/testify/assert"
	"job-throtling-go-poc/pkg/models"
	"strings"
	"testing"
)

func TestNewReaderChannel(t *testing.T) {
	assert.NotPanics(t, func() {
		NewReaderChannel[any]()
	})
}

func TestReaderChannel_Read(t *testing.T) {
	t.Run("Reads multiple lines to a channel", func(t *testing.T) {
		reader := NewReaderChannel[models.Address]()
		records := make(chan models.Address, 100)

		fileReader := strings.NewReader(`John,Doe,120 jefferson st.,Riverside, NJ, 08075
Jack,McGinnis,220 hobo Av.,Phila, PA,09119
"John ""Da Man""",Repici,120 Jefferson St.,Riverside, NJ,08075
Stephen,Tyler,"7452 Terrace ""At the Plaza"" road",SomeTown,SD, 91234
,Blankman,,SomeTown, SD, 00298
"Joan ""the bone"", Anne",Jet,"9th, at Terrace plc",Desert City,CO,00123
`)

		reader.Read(fileReader, records)
		names := []string{"John", "Jack", `John "Da Man"`, "Stephen", "", `Joan "the bone", Anne`}
		for i := 0; i < 6; i++ {
			record := <-records
			assert.Equal(t, names[i], record.Name)
		}
	})
}
