FROM golang:1.25.4

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o short-url ./cmd/short-url

# Экспортируем порт
EXPOSE 8084

# Запускаем приложение
CMD ["./short-url"]