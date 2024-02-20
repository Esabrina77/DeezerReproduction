package main

import (
	"fmt"
	"net/http"

	controller "GROUPIE/controller" // Mise en place du package controller
)

func main() {

	colorGreen := "\033[32m" // Mise en place de couleur pour la lisibilité dans le terminal
	colorBlue := "\033[34m"
	colorYellow := "\033[33m"

	fmt.Println(string(colorBlue), "[SERVER_INFO] : Starting local Server...") // Information lorsque que le serveur est lancé

	fs := http.FileServer(http.Dir("static")) // Mise en place d'un dossier qui contient le CSS/IMG/JS
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Création des différentes routes, pour y accéder le serveur va chercher le package "controller" dans les dossiers et exécuter la fonction demandé.
	http.HandleFunc("/api/relation/", controller.RelationData)
	http.HandleFunc("/map", controller.Map)
	http.HandleFunc("/search", controller.Search)
	http.HandleFunc("/api/artists", controller.Artists)
	http.HandleFunc("/api/locations", controller.Locations)
	http.HandleFunc("/api/dates", controller.Dates)
	http.HandleFunc("/api/relation", controller.Relation)
	http.HandleFunc("/", controller.Accueil)

	fmt.Println(string(colorGreen), "[SERVER_READY] : on http://localhost:8000 ✅ ")    // Mise en place de l'URL pour l'utilisateur
	fmt.Println(string(colorYellow), "[SERVER_INFO] : To stop the program : Ctrl + c") // Information pour couper le serveur
	http.ListenAndServe(":8000", nil)                                                  // Mise en place du port
}
