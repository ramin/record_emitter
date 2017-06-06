package record_emitter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Record struct {
	headerMap map[string]int
	Row       []string
}

func (record *Record) Get(field string) string {
	if index, present := record.headerMap[field]; present {
		if len(record.Row) >= index {
			return strings.Trim(record.Row[index], " ")
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (record *Record) ProtectedGet(field string) (string, error) {
	if index, present := record.headerMap[field]; present {
		if len(record.Row) >= index {
			return record.Row[index], nil
		} else {
			return "", errors.New("Index  out of bounds of HeaderMap")
		}
	} else {
		return "", errors.New("Value " + field + " not found")
	}
}

func NewEmitter(filepath string) RecordEmitter {
	recordEmitter := RecordEmitter{Filepath: filepath}

	err := recordEmitter.openFile()

	if err != nil {
		fmt.Println(err)
	}

	return recordEmitter
}

type RecordEmitter struct {
	Filepath  string
	file      *os.File
	Reader    *csv.Reader
	HeaderMap map[string]int
	cursor    int
}

func (re *RecordEmitter) openFile() error {
	file, err := os.Open(re.Filepath)
	if err != nil {
		return err
	}

	re.file = file

	if err != nil {
		return err
	}

	re.Reader = csv.NewReader(NewCRLFReader(file))
	re.Reader.LazyQuotes = true

	return nil
}

func (re *RecordEmitter) CloseFile() bool {
	re.file.Close()
	return true
}

func (re *RecordEmitter) buildHeaderMap(headerRow []string) {
	re.HeaderMap = make(map[string]int)

	for index, value := range headerRow {
		re.HeaderMap[value] = index
	}
}

func (re *RecordEmitter) Start() <-chan Record {
	ch := make(chan Record)
	go func() {
		for {
			row, err := re.Reader.Read()

			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}

			if re.cursor == 0 {
				re.cursor += 1
				re.buildHeaderMap(row)
				continue
			}

			ch <- Record{Row: row, headerMap: re.HeaderMap}
			re.cursor += 1
		}

		fmt.Println(re.cursor)

		re.CloseFile()
		close(ch)
	}()
	return ch
}
