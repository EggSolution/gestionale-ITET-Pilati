package webserver

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"net/http"
	"io/ioutil"
	"database/sql"
    _"github.com/go-sql-driver/mysql"
	"html/template"
	"github.com/EggSolution/gestionale-ITET-Pilati/moduli/database"
)

// variabile globale per al connessione al database
var InfoDB string
var Cwd string
// variabile per il numero di elaborati
var Nelaborati string
// struct per dahboard
type ElabStruct struct {
	Id          string
	Name        string
	Creator     string
	FilePath    string
	UploadDate  string
}
type ElabStructDash struct {
	Id          string
	Name        string
	Creator     string
	FilePath    string
	UploadDate  string
	Preferito   bool
}
type UserStruct struct {
	Id          string
	Name        string
	Privileges  string
	Date        string
	Password    string
	Email       string
	Nuovo       string
	Preferiti   string
}
type HomeStruct struct {
	Sezione      string
	Errore       string
}
type DashStruct struct {
	TitoloPag    string
	NuovoAcc     string
	IdUtente     string
	NomeUtente   string
	EmailUtente  string
	PassUtente   string
	Elaborati    []ElabStructDash
	Sezione      string
}

func main(){

}

func Routes(infoDB string){
	fileNelaborati, _ := os.Open("\\webserver\\var\\Nelaborati.txt")
	scanner := bufio.NewScanner(fileNelaborati)
	InfoDB = infoDB
	for scanner.Scan() {
		Nelaborati = scanner.Text()
	}
	// static file handling
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	// routes
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/passReset", passReset)
	http.HandleFunc("/uploadFile", uploadFile)
	http.HandleFunc("/cambioImpostazioni", cambioImpostazioni)
}

// all page function
func home(w http.ResponseWriter, r *http.Request){
	// divido l'url
	sezione := r.URL.Query().Get("sez")
	errore := r.URL.Query().Get("err")
	HomeVar := new(HomeStruct)
	HomeVar.Sezione = ""
	HomeVar.Errore = ""
	if sezione != "" {
		HomeVar.Sezione = sezione
	}
	if errore != "" {
		HomeVar.Errore = errore
	}
	// get current working directory
	Cwd, _ =  os.Getwd()
	// execute html template
	template, _ := template.ParseFiles(Cwd + "\\pagine\\home.html")
	template.Execute(w, HomeVar)
}

