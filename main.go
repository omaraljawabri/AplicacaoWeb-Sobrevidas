package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request){
	var fileName = "login.html"
	t, err := template.ParseFiles(fileName)
	if err != nil{
		fmt.Println("Erro ao fazer o parse files", err)
		return
	}
	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil{
		fmt.Println("Erro ao executar a template", err)
		return
	}
}

func loginSubmit(w http.ResponseWriter, r *http.Request){

}

func register(w http.ResponseWriter, r *http.Request){

}

func handler(w http.ResponseWriter, r *http.Request){
	switch r.URL.Path{
	case "/login":
		login(w, r)
	case "/register":
	case "/login-submit":
	default:
	}
}

func main(){
	http.HandleFunc("/", handler)
	http.ListenAndServe("", nil)
}