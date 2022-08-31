package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
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

		vault_map = message_receive.Vault_Map

		Threshold = message_receive.T
		peer_index = message_receive.Peer_index

		for i, peer_ip := range peer_details_list {

			if peer_details_list[i] == p2p.Host_ip {
				// this_vault = port
				// vault_map[this_vault] = peer_ip
				peer_map[strings.Split(peer_ip, "/")[len(strings.Split(peer_ip, "/"))-1]] = peer_ip
				my_index = i
				continue
			}
			connect_to, err := peer.AddrInfoFromString(peer_ip)
			if err != nil {
				log.Println(err)
			}
			if err := p2p.Host.Connect(p2p.Ctx, *connect_to); err != nil {
				log.Println("Connection failed:", peer_ip)
				all_ok = false

			} else {
				log.Println("Connected to: ", peer_ip)
			}
			peer_map[strings.Split(peer_ip, "/")[len(strings.Split(peer_ip, "/"))-1]] = peer_ip
		}

		keygen()

	}
	if message_receive.Type == 3 {
		connect_to, err := peer.AddrInfoFromString(message_receive.Sender)
		if err != nil {
			log.Println(err)
		}
		if err := p2p.Host.Connect(p2p.Ctx, *connect_to); err != nil {
			all_ok = false
		}
		message_send := message_conn{
			Type:       4,
			Vault_name: this_vault,
			Sender:     p2p.Host_ip,
		}
		s, err := p2p.Host.NewStream(p2p.Ctx, connect_to.ID, "/conn/0.0.1")
		if err != nil {
			log.Println(peer_map[message_receive.Sender])
			log.Println(err, "Connecting to send message error")

		}

		b_message, err := json.Marshal(message_send)
		if err != nil {
			log.Println(err, "Error in jsonifying data")

		}
		s.Write(append(b_message, '\n'))

	}
	if message_receive.Type == 4 {
		vault_map[message_receive.Vault_name] = message_receive.Sender
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
		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1+1) + "/" + message_receive.Name + ".txt")

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

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1+1) + "/" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		// fmt.Println("Encrypted Shares received")
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 4 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1+1) + "/" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		// fmt.Println("Encrypted Shares received")
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 5 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1+1) + "/" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		// fmt.Println("Encrypted Shares received")
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 6 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Alphas/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + "alpha" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		fmt.Println("Alpha Recieved from peer :", fmt.Sprint(res1+1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 7 {

		fmt.Println("INSIDE RECIEVE:", message_receive.Value)
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		ij := strings.Split(message_receive.Name, ",")
		//"C's/" + fmt.Sprint(i) + "/" + peer_number + "/C1.txt"
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Shares/To" + ij[0]
		os.MkdirAll(path, os.ModePerm)
		_f, _ := os.Create(path + "/C" + ij[1] + ".txt")
		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 8 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Signing/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 9 || message_receive.Phase == 10 || message_receive.Phase == 11 || message_receive.Phase == 12 {

		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]

		_f, _ := os.Create("Broadcast/" + fmt.Sprint(res1+1) + "/Signing/" + message_receive.Name + ".txt")

		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 13 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Signing/Alphas/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + "alpha" + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)
		fmt.Println("Sign Alpha Recieved from peer :", fmt.Sprint(res1+1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 14 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Signing/Shares/"
		os.MkdirAll(path, 0755)
		file, _ := os.OpenFile(path+"shareTo"+message_receive.Name+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		//_f, _ := os.Create(path + "shareTo" + message_receive.Name + ".txt")
		//_f.WriteString(message_receive.Value)
		_, _ = fmt.Fprint(file, message_receive.Value)

		fmt.Println("Sign Shares Recieved from peer :", fmt.Sprint(res1+1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	} else if message_receive.Phase == 16 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Signing/V_i.txt"
		_f, _ := os.Create(path)
		_f.WriteString(message_receive.Value)
		fmt.Println("V Broadcasted By peer :", fmt.Sprint(res1+1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 17 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Signing/U.txt"
		_f, _ := os.Create(path)
		_f.WriteString(message_receive.Value)
		fmt.Println("U Broadcasted By peer :", fmt.Sprint(res1+1))
		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)

	} else if message_receive.Phase == 15 {
		res1 := peer_index[peer_map[s.Conn().RemotePeer().String()]]
		path := "Broadcast/" + fmt.Sprint(res1+1) + "/Signing/"
		os.MkdirAll(path, 0755)
		_f, _ := os.Create(path + message_receive.Name + ".txt")
		_f.WriteString(message_receive.Value)

		acknowledge(s.Conn().RemotePeer().String(), message_receive.Phase, h)
	}

	_, err = s.Write([]byte(""))
	return err
}
