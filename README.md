# 🚀 DevOps Control Center & NTS Production Infrastructure

> Инфраструктурный проект Production уровня, реализующий микросервисную архитектуру, CI/CD пайплайн, балансировку нагрузки и мониторинг.

---

## 🏗️ Архитектура проекта

Проект развернут на изолированном Linux-сервере и управляется с помощью **Docker Compose**.

* **Reverse Proxy & Load Balancer (`Nginx`):** Прием входящего HTTP-трафика по 80 порту и его редиректирование между двумя репликами приложения (`devops-app-1`, `devops-app-2`) с проверкой здоровья (`healthcheck`).
* **Backend application (`Go`):** Высокопроизводительное веб-приложение на языке Go, которое отдает ответы в JSON-формате, эндпоинты для проверки работы и метрики Prometheus (`/metrics`).
* **Monitoring & Observability (`Prometheus` + `Grafana` + `Node Exporter`):**
  * **Prometheus** выполняет автоматический сбор метрик приложения и системных метрик хоста.
  * **Node Exporter** собирает детализированную телеметрию аппаратных ресурсов сервера (CPU, RAM, Disk, Network).
  * **Grafana** визуализирует собранные данные в реальном времени.
* **CI/CD pipeline (`GitHub Actions` + `Self-Hosted Runner`):** Автоматизированная доставка кода по модели GitOps.

---

## 🛠️ Технологический стек

* **Язык программирования:** Go (Golang)
* **Контейнеризация:** Docker, Docker Compose
* **Веб-сервер / Балансировщик:** Nginx (Alpine)
* **Системы мониторинга:** Prometheus, Grafana, Node Exporter
* **CI/CD и автоматизация:** GitHub Actions (Self-Hosted Runner)

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
