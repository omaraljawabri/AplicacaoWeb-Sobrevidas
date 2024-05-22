package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var db = fazConexaoComBanco()

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/quote", createQuote)

	log.Println("Server rodando na porta 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Fatal(err)
	}
}

func fazConexaoComBanco() *sql.DB{
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("Erro ao carregar arquivo .env")
	}

	usuarioBancoDeDados := os.Getenv("USUARIO")
	senhaDoUsuario := os.Getenv("SENHA")
	nomeDoBancoDeDados := os.Getenv("NOME_BANCO_DE_DADOS")
	dadosParaConexao := "user=" + usuarioBancoDeDados + " dbname=" + nomeDoBancoDeDados + " password=" + senhaDoUsuario + " host=localhost port=5432 sslmode=disable"
	database, err := sql.Open("postgres", dadosParaConexao)
	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Query("CREATE TABLE IF NOT EXISTS cadastro(id SERIAL PRIMARY KEY, nome_completo VARCHAR(255) UNIQUE NOT NULL, cpf VARCHAR(15) UNIQUE NOT NULL, cns VARCHAR(15), cbo VARCHAR(15), cnes VARCHAR(15), ine VARCHAR(15), senha VARCHAR(20))")
	if err != nil{
		log.Fatal(err)
	}

	return database
}

func createQuote(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Redirect(w, r, "/quote", http.StatusSeeOther)
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
		http.Redirect(w, r, "/quote", http.StatusSeeOther)
		return
	}
}