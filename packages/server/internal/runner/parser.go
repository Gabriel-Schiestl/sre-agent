package runner

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

type JTLRecord struct {
	TimeStamp       int64
	Elapsed         int
	Label           string
	ResponseCode    string
	ResponseMessage string
	ThreadName      string
	DataType        string
	Success         bool
	FailureMessage  string
	Bytes           int
	SentBytes       int
	GrpThreads      int
	AllThreads      int
	Latency         int
	IdleTime        int
	Connect         int
}

// TODO: implement streaming parsing for large files instead of loading everything into memory.
// TODO: add safe acessor helper to avoid index out of range and not to enforce 16 fields per record if not needed.
func parse(content []byte) ([]JTLRecord, error) {
	buff := bytes.NewReader(content)
	reader := csv.NewReader(buff)
	reader.FieldsPerRecord = 16

	var records []JTLRecord
	rec, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, row := range rec {
		if i == 0 {
			continue
		}
		record := JTLRecord{
			TimeStamp:       parseInt64(row[0]),
			Elapsed:         parseInt(row[1]),
			Label:           row[2],
			ResponseCode:    row[3],
			ResponseMessage: row[4],
			ThreadName:      row[5],
			DataType:        row[6],
			Success:         parseBool(row[7]),
			FailureMessage:  row[8],
			Bytes:           parseInt(row[9]),
			SentBytes:       parseInt(row[10]),
			GrpThreads:      parseInt(row[11]),
			AllThreads:      parseInt(row[12]),
			Latency:         parseInt(row[13]),
			IdleTime:        parseInt(row[14]),
			Connect:         parseInt(row[15]),
		}
		records = append(records, record)
	}

	return records, nil
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseInt64(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseBool(s string) bool {
	return s == "true"
}