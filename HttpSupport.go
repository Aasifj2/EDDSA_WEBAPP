package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

var tmpl *template.Template

type Dealer_Data struct {
	IP_addreses string
	Threshold_T int
}

type MyInfoStruct struct {
	VaultID string
	MyIp    string
	//Make Every Property start with Capital letter
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	//For accessing the assets folder in index.html file . we need to specify that ( IN HTTP ROUTER)
	// fs := http.FileServer(http.Dir("assets"))
	// http.Handle("/assets/", http.StripPrefix("/assets", fs))

	//SAME AS ABOVE FOR MUX ROUTER ( SHIFTED TO MAIN.GO)
	// fs := http.FileServer(http.Dir("assets"))
	// r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to My world %s", r.URL.Path)

}

func Dummy_api() string {
	return "HELLO THIS IS YOUR GUILTY CONSCIENCE"
}

func Index(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "index.html", nil)

}
func CSS2(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}
func DisplayForm(w http.ResponseWriter, r *http.Request) {
	//P2p_func()
	this_vault = r.FormValue("vaultID")

	recieved_data := MyInfoStruct{
		VaultID: r.FormValue("vaultID"),
		MyIp:    p2p.Host_ip,
	}

	// recieved_data := Data2{
	// 	Item1: p2p.Host_ip,
	// 	Item2: "hello",
	// 	Item3: Dummy_api(),
	// }
	tmpl.ExecuteTemplate(w, "form.html", struct {
		Success bool
		Mydata  MyInfoStruct
	}{true, recieved_data})
}

func DisplayData(w http.ResponseWriter, r *http.Request) {
	tempT, _ := strconv.Atoi(r.FormValue("T"))
	recieved_data := Dealer_Data{
		IP_addreses: r.FormValue("ip"),
		Threshold_T: tempT,
	}
	Threshold = tempT
	fmt.Println("ISIDEEEEEEE::::", Threshold)
	gen_keyshares(recieved_data.IP_addreses)
	tmpl.ExecuteTemplate(w, "display.html", struct {
		Success bool
		Mydata  Dealer_Data
	}{true, recieved_data})
	// tmpl.ExecuteTemplate(w, "display.html", struct {
	// 	Mydata Data2
	// }{recieved_data})
}

// func Init_vault(){

// }
