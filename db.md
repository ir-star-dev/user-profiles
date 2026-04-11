# Выбор базы данных

Нам для проекта "Профили пользователей" нужна база данных, которая будет хранить пользователей.
- Самая популярная БД в продакшене: PostgreSQL
- Второе место: MySQL
- NoSQL (MongoDB) — под специфические задачи

Мы будем использовать PostgreSQL. 
Для тестирования работы с JWT будем также  использовать ее, но в продакшене для хранения токенов обычно используют Redis.

# Пакеты для работы с БД

В go уже есть встроенный пакет для работы с БД - `database/sql`.
`database/sql` — это база, но в продакшене её редко используют напрямую. 
У нее есть преимущества и недостатки. 

Плюсы:
- полный контроль,
- высокая производительность,
- стандарт Go,
- минимальные зависимости.

Минусы:
- много ручного кода,
- ручной маппинг,
- нет работы со структурами "из коробки",
- нет удобных helpers (get, select),
- нет автоматического биндинга структур в sql,
- и пр.

Можно использовать ORM, например `gorm`. В использовании ORM есть минусы:
- скрывает SQL (магия),
- генерирует неочевидные запросы,
- может ухудшать производительность.

Поэтому будем работать с пакетом `sqlx`, который расширяет возможности встроенной Go библиотеки.
Для установки пакета выполним:

```go
go get github.com/jmoiron/sqlx
```

# Поднимем контейнер postgres

Мы будем использовать docker, чтобы поднять контейнер с базой данных. 
В качестве образа возьмем фиксированную версию postgres 16.4
В файле `.env` определим переменные DB_USER, DB_PASS. 
БД будем прокидывать локально в директорию `./postgres-data`, не забывая при этом добавить ее в файл .gitignore.
Порт для postgres дефолтный: 5432
Напишем в файле docker-compose.yml следующее:

```yml
version: '3'

services:
    postgres:
        container_name: postgres_profiles
        image: postgres:16.4
        environment:
            POSTGRES_USER: ${DB_USER}
            POSTGRES_PASSWORD: ${DB_PASS}
            PGDATA: /data/postgres
        volumes: 
            - ./postgres-data:/data/postgres
        ports:
            - "5432:5432"
```

Запустим в терминале `docker compose up`, чтобы понять контейнер.

Если вы работаете в VS code, то установим расширение `PostgreSQL Explorer`.
Через это расширение создадим подключение к нашей БД (нажимая на + Add Connection).
После чего, по шагам, будем вводить данные:
- Host: localhost
- User: postgres (то, что указывали в .env)
- Pass: pass (то, что указывали в .env)
- Port: 5432
- SSL connection: standart
- Show All Databases

Если все настроено правильно, то вы увидим нашу БД, под названием `localhost`.
Далее в созданном подключении, создадим нашу тестовую БД с названием `demo`.
Для этого правой кнопкой на соединении `localhost` выберем `New Query` и напишем запрос:

```sql
CREATE DATABASE user_profiles;
```

Выполним, нажав F5. Также правой кнопкой на соединении `localhost` выберем `Refresh Items`, должна появится наша база данных.

# Подключение к БД

Для подключение к БД нам потребуется драйвер для работы с postgres. 
Установим пакет 

```go
    go get github.com/lib/pq
```

## Переменные окружения для БД

Нам нужен ряд переменных окружения, которые мы определяем в `.env` файле.

```.env
# DB connection
DRIVER_NAME="postgres"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="user_profiles"
DB_USER="postgres"
DB_PASS="pass"
SSL_MODE="disable" (только для локальной разработки)

# Migrations
MIGRATIONS_PATH="././migrations"
```
## Загрузка переменных окружения в приложение

Установим пакет, который будет считывать наш файл с переменными окружения и загружать их в наше приложение.

```go
go get github.com/joho/godotenv
```

## Пакет configs

Мы создали пакет `configs` файл `configs.go`, в котором функция `Load` возвращает указатель на нашу структуру `Config`.
В структуре будут необходимые данные для подключения к БД: Dsn, Driver, MigrationsPath (об этом позже).

```go
package configs

import (
    "os"
    "fmt"
    "github.com/joho/godotenv"
)

type Config struct {
    Db     DbConfig
    Secret string
}

type DbConfig struct {
    Dsn            string
    Driver         string
    MigrationsPath string
}

func Load() (*Config, error) {
    err := godotenv.Load(".env")
    if err != nil {
        return nil, err
    }

    return &Config{
        Secret: os.Getenv("SECRET"),
        Db: DbConfig{
            Dsn:    buildDSN(),
            Driver: os.Getenv("DRIVER_NAME"),
            MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
        },
    }, nil
}

func buildDSN() string {
    return fmt.Sprintf(
        "%s://%s:%s@%s:%s/%s?sslmode=%s",
        os.Getenv("DRIVER_NAME"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
        os.Getenv("SSL_MODE"),
    )
}
```

## Пакет db

Создадим пакет `pkg/db/postgres.go`, в котором напишем подключение к БД.
В нем мы импортируем наш пакет `configs`, пакет `log` (опционально), пакет `sqlx` и драйвер `pg`.
Функция `Connect` выполняет соединение с БД и пингует его. Возвращает указатель на структуру DB и ошибку, если подключиться не удалось.
Принимает строку - имя драйвера и dsn строку, которая выглядит так:

`"postgres://user:pass@host:port/db_name?sslmode=disable"`

