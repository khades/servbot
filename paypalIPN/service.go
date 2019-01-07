package paypalIPN

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/khades/servbot/httpAPI"
)

func ipn(w http.ResponseWriter, req *http.Request) {
	buf, _ := ioutil.ReadAll(req.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))

	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	incomingBuffer, _ := ioutil.ReadAll(rdr1)
	valuesToSend := [][]byte{[]byte("_notify-validate"), incomingBuffer}
	result := bytes.Join(valuesToSend, []byte("&"))
	log.Printf("To send: %s", string(result))

	client := &http.Client{Timeout: 5 * time.Second}
	r, _ := http.NewRequest("POST", "https://ipnpb.sandbox.paypal.com/cgi-bin/webscr", bytes.NewReader(result))
	r.Header.Set("Content-Type", req.Header.Get("Content-Type"))
	r.Header.Set("Content-Length", req.Header.Get("Content-Length"))
	resp, err := client.Do(r)
	if err != nil {
		log.Println("Error sending request")
		return
	}
	htmlData, err := ioutil.ReadAll(resp.Body) //<--- here!
	if err != nil {
		log.Println("Error getting data from paypal")
		return
	}
	log.Printf("GOT: %s", string(htmlData))
	if string(htmlData) == "VERIFIED" {
		log.Println("ALL GOODMAN")
	}
	if string(htmlData) == "INVALID" {
		log.Println("NOT GOOD MAN")
		return
	}
	postBuffer, postBufferError := ioutil.ReadAll(rdr2)

	if postBufferError != nil {
		log.Printf("Error reading data: %s", postBufferError.Error())
	}

	values, valuesError := url.ParseQuery(string(postBuffer))
	if valuesError != nil {
		log.Printf("Error reading values: %s", valuesError.Error())
	}
	log.Printf("Values: %+v", values)

	paymentStatus := values.Get("payment_status")
	if paymentStatus == "Completed" {
		objectID := values.Get("custom")
		log.Printf("DonationID: %s", objectID)
	}
}

func ipnGet(w http.ResponseWriter, r *http.Request) {
	httpAPI.WriteJSONError(w, "Method not allowed", 405)
}
