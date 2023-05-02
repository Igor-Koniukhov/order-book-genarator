package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"generatorCandleChartData/internal/models"
	"golang.org/x/net/websocket"
)

type WebsocketHandler struct {
	conn *websocket.Conn
}

func (wh *WebsocketHandler) Connect(ws *websocket.Conn) {
	wh.conn = ws
}
func (wh *WebsocketHandler) SendData(data models.Data) error {
	if wh.conn == nil {
		return errors.New("websocket connection is not established")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if _, err := wh.conn.Write(jsonData); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	return nil
}
func (wh *WebsocketHandler) SendOrderBookData(data models.OrderBookCup) error {
	if wh.conn == nil {
		return errors.New("websocket connection is not established")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	fmt.Println(jsonData)

	if _, err := wh.conn.Write(jsonData); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	return nil
}
