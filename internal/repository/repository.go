package repository

import (
	"database/sql"
	"fmt"
	"generatorCandleChartData/internal/models"
	"math/rand"
	"time"
)

type DataService struct {
	Db *sql.DB
}

func (ds *DataService) GenerateRandomOrders() models.OrderBookCup {
	orderBookCup := models.OrderBookCup{
		Bids: make([]models.Order, 0),
		Asks: make([]models.Order, 0),
	}

	for i := 0; i < 5; i++ {
		price := rand.Float64() + 0.100
		amount := rand.Float64() * 10

		orderBookCup.Bids = append(orderBookCup.Bids, models.Order{
			Price:  price,
			Amount: amount,
			Total:  price * amount,
		})

		price = rand.Float64() + 0.100
		amount = rand.Float64() * 10

		orderBookCup.Asks = append(orderBookCup.Asks, models.Order{
			Price:  price,
			Amount: amount,
			Total:  price * amount,
		})
	}

	return orderBookCup
}

func (ds *DataService) SaveOrdersToDatabase(orderBookCup models.OrderBookCup) error {

	_, err := ds.Db.Exec(`
		CREATE TABLE IF NOT EXISTS bids (
			id SERIAL PRIMARY KEY,
			price NUMERIC(12, 2) NOT NULL,
			amount NUMERIC(12, 2) NOT NULL,
			total NUMERIC(12, 2) NOT NULL,
			timestamp TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = ds.Db.Exec(`
		CREATE TABLE IF NOT EXISTS asks (
			id SERIAL PRIMARY KEY,
			price NUMERIC(12, 2) NOT NULL,
			amount NUMERIC(12, 2) NOT NULL,
			total NUMERIC(12, 2) NOT NULL,
			timestamp TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	_, err = ds.Db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			price NUMERIC(15, 6) NOT NULL,
			amount NUMERIC(15, 6) NOT NULL,
			total NUMERIC(15, 6) NOT NULL,
			order_type VARCHAR(10) NOT NULL,
		    timestamp TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	timestamp := time.Now()

	for _, bid := range orderBookCup.Bids {
		_, err := ds.Db.Exec("INSERT INTO bids (price, amount, total, timestamp) VALUES ($1, $2, $3, $4)",
			bid.Price, bid.Amount, bid.Total, timestamp)
		if err != nil {
			return err
		}
		_, err = ds.Db.Exec("INSERT INTO orders ( price, amount, total, order_type, timestamp) VALUES ($1, $2, $3, $4, $5)",
			bid.Price, bid.Amount, bid.Total, "buy", timestamp)
		if err != nil {
			return err
		}
	}

	for _, ask := range orderBookCup.Asks {
		_, err := ds.Db.Exec("INSERT INTO asks (price, amount, total, timestamp) VALUES ($1, $2, $3, $4)",
			ask.Price, ask.Amount, ask.Total, timestamp)
		if err != nil {
			return err
		}
		_, err = ds.Db.Exec("INSERT INTO orders ( price, amount, total, order_type, timestamp) VALUES ($1, $2, $3, $4, $5)",
			ask.Price, ask.Amount, ask.Total, "sell", timestamp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DataService) GenerateData() (models.Data, error) {
	now := time.Now()
	return models.Data{
		Date:   now,
		Open:   rand.Float64()*100 - 50,
		High:   rand.Float64()*100 - 50,
		Low:    rand.Float64()*100 - 50,
		Close:  rand.Float64()*100 - 50,
		Volume: rand.Intn(10000000),
	}, nil
}

func (ds *DataService) SaveData(data models.Data) error {
	return ds.InsertData(data)
}
func (ds *DataService) InsertData(data models.Data) error {
	_, err := ds.Db.Exec(`
		CREATE TABLE IF NOT EXISTS stock_chart (
			id SERIAL PRIMARY KEY,
			date TIMESTAMP NOT NULL,
			open NUMERIC(10, 2) NOT NULL,
			high NUMERIC(10, 2) NOT NULL,
			low NUMERIC(10, 2) NOT NULL,
		    close NUMERIC(10, 2) NOT NULL,
		    volume INTEGER NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	query := `INSERT INTO stock_chart (date, open, high, low, close, volume) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = ds.Db.Exec(query, data.Date, data.Open, data.High, data.Low, data.Close, data.Volume)
	return err
}
func (ds *DataService) GetAllData() ([]models.Data, error) {
	query := `SELECT date, open, high, low, close, volume FROM stock_chart`
	rows, err := ds.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.Data
	for rows.Next() {
		var d models.Data
		if err := rows.Scan(&d.Date, &d.Open, &d.High, &d.Low, &d.Close, &d.Volume); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}
func (ds *DataService) FetchOrdersByType(orderType string) ([]models.Order, error) {

	query := `
		SELECT id, order_type, price, amount, total
		FROM orders
		WHERE order_type = $1
		ORDER BY id
	`

	rows, err := ds.Db.Query(query, orderType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.OrderType, &order.Price, &order.Amount, &order.Total)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
func (ds *DataService) FetchLastOrdersByType(orderType string, limit int) ([]models.Order, error) {

	if orderType != "sell" && orderType != "buy" {
		return nil, fmt.Errorf("Invalid order type: %s", orderType)
	}

	query := `
        SELECT id, order_type, price, amount, total
        FROM orders
        WHERE order_type = $1
        ORDER BY id DESC
        LIMIT $2
    `
	rows, err := ds.Db.Query(query, orderType, limit)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch last orders by type: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.OrderType, &order.Price, &order.Amount, &order.Total)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan order: %v", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to iterate over rows: %v", err)
	}

	return orders, nil
}
func (ds *DataService) FetchStockDataForPeriod(startTime, endTime int64) ([]models.Data, error) {
	startTimeParsed := time.Unix(0, startTime*int64(time.Millisecond))
	endTimeParsed := time.Unix(0, endTime*int64(time.Millisecond))

	query := `
		SELECT Date, Open, High, Low, Close, Volume
		FROM stock_chart
		WHERE Date >= $1 AND Date <= $2
		ORDER BY Date
	`

	rows, err := ds.Db.Query(query, startTimeParsed, endTimeParsed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stockData := []models.Data{}
	for rows.Next() {
		var data models.Data
		err := rows.Scan(&data.Date, &data.Open, &data.High, &data.Low, &data.Close, &data.Volume)
		if err != nil {
			return nil, err
		}
		stockData = append(stockData, data)
	}

	return stockData, nil
}
func (ds *DataService) FetchAllOrders() ([]models.Order, error) {

	query := `
		SELECT id, price, amount, total, order_type
		FROM orders
		ORDER BY id
	`

	rows, err := ds.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.Price, &order.Amount, &order.Total, &order.OrderType)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
