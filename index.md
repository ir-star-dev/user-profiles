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


=========================================================================

## Что такое JWT?

JWT (JSON Web Token) — это строка, которая содержит закодированную (в base64) информацию о пользователе и подписана сервером. 
Используется для аутентификации и передачи данных между сервисами в виде JSON-объекта.

## Основное применение

- Аутентификация
1. Ты вводишь логин и пароль в приложение
2. Система проверяет, совпадают ли они с тем, что хранится в базе
3. Если да — тебя "пускают внутрь"
4. Сервер создаёт JWT и отправляет его клиенту
5. Клиент хранит токен (например, в localStorage или cookie)
6. Этот токен используется для доступа к API:

```go
    Authorization: Bearer <JWT>
```                
Сервер проверяет лишь подпись токена и не лезет в базу.

## Когда использовать JWT?

- Используем, если:
- делаем REST API,
- есть фронтенд отдельно (React/Vue),
- микросервисная архитектура,
- есть мобильное приложение,
- планируется масштабирование.

❗Не лучший вариант, если:
- нужна мгновенная инвалидизация сессий,
- сложная система ролей/доступов с частыми изменениями.

### Преимущества

- Stateless
- сервер не хранит состояние о сессиях или токенах,
- сервер не хранит список активных токенов,
- любой сервер может проверить токен только по подписи и содержимому payload.

- Быстро:
- не нужно каждый раз ходить в БД, вся информация уже внутри токена.

- Безопасность (при правильном использовании):
- токен подписан сервером → нельзя подделать,
- можно задать срок жизни токена (expiration),
- HTTPS обязателен для передачи токена, иначе его могут перехватить.

- Удобен для API:
- идеально для REST / GraphQL,
- широко используется в SPA и мобильных приложениях. 


## Из чего состоит JWT токен?

### Header

Всегда содержит 2 обязательных поля:
`type` - тип токена (например, JWT)
`alg` - алгоритм подписи (например, HS256)

Существуют:
Симметричные алгоритмы - `HS256/384/512`. 
То есть для подписи и проверки подписи используется один и тот же секретный ключ. 
Используют чаще всего в простых API / монолитах.
- Плюсы: быстрый, простой.
- Минусы: секретный ключ должен быть у всех серверов ❗
    
Ассиметричные алгоритмы - `RSA, ECDSA`. 
Приватный ключ подписывает токен, публичный проверяет подпись.
Используют: микросервисы, кластер баз данных, CDN, Kubernetes и пр.
- Плюсы: можно раздавать public key всем сервисам, private key хранится в одном месте.
- Минусы: медленнее, сложнее.
        
Мы будем в проекте использовать симметричный алгоритм - `HS256`.
        
### Payload - тело токена или claims

Сюда можно положить любые данные, которые серверу нужно читать для аутентификации или авторизации.

payload токена часто используется для двух разных задач одновременно:
1. 🔐 Аутентификация (кто ты)
Payload может содержать данные, которые помогают понять личность пользователя:
- user_id
- email
- username

2. 🛂 Авторизация (что тебе можно)
Payload также может содержать права доступа:
- role: admin
- permissions: ["read", "write"]
- is_premium: true

- Стандартные клэймы (рекомендованные)
| Claim | Описание                                 |
| ----- | ---------------------------------------- |
| `iss` | Issuer — кто выпустил токен              |
| `sub` | Subject — идентификатор пользователя     |
| `aud` | Audience — для кого токен                |
| `exp` | Expiration — время жизни токена          |
| `nbf` | Not before — токен действителен с…       |
| `iat` | Issued at — время выпуска                |
| `jti` | JWT ID — уникальный идентификатор токена |

- Кастомные клэймы: user_id, role, email, permissions. Любая информация, которая нужна серверу для проверки прав.

⚠️ Важные нюансы:
- Нельзя хранить конфиденциальные данные: пароли, пин-коды, кредитные карты и т.д. payload читается любым base64.
- payload должен быть стабильным, т.е. не стоит класть туда данные, которые часто меняются. Если данные изменились, то токен устарел.
- payload не должен быть слишком большим. Рекомендация: < 1–2 KB. Причина: токен передаётся в заголовке HTTP → большие токены замедляют каждый запрос.

### Signature - подпись

Создаётся сервером с помощью секретного ключа. Гарантирует, что токен не подделан.
❗Секрет никогда не хранят в коде.

Для симметричных алгоритмов `HS256/384/512` размер ключа "256/384/512 бит". 
Секрет должен быть достаточно длинным, иначе подпись будет уязвима.

Мы будем использовать в проекте алгоритм HS256, поэтому секрет (64 байта) для него можно сгенерировать простой функцией:

```go
package key

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func Generate() {
    key := make([]byte, 64)
    _, err := rand.Read(key)
    if err != nil {
        panic(err)
    }

    encoded := base64.StdEncoding.EncodeToString(key)
    fmt.Println(encoded)
}
```

Для ассиметричный алгоритмов это пара ключей: private.pem, public.pem. 
Для RSA обычно 2048–4096 бит (256–512 байт), для ECDSA меньше (например, ES256 → 256 бит / 32 байта).


## Как работать в продакшене с JWT?

### Access токен

