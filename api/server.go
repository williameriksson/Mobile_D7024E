package api

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "Mobile_D7024E/d7024e"
    "os"
    "Mobile_D7024E/common"
)

//const addr string = ":8080"
var default_dir string = "C:/Users/David/go/src/Mobile_D7024E/files/"
var kademlia *d7024e.Kademlia

func StartServer(kad *d7024e.Kademlia) {
    kademlia = kad
    default_dir = default_dir+kademlia.RoutingTable.GetMyID()

    myIP := kademlia.RoutingTable.GetMyIP()
    myIP = convertIP(myIP)
    myPort := myIP[len(myIP) - 5 :]


    if _, err := os.Stat(default_dir); os.IsNotExist(err) {
        mk := os.Mkdir(default_dir, os.ModePerm)
        if mk != nil {
            log.Fatal(mk)
        }
    }
    
    log.Print(default_dir)
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
    go func() {
        for {
            handle := <-kademlia.ServerChannel
            switch handle.Command {
                case common.CMD_FOUND_FILE:
                case common.CMD_RETRIEVE_FILE:
                    log.Println("CMD_RETRIEVE_FILE")
                    GetFile(handle.Hash, convertIP(handle.Ip))
                case common.CMD_REMOVE_FILE:
                default:
            }    
        }
        
    }()
    log.Fatal(http.ListenAndServe(myPort, router))

}

func convertIP(ip string) string{
    temp := []byte(ip)
    temp[len(temp) - 4] = temp[len(temp) - 4] + 1
    new_ip := string(temp)
    return new_ip
}