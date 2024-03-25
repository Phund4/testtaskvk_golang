1. Клонировать репозиторий командой
```bash
git clone https://github.com/Phund4/testtaskvk_golang.git
```
2. Перейти в корневую директорию проекта.
3. В pgadmin4 создать БД и вставить содержимое файла test.dumb в скрипт.
либо написать команду
```bash
pg_restore [параметры для подключения] [параметры восстановления] [дамп базы данных]
```
4. В корневой директории проекта создать файл .env
5. Написать в нем три переменной окружения:
   PASSWORD - пароль при подключении к БД.
   USER - пользователь postgres.
   DBNAME - название БД.
6. В корневой директории проекта написать команду
```bash
go run main.go
```
7. Сервер должен запуститься.
8. Для запуска тестов команда
```bash
go test [директория client или quest, где находятся тесты]
```
9. Для запуска RabbitMQ сперва установить и запустить контейнер с Rabbit:
```bash
docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```
10. Затем запустить из корня проекта брокер командой
```bash
go run rabbit/rabbit-consumer.go
```
URL для запросов:
1.http://localhost:8080/addclient
2.http://localhost:8080/addquest
3.http://localhost:8080/completequest
4.http://localhost:8080/getclientinfo
