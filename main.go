package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Pacientes struct{
	Nome string
	Rua string
	Numero string
	Sexo string
	Bairro string
	Complemento string
	Telefone string
	DataNasc string
	Homem string
	Etilista string
	Tabagista string
	LesaoBucal string
	Endereco string
	Idade int
	Fatores string
	Usuario string
	Primeira string
}

type validarlogin struct{
	Usuario string
	Cpf string
	Senha string
	PrimeiraLetra string
	QtdBaixo int
	QtdMedio int
	QtdAlto int
	PorcBaixo float64
	PorcMedio float64
	PorcAlto float64
	Cns string
	Cbo string
	Cnes string
	Ine string
}

type PegarDados struct{
	Homem string
	Etilista string
	Tabagista string
	LesaoBucal string
}

type UsuarioNoDashboard struct{
	Usuario string
	Primeira string
	QtdBaixo int
	QtdMedio int
	QtdAlto int
	PorcBaixo float64
	PorcMedio float64
	PorcAlto float64
}

type validarCpf struct{
	Cpf string
}

type ACS struct{
	User string
	NomeCompleto string
	CPF string
	CNS string
	CBO string
	CNES string
	INE string
	SenhaACS string
	PrimeiraLetra string
}

type DadosForm struct{
	Cns []string
	Cbo1 string
	Cbo2 string
	Cbo3 string
	Cbo4 string
	Cbo5 string
	Cbo6 string
	Cnes []string
	Ine []string
}


var db = fazConexaoComBanco()
var Cns, Cbo, Cnes, Ine, cpfLogin, senhaLogin, usuarioLogin, primeiraletraLogin string
var qtdBaixo, qtdMedio, qtdAlto, qtdTotal int
var templates = template.Must(template.ParseFiles("./index.html", "./templates/telalogin/login.html", "./templates/telaesqueceusenha/esqueceusenha.html", "./templates/dashboard/dashboard.html", "./templates/telaesqueceusenha/cpfinvalido.html", "./templates/telalogin/logininvalido.html", "./templates/formulario/formulario.html", "./templates/formulario/cadastroinvalido1.html", "./templates/formulario/formulariofeito.html", "./templates/central-usuario/centralusuario.html", "./templates/pacientesgerais/indexPacGerais.html", "./templates/pg-baixo/pg-baixo.html", "./templates/dashboard/dashboardv2.html", "./templates/pg-medio/pg-medio.html", "./templates/pg-alto/pg-alto.html", "./templates/pg-absenteista/pg-absenteista.html", "./templates/pag-Faq/indexFaq.html", "./templates/formulario-preenchido/formpreenchido.html"))

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/login", autenticaCadastroELevaAoLogin)
	http.HandleFunc("/logininvalido", loginInvalido)
	http.HandleFunc("/dashboard", autenticaLoginELevaAoDashboard)
	http.HandleFunc("/dashboard/voltar", dashboard)
	http.HandleFunc("/esqueceusenha", executarEsqueceuSenha)
	http.HandleFunc("/atualizarinvalido", atualizarSenhaInvalido)
	http.HandleFunc("/telalogin", atualizarSenha)
	http.HandleFunc("/cadastrar-paciente", executarFormulario)
	http.HandleFunc("/paciente-cadastrado", cadastrarPaciente)
	colocarDados()
	http.HandleFunc("/central-usuario", executarCentralUsuario)
	http.HandleFunc("/pagina-faq", executarPagFaq)
	http.HandleFunc("/pacientesgerais", executarPacGerais)
	http.HandleFunc("/pagina-baixo-risco", executarPgBaixo)
	http.HandleFunc("/pagina-medio-risco", executarPgMedio)
	http.HandleFunc("/pagina-medio-risco/filtrar", executarPgMedioFiltro)
	http.HandleFunc("/pagina-alto-risco", executarPgAlto)
	http.HandleFunc("/pagina-alto-risco/filtrar", executarPgAltoFiltro)
	http.HandleFunc("/pagina-absenteista", executarPgAbsenteista)
	http.HandleFunc("/formulario/preenchido", executarFormPreenchdio)

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

	_, err = database.Query("CREATE TABLE IF NOT EXISTS pacientes(id SERIAL PRIMARY KEY, nome_completo VARCHAR(255) UNIQUE NOT NULL, data_nasc VARCHAR(30), cpf VARCHAR(15) UNIQUE NOT NULL, nome_mae VARCHAR(255) UNIQUE NOT NULL, sexo VARCHAR(30), cartao_sus VARCHAR(55) UNIQUE NOT NULL, telefone VARCHAR(55) UNIQUE NOT NULL, email VARCHAR(255) UNIQUE NOT NULL, cep VARCHAR(15) UNIQUE NOT NULL, bairro VARCHAR(255), rua VARCHAR(255), numero VARCHAR(255), complemento VARCHAR(255), homem VARCHAR(15) NOT NULL, etilista VARCHAR(15) NOT NULL, tabagista VARCHAR(15) NOT NULL, lesao_bucal VARCHAR(15) NOT NULL)")
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

