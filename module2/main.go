package main

import (
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
	"golang.org/x/crypto/blowfish"
)

// CS2Skin представляет собой структуру данных для скина в CS2.
type CS2Skin struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Wear       string `json:"wear"`
	Pattern    string `json:"pattern"`
	Side       string `json:"side"`
	WeaponType string `json:"weapon_type"`
}

func main() {
	// Подключаемся к серверу NATS-Streaming
	sc, err := stan.Connect("test-cluster", "client-2")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Создаем новый скин
	newSkin := CS2Skin{
		Name:       "Example Skin",
		Price:      100,
		Wear:       "Factory New",
		Pattern:    "Custom Pattern",
		Side:       "Terrorist",
		WeaponType: "AK-47",
	}

	// Преобразуем скин в JSON
	skinJSON, err := json.Marshal(newSkin)
	if err != nil {
		log.Fatal(err)
	}

	// Хэшируем данные с использованием SHA-256
	hashedData := sha256.Sum256(skinJSON)

	// Шифруем данные с использованием Blowfish
	key := []byte("example-key")
	block, err := blowfish.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	encryptedData := make([]byte, len(hashedData))
	block.Encrypt(encryptedData, hashedData[:])

	// Отправляем зашифрованные данные через NATS-Streaming
	err = sc.Publish("skins", encryptedData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encrypted and published skin data successfully.")

	// Подписываемся на канал "skins" и получаем сообщения
	sub, err := sc.Subscribe("skins", func(msg *stan.Msg) {
		// Дешифруем данные с использованием Blowfish
		decryptedData := make([]byte, len(msg.Data))
		block.Decrypt(decryptedData, msg.Data)

		// Проверяем хэш с использованием SHA-256
		receivedHash := sha256.Sum256(decryptedData)
		if base64.StdEncoding.EncodeToString(receivedHash[:]) == base64.StdEncoding.EncodeToString(hashedData[:]) {
			fmt.Println("Data integrity verified: Skin data received successfully.")
		} else {
			fmt.Println("Data integrity check failed: Skin data corrupted or tampered.")
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	// Ожидаем сообщения в течение 5 секунд
	time.Sleep(5 * time.Second)
}