package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	// "gopkg.in/dedis/kyber.v2/util/encoding"
)

func connection_Stream_listener(h host.Host) {
	//fmt.Println("Got a new stream!")

	h.SetStreamHandler("/conn/0.0.1", func(s network.Stream) {
		//log.Println("sender received new stream")
		if err := process_connection(s, h); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}

	})

}
func process_connection(s network.Stream, h host.Host) error {

	buf := bufio.NewReader(s)
	//log.Println(s)
	str, err := buf.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return err
	}
	bytes := []byte(str)
	var message_receive message_conn
	json.Unmarshal(bytes, &message_receive)
	if message_receive.Type == 1 {
		peer_details_list = message_receive.Peers
		vault := message_receive.Vault_name
		vault_map[vault] = message_receive.Sender
		if Threshold == 0 {
			Threshold = message_receive.T
		}
	}
	if execute_send == 0 {
		execute_send = 1
		test()

	} else {

		log.Println(len(vault_map), len(peer_details_list))
		if len(vault_map) == len(peer_details_list) {
			keys := make([]string, 0)
			for k := range vault_map {

				keys = append(keys, k)
			}
			sort.Strings(keys)

			for i, k := range keys {

				log.Println(k, "->key")
				copy = append(copy, vault_map[k])
				// peer_details_list[i] = vault_map[k]
				log.Println(string(vault_map[k]), string(p2p.Host_ip))
				if string(vault_map[k]) == string(p2p.Host_ip) {
					my_index = i

				}
				peer_index[vault_map[k]] = i

			}

			peer_details_list = copy
			log.Println(peer_details_list, my_index)
		}
		// test()
		return nil
	}
	return nil
}

func keygen_Stream_listener(h host.Host) {
	//fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	//Return Channel details
	h.SetStreamHandler("/keygen/0.0.1", func(s network.Stream) {
		//log.Println("sender received new stream")
		if err := process_input(s, h); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}

	})
	// 'stream' will stay open until you close it (or the other side closes it).

}

func process_input2(s network.Stream, h host.Host) error {

	//log.Println(s)
	buf := bufio.NewReader(s)
	//log.Println(s)
	str, err := buf.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return err
	}
	bytes := []byte(str)
	var message_receive message
	json.Unmarshal(bytes, &message_receive)
	//log.Println(s.Conn().RemotePeer())

	//Check and rediect :
	//sender_id := s.ID()[1 : len(s.ID())-2]

	if message_receive.Phase == 1 {
		//Index peer_index[s.Conn().RemotePeer().String()] use this instead of sort.Search()
		log.Println("Got ppk_j: ", message_receive.Value, " from ", s.Conn().RemotePeer())
		// log.Println(s.Conn().RemotePeer().String())
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 2 {
		log.Println("Got Kgc_j: ", message_receive.Value, " from ", s.Conn().RemotePeer())
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 3 {
		log.Println("Got Kgd_j: ", message_receive.Value, " from ", s.Conn().RemotePeer())
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	}

	_, err = s.Write([]byte(""))
	return err
}

func process_input(s network.Stream, h host.Host) error {

	//log.Println(s)
	buf := bufio.NewReader(s)
	//log.Println(s)
	str, err := buf.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return err
	}
	bytes := []byte(str)
	var message_receive message
	json.Unmarshal(bytes, &message_receive)
	//log.Println(s.Conn().RemotePeer())

	//Check and rediect :
	//sender_id := s.ID()[1 : len(s.ID())-2]

	if message_receive.Phase == 2 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]] //use this instead of sort.Search()

		// res1:= peer_index[s.Conn().RemotePeer().String()]
		// log.Println("Hey look at me IN PHASE 2", res1, message_receive.Name, peer_map, my_index)
		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1) + "/" + message_receive.Name + ".txt")

		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 1 {
		//log.Println("Got public key - ", message_receive.Value, " from ", s.Conn().RemotePeer())
		os.MkdirAll("Broadcast/"+message_receive.Name, 0755)
		_f, _ := os.OpenFile("Broadcast/"+message_receive.Name+"/EPK.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		_f.WriteString(message_receive.Value)
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 3 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1) + "/" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		// fmt.Println("Encrypted Shares received")
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 4 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1) + "/" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		// fmt.Println("Encrypted Shares received")
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 5 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1) + "/" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		// fmt.Println("Encrypted Shares received")
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 6 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1) + "/Alphas/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + "alpha" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		fmt.Println("Alpha Recieved from peer :", fmt.Sprint(res1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 7 {

		fmt.Println("INSIDE RECIEVE:", message_receive.Value)
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		ij := strings.Split(message_receive.Name, ",")
		//"C's/" + fmt.Sprint(i) + "/" + peer_number + "/C1.txt"
		path := "Broadcast/" + fmt.Sprint(res1) + "/Shares/To" + ij[0]
		os.MkdirAll(path, os.ModePerm)
		_f, _ := os.Create(path + "/C" + ij[1] + ".txt")
		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 8 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1) + "/Signing/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 9 || message_receive.Phase == 10 || message_receive.Phase == 11 || message_receive.Phase == 12 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1) + "/Signing/" + message_receive.Name + ".txt")

		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 13 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1) + "/Signing/Alphas/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + "alpha" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		fmt.Println("Sign Alpha Recieved from peer :", fmt.Sprint(res1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 14 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1) + "/Signing/Shares/"
		os.MkdirAll(path, 0755)
		file, _ := os.OpenFile(path+"shareTo"+message_receive.Name+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		//_f, _ := os.Create(path + "shareTo" + message_receive.Name + ".txt")
		//_f.WriteString(message_receive.Value)
		_, _ = fmt.Fprint(file, message_receive.Value)

		fmt.Println("Sign Shares Recieved from peer :", fmt.Sprint(res1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 15 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1) + "/Signing/V_i.txt"
		_f, _ := os.Create(path)
		_f.WriteString(message_receive.Value)
		fmt.Println("V Broadcasted By peer :", fmt.Sprint(res1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 16 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1) + "/Signing/U.txt"
		_f, _ := os.Create(path)
		_f.WriteString(message_receive.Value)
		fmt.Println("U Broadcasted By peer :", fmt.Sprint(res1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	}

	_, err = s.Write([]byte(""))
	return err
}
