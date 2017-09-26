package api

import (
	//"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"strconv"

	//"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Cat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Cat endpoint")
}

func Store(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Store endpoint")
}

func Pin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pin endpoint")
}

func Unpin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Unpin endpoint")
}
