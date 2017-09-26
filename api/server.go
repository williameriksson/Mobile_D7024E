package api

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
)


func StartNewHttpServer() *http.Server {
    router := mux.NewRouter().StrictSlash(true)

    srv := &http.Server{Addr: ":8080", Handler: router}

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