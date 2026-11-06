# Тестовое задание на Go

Небольшой веб сервер с REST-интерфейсом для обмена информацией про угрозы.

## Запуск

```sh
go run ./cmd/server
```

Или в Docker:

```sh
docker build -t threats-info-exchange .
docker run --rm -p 8080:8080 threats-info-exchange
```

## API

| Метод | Путь | Результат |
| --- | --- | --- |
| `POST` | `/info` | Создаёт запись об угрозе (`201`), повторный MD5 - `409` |
| `GET` | `/info/{md5}` | Возвращает угрозу (`200`), отсутствие записи - `404` |
| `PUT` | `/info/{md5}` | Заменяет данные существующей угрозы (`200`), отсутствие записи - `404` |

MD5 соответвует паттерну '^[a-fA-F0-9]{32}$'. Модель угрозы: поля `md5`, `name`, `description`, `category`, `severity` (`unknown`, `low`, `medium`, `high`, `critical` или пустая строка), `createdAt`, `updatedAt` в UTC.

Пример запроса:

```sh
curl -i -X POST http://localhost:8080/info \
  -H 'Content-Type: application/json' \
  -d '{"md5":"d41d8cd98f00b204e9800998ecf8427e","name":"Trojan.Win32.Generic","severity":"high"}'
```


## Код

- `cmd/server` собирает зависимости и запускает HTTP-сервер.
- `internal/router` регистрирует маршруты; `internal/handler` разбирает запросы и формирует ответы.
- `internal/domain` содержит модели и ошибки; `internal/repo` хранит записи об угрозах.
- `internal/middleware` логирует запросы и обрабатывает паники.

