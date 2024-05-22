package main

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

var db = fazConexaoComBanco()
var templates = template.Must(template.ParseFiles("./index.html", "./templates/telalogin/login.html", "./templates/telaesqueceusenha/esqueceusenha.html", "./templates/dashboard/dashboard.html"))

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/login", autenticaCadastroELevaAoLogin)
	http.HandleFunc("/dashboard", autenticaLoginELevaAoDashboard)
	http.HandleFunc("/senhaatualizada", atualizarSenha)

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

func autenticaCadastroELevaAoLogin(w http.ResponseWriter, r *http.Request){
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

type validarlogin struct{
	Cpf string
	Senha string
}

func autenticaLoginELevaAoDashboard(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	cpf := r.FormValue("cpf")
	senha := r.FormValue("senha")
	cpfsenha, err := db.Query("SELECT cpf, senha FROM cadastro")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cpfsenha.Close()
	armazenamento := make([]validarlogin, 0)

	for cpfsenha.Next(){
		armazenar := validarlogin{}
		err := cpfsenha.Scan(&armazenar.Cpf, &armazenar.Senha)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = cpfsenha.Err(); err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, armazenado := range armazenamento{
		if armazenado.Cpf == cpf && armazenado.Senha == senha{
			err = templates.ExecuteTemplate(w, "dashboard.html", "a")
			if err != nil{
				return
			}
		}
	}
}

type validarCpf struct{
	Cpf string
}

func atualizarSenha(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	cpf := r.FormValue("cpf")
	senha := r.FormValue("senha")
	confirmarsenha := r.FormValue("confirmpassword")
	pegarcpf, err := db.Query("SELECT cpf FROM cadastro")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	defer pegarcpf.Close()
	armazenamento := make([]validarCpf, 0)

	for pegarcpf.Next(){
		armazenar := validarCpf{}
		err := pegarcpf.Scan(&armazenar.Cpf)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = pegarcpf.Err(); err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, armazenado := range armazenamento{
		if armazenado.Cpf == cpf && senha==confirmarsenha{
			_, err := db.Exec(`UPDATE cadastro SET senha=$1 WHERE cpf=$2`, senha, cpf)
			if err != nil{
				return
			}
			templates.ExecuteTemplate(w, "login.html", "a")
		}
	}
}