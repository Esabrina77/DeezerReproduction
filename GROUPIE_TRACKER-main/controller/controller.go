package controller

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func Accueil(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		codeErreur(w, r, 404, "Page not found")
		return
	}

	//Analyse du fichier accueil.html
	custTemplate, err := template.ParseFiles("./templates/accueil.html")

	if err != nil {
		//Gestion d'erreur
		codeErreur(w, r, 500, "Template not found : accueil.html")
		return
	}

	// Exécution de la d'accueil.html si il n'y a aucune erreur
	err = custTemplate.Execute(w, nil)
}

func Map(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/map" {
		codeErreur(w, r, 404, "Page not found")
		return
	}

	custTemplate, err := template.ParseFiles("./templates/map.html")

	if err != nil {
		codeErreur(w, r, 500, "Template not found : map.html")
		return
	}

	err = custTemplate.Execute(w, nil)
}

func Search(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/search" {
		codeErreur(w, r, 404, "Page not found")
		return
	}

	custTemplate, err := template.ParseFiles("./templates/search.html")

	if err != nil {
		codeErreur(w, r, 500, "Template not found : search.html")
		return
	}

	err = custTemplate.Execute(w, nil)
}

func loadApi(w http.ResponseWriter, r *http.Request, endpoint string) {
	// Tableau contenant le dernier élément de l'URL de chaque API
	tab := [4]string{"artists", "locations", "dates", "relation"}

	endpointIsValid := false

	// Vérification de la similarité du paramètre et des éléments du tableau
	for i := 0; i < len(tab); i++ {
		if endpoint == tab[i] {
			endpointIsValid = true
			break
		}
	}

	// Si le paramètres 'endpoint' ne correspond pas à un élément du tableau, alos on gère l'erreur
	if !endpointIsValid {
		// Gestion d'erreur 400
		codeErreur(w, r, 400, "Invalid endpoint")
		return
	}

	// Récupération de l'API voulu
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/" + endpoint)

	if err != nil {
		//Gestion d'erreur 500
		codeErreur(w, r, 500, "Server API is not responding")
		return
	}

	// Lecture de l'API
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		//Gestion d'erreur 500
		codeErreur(w, r, 500, "No data to sent")
		return
	}

	// On va donc ajouté l'API, en format JSON
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseData)
}

// Création de fonction pour chaque API
func Artists(w http.ResponseWriter, r *http.Request) {
	loadApi(w, r, "artists")
}

func Locations(w http.ResponseWriter, r *http.Request) {
	loadApi(w, r, "locations")
}

func Dates(w http.ResponseWriter, r *http.Request) {
	loadApi(w, r, "dates")
}

func Relation(w http.ResponseWriter, r *http.Request) {
	loadApi(w, r, "relation")
}

// Méthode qui récupère les caractéristiques d'un groupe en particulier
func getId(w http.ResponseWriter, r *http.Request, id string) {

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)

	if err != nil {
		// Gestion d'erreur 500
		codeErreur(w, r, 500, "Server API is not responding")
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		// Gestion d'erreur 500
		codeErreur(w, r, 500, "No data to sent")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseData)
}

// Methode qui va compléter la méthode getId
func RelationData(w http.ResponseWriter, r *http.Request) {
	pathPart := strings.Split(r.URL.Path, "/")
	getId(w, r, pathPart[len(pathPart)-1])
}

func codeErreur(w http.ResponseWriter, r *http.Request, status int, message string) {

	colorRed := "\033[31m" // Mise en place d'une couleur pour les erreurs

	// Mise en place de condition pour les différents types d'erreurs
	if status == 404 {
		http.Error(w, "404 not found", http.StatusNotFound)                                            // Mise en place d'un message qui sera afficher lors de l'erreur
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 404 : File not found, or missing...", message) // Message qui sera afficher sur le terminal avec une précision de l'erreur
	}
	if status == 400 {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 400 : Bad request", message)
	}
	if status == 500 {
		http.Error(w, "500 Internal server error", http.StatusInternalServerError)
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 500 : Internal server error", message)
	}

}
