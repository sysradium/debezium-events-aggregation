package redis

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type RedisStreamUnmarshaler struct{}

func (RedisStreamUnmarshaler) Unmarshal(values map[string]interface{}) (*message.Message, error) {
	payload := values["value"].(string)
	return message.NewMessage(uuid.NewString(), []byte(payload)), nil
}
