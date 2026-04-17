
# SportShop – интернет-магазин спортивных костюмов

Полноценное веб-приложение с фронтендом на **Vue 3** и бэкендом на **Go**.  
Поддерживает регистрацию, авторизацию по сессиям с токенами, корзину, избранное, оформление заказов, статистику пользователя, фильтрацию товаров и платёжную страницу.

## 🚀 Технологии

### Backend
- Go 1.22
- gorilla/mux – маршрутизация
- PostgreSQL (через lib/pq)
- bcrypt – хеширование паролей
- Сессии на основе токенов, хранящихся в БД (с автоматическим продлением)

### Frontend
- Vue 3 (Composition API)
- Vue Router
- Axios (с `withCredentials`)
- Tailwind CSS
- Vite

## 📦 Установка и запуск

### Требования
- Docker и Docker Compose
- Git

### Шаги

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/your-username/sportshop.git
   cd sportshop
   ```

2. **Запустите через Docker Compose**
   ```bash
   docker-compose up --build
   ```

   После успешной сборки:
   - Фронтенд: `http://localhost:5173`
   - Бэкенд API: `http://localhost:8080/api`
   - PostgreSQL: `localhost:5432` (логин `sportuser`, пароль `sportpass`, БД `sportshop`)

3. **Остановка**
   ```bash
   docker-compose down
   ```

### Альтернативный ручной запуск (без Docker)

#### Backend (Go)
```bash
cd backend
go mod tidy
go run main.go
# сервер на http://localhost:8080
```

#### Frontend (Vue)
```bash
cd frontend
npm install
npm run dev
# приложение на http://localhost:5173
```

#### База данных (PostgreSQL)
Убедитесь, что PostgreSQL запущен локально и создана БД `sportshop` с пользователем `sportuser`.  
Выполните скрипт `init.sql` для создания таблиц и наполнения товарами.

## ⚙️ Переменные окружения

### Backend (в Docker или локально)
| Переменная     | Значение по умолчанию |
|----------------|----------------------|
| DB_HOST        | postgres             |
| DB_PORT        | 5432                 |
| DB_USER        | sportuser            |
| DB_PASSWORD    | sportpass            |
| DB_NAME        | sportshop            |
| PORT           | 8080                 |
| SESSION_KEY    | supersecretkey123    |

### Frontend
В файле `frontend/.env` (создайте при необходимости):
```
VITE_API_URL=http://localhost:8080/api
```

## 📁 Структура проекта

```
SportStore_Vue_Go/
├── backend/
│   ├── db/                # инициализация БД и работа с сессиями
│   ├── handlers/          # обработчики API (auth, cart, favorites, orders, products, profile, statistics)
│   ├── middleware/        # AuthRequired, GetUserID
│   ├── models/            # структуры данных
│   ├── utils/             # генерация токенов
│   ├── go.mod, go.sum
│   ├── main.go
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/    # Header, Card, CardList, CartDrawer, CartItem...
│   │   ├── pages/         # Home, Favorites, Login, Register, Profile, Payment
│   │   ├── utils/         # API-вызовы (auth, cart, favorites, orders, products, statistics)
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css
│   ├── public/            # иконки, logo.png
│   ├── package.json
│   ├── vite.config.js
│   └── Dockerfile
├── init.sql               # схема БД и тестовые товары
├── docker-compose.yml
├── nginx.conf (опционально)
└── README.md
```

## 🔌 API Endpoints

### Публичные
| Метод | Эндпоинт                | Описание                   |
|-------|-------------------------|----------------------------|
| POST  | /api/register           | Регистрация пользователя   |
| POST  | /api/login              | Вход (установка cookie)    |
| POST  | /api/logout             | Выход (удаление сессии)    |
| GET   | /api/products           | Список товаров (фильтры)   |
| GET   | /api/products/{id}      | Детали товара              |

### Защищённые (требуют cookie session_token)
| Метод   | Эндпоинт                    | Описание                       |
|---------|-----------------------------|--------------------------------|
| GET     | /api/profile                | Профиль пользователя           |
| PUT     | /api/profile                | Обновление профиля             |
| GET     | /api/favorites              | Список избранного              |
| POST    | /api/favorites              | Добавить в избранное           |
| DELETE  | /api/favorites/{productId}  | Удалить из избранного          |
| GET     | /api/cart                   | Товары в корзине               |
| POST    | /api/cart                   | Добавить товар                 |
| PUT     | /api/cart/{itemId}          | Изменить количество            |
| DELETE  | /api/cart/{itemId}          | Удалить товар из корзины       |
| GET     | /api/orders                 | История заказов                |
| POST    | /api/orders                 | Создать заказ                  |
| GET     | /api/statistics             | Статистика пользователя        |

## 🧪 Тестовые данные

В `init.sql` уже добавлено 17 товаров (мужские, женские, детские костюмы) с разными размерами, ценами и изображениями.  
При первом запуске через Docker Compose база инициализируется автоматически.

## 📝 Примечания

- Сессии хранятся в таблице `sessions`, токен генерируется случайно (32 байта), срок жизни 15 минут, продлевается при каждом запросе к защищённому маршруту.
- Пароли хешируются с помощью bcrypt.
- Корзина и избранное привязаны к пользователю.
- Страница оплаты (`/payment`) имитирует ввод данных карты и после успеха показывает сообщение и перенаправляет на главную.
- Валидация пароля при регистрации: минимум 8 символов, заглавная буква, цифра, спецсимвол.

## 🐛 Возможные проблемы и решения

- **Ошибка 401 при запросах к API**  
  Убедитесь, что cookie `session_token` установлена (смотреть вкладка Application в браузере) и что запросы отправляются с `withCredentials: true` (в проекте настроено).

- **Фронтенд не видит бэкенд при Docker**  
  Проверьте, что в `frontend/.env` указан `VITE_API_URL=http://localhost:8080/api` (или `http://backend:8080/api`, если используется nginx). При `docker-compose` фронтенд использует `VITE_API_URL` из Docker-окружения.

- **Сборка бэкенда падает с ошибками импорта**  
  Убедитесь, что `go.mod` и `go.sum` синхронизированы (`go mod tidy` внутри контейнера уже выполняется).

