package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"task/floodcontroller"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Замените на ваш адрес Redis
		Password: "",               // Если пароль не пустой, укажите его здесь
		DB:       0,                // По умолчанию используется DB 0
	})

	var flood FloodControl
	flood = floodcontroller.NewFloodConfig(client, 20, 10*time.Second)

	// Проверка алгоритма на 10 запросах
	for i := 0; i < 10; i++ {
		allowed, err := flood.Check(context.Background(), 222)
		if err != nil {
			log.Fatal(err)
			return
		}
		if allowed {
			fmt.Println("Request allowed")
		} else {
			fmt.Println("Request denied for")
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
