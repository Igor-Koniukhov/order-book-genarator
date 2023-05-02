package main

import (
	"database/sql"
	"fmt"
	"generatorCandleChartData/internal/handlers"

	"generatorCandleChartData/internal/repository"
	"github.com/rs/cors"
	"io"
	"log"

	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
	"golang.org/x/net/websocket"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "chart_data"
)

func main() {
	fmt.Println("started")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)

	} else {
		fmt.Println("db connected")
	}
	defer db.Close()

	service := &repository.DataService{Db: db}
	schema, err := service.CreateGraphQLSchema()
	if err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	websocketHandler := websocket.Handler(func(ws *websocket.Conn) {
		handler := &handlers.WebsocketHandler{}
		handler.Connect(ws)
		generateDataTicker := time.NewTicker(30 * time.Second)
		listenWSQueries := time.NewTicker(100 * time.Millisecond)
		cupDataTicker := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-generateDataTicker.C:

				data, err := service.GenerateData()

				if err != nil {
					log.Printf("Error generating data: %v", err)
					continue
				}

				if err := service.SaveData(data); err != nil {
					log.Printf("Error saving data: %v", err)
					continue
				}
				if ws != nil && ws.IsServerConn() {
					if err := handler.SendData(data); err != nil {
						log.Printf("Error sending data: %v", err)
						continue
					}
				}

			case <-listenWSQueries.C:
				var msg string
				if err := websocket.Message.Receive(ws, &msg); err != nil {
					if err != io.EOF {
						log.Printf("Error receiving message: %v", err)
					}
					continue
				}

				result := graphql.Do(graphql.Params{
					Schema:        schema,
					RequestString: msg,
				})
				fmt.Println(msg, " msg")

				if len(result.Errors) > 0 {
					log.Printf("GraphQL errors: %v", result.Errors)
				}

				err = websocket.JSON.Send(ws, result)
				if err != nil {
					log.Printf("Error sending GraphQL response: %v", err)
					break
				}
			case <-cupDataTicker.C:
				orderBookCup := service.GenerateRandomOrders()
				if err := service.SaveOrdersToDatabase(orderBookCup); err != nil {
					log.Printf("Error saving orders: %v", err)
					continue
				}
				if err := handler.SendOrderBookData(orderBookCup); err != nil {
					log.Printf("Error sending orders data: %v", err)
					continue
				}
			}
		}
	})
	handler := c.Handler(websocketHandler)
	http.Handle("/websocket", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
