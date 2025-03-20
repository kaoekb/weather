package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const city = "Saint Petersburg"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	WeatherCode float64 `json:"weathercode"`
}

type WeatherResponse struct {
	CurrentWeather CurrentWeather `json:"current_weather"`
}

func getWeather() (string, error) {
	// Запрос к Open-Meteo API
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=59.9343&longitude=30.3351&current_weather=true")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return "", err
	}

	return fmt.Sprintf("Температура: %.1f°C", weatherResp.CurrentWeather.Temperature), nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		if string(message) == "weather" {
			weather, err := getWeather()
			if err != nil {
				log.Println(err)
				conn.WriteMessage(websocket.TextMessage, []byte("Ошибка получения погоды"))
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte(weather))
			}
		}
	}
}

func main() {
	// Принтуем тестовое 'ja workau'
	fmt.Println("ja workau")
	// Создаем WebSocket-соединение на порт 8080
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
