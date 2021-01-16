package app

import (
	"banking/domain"
	"banking/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func checkVariables() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	if dbUser == "" ||
		dbPassword == "" ||
		dbAddress == "" ||
		dbPort == "" ||
		dbName == "" {
		log.Fatal("Envoriment variables not defined...")
	}
}

func Start() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	checkVariables()

	// Create a new gorilla multiplexer
	router := mux.NewRouter()

	// wiring
	dbClient := getDBClient()

	customerRepositoryDB := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDB := domain.NewAccountRepositoryDB(dbClient)
	transactionRepositoryDB := domain.NewTransactionRepositoryDB(dbClient)

	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDB)}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDB)}
	th := TransactionHandlers{service.NewTransactionService(transactionRepositoryDB)}

	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.newAccount).Methods(http.MethodPost)
	router.HandleFunc("/transaction", th.newTransaction).Methods(http.MethodPost)

	// Listen on port 8000, we can pass a router as a seccond param but we are using the default one
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func getDBClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbAddress, dbPort, dbName)

	// Create a database client
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
