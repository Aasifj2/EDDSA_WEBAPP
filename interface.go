package main

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"gopkg.in/dedis/kyber.v2"
)

//Discovered values
// type peer_details struct {
// 	id   peer.ID
// 	addr peer.AddrInfo
// }

// var peer_details_list []peer_details = make([]peer_details, 0, 10)
var peer_details_list []string

// var round = make(map[string]int)
var peer_map = make(map[string]string)

var sorted_peer_id []string
var my_index int = 0

//Store the current phase value received
var sent_peer_phase = make(map[string]int)

//Store the phase value of peer acknowledgement
var receive_peer_phase = make(map[string]int)

//Lock map to avoid concurrent map writes
var l = sync.Mutex{}

var m = sync.Mutex{}

type P2P struct {

	// Represents the libp2p host
	Host    host.Host
	Host_ip string
	Ctx     context.Context
	Peers   []string
}

var p2p P2P

//Rework flags and channels to conform to this struct
type Status struct {
	Phase     int
	Chan      string
	Num_peers int
}

var status_struct Status
var all_ok = true
var peer_index = make(map[string]int)

var copy []string
var this_vault string
var vault_map = make(map[string]string)
var execute_send = 0

var n = 0
var Threshold int

var Group_key kyber.Point

var port = ":8082"

var keygenFlag = false
