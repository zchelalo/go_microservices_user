package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zchelalo/go_microservices_user/internal/user"
	"github.com/zchelalo/go_microservices_user/pkg/bootstrap"
)

func main() {
	logger := bootstrap.InitLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("error loading .env file")
	}

	db, err := bootstrap.DBConnection()
	if err != nil {
		logger.Fatal(err)
	}

	router := http.NewServeMux()

	pageLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pageLimDef == "" {
		logger.Fatal("paginator limit default is required")
	}

	config := user.Config{
		LimPageDef: pageLimDef,
	}

	userRepository := user.NewRepository(logger, db)
	userService := user.NewService(logger, userRepository)
	userEndpoints := user.MakeEndpoints(userService, config)

	router.HandleFunc("GET /users", userEndpoints.GetAll)
	router.HandleFunc("GET /users/{id}", userEndpoints.Get)
	router.HandleFunc("POST /users", userEndpoints.Create)
	router.HandleFunc("PATCH /users/{id}", userEndpoints.Update)
	router.HandleFunc("DELETE /users/{id}", userEndpoints.Delete)

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := &http.Server{
		// Handler:      http.TimeoutHandler(router, 5*time.Second, "Timeout!"),
		Handler:      router,
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
