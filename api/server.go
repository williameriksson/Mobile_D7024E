package api

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "Mobile_D7024E/d7024e"
)

const privateAddr string = "127.0.0.1:8080"
const publicAddr string = ":8081"
var kademlia *d7024e.Kademlia

func StartServer(kad *d7024e.Kademlia) {
    kademlia = kad
    privateRouter := mux.NewRouter().StrictSlash(true)
    publicRouter := mux.NewRouter().StrictSlash(true)

    for _, route := range privateRoutes {
        var handler http.Handler

        handler = route.HandlerFunc
        handler = Logger(handler, route.Name)

        privateRouter.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }

    for _, route := range publicRoutes {
        var handler http.Handler

        handler = route.HandlerFunc
        handler = Logger(handler, route.Name)

        publicRouter.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }

    go func(){
        log.Fatal(http.ListenAndServe(privateAddr, privateRouter))
    }()
    log.Fatal(http.ListenAndServe(publicAddr, publicRouter))
}