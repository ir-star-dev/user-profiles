# Frontend и Go
## API + SPA (самый распространённый способ)
Стек:
Backend: Go (Gin / Fiber / Echo)
Frontend: React / Vue.js / Angular

Где используется:
- SaaS продукты
- админки
- дашборды
- сложные интерфейсы
Плюсы:
- масштабируемость
- богатый UI
- независимые команды
Минусы
- 2 проекта
- сложный деплой
- CORS, auth сложнее

Реальность
👉 Это де-факто стандарт индустрии

## Go + SSR (html/template)
Без JS или почти без него.

Где используется:
- внутренние инструменты
- простые сервисы
- MVP
Плюсы:
- простота
- один бинарник
- минимум зависимостей
Минусы
- слабая интерактивность
- тяжело масштабировать UI

Реальность
👉 Используется, но реже в современных продуктах

## Go + HTMX (новый тренд)
"SSR + интерактивность без SPA"

Где используется:
- стартапы
внутренние панели
- CRUD-приложения
- admin UI
Плюсы:
- почти нет JS
- быстрее разработка
- проще, чем SPA
Минусы:
- не для сложных интерфейсов
- меньше экосистема

Реальность
👉 очень быстро набирает популярность
👉 особенно среди Go-разработчиков

## Hybrid (SSR + SPA части)
Самый “умный” вариант

Пример:
- страница профиля → SSR / HTMX
- админ-дашборд → React
Где используется:
- средние/крупные проекты

Реальность
👉 часто в проде, но сложнее в реализации

## Static frontend + embed в Go
Собираешь фронт и встраиваешь

Где используется:
- небольшие сервисы
- self-hosted приложения
Плюсы:
- один бинарник
- удобно деплоить
Минусы:
- билд фронта отдельно

# Резюме в одной строке
👉 Если большой продукт → SPA + API
👉 Если простой/средний → HTMX + Go — очень сильный выбор

# Что такое HTMX?
HTMX — это маленькая JS-библиотека, которая позволяет:
- делать HTTP-запросы прямо из HTML
- и вставлять HTML-ответ в страницу

## Как это выглядит?
Без HTMX (чистый JS)
```js
fetch("/api/profile")
  .then(res => res.json())
  .then(data => {
    document.getElementById("result").innerText = data.email
  })
```

С HTMX
```html
<button hx-get="/hx/profile" hx-target="#result">
  Load Profile
</button>

<div id="result"></div>
```

## Модель HTMX

HTMX строится на 3 идеях:

1. HTML — это API
Ты описываешь поведение прямо в HTML:
```html
<button hx-get="/users">
```

2. Сервер возвращает HTML, а не JSON
```html
<tr>
  <td>user@mail.com</td>
</tr>
```

3. Частичное обновление DOM
обновляется не вся страница, а кусок

4. Триггеры (очень мощно)
Ты можешь управлять, когда отправлять запрос:
🔹 Базовые
```html
hx-trigger="click"
hx-trigger="load"
```
🔹 Продвинутые
```html
hx-trigger="keyup changed delay:500ms"
```
🔹 Несколько триггеров
```html
hx-trigger="click, keyup"
```
🔹 Polling (реалтайм)
```html
hx-trigger="every 5s"
```
👉 автообновление данных

4. Управление DOM (очень гибко)

🔹 Куда вставлять (hx-target)
```html
hx-target="#result"
```
🔹 Как вставлять (hx-swap)
```html
hx-swap="innerHTML"
hx-swap="outerHTML"
hx-swap="beforeend"
hx-swap="afterbegin"
```

5. Работа с формами (очень удобно)

HTMX автоматически:
- сериализует форму
- отправляет её
- обрабатывает ответ
Пример:
```html
<form hx-post="/hx/login">
```
Валидация с сервера:
```go
if invalid {
    fmt.Fprint(w, "<p>Error</p>")
}
```
👉 сразу отображается в UI

6. Загрузка при открытии страницы
```html
<div hx-get="/hx/profile" hx-trigger="load">
```

