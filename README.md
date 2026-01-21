# Веб-сервис бронирования аудиторий

Данный проект представляет собой реализацию сайта для бронирования аудиторий, созданный в стиле дизайна LMS Центрального университета. 

## Содержание
- [Технические особенности](#технические-особенности)
- [Запуск проекта](#запуск-проекта)
- [Система ролей](#система-ролей)
- [Booking API (Java)](#booking-api-java)
- [User API (Golang)](#user-api-golang)

## Технические особенности 
- На `Golang` реализован REST-сервис с помощью `net/http` для авторизации и управления пользователями
- Авторизация происходит с помощью `JWT-токенов`
- Пароли сохраняются с помощью `crypto/bcrypt`
- На `Java` реализован REST-сервис на фреймворке `Spring` для управления аудиториями и бронированиями
- Используемая база данных – `PostgreSQL`. На Golang обращение происходит в явном виде с использованием `sqlx`, а на Java – с применением `ORM`
- Frontend полностью реализован на `React`
- Сервисы изолированы и поднимаются с помощью `docker compose`
- Имеется несколько `Unit-тестов` на Java для проверки корректности валидации времени

## Запуск проекта
Перед запуском проекта необходимо создать файл `.env` со всеми необходимыми переменными окружения. Пример приведен в файле `.env.dist` и для тестового запуска достаточно скопировать его.

Затем необходимо поднять приложение с помощью `docker compose`:
```bash
docker compose up -d --build
```
Прослушиваемые порты:
- Сервис на Golang – `7070`
- Сервис на Java – `8080`
- Frontend – `3000`

Таким образом, чтобы открыть сайт локально, необходимо перейти по ссылке – [http://localhost:3000](http://localhost:3000)

Чтобы обратиться к API сервисам Java и Go, необходимо отправлять запросы на `http://localhost:8080/api/v1/...` и `http://localhost:7070/api/v1/...` соответственно

Также нужно установить следующие переменные окружения в файле `.env`:
- POSTGRES_USER – ник пользователя бд
- POSTGRES_PASSWORD – пароль пользователя бд
- POSTGRES_DATABASE – название бд
- JAVA_SOURCE – строка вида `jdbc:postgresql://db:5432/{POSTGRES_DATABASE}`
- JWT_SECRET_KEY – секретный ключ шифрования JWT-токенов, рекомендуется указывать длинный набор случайных символов
- JWT_ADMIN_DURATION – длительность действия JWT-токена для администраторов в **минутах** (необязательное поле; по умолчанию 60)
- JWT_USER_DURATION – длительность действия JWT-токена для пользователей в **минутах** (необязательное поле; по умолчанию 60)

Если не указать все выше перечисленные поля, не считая необязательных, то приложение не запустится. Если не указывать необязательные поля, то высветится предупреждение.

## Система ролей
На сайте присутствует простая система ролей:
- MAIN_ADMIN – Пользователь с полными правами на сайте, но такой ролью может обладать **только один** пользователь. Планируется реализовать возможность передачи прав.
- ADMIN – Пользователь, который может редактировать других пользователей с ролью `USER` (будет в будущем), управлять аудиториями и чужими записями. Права на редактирование других `ADMIN` не имеет.
- USER – Обычный пользователь, который может редактировать только свои данные.


> УКАЗАНИЕ: Первый зарегистрировавшийся пользователь сайта автоматически получает роль MAIN_ADMIN.

> УКАЗАНИЕ: Управление от лица администратора происходит только через API. Изменение пользователей через администратора не реализовано.

## Booking API (Java)
Swagger API доступен по [ссылке](http://localhost:8080/swagger-ui/index.html) после успешного поднятия docker compose

Отправление запросов через Swagger может происходить некорректно. В таком случае можно использовать [Postman](https://www.postman.com/downloads/), [Insomnia](https://insomnia.rest/download) или подобные.

## User API (Golang)

### localhost:7070/api/v1/signup
Регистрирует нового пользователя
```
METHOD: POST
HEADERS: 
- Content-Type: application/json; charset=utf-8
BODY:
- email (Почта)
- first_name (Имя)
- second_name (Фамилия)
- patronymic (Отчество, необязательное поле)
- password (Пароль)

EXAMPLE:
{
    	"email": "admin@example.com",
    	"first_name": "Test",
	"second_name": "Test",
	"patronymic": "Test",
	"password": "12345" 
}

RESULT:
[201 CREATED]
```

### localhost:7070/api/v1/login
Генерация JWT-токена
```
METHOD: POST
HEADERS: 
- Content-Type: application/json; charset=utf-8
BODY:
- email (Почта)
- password (Пароль)

EXAMPLE:
{
    	"email": "admin@example.com",
	"password": "12345" 
}

RESULT:
[201 CREATED]
HEADERS:
- Authorization (JWT-токен)
```

### localhost:7070/api/v1/user
Парсинг JWT-токена (необходим для бизнес-логики в некоторых сервисах)
```
METHOD: GET
HEADERS: 
- Authorization (JWT-токен)
BODY: NONE

RESULT:
[200 OK]
{
    	"id": 1,
	"email": "admin@example.com",
	"first_name": "Test",
	"second_name": "Test",
	"patronymic": "Test",
	"role": "MAIN_ADMIN",
	"exp": 1751342317
}
```
Изменение пользователя
```
METHOD: PUT
HEADERS: 
- Content-Type: application/json; charset=utf-8
- Authorization (JWT-токен)
BODY:
- email (Почта, необязательное поле)
- first_name (Имя, необязательное поле)
- second_name (Фамилия, необязательное поле)
- patronymic (Отчество, необязательное поле)

EXAMPLE:
{
    	"email": "newadmin@example.com",
    	"first_name": "newTest",
	"second_name": "newTest",
	"patronymic": "newTest",
}

RESULT:
[204 NO CONTENT]
HEADERS: 
- Authorization (новый JWT-токен)
```
Удаление аккаунта
```
METHOD: DELETE
HEADERS: 
- Authorization (JWT-токен)
BODY: NONE

RESULT:
[204 NO CONTENT]
```
### localhost:7070/api/v1/user/password
Изменение пароля
```
METHOD: PUT
HEADERS: 
- Content-Type: application/json; charset=utf-8
- Authorization (JWT-токен)
BODY:
- new_password (новый пароль)

EXAMPLE:
{
	"old_password": "12345",
    	"new_password": "1234567890"
}

RESULT:
[204 NO CONTENT]
```