func dashboard(w http.ResponseWriter, r *http.Request){
	sezione := r.URL.Query().Get("sez")
	// get current working directory
	Cwd, _ =  os.Getwd()
	emailForm := ""
	passForm := ""
	switch r.Method {
		// filtro richieste
		case "POST":
			emailForm = r.FormValue("email")
			passForm = r.FormValue("password")
		case "GET":
			// execute html template
			http.Redirect(w, r, "http://localhost/?sez=1", http.StatusSeeOther)

			return
	}

	// connessione database
	DBconn, _ := sql.Open("mysql", InfoDB)
	// query al database
	credenziali, _ := DBconn.Query("SELECT * FROM user WHERE email='"+string(emailForm) + "' AND password='"+string(passForm)+"';")
	// divido la query
	credVar := database.QueryUser()
	// for che controlla tutti i risultati
	for credenziali.Next() {
		err := credenziali.Scan(&credVar.Id, &credVar.Name, &credVar.Privileges, &credVar.Date, &credVar.Password, &credVar.Email, &credVar.Nuovo, &credVar.Preferiti)
		if err != nil {
			panic(err)
		}
	}
	// creo gli arrey degli elaborati preferiti
	var preferiti []string
	preferiti = strings.Split(credVar.Preferiti, ",")
	// creo la struct da mettere nell HTML
	titoloP := "Dashboard - " + credVar.Name
	var ElaboratiStructData []ElabStructDash
	// raccolgo gli elaborati per renderizzarli nella dash
	var ElabStruct1 ElabStructDash
	elaboratiQueryData, _ := DBconn.Query("SELECT * FROM elaborati")
	for elaboratiQueryData.Next(){
		err := elaboratiQueryData.Scan(&ElabStruct1.Id, &ElabStruct1.Name, &ElabStruct1.Creator, &ElabStruct1.FilePath, &ElabStruct1.UploadDate)
		if err != nil {
			fmt.Println(err)
		}
		for i := 0; i < len(preferiti); i++ {
			fmt.Print(i)
			if preferiti[i] == ElabStruct1.Id {
				ElabStruct1.Preferito = true
			} else {
				ElabStruct1.Preferito = false
			}
		}
		ElaboratiStructData = append(ElaboratiStructData, ElabStruct1)
	}

	elaboratiHTML := DashStruct{
		TitoloPag: titoloP,
		NuovoAcc: credVar.Nuovo,
		IdUtente: credVar.Id,
		NomeUtente: credVar.Name,
		EmailUtente: credVar.Email,
		PassUtente: credVar.Password,
		Elaborati: ElaboratiStructData,
		Sezione: sezione,
	}

	// controlle se le credenziali esistono
	if credVar.Id == "" {
		// execute html template
			http.Redirect(w, r, "http://localhost/?sez=1&err=1", http.StatusSeeOther)
	} else {
		fmt.Println("Utente loggato:")
		fmt.Println("  -id: " + credVar.Id)
		fmt.Println("  -nome: " + credVar.Name)
		fmt.Println("  -email: " + credVar.Email + "\n")
		// aggiorno il profilo non più nuovo
		if credVar.Nuovo == "si" {
			DBconn.Query("UPDATE user SET nuovo = 'no'")
		}
		// execute html template
		template, _ := template.ParseFiles(Cwd + "\\pagine\\dashboard.html")
		template.Execute(w, elaboratiHTML)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			http.Redirect(w, r, "http://localhost/?sez=2", http.StatusSeeOther)
		case "POST":
			DBconn, _ := sql.Open("mysql", InfoDB)
			// prendo dati
			NomeUtente := r.FormValue("nome")
			MailUtente := r.FormValue("email")
			PasswordUtente := r.FormValue("password")
			// query per controllare omonimi
			utenti, _ := DBconn.Query("SELECT * FROM user WHERE name='"+NomeUtente+"' OR email='"+MailUtente+"';")
			for utenti.Next(){
				utStr := new(UserStruct)
				err := utenti.Scan(&utStr.Id, &utStr.Name, &utStr.Privileges, &utStr.Date, &utStr.Password, &utStr.Email, &utStr.Nuovo, &utStr.Preferiti)
				if err != nil {
					fmt.Println(err)
				}
				if utStr.Id != "" {
					// con errore
					http.Redirect(w, r, "http://localhost/?sez=2&err=2", http.StatusSeeOther)
					return
				}
			}
			// registrazione account
			_, err := DBconn.Query("INSERT INTO user (name, privileges, password, email, nuovo, preferiti) VALUES ('"+NomeUtente+"', '"+"3"+"', '"+PasswordUtente+"', '"+MailUtente+"', 'si', '');")
			if err != nil {
				fmt.Println(err)
				// con errore
				http.Redirect(w, r, "http://localhost/?sez=2&err=3", http.StatusSeeOther)
				return
			}
			fmt.Println("Utente registrato:")
			fmt.Println("  -nome: " + NomeUtente)
			fmt.Println("  -email: " + MailUtente + "\n")
			// redirect alla home
			http.Redirect(w, r, "http://localhost/?sez=1", http.StatusSeeOther)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// get current working directory
	Cwd, _ =  os.Getwd()
	switch r.Method {
		case "POST":
			// prendo il file
			file, handler, err := r.FormFile("file")
			if err != nil {
				fmt.Println("errore nel caricamento del file")
				fmt.Println(err)
			}
			defer file.Close()
			// aggiungo l'elaborato nel database
			name := r.FormValue("nomeElaborato")
			// VA MODIFICATO DOPO LA AGGIUNTA DELLE SESSIONI
			creator := "admin"
			filePath := "/elaborati/" + string(handler.Filename)
			DBconn, _ := sql.Open("mysql", InfoDB)
			query, _ := DBconn.Query("INSERT INTO elaborati (name, creator, filePath) VALUES ('"+name+"','"+creator+"','"+filePath+"');")
			fmt.Println(query)
			// creo l'elaborato
			tempFile, _ := ioutil.TempFile(Cwd + "\\elaborati\\", "elaborato-*.pdf")
			defer tempFile.Close()
			fileByte, _ := ioutil.ReadAll(file)
			tempFile.Write(fileByte)
			http.Redirect(w, r, "http://localhost/dashboard", http.StatusSeeOther)
		case "GET":
			// renderizzo la pagina (con errori)
			http.Redirect(w, r, "http://localhost/dashboard", http.StatusSeeOther)
	}
}

func passReset(w http.ResponseWriter, r *http.Request) {
	// get current working directory
	Cwd, _ =  os.Getwd()
	switch r.Method {
		case "POST":
			// execute html template
			template, _ := template.ParseFiles(Cwd + "\\pagine\\passDimenticata.html")
			template.Execute(w,"")
		case "GET":
			// execute html template
			template, _ := template.ParseFiles(Cwd + "\\pagine\\home.html")
			template.Execute(w,"")
	}

}

func cambioImpostazioni(w http.ResponseWriter, r *http.Request) {
	Cwd, _ =  os.Getwd()
	var VecchioEmail string
	var VecchioPass string
	var NuovoNome string
	var NuovoEmail string
	var NuovoPass string
	switch r.Method {
		case "POST":
			VecchioEmail = r.FormValue("emailOriginale")
			VecchioPass = r.FormValue("passOriginale")
			NuovoNome = r.FormValue("nomeUtente")
			NuovoEmail = r.FormValue("emailUtente")
			NuovoPass = r.FormValue("passUtente")
		case "GET":
			http.Redirect(w, r, "http.//localhost/dashboard", http.StatusSeeOther)
	}
	DBconn, _ := sql.Open("mysql", InfoDB)

	QueryAggiornamento, _ := DBconn.Query("UPDATE user SET name = '"+NuovoNome+"', email = '"+NuovoEmail+"', password = '"+NuovoPass+"' WHERE email = '"+VecchioEmail+"' and password = '"+VecchioPass+"';")
	fmt.Println(QueryAggiornamento)
	http.Redirect(w, r, "http://localhost/dashboard", http.StatusSeeOther)
}