func loginInvalido(w http.ResponseWriter, _ *http.Request){
	err := templates.ExecuteTemplate(w, "logininvalido.html", "a")
	if err != nil{
		return
	}
}

func colocarDados(){
	pegardados, err := db.Query("SELECT homem, etilista, tabagista, lesao_bucal FROM pacientes")
	if err != nil{
		return
	}
	defer pegardados.Close()
	armazenamento := make([]PegarDados, 0)

	for pegardados.Next(){
		armazenar := PegarDados{}
		err := pegardados.Scan(&armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil{
			log.Println(err)
			return
		}
		armazenamento = append(armazenamento, armazenar)
	}
	if err = pegardados.Err(); err != nil{
		return
	}
	pgbaixo := &qtdBaixo
	pgmedio := &qtdMedio
	pgalto := &qtdAlto
	pgtotal := &qtdTotal
	for _, armazenado := range armazenamento{
		if armazenado.Tabagista == "Não" && armazenado.LesaoBucal == "Não"{
			*pgbaixo++
		} else if armazenado.Tabagista == "Sim" && armazenado.LesaoBucal == "Não"{
			*pgmedio++
		} else if armazenado.LesaoBucal == "Sim"{
			*pgalto++
		}
	}
	*pgtotal = *pgbaixo + *pgmedio + *pgalto
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
	cpf := &cpfLogin
	senha := &senhaLogin
	*cpf = r.FormValue("cpf")
	*senha = r.FormValue("senha")
	cpfsenha, err := db.Query("SELECT nome_completo, cpf, senha, cns, cbo, cnes, ine FROM cadastro")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cpfsenha.Close()
	armazenamento := make([]validarlogin, 0)

	for cpfsenha.Next(){
		armazenar := validarlogin{}
		err := cpfsenha.Scan(&armazenar.Usuario, &armazenar.Cpf, &armazenar.Senha, &armazenar.Cns, &armazenar.Cbo, &armazenar.Cnes, &armazenar.Ine)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		armazenar.QtdBaixo = qtdBaixo
		armazenar.QtdMedio = qtdMedio
		armazenar.QtdAlto = qtdAlto
		porcbaixo := (float64(qtdBaixo)/float64(qtdTotal))*100
		porcmedio := (float64(qtdMedio)/float64(qtdTotal))*100
		porcalto := (float64(qtdAlto)/float64(qtdTotal))*100
		porcbaixo = float64(int(porcbaixo * 100)) / 100
		porcmedio = float64(int(porcmedio * 100)) / 100
		porcalto = float64(int(porcalto * 100)) / 100
		armazenar.PorcBaixo = porcbaixo
		armazenar.PorcMedio = porcmedio
		armazenar.PorcAlto = porcalto
		armazenamento = append(armazenamento, armazenar)
	}
	if err = cpfsenha.Err(); err != nil{
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, armazenado := range armazenamento{
		if armazenado.Cpf == cpfLogin && armazenado.Senha == senhaLogin{
			armazenador := &primeiraletraLogin
			armazenador2 := &usuarioLogin
			armazenado.PrimeiraLetra = string(armazenado.Usuario[0])
			*armazenador = string(armazenado.Usuario[0])
			quebrado := strings.Split(armazenado.Usuario, " ")
			armazenado.Usuario = quebrado[0]
			*armazenador2 = armazenado.Usuario
			cns := &Cns
			cbo := &Cbo
			cnes := &Cnes
			ine := &Ine
			*ine = armazenado.Ine
			*cns = armazenado.Cns
			*cnes = armazenado.Cnes
			*cbo = armazenado.Cbo
			err = templates.ExecuteTemplate(w, "dashboard.html", armazenado)
			if err != nil{
				return
			}
			return
		}
	}
	http.Redirect(w, r, "/logininvalido", http.StatusSeeOther)
}

func dashboard(w http.ResponseWriter, r *http.Request){
	porcbaixo := (float64(qtdBaixo)/float64(qtdTotal))*100
	porcmedio := (float64(qtdMedio)/float64(qtdTotal))*100
	porcalto := (float64(qtdAlto)/float64(qtdTotal))*100
	porcbaixo = float64(int(porcbaixo * 100)) / 100
	porcmedio = float64(int(porcmedio * 100)) / 100
	porcalto = float64(int(porcalto * 100)) / 100
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin, QtdBaixo: qtdBaixo, QtdMedio: qtdMedio, QtdAlto: qtdAlto, PorcBaixo: porcbaixo, PorcMedio: porcmedio, PorcAlto: porcalto}
	err := templates.ExecuteTemplate(w, "dashboardv2.html", u)
	if err != nil{
		return
	}
}

func executarEsqueceuSenha(w http.ResponseWriter, _ *http.Request){
	err := templates.ExecuteTemplate(w, "esqueceusenha.html", "a")
	if err != nil{
		return
	}
}

func atualizarSenhaInvalido(w http.ResponseWriter, _ *http.Request){
	err := templates.ExecuteTemplate(w, "cpfinvalido.html", "a")
	if err != nil{
		return
	}
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
			err = templates.ExecuteTemplate(w, "login.html", "a")
			if err != nil{
				return
			}
			return
		}	
	}
	http.Redirect(w, r, "atualizarinvalido", http.StatusSeeOther)
}

