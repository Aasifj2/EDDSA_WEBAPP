package main

type test_struct struct {
	Peer_list string `json:"peer_list"`
}

type gen_share struct {
	Pvt  string `json:"pvt"`
	List string `json:"peer_list"`
	T    int    `json:"t"`
	// V    string `json:"vault"`
}

type vault struct {
	V string `json:"vault"`
}

var start_p2p_flag = 0

var debug = false
