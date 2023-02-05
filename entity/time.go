package entity

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05Z"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	if strings.Contains(string(b), "Z") {
		ct.Time, err = time.Parse(ctLayout, string(b))
		return
	}
	if len(b) > 22 {
		b = b[0:22]
	}
	ct.Time, err = time.Parse(ctLayout, string(b))
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format("2006-01-02T15:04:05")), nil
}

type SummaryTime struct {
	time.Time
}

const summaryLayout = "2006-01-02T15:04:05"

func (ct *SummaryTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(ctLayout)), nil
}

func (ct *SummaryTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	if len(b) > 22 {
		b = b[0:22]
	}
	ct.Time, err = time.Parse(summaryLayout, string(b))
	return
}
