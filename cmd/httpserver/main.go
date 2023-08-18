package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "net/http/pprof" // enable to trace

	"github.com/lib/pq"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/configs"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/adapter/in/web"
	inmemory "github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/adapter/out/in-memory"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/adapter/out/postgres"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/out"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/service"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB

	personPersistenceAdapter out.PersistencePort

	personMemoryAdapter out.CachePort

	personChan chan domain.Person

	err error
)

func main() {
	cfg := configs.GetConfig()
	ctx := context.Background()

	logLevel := logger.Silent
	if cfg.EnableDebug {
		logLevel = logger.Info
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Port)
	db, err = gorm.Open(driver.New(driver.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Hour)

	personChan = make(chan domain.Person)

	personPersistenceAdapter = postgres.NewPersonPersistenceAdapter(db)
	personMemoryAdapter = inmemory.NewPersonMemoryAdapter()

	personService := service.NewPersonService(personChan, personPersistenceAdapter, personMemoryAdapter)

	personCreateHandler := web.NewPersonCreateHandler(ctx, personService)
	personGetByIdHandler := web.NewPersonGetByIdHandler(ctx, personService)
	personCountHandler := web.NewPersonCountHandler(ctx, personService)
	personSearchHandler := web.NewPersonSearchHandler(ctx, personService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/pessoas" {
			personCreateHandler.HandleCreatePerson(w, r)
			return
		}

		if r.Method == http.MethodGet && strings.Index(r.URL.Path, "/pessoas") == 0 {
			path := r.URL.Path
			var id string
			if strings.Contains(path, "/pessoas/") {
				id = path[len("/pessoas/"):]
			}

			if id != "" {
				personGetByIdHandler.HandleGetPersonById(w, r)
				return
			}

			personSearchHandler.HandleSearchPersonsByTerm(w, r)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == "/contagem-pessoas" {
			personCountHandler.HandleCountPerson(w, r)
			return
		}

		w.Header().Set("Allow", "GET,POST")
		w.WriteHeader(http.StatusNotFound)
	})

	port := fmt.Sprintf(":%s", configs.GetConfig().Server.Port)

	if err != nil {
		fmt.Println(err)
	}

	_, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("person_created")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	go waitForNotification(ctx, listener)

	fmt.Printf("Starting server on port %s...\n", port)
	go http.ListenAndServe(port, nil)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	sigs := make(chan os.Signal, 1)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	fmt.Printf("Server Shutdown [%s]...\n", sig)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	<-ctx.Done()
	close(personChan)

	fmt.Println("Server Shutdown: Done!")
}

type PersonFromPostgres struct {
	ID       string `json:"id"`
	Name     string `json:"nome"`
	Nickname string `json:"apelido"`
	Birthday string `json:"nascimento"`
	Stack    string `json:"stack"`
}

func (p *PersonFromPostgres) toDomain() domain.Person {
	var stack []string
	json.Unmarshal([]byte(p.Stack), &stack)
	person := domain.Person{
		ID:       p.ID,
		Name:     p.Name,
		Nickname: p.Nickname,
		// Birthday: p.Birthday, // TODO: return correct date
		Birthday: &time.Time{},
		Stack:    stack,
	}
	return person
}

func waitForNotification(ctx context.Context, l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			person := PersonFromPostgres{}
			err := json.Unmarshal([]byte(n.Extra), &person)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			p := person.toDomain()

			personMemoryAdapter.StorePerson(ctx, p)
			personMemoryAdapter.StoreNickname(ctx, p.Nickname, p.ID)
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection...")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}
