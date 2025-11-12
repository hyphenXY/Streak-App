// cache/redis.go
package cache

import (
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Reminder struct {
	UserID  int       `json:"userId"`
	ClassID int       `json:"classId"`
	Time    time.Time `json:"time"`
}

func AddReminder(reminder Reminder) error {
    // sendAt = classTime - 30 mins
    sendAt := reminder.Time.Add(-30 * time.Minute).Unix()

    data, _ := json.Marshal(reminder)
    return Rdb.ZAdd(Ctx, "reminders", redis.Z{
        Score:  float64(sendAt),
        Member: data,
    }).Err()
}