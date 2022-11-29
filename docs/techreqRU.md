Russian

## Инкремент 14

Сделайте в таблице базы данных с сокращёнными URL дополнительное поле с флагом, указывающим на то, что URL должен считаться удалённым.
Далее добавьте в сервис новый асинхронный хендлер DELETE /api/user/urls, который принимает список идентификаторов сокращённых URL для удаления в формате:
```
[ "a", "b", "c", "d", ...]
```
В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted. Фактический результат удаления может происходить позже — каким-либо образом оповещать пользователя об успешности или неуспешности не нужно.
Успешно удалить URL может пользователь, его создавший. При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone.
Совет:
Для эффективного проставления флага удаления в базе данных используйте множественное обновление (batch update).
Используйте паттерн fanIn для максимального наполнения буфера объектов обновления.

## Инкремент 13

Сделайте в таблице базы данных с сокращёнными URL уникальный индекс для поля с исходным URL. Это позволит избавиться от дублирующих записей в базе данных.
При попытке пользователя сократить уже имеющийся в базе URL через хендлеры POST / и POST /api/shorten сервис должен вернуть HTTP-статус 409 Conflict, а в теле ответа — уже имеющийся сокращённый URL в правильном для хендлера формате.

## Инкремент 12
Задание для трека «Сервис сокращения URL»

Добавьте новый хендлер POST /api/shorten/batch, принимающий в теле запроса множество URL для сокращения в формате:
Скопировать код
```
[
    {
        "correlation_id": "<строковый идентификатор>",
        "original_url": "<URL для сокращения>"
    },
    ...
]
```
В качестве ответа хендлер должен возвращать данные в формате:
Скопировать код
```
[
    {
        "correlation_id": "<строковый идентификатор из объекта запроса>",
        "short_url": "<результирующий сокращённый URL>"
    },
    ...
]
```

## Инкремент 11

Перепишите сервис так, чтобы СУБД PostgreSQL стала хранилищем сокращённых URL вместо текущей реализации.

Сервису нужно самостоятельно создать все необходимые таблицы в базе данных. Схема и формат хранения остаются на ваше усмотрение.

При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d или при их пустых значениях вернитесь последовательно к:
- хранению сокращённых URL в файле при наличии соответствующей переменной окружения или флага командной строки;
- хранению сокращённых URL в памяти.

## Инкремент 10
Задание для трека «Сервис сокращения URL»

- Добавьте в сервис функциональность подключения к базе данных. В качестве СУБД используйте PostgreSQL не ниже 10 версии.
- Добавьте в сервис хендлер GET /ping, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
- Строка с адресом подключения к БД должна получаться из переменной окружения DATABASE_DSN или флага командной строки -d.

## Инкремент 9

Добавьте в сервис функциональность аутентификации пользователя.
Сервис должен:
- Выдавать пользователю симметрично подписанную куку, содержащую уникальный идентификатор пользователя, если такой куки не существует или она не проходит проверку подлинности.
- Иметь хендлер GET /api/user/urls, который сможет вернуть пользователю все когда-либо сокращённые им URL в формате:
  Скопировать код
````
[
    {
        "short_url": "http://...",
        "original_url": "http://..."
    },
    ...
]
````
- При отсутствии сокращённых пользователем URL хендлер должен отдавать HTTP-статус 204 No Content.

## Инкремент 8

Добавьте поддержку gzip в ваш сервис. Научите его:
- принимать запросы в сжатом формате (HTTP-заголовок Content-Encoding);
- отдавать сжатый ответ клиенту, который поддерживает обработку сжатых ответов (HTTP-заголовок Accept-Encoding).
  Вспомните middleware из урока про HTTP-сервер, это может вам помочь.

## Инкремент 7

Поддержите конфигурирование сервиса с помощью флагов командной строки наравне с уже имеющимися переменными окружения:
- флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
- флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
- флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).

## Инкремент 6

Сохраняйте все сокращённые URL на диск в виде файла. При перезапуске приложения все URL должны быть восстановлены.

Путь до файла должен передаваться в переменной окружения FILE_STORAGE_PATH.

При отсутствии переменной окружения или при её пустом значении вернитесь к хранению сокращённых URL в памяти.

## Инкремент 5

Добавьте возможность конфигурировать сервис с помощью переменных окружения:
- адрес запуска HTTP-сервера с помощью переменной SERVER_ADDRESS.
- базовый адрес результирующего сокращённого URL с помощью переменной BASE_URL.

## Инкремент 4

Добавьте в сервер новый эндпоинт POST /api/shorten, принимающий в теле запроса JSON-объект {"url":"<some_url>"} и возвращающий в ответ объект {"result":"<shorten_url>"}.

Не забудьте добавить тесты на новый эндпоинт, как и на предыдущие.

Помните про HTTP content negotiation, проставляйте правильные значения в заголовок Content-Type.

## Инкремент 3

Вы написали приложение с помощью стандартной библиотеки net/http. Используя любой пакет (роутер или фреймворк), совместимый с net/http, перепишите ваш код.

Задача направлена на рефакторинг приложения с помощью готовой библиотеки.

Обратите внимание, что необязательно запускать приложение вручную: тесты, которые вы написали до этого, помогут вам в рефакторинге.

## Инкремент 2

Покройте сервис юнит-тестами. Сконцентрируйтесь на покрытии тестами эндпоинтов, чтобы защитить API сервиса от случайных изменений.

## Инкремент 1

Напишите сервис для сокращения длинных URL. Требования:
- Сервер должен быть доступен по адресу: http://localhost:8080.
- Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
- Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
- Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
- Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.

---
# Other Instructions:

## go-musthave-shortener-tpl
Шаблон репозитория для практического трека «Go в веб-разработке».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.