Dsn строку мы уже сформировали в нашем пакете `configs`.

```go
package db

import (
    "user-profiles/configs"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func Connect(conf *configs.Config) (*sqlx.DB, error) {
    dsn := conf.Db.Dsn
    db, err := sqlx.Connect(conf.Db.Driver, dsn)
    if err != nil {
        return nil, err
    }
    log.Println("Connected to PostgreSQL ✔️")
    return db, nil
}
```

## Соединение с БД

В пакете `internal/app/app.go` мы вызываем функцию `Load()` из пакета `configs`.
    
```go
conf, err := configs.Load()
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// DB
dbConn, err := db.Connect(conf)
if err != nil {
    return fmt.Errorf("failed to connect db: %w", err)
}
defer dbConn.Close()
```

Когда запустим наше приложение `go run cmd/main.go`, то при успешном соединении с БД в консоли мы должны увидеть

`"Connected to PostgreSQL ✔️"`

# Миграции

На продакшене миграции - стандарт. 
Позволяют:
- версионировать БД,
- безопасно изменять схему,
- работать в команде.

Самый правильный вариант на продакшене:
- делают деплой приложения,
- выполняют миграции,
- запускают приложение.

Для выполнения миграций, мы создадим отдельный пакет `cmd/migrate/auto.go`
И будем запускать в терминале до запуска приложения.

```go
    go run cmd/migrate/auto.go -up
```

Только для удобства локальной разработки добавлена команда:

```go
    go run cmd/migrate/auto.go -down
```

Она отменит все миграции, то есть удалит все таблицы из БД.

Для работы с миграциями нам нужен пакет

```go
    go get github.com/golang-migrate/migrate/v4
```
    
## Форматы миграций

1. SQL файлы в директории migrations (наиболее распространенное название директории).
2. SQL миграции прямо в Go.
3. Миграции через ORM.

Первый вариант самый популярный, но можно встретить и второй, но реже.

### SQL файлы для миграции

SQL файлы для миграции - это последовательность версионных изменений БД.
Файлы имеют четкий формат именования.

- 0001_название.up.sql
- 0001_название.down.sql

Номер `0001` обязателен, и каждое последующее действие для БД, должно быть описано в файле с номерами 0002, 0003 и т.д.
Номера представляют собой историю изменений.

Название обычно описывает то, что делает эта миграция, например `create_users_table`, `alter_column_code_to_users_table`.

Расширения `.up.sql` и `.down.sql` также обязательны и обозначают:
- up - применить изменения, 
- down - откатить.

На практике, на продакшене файлы `.down.sql` могут отсутствовать, то есть откат через `down` не делают.
- Так как можно потерять данные, они безвозвратно удаляются.
- Откат может сломать приложение, если колонки из БД уже используются в коде.

Вместо этого делают новую миграцию, с необходимыми изменениями.

Пример файла `0001_create_users_and_tokens_tables.up.sql` миграции для создания таблиц `users` и `tokens`.

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    refresh_token TEXT NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expired_at TIMESTAMPTZ NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL
);
```

Новая миграция `0002_add_more_columns_to_tokens.up` для таблицы `tokens`.

```sql
ALTER TABLE tokens 
ADD COLUMN family_id UUID,
ADD COLUMN device VARCHAR(50),
ADD COLUMN ip TEXT, 
ADD COLUMN user_agent VARCHAR(50);
```

В результате применения миграций, создадутся 2 таблицы. Потом в таблицу `tokens` добавятся 4 новых поля.

## Реализация миграций

Мы уже установили пакет для миграций, который мы импортируем. Также мы импортируем пакеты:

```go
_ "github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
```

Они напрямую не используются, но необходимы для выполнения миграций и будут работать "за кулисами".

Функция `MigrationsUp` принимает `dsn` - подключение к БД и `mPath` - путь к директории с SQL файлами миграций, который мы определили в файле `.env`.

`main` функция загрузит наш конгфиг, и вызовет функцию для запуска миграций. Если миграции успешно выполнены, то в терминале мы увидим: 

`Migrations applied successfully ✔️`


```go
package main

import (
	"user-profiles/configs"
	"log"
	"flag"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	up := flag.Bool("up", false, "Apply all migrations")
	down := flag.Bool("down", false, "Rollback all migrations")
	flag.Parse()

	if !*up && !*down {
		log.Println("Usage:")
		log.Println("  -up     Apply migrations")
		log.Println("  -down   Rollback all migrations")
		os.Exit(1)
	}

    conf, err := configs.Load()
	if err != nil {
        log.Fatal(err)
	}
	dsn := conf.Db.Dsn
    mPath := conf.Db.MigrationsPath

	
	if *up {
		if err := MigrationsUp(dsn, mPath); err != nil {
			log.Fatal(err)
		}
	}

	if *down {
		if err := MigrationsDownAll(dsn, mPath); err != nil {
			log.Fatal(err)
		}
	}
}

func MigrationsUp(dsn string, mPath string) error {
	m, err := migrate.New(
		"file://"+mPath,
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully ✔️")
	return nil
}

func MigrationsDownAll(dsn string, mPath string) error {
    m, err := migrate.New(
        "file://"+mPath,
        dsn,
    )
    if err != nil {
        return err
    }

    if err := m.Down(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    log.Println("All migrations rolled back ✔️")
	return nil
}
```

После выполнения миграций, можно запускать приложение.
Пожалуй этих знаний достаточно для работы с БД.
