# Command-executor service
Сервис для создания и удаленного исполнения bash-команд/скриптов

### Инструкция по запуску
```sh
git clone https://github.com/timohahaa/executor && cd executor
```
Затем создать файл `.env` и указать в нем все нужные переменные окружения (для примера см. `.env.example`)
После выполнить команду:
```sh
docker-compose up
```

Сервер будет доступен для взаимодействия на порту, указанном в `.env` файле.

### Тесты
Для тестов были использованы Postman, стандартная библиотека Go и https://github.com/stretchr/testify.
Запустить тесты можно командой `go test -v ./...`. 
Так же в файле `Executor service tests.postman_collection.json` представлена коллекция тестовых запросов с описанием. Их можно импортировать в Postman и позапускать (не забыть поменять порт в запросах на свой).

### Используемые технологии и библиотеки
- Linux Ubuntu 22.04 
- Golang версии 1.22
- PostgreSQL версии 14
- Redis версии 7.2
- Docker и docker-compose
- **Самописная** библиотека для работы с Postgres (https://github.com/timohahaa/postgres)
- Клиент для работы с Redis (https://github.com/redis/go-redis)
- Роутер Echo (https://github.com/labstack/echo)
- Логгер Logrus (https://github.com/sirupsen/logrus)
- Милая библиотека для работы с конфигами (https://github.com/ilyakaznacheev/cleanenv)

### Выполненные задачи
- API для создания команды, получения команды, получения списка команд
- API для запуска и остановки команды
- Поддержка долгих команд и сохранение вывода команды
- Сборка и запуск приложения с помощью Docker

### Ход мыслей
В целом сложности возникли только с задачей параллельного запуска множества команд и остановки произвольной команды.
Команды хранятся в БД по id, а останавливаются по pid (id процесса запущенной команды). Поэтому было принято решение не позвалять запускать команду с одним и тем же id до того, как завершится выполнение уже запущенной команды с этим id.
Соответсвтенно можно запускать параллельно множество команд с разными id, но нет возможности запускать параллельно несколько команд с одинаковым id.
Я не вижу это большой проблемой, так как требование параллельного запуска команд все еще выполнено. Потенциально можно создать очередь на запуск и тогда можно не дожидаться выполнения предыдущей команды для запуска копии...
