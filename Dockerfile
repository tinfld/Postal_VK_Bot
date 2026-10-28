# Упращенный образ меньше весит 
FROM golang:1.25-alpine
# Рабочая директория внутри контейнера
WORKDIR /app

COPY go.mod go.sum ./

# Загружаем Go-модули указанные в go.mod/go.sum
RUN go mod download

# Копируем исходный код проекта в контейнер
COPY . .

# Компилируем Go-приложение в исполняемый файл bot
RUN go build -o bot .

# Команда, которую Docker запускает при старте контейнера
CMD ["./bot"]