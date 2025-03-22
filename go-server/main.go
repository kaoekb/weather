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
	log.Println("Запрос к Open-Meteo API для получения текущей погоды...")
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=59.9343&longitude=30.3351&current_weather=true")
	if err != nil {
		log.Printf("Ошибка при запросе погоды: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа от API: %v\n", err)
		return "", err
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		log.Printf("Ошибка при разборе ответа JSON: %v\n", err)
		return "", err
	}

	weather := fmt.Sprintf("Температура: %.1f°C", weatherResp.CurrentWeather.Temperature)
	log.Printf("Получены данные о погоде: %s\n", weather)
	return weather, nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Попытка установления WebSocket-соединения...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при установке WebSocket-соединения: %v\n", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket-соединение установлено успешно.")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Ошибка при чтении сообщения от клиента: %v\n", err)
			break
		}

		log.Printf("Получено сообщение от клиента: %s\n", string(message))
		if string(message) == "weather" {
			weather, err := getWeather()
			if err != nil {
				log.Println("Ошибка получения погоды.")
				conn.WriteMessage(websocket.TextMessage, []byte("Ошибка получения погоды"))
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte(weather))
			}
		}
	}
}

func main() {
	log.Println("Запуск сервера на порту 8080...")
	// Создаем WebSocket-соединение на порт 8080
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
