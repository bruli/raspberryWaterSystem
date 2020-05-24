package jsontime

import (
	"strconv"
	"time"
)

const Layout = "2006-01-02 15:04:05"

type JsonTime time.Time

func (j *JsonTime) MarshalJSON() ([]byte, error) {
	t := strconv.Quote(time.Time(*j).Format(Layout))
	return []byte(t), nil
}

func (j *JsonTime) ToTime() time.Time {
	return time.Time(*j)
}
func (j *JsonTime) ToString() string {
	return j.ToTime().Format(Layout)
}

func (j *JsonTime) UnmarshalJSON(s []byte) error {

	unquote, err := strconv.Unquote(string(s))
	if err != nil {
		return err
	}
	t, err := time.Parse(Layout, unquote)
	if err != nil {
		return err
	}
	*j = JsonTime(t)
	return nil
}
