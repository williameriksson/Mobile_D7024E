package api

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Cat",
		"GET",
		"/cat/{filename}",
		Cat,
	},
	Route{
		"Store",
		"POST",
		"/store",
		Store,
	},
	Route{
		"Pin",
		"GET",
		"/pin/{filename}",
		Pin,
	},
	Route{
		"Unpin",
		"GET",
		"/unpin/{filename}",
		Unpin,
	},
	Route{
		"AddNode",
		"GET",
		"/addnode/{addr}",
		AddNode,
	},
}