JWT токен еще называют `access` токеном.
- Обычно stateless → Сервер не хранит его. Любой сервер может проверить токен сам, без хранения состояния.
- Проверка идёт только по подписи токена (Signature) и по payload (например, user_id, role, exp).
- Содержит минимальную информацию о пользователе.
- Используется для проверки прав на каждом запросе.
- Хранится на фронтенде.

Минусы:
- Короткоживущий (5-15 минут).
- Если украли access token → нельзя просто "удалить его с сервера".
- Данные API уязвимы, так как доступ через access токен закроется, только когда закончится срок жизни токена.
- После истечения срока жизни токена, пользователю нужно снова входить в систему (т.е. каждый 5-15 минут).

В проекте, для генерации JWT будем использовать библиотеку:

```go
go get github.com/golang-jwt/jwt/v5
```

Создадим пакет `internal/infrastructure/security/jwt.go`

```go
package security

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
    Create(uId int) (string, error)
    Parse(token string) (jwt.MapClaims, error)
}

type JWT struct {
    Secret string
}

func NewJWTService(secret string) JWTService {
    return &JWT{
        Secret: secret,
    }
}

func (j *JWT) Create(uId int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iat": jwt.NewNumericDate(time.Now()),
        "exp": jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
        "sub": uId,
    })

    s, err := token.SignedString([]byte(j.Secret))
    if err != nil {
        return "", err
    }
    return s, nil
}

func (j *JWT) Parse(tokenStr string) (jwt.MapClaims, error) {
    claims := jwt.MapClaims{}
    t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
        return []byte(j.Secret), nil
    })
    if err != nil {
        return nil, err
    }
    if !t.Valid {
        return nil, errors.New("Invalid token")
    }
    return claims, nil
}
```

В пакете есть интерфейс с 2 методами: `Create`(создает токен) и `Parse`(проверяет подпись).
Структура, которая зависит от секретного ключа. И функция констуктор.

### Refresh токен 

Практически везде используют JWT (access) + refresh token, особенно для SPA, мобильных приложений и микросервисов:
- REST API
- GraphQL API
- Мобильные приложения

Связка:
- Безопаснее: если access токен украли, злоумышленник может использовать только пока он жив.
- Удобнее для пользователя: не нужно логиниться каждые 5-15 минут.
- Легко обновлять access token автоматически.
    
Только `access` token без `refresh` token встречается редко → только для очень простых внутренних сервисов.

`Refresh` токен — это специальный токен, который используется для получения нового access токена, когда тот истёк.
Он долгоживущий, от нескольких дней до нескольких недель.

Простыми словами:
- Access token → "билет на вход" (короткоживущий).
- Refresh token → "абонемент, по которому тебе дают новый билет".

#### Как выглядит Refresh токен?

По сути это длинная случайная строка, не обязательно в JWT.
1. Opaque token (самый частый вариант). Просто случайная строка: a8f9c2d1e7b4f6...
2. Иногда JWT refresh token. Используется реже, так как сложнее отзывать.

В проекте мы создадим пакет `internal/infrastructure/security/refresh.go`Ж

```go
package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type RefreshTokenService interface {
	Generate() (string, error)
	Hash(token string) []byte
}

type refreshTokenService struct{}

func NewRefreshTokenService() RefreshTokenService {
	return &refreshTokenService{}
}

func (r *refreshTokenService) Generate() (string, error) {
	b := make([]byte, 32) // 256 bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *refreshTokenService) Hash(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}
```
Этот пакет будет генерировать случайную 32 байтовую строку и создавать хэш этой строки.

#### Как работает?

1. Клиент авторизуется.
2. Сервер выдает access токен на 10 минут.
3. Генерирует refresh токен и записывает его хэш в БД, а сам токен, как HttpOnly cookie.
4. Клиент через 10 минут отправляет запрос с access токеном, сервер отвечает "401 Unauthorized".
5. Клиент отправляет запрос на получения нового access токена, например /refresh.
6. Сервер берет куку, в которой записан refresh токен, проверяет его, если все ок - то выдает клиенту новый access токен.
    
#### Безопастность

Хэш `refresh` токена нужно хранить на сервере (БД или Redis). В продакшене обычно используют Redis.
Это нужно чтобы:
- можно было его отозвать при необходимости → удаляют его из базы или помечают как недействительный,
- проверить, что токен не украден или не истёк.

Минус `refresh` токена:
❗Если украли refresh token — злоумышленник может “жить” в системе долго.

На продакшене, лучшей практикой является система - использовать `rotation refresh tokens`:
- каждый раз выдаётся новый `refresh` token,
- старый инвалидируется.

То есть процесс выдачи `access` токена, описанный выше, будет такой:

1. При каждом обновлении access token, сервер проверяет его `refresh`:
- убеждается, что он валидный и не отозван, 
- старый `refresh` token инвалидируется (revoked),
- создаётся новый `refresh` token,
- создаётся новый `access` token.

Клиент получает новую пару токенов.

🔥 Best Practice
- Refresh Token Rotation
- Хранение в БД (или Redis)
- Поле revoked
- Поле expires_at
- Связка с устройством (device_id)

В проекте мы как раз реализуем хранение в БД, с полями revoked и expires_at. А также Refresh Token Rotation.


