package main

import (
	"encoding/json"
	"log"
	"sort"
	"strings"

	peer "github.com/libp2p/go-libp2p-core/peer"
)

func test_conn() {
	peer_details_list = append(peer_details_list, p2p.Host_ip)

	// sort.Strings(peer_details_list)
	// for i, item := range peer_details_list {
	// 	sorted_peer_id = append(sorted_peer_id, strings.Split(item, "/")[len(strings.Split(item, "/"))-1])
	// 	peer_index[strings.Split(item, "/")[len(strings.Split(item, "/"))-1]] = i
	// 	if item == p2p.Host_ip {
	// 		my_index = i
	// 	}
	// }

	for i, peer_ip := range peer_details_list {
		// peer_map[strings.Split(peer_ip, "/")[len(strings.Split(peer_ip, "/"))-1]] = peer_ip
		// fmt.Println(len(sorted_peer_id))
		// if i == my_index {
		// 	continue
		//
		if peer_details_list[i] == p2p.Host_ip {
			// this_vault = port
			vault_map[this_vault] = peer_ip
			peer_map[strings.Split(peer_ip, "/")[len(strings.Split(peer_ip, "/"))-1]] = peer_ip

			continue
		}
		peer_map[strings.Split(peer_ip, "/")[len(strings.Split(peer_ip, "/"))-1]] = peer_ip
		connect_to, err := peer.AddrInfoFromString(peer_ip)
		if err != nil {
			log.Println(err)
		}
		if err := p2p.Host.Connect(p2p.Ctx, *connect_to); err != nil {
			log.Println("Connection failed:", peer_ip)
			all_ok = false
			return
		} else {
			log.Println("Connected to: ", peer_ip)
		}
		message_send := message_conn{
			Type:       1, //Type 1 -> Dealer to Non Dealer
			Peers:      peer_details_list,
			Vault_name: this_vault,
			Sender:     p2p.Host_ip,
		}

		s, err := p2p.Host.NewStream(p2p.Ctx, connect_to.ID, "/conn/0.0.1")
		if err != nil {
			log.Println(peer_map[peer_ip])
			log.Println(err, "Connecting to send message error")
			return
		}

		b_message, err := json.Marshal(message_send)
		if err != nil {
			log.Println(err, "Error in jsonifying data")
			return
		}
		s.Write(append(b_message, '\n'))

	}
	keygen()

}

func test() {
	for i, peer_ip := range peer_details_list {

		// fmt.Println(len(sorted_peer_id))

		peer_map[strings.Split(peer_ip, "/")[len(strings.Split(peer_ip, "/"))-1]] = peer_ip
		if peer_details_list[i] == p2p.Host_ip {
			// this_vault = port

			vault_map[this_vault] = peer_ip

			// my_index = i
			continue
		}

		connect_to, err := peer.AddrInfoFromString(peer_ip)
		if err != nil {
			log.Println(err)
		}
		if err := p2p.Host.Connect(p2p.Ctx, *connect_to); err != nil {
			log.Println("Connection failed:", peer_ip)
			all_ok = false
			return
		} else {
			log.Println("Connected to: ", peer_ip)
		}
		message_send := message_conn{
			Type:       1, //Type 1 -> Dealer to Non Dealer
			Peers:      peer_details_list,
			Vault_name: this_vault,
			Sender:     p2p.Host_ip,
		}

		s, err := p2p.Host.NewStream(p2p.Ctx, connect_to.ID, "/conn/0.0.1")
		if err != nil {
			log.Println(peer_map[peer_ip])
			log.Println(err, "Connecting to send message error")
			return
		}

		b_message, err := json.Marshal(message_send)
		if err != nil {
			log.Println(err, "Error in jsonifying data")
			return
		}
		s.Write(append(b_message, '\n'))

	}

	log.Println(len(vault_map), len(peer_details_list))

	if len(vault_map) == len(peer_details_list) {
		keys := make([]string, 0, len(vault_map))
		for k := range vault_map {

			keys = append(keys, k)
		}
		sort.Strings(keys)
		// peer_details_list = make([]string, len(peer_details_list))
		for i, k := range keys {

			copy = append(copy, vault_map[k])
			// peer_details_list[i] = vault_map[k]
			// fmt.Println("YOOOOOOOOOO ", string(vault_map[k]), string(p2p.Host_ip))
			if string(vault_map[k]) == string(p2p.Host_ip) {
				my_index = i

			}
			peer_index[vault_map[k]] = i
			// global peer_details_list
			// append(peer_details_list,copy)//[len(copy):]
			peer_details_list = copy
			log.Println(peer_details_list, my_index)

		}
	}
	keygen()
}