7. Индикаторы загрузки
```html
<button hx-get="/hx/data" hx-indicator="#loader">

<div id="loader" style="display:none;">Loading...</div>
```
👉 HTMX сам показывает/скрывает

8. Работа с history (как SPA!)
```html
hx-push-url="true"
```
👉 URL меняется без перезагрузки

9. События (очень мощно)
HTMX генерирует события:
```html
document.body.addEventListener("htmx:afterRequest", ...)
Примеры событий:
htmx:beforeRequest
htmx:afterRequest
htmx:responseError
```
👉 можно подключать JS при необходимости

10. Out-of-band updates (🔥 продвинутая фича)
👉 сервер может обновить НЕ только target
```html
HTML:
<div id="notifications"></div>
```
Ответ сервера:
```html
<div id="notifications" hx-swap-oob="true">
  New message!
</div>
```
👉 обновится в любом месте страницы

11. SSE и WebSocket поддержка
HTMX умеет:
```
Server-Sent Events
WebSockets
```
👉 почти realtime без фронтенд-фреймворка

12. Lazy loading
```html
<div hx-get="/hx/data" hx-trigger="revealed">
```
👉 загружается при прокрутке

13. Кэширование
hx-cache="true"
👉 уменьшает количество запросов

14. Композиция интерфейса
Ты можешь строить UI из кусочков:
/navbar
/profile
/users/list
/users/row

👉 как компоненты, но на сервере

15. Интеграция с Go templates
```html
{{ define "userRow" }}
<tr>
  <td>{{ .Email }}</td>
</tr>
{{ end }}
```
👉 идеально ложится на HTMX


# Документация по HTMX
```
https://htmx.org/reference/
```

# Как работает под капотом
1. Пользователь нажимает кнопку
2. HTMX делает HTTP-запрос
3. Go возвращает HTML
4. HTMX вставляет HTML в DOM
```go
click → request → HTML → replace DOM
```

# Ограничения HTMX

Важно понимать, где он НЕ подходит:
Сложный UI
- drag & drop
- графики
- сложные таблицы
Большие SPA
- много состояния
- клиентская логика

# Go templates

В Go есть встраенный пакет для работы с html "html/template"
Чтобы собирать страницы с помощью него нужно:
1. Иметь страницу каркас:
```html
base.html

<!doctype html>
<html>
<head>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js" integrity="sha384-/TgkGk7p307TH7EXJDuUlgG3Ce1UVolAOFopFekQkkXihi5u/6OCvVKyz1W+idaz" crossorigin="anonymous"></script>
</head>

<body id="container">

  {{ block "content" . }}{{ end }}

</body>
</html>
```
2. Часть страницы, которую надо вставить в каркас
```html
register.html

{{ define "content" }}

<h2>Register</h2>
<div id="result"></div>

<form hx-post="/store">
    <div class="mb-3">
        <label for="name" class="form-label">Name</label>
        <input type="text" class="form-control" id="name" required>
    </div>
    <div class="mb-3">
        <select class="form-select" aria-label="Default select example" required>
            <option selected>Role</option>
            <option value="user">User</option>
            <option value="admin">Admin</option>
        </select>
    </div>
    <div class="mb-3">
        <label for="email" class="form-label">Email</label>
        <input type="email" class="form-control" id="email" required>
    </div>
    <div class="mb-3">
        <label for="password" class="form-label">Password</label>
        <input type="password" id="password" class="form-control" aria-describedby="passwordHelpBlock" required>
        <div id="passwordHelpBlock" class="form-text">
            Your password must be 5-12 characters long.
        </div>
    </div>
    <button type="submit" class="btn btn-primary">Register</button>
</form>

{{end}}

```
Части страница вставляются через специальный синтаксис Go
```go
{{ block "content" . }}{{ end }}
```
block = "место, куда вставится контент"

`Точка` обязательна. Это передача контента внутрь блока.

```go
{{ define "content" }}
...
{{ end }}
```
define = "контент, который вставится на место в block"

👉 block и define работают только вместе
👉 это механизм “layout + content”