package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IpApiResponse struct {
	Address string `json:"query"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	fmt.Println("Looking for your public IP")

	ipResponse := &IpApiResponse{}
	err := getJson("http://ip-api.com/json", ipResponse)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	fmt.Printf("Your public Ip is %s\n", ipResponse.Address)
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
