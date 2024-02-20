package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"groupieTrack/manager"
	inittemplate "groupieTrack/templates"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

const (
	Port            = "localhost:8080"
	deezerAPI       = "https://api.deezer.com"
	StripePublicKey = "pk_test_51OkP9QKa0BEOKwek4AcHZOLCTI4gsDDZSCzWGrRjQt8hHy8sCueAiNxxwnjbUAPfEEtOXRiJ72nF2oO5puW0G8oW00efoSjW1x"
	StripeSecretKey = "sk_test_51OkP9QKa0BEOKwekOTKYJaJQPDSwfMmT4Fb8PtYYKgixcOyL5II3106UbajXitNMxy4MUAs767XG21ZE8JId4wKt00El13BkiO"
)

var store = sessions.NewCookieStore([]byte(SecretKey()))

func SecretKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

// Sécurisation des routes/gestions des erreurs de chargement de pages
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	inittemplate.Temp.ExecuteTemplate(w, "404", nil)
}
func RessourceNotFoundHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "notFound", nil)
}
func ConnexionHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "connexion", nil)
}
func InscriptionHandler(w http.ResponseWriter, r *http.Request) {

	inittemplate.Temp.ExecuteTemplate(w, "inscription", nil)
}
func TreatInscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	//recupérer les données du formulaire d'enregistrement
	email := r.FormValue("email")
	password := r.FormValue("password")
	pseudo := r.FormValue("pseudo")

	//Enregistrer le nouvel Utilisateur
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if user.Email == email && user.Password == password && user.Pseudo == pseudo {
			//verifier si le login est déjà enregistré
			login = true
		}
	}
	if login {

		http.Redirect(w, r, "/connexion?error=already_registred", http.StatusFound)
	} else {
		//IL S AGIT D'UNE PREMIERE CONNEXION !
		//rediriger vers la page dc'acceuil & enregistrer le login
		manager.MarkLogin(email, password, pseudo)

		i := 0
		//Creer une nouvelle session & stocker l'email
		var err error
		session, err = store.Get(r, "session-name")
		for i > 1 {
			if err != nil {
				http.Error(w, "ERREUR DE SESSION_1", http.StatusInternalServerError)
				return
			}
		}

		session.Values["email"] = email
		fmt.Println("EMAIL RECU", email)
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERREUR DE SESSION_2", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home?success=Login_registred", http.StatusFound)
	}
}
func TreatConnexionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	//recupérer les données du formulaire de connexion
	// email := r.FormValue("email")
	password := r.FormValue("password")
	pseudo := r.FormValue("pseudo")

	// fmt.Println("l' email:", email)
	fmt.Println("le password:", password)
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if /*user.Email == email &&*/ user.Password == password && user.Pseudo == pseudo {
			//verifier si le login est correcte
			login = true
			break
		}
	}
	if login {
		i := 0
		//Creer une nouvelle session & stocker l'email
		var err error
		session, err = store.Get(r, "session-name")
		for i > 1 {
			if err != nil {
				http.Error(w, "ERREUR DE SESSION_1", http.StatusInternalServerError)
				return
			}
		}

		session.Values["pseudo"] = pseudo
		fmt.Println("PSEUDO RECU", pseudo)
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "ERREUR DE SESSION_2", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		//rediriger vers la page de connexion avec un message d'erreur
		http.Redirect(w, r, "/connexion?error=invalid_login_try_again", http.StatusFound)
	}
}

func RadiosHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(deezerAPI + "/radio")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var radioResp manager.RadioResponse
	if err := json.NewDecoder(response.Body).Decode(&radioResp); err != nil {
		log.Fatal(err)
	}
	inittemplate.Temp.ExecuteTemplate(w, "radio", radioResp.Data)
}
func EditorialsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(deezerAPI + "/editorial")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var editorialResp manager.EditorialResponse
	if err := json.NewDecoder(response.Body).Decode(&editorialResp); err != nil {
		log.Fatal(err)
	}
	inittemplate.Temp.ExecuteTemplate(w, "editorial", editorialResp.Data)
}
func GenreHandler(w http.ResponseWriter, r *http.Request) {
	inittemplate.Temp.ExecuteTemplate(w, "genre", nil)
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	inittemplate.Temp.ExecuteTemplate(w, "artist", nil)
}
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	inittemplate.Temp.ExecuteTemplate(w, "search", nil)
}

// page d'acceuil
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("https://api.deezer.com/chart")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var charts manager.Charts
	err = json.NewDecoder(response.Body).Decode(&charts)
	if err != nil {
		log.Fatal(err)
	}
	//Récuperer la musique la plus ecouter
	topTrack := charts.Tracks.Data[0]
	//Récuperer les 10 artistes les plus écouter en france
	topArtists := charts.Artist.Data[:10]
	for _, artist := range topArtists {
		fmt.Println("Nom de l'artiste:", artist.Name)
		fmt.Println("Photo de l'artiste:", artist.PictureMedium)
	}
	//Lecture du fichier Description.txt
	descript, err := os.ReadFile("Description.txt")
	if err != nil {
		log.Fatal(err)
	}
	description := manager.Description{
		Phrases: strings.Split(string(descript), "\n"),
	}
	//génération d'un nombre aléatiore
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(description.Phrases))

	data := struct {
		TopTrack    manager.Tracks
		TopArtists  []manager.Artist
		Description string
	}{
		TopTrack:    topTrack,
		TopArtists:  topArtists,
		Description: description.Phrases[randomIndex],
	}
	inittemplate.Temp.ExecuteTemplate(w, "home", data)
}

func SuccesHandler(w http.ResponseWriter, r *http.Request) {

	// Afficher la page HTML de confirmation
	fmt.Fprintf(w, "<h1>Paiement réussi</h1>")
	fmt.Fprintf(w, "<p>Merci pour votre paiement ! Votre commande a été traitée avec succès.</p>")

	inittemplate.Temp.ExecuteTemplate(w, "success", nil)
}

func CancelHandler(w http.ResponseWriter, r *http.Request) {

	// Afficher la page HTML d'annulation
	fmt.Fprintf(w, "<h1>Paiement annulé</h1>")
	fmt.Fprintf(w, "<p>Vous avez annulé votre paiement. Si vous avez des questions, veuillez nous contacter.</p>")
	inittemplate.Temp.ExecuteTemplate(w, "cancel", nil)
}
