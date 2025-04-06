package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toramanomer/passwd-auth-go/core/emailverification"
	"github.com/toramanomer/passwd-auth-go/core/mailer"
	"github.com/toramanomer/passwd-auth-go/core/repository"
	"github.com/toramanomer/passwd-auth-go/handlers"
)

func init() {
	filename := ".env"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open %s file: %v\n", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name, value, found := strings.Cut(line, "=")
		if !found {
			log.Fatalf("Invalid line %d in file %s: %s\n", lineNumber, filename, line)
		}

		if err := os.Setenv(name, value); err != nil {
			log.Fatalf("Failed to set %s environment variable: %v\n", name, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning %s file: %v\n", filename, err)
	}
}

func main() {
	var (
		connString = os.Getenv("POSTGRESQL_CONNECTION")
		serverPort = os.Getenv("SERVER_PORT")
		serverAddr = fmt.Sprintf(":%s", serverPort)
	)

	var db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var (
		mux    = http.NewServeMux()
		server = &http.Server{
			Addr:              serverAddr,
			Handler:           mux,
			ReadTimeout:       3 * time.Second,
			WriteTimeout:      3 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
		}
	)

	var (
		userManagementRepo = repository.NewUserManagementRepository(db)
		mailer             = mailer.NewMailer()
		evStrategy         = emailverification.NewEmailVerificationStrategy()
	)

	signupHandler := &handlers.SignupController{
		UserManagementRepo:        userManagementRepo,
		EmailVerificationStrategy: evStrategy,
		Mailer:                    mailer,
	}
	mux.Handle("/signup", signupHandler)

	signoutHandler := &handlers.SignoutController{
		UserManagementRepo: userManagementRepo,
	}
	mux.Handle("/signout", signoutHandler)

	verifyHandler := &handlers.VerifyController{
		UserManagementRepo: userManagementRepo,
	}
	mux.Handle("/verify", verifyHandler)

	log.Fatalln(server.ListenAndServe())
}
