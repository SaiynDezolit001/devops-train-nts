package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>DevOps Control Center - NTS</title>
    <style>
        :root {
            --bg-color: #0f172a;
            --card-bg: #1e293b;
            --accent-green: #22c55e;
            --accent-blue: #38bdf8;
            --text-main: #f8fafc;
            --text-muted: #94a3b8;
            --border-color: #334155;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: var(--bg-color);
            color: var(--text-main);
            margin: 0;
            padding: 20px;
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        .container {
            width: 100%;
            max-width: 900px;
        }
        header {
            text-align: center;
            margin-bottom: 30px;
        }
        header h1 {
            color: var(--accent-blue);
            margin-bottom: 5px;
        }
        header p {
            color: var(--text-muted);
        }
        .grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
        }
        @media (max-width: 768px) {
            .grid { grid-template-columns: 1fr; }
        }
        .card {
            background-color: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }
        .card h2 {
            margin-top: 0;
            font-size: 1.2rem;
            border-bottom: 1px solid var(--border-color);
            padding-bottom: 10px;
        }
        button {
            background-color: var(--accent-blue);
            color: #0f172a;
            border: none;
            padding: 10px 16px;
            font-weight: bold;
            border-radius: 6px;
            cursor: pointer;
            transition: opacity 0.2s;
            margin: 5px 5px 5px 0;
        }
        button:hover {
            opacity: 0.9;
        }
        button.secondary {
            background-color: #475569;
            color: var(--text-main);
        }
        .output {
            background-color: #020617;
            border: 1px solid var(--border-color);
            border-radius: 6px;
            padding: 12px;
            font-family: 'Courier New', Courier, monospace;
            font-size: 0.9rem;
            color: var(--accent-green);
            min-height: 80px;
            max-height: 150px;
            overflow-y: auto;
            margin-top: 15px;
            white-space: pre-wrap;
        }
        .links a {
            display: inline-block;
            color: var(--accent-blue);
            text-decoration: none;
            margin-right: 15px;
            margin-top: 10px;
        }
        .links a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🚀 DevOps Control Center</h1>
            <p>Production Infrastructure Dashboard & CI/CD Demo</p>
        </header>

        <div class="grid">
            <!-- Секция тестирования API -->
            <div class="card">
                <h2>API Endpoints</h2>
                <p>Проверка работоспособности сервисов:</p>
                <button onclick="testEndpoint('/health')">GET /health</button>
                <button onclick="testEndpoint('/api/v1/hello')">GET /api/v1/hello</button>
                <button class="secondary" onclick="testEndpoint('/metrics')">GET /metrics</button>
                
                <div class="output" id="api-output">Нажмите кнопку для теста...</div>
            </div>

            <!-- Секция мониторинга и инфраструктуры -->
            <div class="card">
                <h2>Infrastructure & Links</h2>
                <p>Внешние панели мониторинга и инструменты:</p>
                <div class="links">
                    <a href="/metrics" target="_blank">📊 Prometheus Metrics</a>
                    <a href="https://hub.docker.com/r/justsaiyn/devops-app" target="_blank">🐳 Docker Hub</a>
                </div>
                <hr style="border: 0; border-top: 1px solid var(--border-color); margin: 15px 0;">
                <p style="font-size: 0.9rem; color: var(--text-muted);">
                    <strong>Status:</strong> <span style="color: var(--accent-green);">● Operational (Zero-Downtime)</span><br>
                    <strong>CI/CD:</strong> GitHub Actions + Self-Hosted Runner
                </p>
            </div>
        </div>
    </div>

    <script>
        async function testEndpoint(url) {
            const output = document.getElementById('api-output');
            output.textContent = "Запрос отправлен...";
            try {
                const response = await fetch(url);
                const text = await response.text();
                // Пробуем отформатировать если это JSON
                try {
                    const json = JSON.parse(text);
                    output.textContent = JSON.stringify(json, null, 2);
                } catch {
                    output.textContent = text;
                }
            } catch (error) {
                output.textContent = "Ошибка запроса: " + error.message;
            }
        }
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	mux := http.NewServeMux()

	// Главная страница - Дашборд
	mux.HandleFunc("/", dashboardHandler)

	// Технические эндпоинты
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"message":  "Hello, CI/CD is working smoothly!",
			"hostname": hostname,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w.Encode(response))
	})

	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("Starting server", "port", "8080", "hostname", hostname)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")
	ctx, cancel := contextWithTimeout()
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server stopped")
}

func contextWithTimeout() (time.Duration, func()) {
	// Вспомогательная заглушка для таймаута контекста
	return 5 * time.Second, func() {}
}
