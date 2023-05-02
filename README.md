## The app, for order generation, is devoted to the creation of dummy order-book data and sending via WebSocket and by graphql query.
#### queries:
```azure
query {
      stockDataForPeriod(startTime: ${startTime}, endTime: ${endTime}) {
        Date
        Open
        High
        Low
        Close
        Volume
      }
    }
```
```azure
query {
  lastBuyOrders( limit: 50) {
    ID
    OrderType
    Price
    Amount
    Total
  }
}
```
```azure
query {
  lastSellOrders( limit: 50) {
    ID
    OrderType
    Price
    Amount
    Total
  }
}
```
