package timestamp

import "time"

func Now() int64 {
	return time.Now().UnixMicro()
}
