# Telegram Bot Project

## Описание

Этот проект представляет собой Telegram-бота, который взаимодействует с базой данных PostgreSQL и RabbitMQ для обработки запросов на вычисление MD5-хэшей. Бот получает сообщения от пользователей, отправляет задачи в очередь RabbitMQ для обработки, и возвращает результаты пользователям.

## Структура проекта

- **telegram_bot**: Основной контейнер с Telegram-ботом.
- **md5worker**: Контейнер для обработки задач по вычислению MD5-хэшей.
- **my_bot_db**: Контейнер с базой данных PostgreSQL.
- **rabbitmq**: Контейнер с RabbitMQ.

## Установка и запуск

### Предварительные требования

- Установите Docker и Docker Compose.

### Конфигурация

1. Создайте файл `.env` в корневом каталоге проекта с содержимым:

    ```env
    TELEGRAM_BOT_TOKEN = Ваш_токен_бота
    TIMEOUT_BOT = 60

    DB_USER =
    DB_PASSWORD =
    DB_NAME =
    DB_HOST =
    DB_PORT =

    RABBITMQ_URL =
    RABBITMQ_QUEUE =

    MAX_ATTEMPTS=10
    PERIOD_ATTEMPTS=1
    ```

### Запуск проекта

1. Сборка и запуск контейнеров:

    ```bash
    docker-compose up --build
    ```

2. Для остановки и удаления контейнеров:

    ```bash
    docker-compose down
    ```

### Описание сервисов

- **Telegram Bot**: Обрабатывает сообщения от пользователей, отправляет запросы на вычисление MD5-хэшей в очередь RabbitMQ и возвращает результаты пользователям.
- **md5worker**: Работает с RabbitMQ для получения задач по вычислению MD5-хэшей и возвращает результаты.
- **my_bot_db**: База данных PostgreSQL для хранения информации о запросах.
- **rabbitmq**: Система очередей сообщений RabbitMQ для передачи задач между ботом и обработчиком.

## Примечания

- Убедитесь, что все контейнеры запущены и работают корректно перед использованием бота.
- При необходимости измените параметры в файле `.env` в соответствии с вашими требованиями и настройками.

## Лицензия

Этот проект лицензируется под MIT License - см. файл [LICENSE](LICENSE) для подробностей.

