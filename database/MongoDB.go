package database

import (
	"context"
	"log"
	"os/exec"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitMongoDB() {

	cmd := exec.Command("ssh", "-f", "-N", "-L", "27017:localhost:27017", "root@209.38.33.161")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Ошибка при запуске SSH-туннеля: %v", err)
	}
	log.Println("SSH-туннель успешно запущен")

	mongoURL, err := beego.AppConfig.String("mongo_url")
	if err != nil {
		log.Fatalf("Не удалось получить mongo_url из конфигурации: %v", err)
	}

	clientOptions := options.Client().ApplyURI(mongoURL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Не удалось подключиться к MongoDB: %v", err)
	}

	if err = Client.Ping(ctx, nil); err != nil {
		log.Fatalf("Не удалось подключиться к серверу MongoDB: %v", Client)
	}

	log.Println("Подключение к MongoDB успешно!")
}
