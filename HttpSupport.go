package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var tmpl *template.Template

type Data2 struct {
	Item1 string
	Item2 string
	Item3 string
}

type MyInfoStruct struct {
	MyIp string
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
	recieved_data := MyInfoStruct{
		MyIp: p2p.Host_ip,
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
	recieved_data := Data2{
		Item1: r.FormValue("1"),
		Item2: r.FormValue("2"),
		Item3: Dummy_api(),
	}
	tmpl.ExecuteTemplate(w, "display.html", struct {
		Success bool
		Mydata  Data2
	}{true, recieved_data})
	// tmpl.ExecuteTemplate(w, "display.html", struct {
	// 	Mydata Data2
	// }{recieved_data})
}