func executarCentralUsuario(w http.ResponseWriter, r *http.Request){
	cpfsenha, err := db.Query("SELECT nome_completo, cpf, cns, cbo, cnes, ine, senha FROM cadastro")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer cpfsenha.Close()
	armazenamento := make([]ACS, 0)

	for cpfsenha.Next(){
		armazenar := ACS{}
		err := cpfsenha.Scan(&armazenar.NomeCompleto, &armazenar.CPF, &armazenar.CNS, &armazenar.CBO, &armazenar.CNES, &armazenar.INE, &armazenar.SenhaACS)
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
		if armazenado.CPF == cpfLogin && armazenado.SenhaACS == senhaLogin{
			armazenado.PrimeiraLetra = string(armazenado.NomeCompleto[0])
			armazenado.CPF = strings.ReplaceAll(armazenado.CPF, armazenado.CPF[:5], "*****")
			armazenado.CNS = strings.ReplaceAll(armazenado.CNS, armazenado.CNS[:5], "*****")
			armazenado.CNES = strings.ReplaceAll(armazenado.CNES, armazenado.CNES[:3], "***")
			quebrado2 := strings.Split(armazenado.NomeCompleto, " ")
			armazenado.User = quebrado2[0]
			quebrado := strings.Split(armazenado.SenhaACS, "")
			for i := 0; i < len(quebrado); i++ {
				armazenado.SenhaACS = strings.Replace(armazenado.SenhaACS, quebrado[i], "*", -1)
			}
			err = templates.ExecuteTemplate(w, "centralusuario.html", armazenado)
			if err != nil{
				return
			}
			return
		}
	}

}

func executarFormulario (w http.ResponseWriter, _ *http.Request){
	cnsq := strings.Split(Cns, "")
	cboq := strings.Split(Cbo, "")
	cnesq := strings.Split(Cnes, "")
	ineq := strings.Split(Ine, "")
	cbo1 := cboq[0]
	cbo2 := cboq[1]
	cbo3 := cboq[2]
	cbo4 := cboq[3]
	cbo5 := cboq[4]
	cbo6 := cboq[5]
	d := DadosForm{Cns: cnsq, Cbo1: cbo1, Cbo2: cbo2, Cbo3: cbo3, Cbo4: cbo4, Cbo5: cbo5, Cbo6: cbo6, Cnes: cnesq, Ine: ineq}
	err := templates.ExecuteTemplate(w, "formulario.html", d)
	if err != nil{
		return
	}
}

