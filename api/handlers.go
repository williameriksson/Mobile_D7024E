package api

import (
	//"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Kademlia API:\nGET: /cat/{hash}\nPOST: /store/\nGET: /pin/{hash}\nGET: /unpin/{hash}")
}

func Cat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Cat endpoint")
	vars := mux.Vars(r)
	fmt.Fprintln(w, vars["hash"])
}

func Store(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Store endpoint")
}

func Pin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pin endpoint")
	vars := mux.Vars(r)
	fmt.Fprintln(w, vars["hash"])
}

func Unpin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Unpin endpoint")
	vars := mux.Vars(r)
	fmt.Fprintln(w, vars["hash"])
}
