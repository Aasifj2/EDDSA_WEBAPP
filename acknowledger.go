package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func acknowledge(peerID string, phase int, h host.Host) {

	//Save sent value
	m.Lock()
	sent_peer_phase[peerID] = phase
	m.Unlock()

	//Send acknowledgement
	peer_ip, _ := peer_map[peerID]
	addr, _ := multiaddr.NewMultiaddr(peer_ip)
	peer_info, err := peer.AddrInfoFromP2pAddr(addr)
	s, err := h.NewStream(context.Background(), peer_info.ID, "/ack/0.0.1")
	if err != nil {
		log.Println(err, "Error in creating connection")
		return
	}
	message_send := ack_message{Phase: phase}
	b_message, err := json.Marshal(message_send)
	if err != nil {
		log.Println(err, "Error in jsonifying data")
		return
	}
	//fmt.Println(b_message)
	_, err = s.Write(append(b_message, '\n'))
	if err != nil {
		log.Println(err, "Error in jsonifying data")
		return
	}
	rand.Seed(time.Now().UnixNano())
	rand_time := rand.Intn(10)
	time.Sleep(100 * time.Duration(rand_time))
	// receive_peer_phase[peerID] = phase
}

func host_acknowledge(h host.Host) {

	h.SetStreamHandler("/ack/0.0.1", func(s network.Stream) {
		//log.Println("sender received new stream")
		if err := process_ack(s, h); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})
}
func process_ack(s network.Stream, h host.Host) error {
	buf := bufio.NewReader(s)
	str, _ := buf.ReadBytes('\n')
	bytes := []byte(str)
	var message_receive1 message
	json.Unmarshal(bytes, &message_receive1)

	// append_new_phase(s.Conn().RemotePeer().String(), message_receive1.Phase)

	l.Lock()
	receive_peer_phase[s.Conn().RemotePeer().String()] = message_receive1.Phase
	l.Unlock()
	_, err := s.Write([]byte(""))
	return err

}

// func append_new_phase(s string, i int) {
// 	receive_peer_phase[s] = i
// 	defer func() {
// 		if p := recover(); p != nil {
// 			time.Sleep(time.Millisecond * 2)
// 			append_new_phase(s, i)
// 		}
// 	}()
// }
