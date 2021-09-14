package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

var retrier = 0

type shippingInfo struct {
	Customer    string `json:"customer"`
	Transaction string `json:"transaction"`
	Order       []struct {
		Productid int `json:"productID"`
		Quantity  int `json:"quantity"`
	} `json:"order"`
}

func main() {
	direktivapps.StartServer(creditHandler)
}

func creditHandler(w http.ResponseWriter, r *http.Request) {

	cc := new(shippingInfo)

	aid := r.Header.Get(direktivapps.DirektivActionIDHeader)
	lw, err := direktivapps.NewDirektivLogWriter(aid)
	if err != nil {
		direktivapps.RespondWithError(w, "io.direktiv.logger", err.Error())
		return
	}

	log.SetOutput(lw)
	log.SetFlags(0)

	log.Println("shipping request")

	_, err = direktivapps.Unmarshal(&cc, r)
	if err != nil {
		direktivapps.RespondWithError(w, "io.direktiv.data", err.Error())
		return
	}

	log.Printf("executing shipping %s (%v)", cc.Customer, aid)

	if len(cc.Customer) == 0 {
		direktivapps.RespondWithError(w, "io.direktiv.customer", "no customer provided")
		return
	}

	shipping := true
	if cc.Customer == "Johnny Mars" {
		shipping = false
	}
	result, _ := json.Marshal(&shipping)

	log.Printf("shipping request for customer %s: %v", cc.Customer, shipping)

	w.Write(result)

}
