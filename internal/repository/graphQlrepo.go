package repository

import (
	"github.com/graphql-go/graphql"
)

func (ds *DataService) CreateGraphQLSchema() (graphql.Schema, error) {
	queryFields := graphql.Fields{
		"stockDataForPeriod": &graphql.Field{
			Type: graphql.NewList(StockData),
			Args: graphql.FieldConfigArgument{
				"startTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"endTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				startTime := int64(params.Args["startTime"].(int))
				endTime := int64(params.Args["endTime"].(int))
				return ds.FetchStockDataForPeriod(startTime, endTime)
			},
		},
		"allOrders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return ds.FetchAllOrders()
			},
		},
		"sellOrders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return ds.FetchOrdersByType("sell")
			},
		},
		"buyOrders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return ds.FetchOrdersByType("buy")
			},
		},
		"lastSellOrders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Args: graphql.FieldConfigArgument{

				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 50,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				limit, _ := params.Args["limit"].(int)
				return ds.FetchLastOrdersByType("sell", limit)
			},
		},
		"lastBuyOrders": &graphql.Field{
			Type: graphql.NewList(OrderType),
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 50,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				limit, _ := params.Args["limit"].(int)
				return ds.FetchLastOrdersByType("buy", limit)
			},
		},
	}
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: queryFields,
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
