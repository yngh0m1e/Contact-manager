# Contact Manager

Полноценное веб-приложение для управления контактами.

# Функции

Регистрация / авторизация (JWT)
CRUD контактов + поиск
Редактирование профиля
Красивый тёмный интерфейс

![Логин](screenshots/login.png)
![Регистрация](screenshots/new_account.png)
![Контакты](screenshots/saved_contacts.png)
![Новый_Контакт](screenshots/add_new_contact.png)
![Профиль](screenshots/my_profile.png)

## Стек
- Frontend: Vue 3 + Vite + Tailwind CSS + DaisyUI (тёмный дизайн)
- Backend: Go + Gin + GORM
- База: PostgreSQL
- Всё в Docker Compose

## Запуск

```bash
docker compose up --build