package pages

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func AutenticaCadastroELevaAoLogin(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	nomecompleto := r.PostForm.Get("nome_completo")
	cpf := r.PostForm.Get("cpf")
	cns := r.PostForm.Get("cns")
	cbo := r.PostForm.Get("cbo")
	cnes := r.PostForm.Get("cnes")
	ine := r.PostForm.Get("ine")
	senha := r.PostForm.Get("senha")
	confirmsenha := r.PostForm.Get("confirmsenha")

	if confirmsenha == senha{
		_, err = db.Exec("INSERT INTO cadastro(nome_completo, cpf, cns, cbo, cnes, ine, senha) VALUES($1, $2, $3, $4, $5, $6, $7)", nomecompleto, cpf, cns, cbo, cnes, ine, senha)
		if err != nil{
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else{
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = templates.ExecuteTemplate(w, "login.html", "a")
	if err != nil{
		return
	}
}