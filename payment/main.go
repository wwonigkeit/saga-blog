package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

var retrier = 0

type orderInfo struct {
	Customer    string `json:"customer"`
	Transaction string `json:"transaction"`
	Order       []struct {
		Productid int `json:"productID"`
		Quantity  int `json:"quantity"`
	} `json:"order"`
	Action string `json:"action"`
}

type orderResult struct {
	Result        bool `json:"result"`
	TransactionID int  `json:"transactionID"`
}

func main() {
	direktivapps.StartServer(paymentHandler)
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {

	cc := new(orderInfo)

	aid := r.Header.Get(direktivapps.DirektivActionIDHeader)
	lw, err := direktivapps.NewDirektivLogWriter(aid)
	if err != nil {
		direktivapps.RespondWithError(w, "io.direktiv.logger", err.Error())
		return
	}

	log.SetOutput(lw)
	log.SetFlags(0)

	log.Println("payment request")

	_, err = direktivapps.Unmarshal(&cc, r)
	if err != nil {
		direktivapps.RespondWithError(w, "io.direktiv.data", err.Error())
		return
	}

	if cc.Action == "undo" {
		log.Printf("undo payment %s (%v), transaction id: %v", cc.Customer, aid, cc.Transaction)
		return
	}

	log.Printf("running payment for %s (%v), transaction id: %v", cc.Customer, aid, cc.Transaction)

	if len(cc.Customer) == 0 {
		direktivapps.RespondWithError(w, "io.direktiv.customer", "no customer provided")
		return
	}

	processed := true

	// hardcoded error if customer id is 123
	if cc.Customer == "Johnny Patience" {
		time.Sleep(2 * time.Minute)
	}

	// hardcoded error if customer id is 123
	if cc.Customer == "Johnny No-Cash" {
		processed = false
	}

	// hardcoded retry demo. success after 3rd attempt
	if cc.Customer == "Pay Retry" {
		retrier++
		if retrier%3 != 0 {
			direktivapps.RespondWithError(w, "io.direktiv.customer", "payment failed")
			return
		}
	}

	or := orderResult{
		Result:        processed,
		TransactionID: rand.Intn(100),
	}

	result, err := json.Marshal(or)
	if err != nil {
		direktivapps.RespondWithError(w, "io.direktiv.internal", err.Error())
		return
	}

	log.Printf("payment request for customer %s: %v", cc.Customer, processed)

	log.Printf("json result: %v %+v", string(result), or)

	w.Write(result)

}
