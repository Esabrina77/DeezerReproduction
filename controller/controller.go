package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"groupieTrack/manager"
	inittemplate "groupieTrack/templates"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
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

// tous les artistes
func AllArtistsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(deezerAPI + "/artist")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var artistResp manager.ArtistResponse
	if err := json.NewDecoder(response.Body).Decode(&artistResp); err != nil {
		log.Fatal(err)

	}
	inittemplate.Temp.ExecuteTemplate(w, "allArtists", artistResp.Data)
}

// RADIO DEEZER
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

	// Trier les radios par ordre alphabétique
	sort.Slice(radioResp.Data, func(i, j int) bool {
		return strings.ToLower(radioResp.Data[i].Title) < strings.ToLower(radioResp.Data[j].Title)
	})
	inittemplate.Temp.ExecuteTemplate(w, "radio", radioResp.Data)
}
func FilteredRadiosHandler(w http.ResponseWriter, r *http.Request) {
	letter := r.FormValue("letter")

	response, err := http.Get(deezerAPI + "/radio")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var radioResp manager.RadioResponse
	if err := json.NewDecoder(response.Body).Decode(&radioResp); err != nil {
		log.Fatal(err)
	}

	// Filtrer les radios en fonction de la lettre sélectionnée
	filteredRadios := make([]manager.Radio, 0)
	for _, radio := range radioResp.Data {
		if strings.HasPrefix(strings.ToLower(radio.Title), strings.ToLower(letter)) {
			filteredRadios = append(filteredRadios, radio)
		}
	}

	// Trier les radios par ordre alphabétique
	sort.Slice(filteredRadios, func(i, j int) bool {
		return strings.ToLower(filteredRadios[i].Title) < strings.ToLower(filteredRadios[j].Title)
	})

	inittemplate.Temp.ExecuteTemplate(w, "filtered-radios", filteredRadios)
}

// EDITORIAL DEEZER
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
	// for _, artist := range topArtists {
	// 	fmt.Println("Nom de l'artiste:", artist.Name)
	// 	fmt.Println("Photo de l'artiste:", artist.PictureMedium)
	// }
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

// PAGE CORREPONDANT A CHAQUE ARTISTE EN FONCTION DE L4ID
func extractArtistID(urlPath string) int {
	parts := strings.Split(urlPath, "/")
	if len(parts) < 3 {
		return -1
	}
	artistID, err := strconv.Atoi(parts[2])
	if err != nil {
		return -1
	}
	return artistID
}
func getArtistAlbums(artistID int) ([]manager.Album, error) {
	//REQUETE A L API POUR OBTENIR LES ALBUMS DE L'ARTISTE PAR SON ID
	albumURL := fmt.Sprintf("https://api.deezer.com/artist/%d/albums", artistID)

	response, err := http.Get(albumURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	//ANALYSE ET DECRYTAGE DE LA REPONSE JSON
	var albumResponse manager.AlbumResponse
	err = json.Unmarshal(body, &albumResponse)
	if err != nil {
		return nil, err
	}

	//RENVOYER LA LISTE DES ALBUMS DE L'ARTISTE
	return albumResponse.Data, nil
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistID := extractArtistID(r.URL.Path)
	if artistID == -1 {
		http.Error(w, "ID d'artiste invalide", http.StatusBadRequest)
		return
	}

	// Faites une requête HTTP pour récupérer les détails de l'artiste à partir de l'API Deezer
	response, err := http.Get("https://api.deezer.com/artist/" + strconv.Itoa(artistID))
	if err != nil {
		fmt.Println("Erreur lors de la requête :", err)
		http.Error(w, "Erreur lors de la requête", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Vérifiez le code de statut de la réponse
	if response.StatusCode != http.StatusOK {
		fmt.Println("Artiste non trouvé")
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}

	// Analysez la réponse JSON dans la structure `Artist`
	var artist manager.Artist
	err = json.NewDecoder(response.Body).Decode(&artist)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse de la réponse JSON :", err)
		http.Error(w, "Erreur lors de l'analyse de la réponse JSON", http.StatusInternalServerError)
		return
	}

	// Vérifiez si des données d'artiste ont été renvoyées
	if artist.ID == 0 {
		fmt.Println(artistID, "<---- ID")
		fmt.Println("Artiste non trouvé")
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}
	// Analysez la réponse JSON dans la structure `Album`
	albums, err := getArtistAlbums(artistID)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des albums :", err)
		http.Error(w, "Erreur lors de la récupération des albums :", http.StatusInternalServerError)
		return
	}
	// Créez une structure de données pour les données à renvoyer à la page
	data := struct {
		Artist manager.Artist
		Albums []manager.Album
	}{
		Artist: artist,
		Albums: albums,
	}
	// Utilisez votre modèle de page pour afficher les données
	inittemplate.Temp.ExecuteTemplate(w, "artist", data)
}

func extractAlbumID(urlPath string) int {
	parts := strings.Split(urlPath, "/")
	if len(parts) < 3 {
		return -1
	}
	albumID, err := strconv.Atoi(parts[2])
	if err != nil {
		return -1
	}
	return albumID
}
func AlbumHandler(w http.ResponseWriter, r *http.Request) {
	albumID := extractAlbumID(r.URL.Path)
	log.Println(albumID)
	if albumID == -1 {
		http.Error(w, "ID d'album invalide", http.StatusBadRequest)
		return
	}

	// Faites une requête HTTP pour récupérer les détails de l'album à partir de l'API Deezer
	response, err := http.Get("https://api.deezer.com/album/" + strconv.Itoa(albumID))
	if err != nil {
		log.Println("Erreur lors de la requête :", err)
		http.Error(w, "Erreur lors de la requête", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Vérifiez le code de statut de la réponse
	if response.StatusCode != http.StatusOK {
		log.Println("Aucun album trouvé")
		http.Error(w, "Album introuvable", http.StatusNotFound)
		return
	}

	// Analysez la réponse JSON dans la structure `Album`
	var album manager.Album
	err = json.NewDecoder(response.Body).Decode(&album)
	if err != nil {
		log.Println("Erreur lors de l'analyse de la réponse JSON :", err)
		http.Error(w, "Erreur lors de l'analyse de la réponse JSON", http.StatusInternalServerError)
		return
	}

	// Utilisez LE modèle de page pour afficher les données
	inittemplate.Temp.ExecuteTemplate(w, "album", album)
}

// fonctionnalité d e recherche
func search(query string, searchType string) (manager.SearchResult, error) {
	formattedQuery := strings.ReplaceAll(query, " ", "+")
	url := fmt.Sprintf("https://api.deezer.com/search?q=%s:\"%s\"", searchType, formattedQuery)
	response, err := http.Get(url)
	if err != nil {
		return manager.SearchResult{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return manager.SearchResult{}, err
	}

	var result manager.SearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return manager.SearchResult{}, err
	}

	return result, nil
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	searchType := r.FormValue("search_type")
	// Création de la structure SearchResult
	result := manager.SearchResult{
		Query:      query,
		SearchType: searchType,
	}

	// Effectuer la recherche et récupérer les résultats
	searchResult, err := search(query, searchType)
	if err != nil {
		http.Error(w, "Erreur lors de la recherche", http.StatusInternalServerError)
		return
	}
	// Ajouter les résultats à la structure SearchResult
	result.Data = searchResult.Data
	log.Print(result)
	inittemplate.Temp.ExecuteTemplate(w, "search", result)
}
