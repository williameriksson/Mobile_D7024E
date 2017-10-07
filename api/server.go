package api

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "Mobile_D7024E/d7024e"
    "os"
    "path/filepath"
    "time"
)

//const addr string = ":8080"
var default_dir string //= "C:/Users/David/go/src/Mobile_D7024E/files/"
var kademlia *d7024e.Kademlia
var res chan string

const files_dir string = "files"

func StartServer(kad *d7024e.Kademlia) {
    res = make(chan string)
    kademlia = kad

    var err error

    default_dir, err = filepath.Abs("..")
    if err != nil {
        log.Fatal(err)
    }

    default_dir = filepath.Join(default_dir, files_dir, kademlia.RoutingTable.GetMyID())
    log.Println("Kademlia directory: " + default_dir)

    myIP := kademlia.RoutingTable.GetMyIP()
    myIP = convertIP(myIP)
    myPort := myIP[len(myIP) - 5 :]


    if _, err = os.Stat(default_dir); os.IsNotExist(err) {
        mk := os.MkdirAll(default_dir, os.ModePerm)
        if mk != nil {
            log.Fatal(mk)
        }
    }

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
            case d7024e.CMD_FOUND_FILE:
                    log.Println("CMD_FOUND_FILE")
                    select {
                      case  res <- handle.Ip:
                      case <-time.After(time.Second * 2):
                        log.Println("Handler did not listen, timeout")
                    }

                case d7024e.CMD_RETRIEVE_FILE:
                    log.Println("CMD_RETRIEVE_FILE")
                    GetFile(handle.PurgeInfo, convertIP(handle.Ip))
                case d7024e.CMD_REMOVE_FILE:
                    log.Println("CMD_REMOVE_FILE")
                    DeleteFile(handle.Ip)
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
