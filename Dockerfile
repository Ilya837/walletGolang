FROM golang:1.25-alpine3.23 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY main.go ./
COPY dataStorage/dataStorage.go ./dataStorage/
COPY server/server.go ./server/


# Собираем бинарник
RUN go build -o runServer .

# Используем лёгкий образ для запуска
FROM alpine:3.22
WORKDIR /app

# Копируем бинарник из стадии сборки
COPY config.env .
COPY --from=builder /app/runServer .

# Указываем порт, который контейнер будет слушать
EXPOSE 80

# Команда для запуска сервера
CMD ["./runServer"]