package timestamp

import "time"

func Now() int64 {
	return time.Now().UnixMicro()
}

func Add(duration time.Duration) int64 {
	return time.Now().Add(duration).UnixMicro()
}
