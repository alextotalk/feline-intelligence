FROM golang:1.24.1-alpine AS builder

WORKDIR /usr/local/src

# Установка необходимых пакетов
RUN apk --no-cache add bash git make

# Копируем зависимости Go и загружаем модули
COPY go.mod go.sum ./
RUN go mod download

# Создаем директорию для конфигурации
RUN mkdir -p /usr/local/src/config

# Копируем конфигурационный файл в правильное место
COPY config/local.yaml /usr/local/src/config/local.yaml

# Копируем исходный код
COPY cmd/ cmd/
COPY internal/ internal/

# Сборка Go-приложения
RUN go build -o ./bin/app ./cmd/app/main.go

# Финальный этап для создания минимального образа
FROM alpine AS runner

# Копируем бинарный файл из фазы сборки
COPY --from=builder /usr/local/src/bin/app /app

# Копируем конфигурацию
COPY --from=builder /usr/local/src/config /usr/local/src/config

# Указываем рабочую директорию
WORKDIR /usr/local/src

# Указываем команду для запуска контейнера
CMD ["/app"]
