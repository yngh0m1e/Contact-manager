FROM node:20-alpine AS builder
WORKDIR /app

# Копируем package.json и устанавливаем зависимости
COPY package*.json ./
RUN npm install          # ← ИЗМЕНИЛИ с ci на install

# Копируем весь код и собираем проект
COPY . .
RUN npm run build

# Финальная стадия (маленькая)
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80