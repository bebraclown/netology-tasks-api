# REST API задач

## Запуск

Нужен Go 1.20 или новее. Из корня проекта запустить сервер.

```bash
go run ./cmd/server
```

Сервер доступен по адресу `http://127.0.0.1:8080`. Задачи хранятся в памяти до остановки сервера.

## Эндпоинты

| Метод и путь | Действие | Код успеха |
| --- | --- | --- |
| `GET /tasks` | Получить список задач | `200 OK` |
| `POST /tasks` | Создать задачу | `201 Created` |
| `GET /tasks/{id}` | Получить задачу | `200 OK` |
| `PUT /tasks/{id}` | Заменить заголовок и статус задачи | `200 OK` |
| `DELETE /tasks/{id}` | Удалить задачу | `204 No Content` |

`POST` и `PUT` принимают `application/json`. Поле `title` обязательно, `done` по умолчанию равно `false`. ID создаёт сервер. Формат ошибок — `{"error":"сообщение"}`.

## Проверка через curl

Выполняй команды по порядку после запуска сервера с пустым списком задач.

```bash
curl -i -X POST 'http://127.0.0.1:8080/tasks' -H 'Content-Type: application/json' -d '{"title":"Learn Go"}'
# 201 Created
# {"id":1,"title":"Learn Go","done":false}

curl -i -X POST 'http://127.0.0.1:8080/tasks' -H 'Content-Type: application/json' -d '{"title":" "}'
# 400 Bad Request
# {"error":"поле title обязательно"}

curl -i -X GET 'http://127.0.0.1:8080/tasks'
# 200 OK
# [{"id":1,"title":"Learn Go","done":false}]

curl -i -X PATCH 'http://127.0.0.1:8080/tasks'
# 405 Method Not Allowed
# {"error":"метод не поддерживается"}

curl -i -X GET 'http://127.0.0.1:8080/tasks/1'
# 200 OK
# {"id":1,"title":"Learn Go","done":false}

curl -i -X GET 'http://127.0.0.1:8080/tasks/999'
# 404 Not Found
# {"error":"задача не найдена"}

curl -i -X PUT 'http://127.0.0.1:8080/tasks/1' -H 'Content-Type: application/json' -d '{"title":"Learn Go HTTP","done":true}'
# 200 OK
# {"id":1,"title":"Learn Go HTTP","done":true}

curl -i -X PUT 'http://127.0.0.1:8080/tasks/1' -H 'Content-Type: application/json' -d '{"title":"Learn Go HTTP","done":"true"}'
# 400 Bad Request
# {"error":"поле done должно быть логическим значением"}

curl -i -X DELETE 'http://127.0.0.1:8080/tasks/1'
# 204 No Content, тело пустое

curl -i -X DELETE 'http://127.0.0.1:8080/tasks/1'
# 404 Not Found
# {"error":"задача не найдена"}
```
