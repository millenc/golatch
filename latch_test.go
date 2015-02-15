package golatch

import (
	"regexp"
	"testing"
)

func TestGetCurrentDateTime(t *testing.T) {
	l := &Latch{}
	d := l.getCurrentDateTime()

	if l == nil || len(d) == 0 {
		t.Error("Null or empty date time")
	} else {
		match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, d)
		if !match {
			t.Errorf("Date %s doesn't have the expected format.", d)
		}
	}
}
