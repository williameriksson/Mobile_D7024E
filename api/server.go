package api

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "Mobile_D7024E/d7024e"
)

const port string = ":8080"
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

    log.Fatal(http.ListenAndServe(port, router))
}

// Returns a server object so it can be shut down
func StartAndReturnServer() *http.Server {
    router := mux.NewRouter().StrictSlash(true)

    srv := &http.Server{Addr: port, Handler: router}

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

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            // cannot panic, because this probably is an intentional close
            log.Printf("Httpserver: ListenAndServe() error: %s", err)
        }
    }()

    // returning reference so caller can call Shutdown()
    return srv
}