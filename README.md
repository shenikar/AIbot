# :bust_in_silhouette: AIbot

AIbot — Telegram-бот для генерации изображений с использованием DALL·E и других AI-сервисов, разработанный на Go.

## Функции

- Генерация изображений с помощью DALL·E.
- Интеграция с Telegram.
- Простая настройка и конфигурация.

## Установка

### Требования

- Go 1.18 и выше.
- Docker (по желанию для использования Docker Compose).

### Шаги

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/shenikar/AIbot.git
2. Установить зависимости:
    ```bash
    go mod tidy
    ```
3. Настроить переменные окружения:
    - BOT_TOKEN -  токен Telegram API.
    - REDIS_ADDR - Redis port
    - ProxyAPIKey -  ключ API через ProxyAPI
    - ProxyAPIURL - url ProxyAPI
4. Запустить бота:
    ```bash
    go run ./cmd/app/main.go
    ```

## Использование Docker

1. Создайте .env файл с переменными окружения.
2. Запустите через Docker Compose:
    ```bash
    docker-compose up --build
    ```
## Лицензия
   Проект распространяется под лицензией MIT — подробности см. в файле [LICENSE](LICENSE).
