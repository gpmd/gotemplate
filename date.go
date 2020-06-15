package gotemplate

import "time"

func dateFmt(format, datestring string) string {
	if format == "ukshort" {
		format = "02/01/06"
	}
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, datestring)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05-0700", datestring)
		if err != nil {
			return datestring
		}
	}
	return t.Format(format)
}

func dateFmtLayout(format, datestring, layout string) string {
	if format == "ukshort" {
		format = "02/01/06"
	}
	t, err := time.Parse(layout, datestring)
	if err != nil {
		return err.Error()
	}
	return t.Format(format)
}

func formatUKDate(datestring string) string {
	return dateFmt("ukshort", datestring)
}

func datetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ukdate() string {
	return time.Now().Format("02/01/06")
}

func ukdatetime() string {
	return time.Now().Format("02/01/06 15:04:05")
}

func timeFormat(format string) string {
	return time.Now().Format(format)
}

func timeFormatMinus(format string, minus float64) string {
	return time.Now().Add(time.Duration(minus) * -time.Second).Format(format)
}

func unixtimestamp() int32 {
	return int32(time.Now().Unix())
}

func nanotimestamp() int64 {
	return int64(time.Now().UnixNano())
}

func timestamp() string {
	return time.Now().String()
}
