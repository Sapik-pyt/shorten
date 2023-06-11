### Сервис по сокращению ссылок. Тестовое задание для Ozon Fintech.

## Описание проекта
Сервис разработан на языке Go.
Реализована трехуровневая архитектура: `handler` -> `service` -> `storage`.
В качестве хранилища реализованы PostgreSql и In-Memory Storage (на выбор).
Для взаимодействия с БД PostgreSql использовал драйвер `jackc/pgx`.
Реализован обработчик `gRpc`.
Протестировал слои бизнес-логики и хэндлеров с помощью `testify` и `gomock`.

## Пример использования сервиса

## Пример создания сокращения ссылки
```
curl --location --request POST 'http://localhost:8080/post' \
 --header 'Content-Type: application/json' \
    --data-raw '{"original_link": "https://www.youtube.com/"}'
 ```

Ответ:

Content-Type: application/json

```json
{
    "shortLink":"fkkG0hxOP6"
}
```

## Пример получения оригинальной ссылки

```
curl --location --request GET 'http://localhost:8080/get/fkkG0hxOP6' \
 --header 'Content-Type: application/json'
```

```json
{
    "originalLink":"https://www.youtube.com/"
}
```

### Запуск сервиса
Для запуска сервиса с БД в роли хранилища использовать команду 

```
make run_db
```
Для запуска сервиса в In-Memory в роли хранилища использовать команду
```
make run_inmemory
```

### Для запуска тестов 
```
make test
```