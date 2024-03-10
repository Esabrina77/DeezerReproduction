package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"groupieTrack/manager"
	inittemplate "groupieTrack/templates"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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
	inittemplate.Temp.ExecuteTemplate(w, "radio", radioResp.Data)
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

	// Utilisez votre modèle de page pour afficher les données
	inittemplate.Temp.ExecuteTemplate(w, "album", album)
}

// FONCTION  DE RECHERCHE
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	page := getPageQueryParam(r)
	itemsPerPage := 10

	// Effectuer la recherche à l'aide de l'API Deezer
	results, _, err := performSearch(query)
	if err != nil {
		http.Error(w, "Erreur de recherche", http.StatusInternalServerError)
		return
	}

	// Paginer les résultats
	paginatedResults := paginateResults(results, page, itemsPerPage)

	// Construire la requête de redirection avec les résultats de recherche
	redirectURL := fmt.Sprintf("/search?query=%s&page=%d&results=%s", query, page, strings.Join(paginatedResults, ","))

	// Rediriger l'utilisateur vers la page de recherche avec les résultats
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func getPageQueryParam(r *http.Request) int {
	pageParam := r.URL.Query().Get("page")
	if pageParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err == nil {
			return page
		}
	}
	return 1 // Page par défaut si le paramètre est manquant ou invalide
}

func performSearch(query string) ([]string, int, error) {
	url := "https://api.deezer.com/search?q=" + url.QueryEscape(query)
	response, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()
	// Vérification du code de statut de la réponse
	if response.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("la requete a l'api de deezer avec le code de statut : %d", response.StatusCode)
	}
	// Lecture de la réponse JSON
	var deezerResponse struct {
		Data []struct {
			Title string `json:"title"`
		} `json:"data"`
		Total int `json:"total"`
	}

	err = json.NewDecoder(response.Body).Decode(&deezerResponse)
	if err != nil {
		return nil, 0, err
	}

	// Extraction des résultats de recherche et du nombre total de résultats
	results := make([]string, 0, len(deezerResponse.Data))
	for _, item := range deezerResponse.Data {
		results = append(results, item.Title)
	}

	totalResults := deezerResponse.Total

	return results, totalResults, nil
}

func paginateResults(results []string, page, itemsPerPage int) []string {
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage

	// Vérifier les limites de pagination
	if startIndex >= len(results) {
		return []string{} // Aucun résultat pour cette page
	}

	if endIndex > len(results) {
		endIndex = len(results)
	}

	return results[startIndex:endIndex]
}

func getTotalPages(totalResults, itemsPerPage int) int {
	totalPages := totalResults / itemsPerPage
	if totalResults%itemsPerPage != 0 {
		totalPages++
	}
	return totalPages
}

type SearchResponse struct {
	Results     []string `json:"results"`
	TotalPages  int      `json:"total_pages"`
	CurrentPage int      `json:"current_page"`
}
