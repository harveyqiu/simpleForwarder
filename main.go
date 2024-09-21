// main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Info struct {
	URL string `json:"url"`
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var info Info
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &info); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Println(info.URL)

	item := map[string]interface{}{
		"url":  info.URL,
		"tags": []string{"Privacy Law"},
	}

	jsonData, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "https://readwise.io/api/v3/save/", bytes.NewBuffer(jsonData))

	// 从环境变量中获取令牌
	token := os.Getenv("TOKEN")
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Err %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		w.Write([]byte("Created"))
	} else {
		w.Write([]byte("Error Request"))
	}
}

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8999", nil))
}
