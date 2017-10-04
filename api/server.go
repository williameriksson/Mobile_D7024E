package api

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "Mobile_D7024E/d7024e"
)

const addr string = ":8080"
var kademlia *d7024e.Kademlia

func StartServer(kad *d7024e.Kademlia) {
    kademlia = kad
    router := mux.NewRouter().StrictSlash(true)

    for _, route := range routes {
        var handler http.Handler

        handler = route.HandlerFunc
        handler = Logger(handler, route.Name)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }
    log.Fatal(http.ListenAndServe(addr, router))
}