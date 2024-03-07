package main

import (
	"dz_go/handlers"
	"dz_go/proxy"
	"dz_go/storage"
	"flag"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	proxyPort   string
	targetPorts string
)

func init() {
	flag.StringVar(&proxyPort, "proxyPort", "", "Port for the proxy server")
	flag.StringVar(&targetPorts, "targetPort", "", "Comma-separated list of ports for the replicas")
	flag.Parse()
}

func main() {
	// Common SQLite storage for all replicas
	s, err := storage.NewSQLiteStorage("_files/users.db")
	if err != nil {
		log.Fatal("Failed to create storage:", err)
	}

	// Запуск прокси-сервера, если задан proxyPort и список портов для реплик
	if proxyPort != "" && targetPorts != "" {
		go proxy.StartProxyServer(proxyPort, targetPorts)
	}

	// Если задан список портов для реплик, запускаем реплики
	if targetPorts != "" {
		ports := strings.Split(targetPorts, ",")
		var wg sync.WaitGroup

		for _, port := range ports {
			wg.Add(1)
			go func(port string) {
				defer wg.Done()
				port = strings.TrimSpace(port)
				startApp(port, s)
			}(port)
		}

		wg.Wait() // Ожидаем завершения всех реплик
	}
}

func startApp(port string, storage storage.Storage) {
	h := handlers.Handlers{Storage: storage}
	r := chi.NewRouter()

	r.Post("/create", h.CreateUser)
	r.Post("/make_friends", h.MakeFriends)
	r.Get("/friends/{userID}", h.ListFriends)
	r.Get("/user/{userID}", h.GetUser)
	r.Put("/user/{userID}", h.UpdateUser)
	r.Delete("/user", h.DeleteUser)

	log.Printf("Starting application server on :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
