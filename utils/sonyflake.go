package utils

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(
		sonyflake.Settings{
			StartTime: time.Unix(1685548800, 0),
		},
	)
	rand.Seed(time.Now().UnixNano())
}

// GenSeqId 生成自增ID
func GenSeqId() string {
	id, _ := sf.NextID()
	return strconv.FormatUint(id, 36)
}
