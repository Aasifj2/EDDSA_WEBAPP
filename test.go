package main

import (
	"encoding/json"
	"fmt"
)

type Bird struct {
	Phase       string
	Description string
}

type message1 struct {
	Phase string
	Name  string
	Value string
}

func massin() {

	message_send_1 := message1{
		Phase: "current_flag",
		Name:  "u_i",
		Value: "sdsD",
	}
	b_message_1, _ := json.Marshal(message_send_1)
	//birdJson := `{"phase": "pigeon","description": "likes to perch on rocks"}`
	var bird message1
	json.Unmarshal(b_message_1, &bird)
	fmt.Printf("Species: %s, Description: %s", bird.Name, bird.Phase)

}
