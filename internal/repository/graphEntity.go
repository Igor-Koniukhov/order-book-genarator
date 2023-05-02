package repository

import "github.com/graphql-go/graphql"

var StockData = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StockData",
		Fields: graphql.Fields{
			"Date": &graphql.Field{
				Type: graphql.String,
			},
			"Open": &graphql.Field{
				Type: graphql.Float,
			},
			"High": &graphql.Field{
				Type: graphql.Float,
			},
			"Low": &graphql.Field{
				Type: graphql.Float,
			},
			"Close": &graphql.Field{
				Type: graphql.Float,
			},
			"Volume": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
var OrderType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.ID,
			},
			"Price": &graphql.Field{
				Type: graphql.Float,
			},
			"Amount": &graphql.Field{
				Type: graphql.Float,
			},
			"Total": &graphql.Field{
				Type: graphql.Float,
			},
			"OrderType": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
