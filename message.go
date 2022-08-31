package main

//Message Struct
type message struct {
	Phase int
	Name  string
	Value string
	To    string
}

type ack_message struct {
	Phase int
}

var ack_msg ack_message

type message_conn struct {
	Type       int
	Peers      []string
	Vault_name string
	Sender     string
	T          int
	Vault_Map  map[string]string
	Peer_index map[string]int
}
