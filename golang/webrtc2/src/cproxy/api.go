package cproxy

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GenerateUserID() string {
	h, _ := os.Hostname()
	return fmt.Sprintf("%s_%d_%d", h, os.Getuid(), os.Getpid())
}

func SaveToFile(name string, data []byte) {
	file, err := os.Create(name)
	if err != nil {
		log.Println("Cannot create file: ", name, " ", err.Error())
		return
	}

	file.Write(data)
	file.Sync()
	defer file.Close()
}

func AppendToFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(data)
}

func GenerateSessionID() string {
	var id = "0x"

	rand.Seed(time.Now().Unix())
	result := rand.Perm(128)
	for _, i := range result {
		id = id + strconv.Itoa(i)
	}

	return id
}
