package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/arjun1malhotra/armada/data"
	"github.com/arjun1malhotra/armada/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
)

func sendMessage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error reading the input", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var message data.Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		http.Error(res, "invalid data", http.StatusBadRequest)
		return
	}
	ok, err := data.ValidateMessageBody(message)
	if !ok {
		log.Println("error validating data: ", err)
		http.Error(res, fmt.Sprintf("invalid data %v", err), http.StatusBadRequest)
		return
	}

	message.Time = time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("sensor:%s", message.DeviceId)
	err = cache.Set(ctx, key, message)
	if err != nil {
		log.Println("error storing data: ", err)
		http.Error(res, "failed to store the data in cache", http.StatusInternalServerError)
		return
	}

	log.Println("message written: ", message)
	render.JSON(res, req, message)

}

func getMessage(res http.ResponseWriter, req *http.Request) {
	devId := chi.URLParam(req, "device_id")
	fmt.Println("devId: ", devId)
	key := fmt.Sprintf("sensor:%s", devId)

	// TODO: Add data validations

	ctx := context.Background()
	data, err := cache.Get(ctx, key)
	if err != nil {
		log.Println("error retrieving data: ", err)
		http.Error(res, "error retrieving data", http.StatusNotFound)
		return
	}

	log.Println("got data: ", data)
	render.JSON(res, req, data)

}

var cache service.Cache

func main() {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	cache = service.RedisCache{Client: redisClient}
	r := chi.NewMux()
	r.Post("/message", sendMessage)
	r.Get("/message/{device_id}", getMessage)
	http.ListenAndServe(":8086", r)

}
