package data

import "time"

func toTime(v any) time.Time {
	switch t := v.(type) {
	case time.Time:
		return t
	case []byte:
		parsed, _ := time.Parse(time.RFC3339, string(t))
		return parsed
	case string:
		parsed, _ := time.Parse(time.RFC3339, t)
		return parsed
	}
	return time.Time{}
}
