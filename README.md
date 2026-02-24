# Contact Manager

Полноценное веб-приложение для управления контактами.

## Основные возможности

- Регистрация и авторизация пользователей
- Добавление, редактирование, удаление, поиск контактов
- Редактирование профиля (имя, email)
- Современный тёмный интерфейс
- Полностью запускается одной командой

![Логин](screenshots/login.png)
![Регистрация](screenshots/new_account.png)
![Контакты](screenshots/saved_contacts.png)
![Новый_Контакт](screenshots/add_new_contact.png)
![Профиль](screenshots/my_profile.png)

## Стек технологий

- **Frontend** — Vue 3 + Vite + Tailwind CSS + DaisyUI (тёмный стиль)
- **Backend** — Go (Gin + GORM + JWT)
- **База данных** — PostgreSQL
- **Контейнеризация** — Docker + Docker Compose

## Как запустить локально

```bash
git clone https://github.com/yngh0m1e/Contact-manager.git
cd Contact-manager

docker compose up --build