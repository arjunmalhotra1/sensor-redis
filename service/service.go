package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/arjun1malhotra/armada/data"
	"github.com/redis/go-redis/v9"
)

// type cache interface {
// 	set(key string, value string) error
// 	get(key string) (string, bool)
// }

type RedisCache struct {
	Client *redis.Client
}

func (r RedisCache) Set(ctx context.Context, key string, value data.Message) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshalling the data in the Set function")
	}
	// RPush will append to the deviceId list
	res := r.Client.RPush(ctx, key, valueBytes)
	return res.Err()
}

func (r RedisCache) Get(ctx context.Context, key string) ([]data.Message, error) {
	// Retrieves all Ids in the list based on the deviceId
	fmt.Println("key: ", key)
	res, err := r.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Println("error getting the result", res)
		return nil, fmt.Errorf("failed to get a value from redis %v", err)
	} else if len(res) == 0 {
		return nil, fmt.Errorf("get key not found %s", key)
	}

	var messages []data.Message
	// Loop through each item in the Redis list
	for _, item := range res {
		var msg data.Message
		// Unmarshal each item into a Message object
		err := json.Unmarshal([]byte(item), &msg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the redis response item %v: %v", item, err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
