package main
 
import (
    "net/http"
    "fmt"
    "log"
    "os"
    "html/template"
)
 
func main() {
    routes()
    acceso()
    
}
func acceso() {
    fmt.Print("acceso")
    porta := ":90"
    log.Fatal(http.ListenAndServe(string(porta), nil))
    
}
 func routes() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static", fs))
    http.HandleFunc("/pagine", pagine)
 }
 
 func pagine(w http.ResponseWriter, r *http.Request){
    fmt.Print("ciao1")
    Cwd, _ := os.Getwd()
    t, _ := template.ParseFiles(Cwd + "\\pagine\\login.html")
    t.Execute(w,"")
 }