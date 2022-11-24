# Project: URL shortner

## Increment14 - EN
In the database table with short URLS, create an additional field with a flag, indicating that this URL should be considered deleted. Then add an asynchronous handler DELETE /api/user/urls, which accepts a list of short URL identifiers to be deleted in the format:
````
[ "a", "b", "c", "d", ...]
````
If the request is accepted successfully, the handler should return HTTP status 202 Accepted.The actual result of deletion may occur later - there is no need to notify the client of the operation's success or failure in any way.

Only the user who created the URL can successfully delete the URL. When requesting a deleted URL using the handler GET /{id}, a 410 Gone status should be returned.

Advice:
- Use batch update to effectively set the Deleted flag in the database.
- Use the fan-In pattern to maximize buffer load of update objects.

## Increment14 - RU

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

