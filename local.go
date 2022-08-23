package main

import (
	"fmt"
	"strings"
	"time"
)

func local() {

	fmt.Println("Enter Vault Name:")

	fmt.Scan(&this_vault)

	fmt.Println("If Dealer type '1', else input anything else.")
	var deal string
	fmt.Scan(&deal)
	if deal == "1" {
		//To avoid double send test_conn
		execute_send = 1
		fmt.Println("Enter all addresses seperated by ',' and no space: ")
		var inp_strings string
		fmt.Scan(&inp_strings)

		peer_details_list = strings.Split(inp_strings, ",")
	}

	keygen_Stream_listener(p2p.Host)
	//Start Acknowledger
	host_acknowledge(p2p.Host)
	connection_Stream_listener(p2p.Host)

	if deal == "1" {
		test_conn()
		// keygen()
	} else {
		time.Sleep(time.Second * 120)
	}

}