func cadastrarPaciente(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	nome := r.FormValue("nome")
	datanascimento := r.FormValue("datanascimento")
	cpf := r.FormValue("cpfpaciente")
	nomemae := r.FormValue("nomemae")
	sexo := r.FormValue("sexo")
	cartaosus := r.FormValue("cartaosus")
	telefone := r.FormValue("telefone")
	email := r.FormValue("email")
	cep := r.FormValue("cep")
	bairro := r.FormValue("bairro")
	rua := r.FormValue("rua")
	numero, _ := strconv.Atoi(r.FormValue("numero"))
	complemento := r.FormValue("complemento")
	homem := r.FormValue("tipo1")
	etilista := r.FormValue("tipo2")
	tabagista := r.FormValue("tipo3")
	lesao_bucal := r.FormValue("tipo4")
	cnsq := strings.Split(Cns, "")
	cboq := strings.Split(Cbo, "")
	cnesq := strings.Split(Cnes, "")
	ineq := strings.Split(Ine, "")
	cbo1 := cboq[0]
	cbo2 := cboq[1]
	cbo3 := cboq[2]
	cbo4 := cboq[3]
	cbo5 := cboq[4]
	cbo6 := cboq[5]
	d := DadosForm{Cns: cnsq, Cbo1: cbo1, Cbo2: cbo2, Cbo3: cbo3, Cbo4: cbo4, Cbo5: cbo5, Cbo6: cbo6, Cnes: cnesq, Ine: ineq}

	if homem != "" && etilista != "" && tabagista != "" && lesao_bucal != "" && sexo != "" && (homem != "Não" || etilista != "Não" || tabagista != "Não" || lesao_bucal != "Não"){
		_, err := db.Exec("INSERT INTO pacientes(nome_completo, data_nasc, cpf, nome_mae, sexo, cartao_sus, telefone, email, cep, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)", nome, datanascimento, cpf, nomemae, sexo, cartaosus, telefone, email, cep, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal)
		if err != nil{
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		err = templates.ExecuteTemplate(w, "formulariofeito.html", d)
		if err != nil{
			return
		}
	} else{
		err := templates.ExecuteTemplate(w, "cadastroinvalido1.html", d)
		if err != nil{
			return
		}
	}
	pgbaixo := &qtdBaixo
	pgmedio := &qtdMedio
	pgalto := &qtdAlto
	pgtotal := &qtdTotal
	if tabagista == "Não" && lesao_bucal == "Não"{
		*pgbaixo++
	} else if tabagista == "Sim" && lesao_bucal == "Não"{
		*pgmedio++
	} else if lesao_bucal == "Sim"{
		*pgalto++
	}
	*pgtotal = *pgalto + *pgmedio + *pgbaixo
}

func executarPagFaq(w http.ResponseWriter, _ *http.Request){
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin}
	err := templates.ExecuteTemplate(w, "indexFaq.html", u)
	if err != nil{
		return
	}
}

func executarPacGerais(w http.ResponseWriter, _ *http.Request){
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin, QtdBaixo: qtdBaixo, QtdMedio: qtdMedio, QtdAlto: qtdAlto}
	err := templates.ExecuteTemplate(w, "indexPacGerais.html", u)
	if err != nil{
		return
	}
}

