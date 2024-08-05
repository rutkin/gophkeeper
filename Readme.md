Клиент расположен в /cmd/gophkeeper
Сервер расположен в /cmd/server

Использование клиента:

1) Регистрация пользователя
gophkeeper register -u admin
2) Логин
gophkeeper login -u admin
3) Загрузка бинарных данных
gophkeeper set file --path /path/to/file.bin
4) Получение списка всех данных
gophkeeper list
5) Получение бинарных данных
gophkeeper get --id {guid}

Полный список команд gophkeeper --help