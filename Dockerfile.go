#
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum (для загрузки зависимостей)
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Компилируем Go-программу
RUN go build -o weather-server ./go-server/main.go

# Открываем порт, на котором будет работать сервер
EXPOSE 8080

# Запускаем Go-сервер
CMD ["./weather-server"]
