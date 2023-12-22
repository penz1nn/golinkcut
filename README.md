# API
## REST
### POST /new
```
Request body: {"url": "http://example.com/some/long/link"}
Response: {"shortLink": "appdomain.com/shortalias"}
```
При POST запросе на /new будет сгенерирована короткая ссылка длиной 10 символов после "/"
### GET /*shortalias*
```
GET /shortalias (empty request body)
Response: {"url": "http://example.com/some/long/link"}
```
GET запрос на полученный после POST /new запроса адрес вернёт сохранённый ранее url

## GRPC
### CreateLink
```
CreateLink(CreateLinkRequest(Url: "http://example.com/some/long/link"))
-> CreateLinkResponse(Alias: "shortalias")
```
### GetLink
```
GetLink(GetLinkRequest(Alias: "shortalias"))
-> GetLinkResponse(Url: "http://example.com/some/long/link")
```

# Command Line Arguments
 - memory: использовать in-memory хранилище
 - validate: проверять корректность url перед сохранением
 - redirect: (только REST) перенаправлять по ссылке вместо того чтобы просто возвращать её в JSON
 - db-host, db-user, db-port, db-timezone, db-password: параметры для PostgreSQL БД (имеют предустановки в коде)

# Quick Start
```
$ docker compose -f golinkcut_postgres-compose.yml up --build
```
ИЛИ (если будет использовано только in-memory хранилище, тогда контейнер Postgres не нужен)
```
$ docker compose -f golinkcut_standalone-compose.yml up --build
```
# P.S.
- в /for-tests/ находятся скрипты для тестирования API 