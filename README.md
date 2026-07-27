# 🚀 DevOps Control Center & NTS Production Infrastructure

> Комплексный демонстрационный проект инфраструктуры уровня Production, реализующий микросервисную архитектуру, автоматизированный CI/CD пайплайн, балансировку нагрузки и полноценный стек мониторинга.

---

## 🏗️ Архитектура проекта

Проект развернут на изолированном Linux-сервере и управляется с помощью **Docker Compose**. 

* **Reverse Proxy & Load Balancing (`Nginx`):** Принимает входящий HTTP-трафик на 80-м порту и распределяет его между двумя репликами бэкенда (`devops-app-1`, `devops-app-2`) с учетом проверок здоровья (`healthcheck`).
* **Backend Application (`Go`):** Высокопроизводительное веб-приложение на Go, отдающее JSON-ответы, эндпоинты проверки работоспособности и метрики Prometheus (`/metrics`).
* **Monitoring & Observability (`Prometheus` + `Grafana` + `Node Exporter`):**
  * **Prometheus** производит автоматический сбор метрик приложения и системных метрик хоста.
  * **Node Exporter** собирает детальную телеметрию аппаратных ресурсов сервера (CPU, RAM, Disk, Network).
  * **Grafana** визуализирует собранные данные в режиме реального времени.
* **CI/CD Pipeline (`GitHub Actions` + `Self-Hosted Runner`):** Автоматизированная доставка кода по модели GitOps.

---

## 🛠️ Технологический стек

* **Язык программирования:** Go (Golang)
* **Контейнеризация:** Docker, Docker Compose
* **Веб-сервер / Балансировщик:** Nginx (Alpine)
* **Системы мониторинга:** Prometheus, Grafana, Node Exporter
* **CI/CD и Автоматизация:** GitHub Actions (Self-Hosted Runner)

---

## 📂 Структура репозитория

```text
├── .github/
│   └── workflows/
│       └── deploy.yml       # Production CI/CD пайплайн (тесты, сборка, деплой)
├── cmd/
│   └── server/
│       ├── main.go          # Исходный код Go-приложения и веб-интерфейса
│       └── main_test.go     # Юнит-тесты эндпоинтов
├── docker-compose.yml       # Оркестрация всех сервисов стека
├── Dockerfile               # Многоэтапная (multi-stage) сборка приложения
├── nginx.conf               # Конфигурация Nginx и балансировки
├── prometheus.yml           # Конфигурация сбора метрик Prometheus
└── go.mod                   # Зависимости Go