func executarPgBaixo(w http.ResponseWriter, _ *http.Request){
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next(){
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "-")
		if armazenar.Complemento != ""{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Não" && armazenar.LesaoBucal == "Não"{
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" {
				armazenar.Fatores = "Homem/Etilista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não"{
				armazenar.Fatores = "Homem"
			} else if armazenar.Homem == "Não" && armazenar.Etilista == "Sim"{
				armazenar.Fatores = "Mulher/Etilista"
			}
			now := time.Now()
			ano, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			dia, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia){
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	err = templates.ExecuteTemplate(w, "pg-baixo.html", armazenamento)
	if err != nil{
		return
	}
}

func executarPgMedio(w http.ResponseWriter, _ *http.Request){
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next(){
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "-")
		if armazenar.Complemento != ""{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Sim" && armazenar.LesaoBucal == "Não"{
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Etilista/Tabagista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim"{
				armazenar.Fatores = "Homem/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Não"{
				armazenar.Fatores = "Mulher/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Sim"{
				armazenar.Fatores = "Mulher/Etilista/Tabagista"
			}
			now := time.Now()
			ano, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			dia, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia){
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	err = templates.ExecuteTemplate(w, "pg-medio.html", armazenamento)
	if err != nil{
		return
	}
}

func executarPgMedioFiltro(w http.ResponseWriter, r *http.Request){
	sexo := r.FormValue("tipo1")
	idade := r.FormValue("tipo2")
	etilista := r.Form.Has("etilista")
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next(){
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "-")
		if armazenar.Complemento != ""{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.Tabagista == "Sim" && armazenar.LesaoBucal == "Não"{
			if armazenar.Homem == "Sim" && armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim" {
				armazenar.Fatores = "Homem/Etilista/Tabagista"
			} else if armazenar.Homem == "Sim" && armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim"{
				armazenar.Fatores = "Homem/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Não"{
				armazenar.Fatores = "Mulher/Tabagista"
			} else if armazenar.Homem == "Não" && armazenar.Tabagista == "Sim" && armazenar.Etilista == "Sim"{
				armazenar.Fatores = "Mulher/Etilista/Tabagista"
			}
			now := time.Now()
			ano, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			dia, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia){
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	var armazenadoPgMedio []Pacientes
	for _, armazenado := range armazenamento{
		if idade == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50{
			if !etilista{
				if armazenado.Etilista == "Não" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			} else if etilista{
				if armazenado.Etilista == "Sim" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
		if idade == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60{
			if !etilista{
				if armazenado.Etilista == "Não" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			} else if etilista{
				if armazenado.Etilista == "Sim" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
		if idade == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70{
			if !etilista{
				if armazenado.Etilista == "Não" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			} else if etilista{
				if armazenado.Etilista == "Sim" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
		if idade == "70+" && armazenado.Idade > 70{
			if !etilista{
				if armazenado.Etilista == "Não" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			} else if etilista{
				if armazenado.Etilista == "Sim" && armazenado.Sexo == sexo{
					armazenadoPgMedio = append(armazenadoPgMedio, armazenado)
				}
			}
		}
	}
	err = templates.ExecuteTemplate(w, "pg-medio.html", armazenadoPgMedio)
	if err != nil{
		return
	}
}

func executarPgAlto(w http.ResponseWriter, _ *http.Request){
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next(){
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "-")
		if armazenar.Complemento != ""{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.LesaoBucal == "Sim"{
			if armazenar.Homem == "Sim"{
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Homem/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Homem/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não"{
					armazenar.Fatores = "Homem/Etilista/Feridas Bucais"
				}
			}
			if armazenar.Homem == "Não"{
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Mulher/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Mulher/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não"{
					armazenar.Fatores = "Mulher/Etilista/Feridas Bucais"
				}
			}
			now := time.Now()
			ano, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			dia, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia){
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	err = templates.ExecuteTemplate(w, "pg-alto.html", armazenamento)
	if err != nil{
		return
	}
}

func executarPgAltoFiltro(w http.ResponseWriter, r *http.Request){
	sexo := r.FormValue("tipo1")
	idade := r.FormValue("tipo2")
	etilista := r.Form.Has("etilista")
	tabagista := r.Form.Has("tabagista")
	feridasbucais := r.Form.Has("feridasbucais")
	pesquisa, err := db.Query("SELECT nome_completo, data_nasc, sexo, telefone, bairro, rua, numero, complemento, homem, etilista, tabagista, lesao_bucal FROM pacientes ORDER BY nome_completo")
	if err != nil{
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	defer pesquisa.Close()
	var armazenamento []Pacientes
	for pesquisa.Next(){
		armazenar := Pacientes{}
		err := pesquisa.Scan(&armazenar.Nome, &armazenar.DataNasc, &armazenar.Sexo, &armazenar.Telefone, &armazenar.Bairro, &armazenar.Rua, &armazenar.Numero, &armazenar.Complemento, &armazenar.Homem, &armazenar.Etilista, &armazenar.Tabagista, &armazenar.LesaoBucal)
		if err != nil{
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		quebrar := strings.Split(armazenar.DataNasc, "-")
		if armazenar.Complemento != ""{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro + "," + armazenar.Complemento
		} else{
			armazenar.Endereco = armazenar.Rua + "," + armazenar.Numero + "," + armazenar.Bairro
		}
		if armazenar.LesaoBucal == "Sim"{
			if armazenar.Homem == "Sim"{
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Homem/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Homem/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não"{
					armazenar.Fatores = "Homem/Etilista/Feridas Bucais"
				}
			}
			if armazenar.Homem == "Não"{
				if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Mulher/Etilista/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Não" && armazenar.Tabagista == "Sim"{
					armazenar.Fatores = "Mulher/Tabagista/Feridas Bucais"
				} else if armazenar.Etilista == "Sim" && armazenar.Tabagista == "Não"{
					armazenar.Fatores = "Mulher/Etilista/Feridas Bucais"
				}
			}
			now := time.Now()
			ano, _ := strconv.Atoi(quebrar[0])
			mes, _ := strconv.Atoi(quebrar[1])
			dia, _ := strconv.Atoi(quebrar[2])
			armazenar.Idade = now.Year() - ano
			if int(now.Month()) < mes || (int(now.Month()) == mes && now.Day() < dia){
				armazenar.Idade--
			}
			armazenar.Usuario = usuarioLogin
			armazenar.Primeira = primeiraletraLogin
			armazenamento = append(armazenamento, armazenar)
		}
	}
	var armazenadoPgAlto []Pacientes
	for _, armazenado := range armazenamento{
		if idade == "40-50" && armazenado.Idade >= 40 && armazenado.Idade <= 50{
			if !etilista && !tabagista && !feridasbucais && sexo == ""{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && !feridasbucais && armazenado.Sexo == sexo{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == "" && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == armazenado.Sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if etilista && !tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Não"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if !etilista && tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if etilista && tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
		if idade == "51-60" && armazenado.Idade > 50 && armazenado.Idade <= 60{
			if !etilista && !tabagista && !feridasbucais && sexo == ""{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && !feridasbucais && armazenado.Sexo == sexo{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == "" && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == armazenado.Sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if etilista && !tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Não"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if !etilista && tabagista {
				if armazenado.Sexo == sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if etilista && tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
		if idade == "61-70" && armazenado.Idade > 60 && armazenado.Idade <= 70{
			if !etilista && !tabagista && !feridasbucais && sexo == ""{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && !feridasbucais && armazenado.Sexo == sexo{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == "" && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == armazenado.Sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if etilista && !tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Não"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if !etilista && tabagista {
				if armazenado.Sexo == sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if etilista && tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
		if idade == "70+" && armazenado.Idade > 70{
			if !etilista && !tabagista && !feridasbucais && sexo == ""{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && !feridasbucais && armazenado.Sexo == sexo{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == "" && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if !etilista && !tabagista && feridasbucais && sexo == armazenado.Sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
				armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
			}
			if etilista && !tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Não"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if !etilista && tabagista {
				if armazenado.Sexo == sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			} else if etilista && tabagista{
				if armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Sim"{
					armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
				}
			}
		}
		if idade == "" && etilista && !tabagista && armazenado.Sexo == sexo && armazenado.Etilista == "Sim" && armazenado.Tabagista == "Não"{
			armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
		}
		if idade == "" && !etilista && tabagista && armazenado.Sexo == sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Sim"{
			armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
		}
		if idade == "" && !etilista && !tabagista && feridasbucais && armazenado.Sexo == sexo && armazenado.Etilista == "Não" && armazenado.Tabagista == "Não"{
			armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
		}
		if idade == "" && !etilista && !tabagista && !feridasbucais && armazenado.Sexo == sexo{
			armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
		}
		if idade == "" && !etilista && !tabagista && !feridasbucais && sexo == ""{
			armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
		}
		armazenado.Usuario = usuarioLogin
		armazenado.Primeira = primeiraletraLogin
		armazenadoPgAlto = append(armazenadoPgAlto, armazenado)
	}

	err = templates.ExecuteTemplate(w, "pg-alto.html", armazenadoPgAlto)
	if err != nil{
		return
	}
}

func executarPgAbsenteista(w http.ResponseWriter, _ *http.Request){
	u := UsuarioNoDashboard{Usuario: usuarioLogin, Primeira: primeiraletraLogin}
	err := templates.ExecuteTemplate(w, "pg-absenteista.html", u)
	if err != nil{
		return
	}
}

func executarFormPreenchdio(w http.ResponseWriter, _ *http.Request){
	err := templates.ExecuteTemplate(w, "formpreenchido.html", "a")
	if err != nil{
		return
	}
}
