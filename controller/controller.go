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
	// htgotts "github.com/hegedustibor/htgo-tts"
)

const (
	Port            = ":8080"
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

// verif de la session
func Connected(w http.ResponseWriter, r *http.Request) {
	//user non connecté
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["pseudo"] == nil {
		http.Redirect(w, r, "/connexion?error=SESSION_INVALID", http.StatusSeeOther)
		return
	}
}

// commande vocale pour détecter les actions fait sur le site
// func WelcomeUserServe(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
// 	pseudo := session.Values["pseudo"].(string)
// 	welcomeMessage := fmt.Sprintf(" l'utilisateur %s vient de connecter", pseudo)

// 	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("(New-Object -ComObject SAPI.SpVoice).Speak('%s')", welcomeMessage))
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println("Error executing PowerShell command:", err)
// 	}
// }

// AboutUsHandler
func AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	inittemplate.Temp.ExecuteTemplate(w, "aboutUs", nil)
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
	// Vérifier si une erreur est présente dans la requête
	errorMessage := r.URL.Query().Get("error")
	data := struct {
		Error string
	}{
		Error: errorMessage,
	}

	// Exécuter le modèle en passant les données
	inittemplate.Temp.ExecuteTemplate(w, "inscription", data)
}
func TreatInscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	// Récupérer les données du formulaire d'enregistrement
	email := r.FormValue("email")
	password := r.FormValue("password")
	pseudo := r.FormValue("pseudo")

	// Récupérer tous les utilisateurs existants
	users := manager.RetrieveUser()
	// Vérifier l'unicité du pseudo, de l'e-mail et du mot de passe
	if !manager.IsUnique(email, pseudo, users) {
		http.Redirect(w, r, "/inscription?error=already_registered", http.StatusFound)
		return
	}
	// Enregistrer le nouvel utilisateur
	manager.MarkLogin(email, password, pseudo)

	// Créer une nouvelle session pour l'utilisateur
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}

	session.Values["pseudo"] = pseudo
	session.Values["email"] = email
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Erreur lors de l'enregistrement de la session", http.StatusInternalServerError)
		return
	}
	// Rediriger l'utilisateur vers la page d'accueil
	http.Redirect(w, r, "/home?source=inscription", http.StatusFound)
}
func TreatConnexionHandler(w http.ResponseWriter, r *http.Request) {
	var session *sessions.Session
	password := r.FormValue("password")
	pseudo := r.FormValue("pseudo")

	// fmt.Println("l' email:", email)
	fmt.Println("le password:", password)
	users := manager.RetrieveUser()
	var login bool

	for _, user := range users {
		if user.Password == password && user.Pseudo == pseudo {
			//verifier si le login est correcte
			login = true
			break
		}
	}
	if login {
		i := 0
		//Creer une nouvelle session & stocker le pseudo
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
		// Stockez le message d'accueil dans une structure de données accessible au modèle HTML
		http.Redirect(w, r, "/home?source=connexion", http.StatusFound)
	} else {
		//rediriger vers la page de connexion avec un message d'erreur
		http.Redirect(w, r, "/connexion?error=invalid_login_try_again", http.StatusFound)
	}
}

// RADIO DEEZER
func RadiosHandler(w http.ResponseWriter, r *http.Request) {
	//verif de la session
	Connected(w, r)
	//user connecté
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

// EDITORIAL DEEZER
func EditorialsHandler(w http.ResponseWriter, r *http.Request) {
	Connected(w, r)
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
	Connected(w, r)
	inittemplate.Temp.ExecuteTemplate(w, "genre", nil)
}

// page d'acceuil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Connected(w, r)
	// Récupérer le pseudo de la session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}
	pseudo := session.Values["pseudo"].(string)

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
	//bienvenue au client
	data := struct {
		TopTrack    manager.Tracks
		TopArtists  []manager.Artist
		Description string
		Pseudo      string
	}{
		TopTrack:    topTrack,
		TopArtists:  topArtists,
		Description: description.Phrases[randomIndex],
		Pseudo:      pseudo,
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
	Connected(w, r)
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
	Connected(w, r)
	// Récupérer le pseudo de la session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}
	pseudo := session.Values["pseudo"].(string)

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

	// Faites une autre requête HTTP pour récupérer les pistes de l'album
	tracksResponse, err := http.Get(album.Tracklist)
	if err != nil {
		log.Println("Erreur lors de la requête pour les pistes de l'album :", err)
		http.Error(w, "Erreur lors de la récupération des pistes de l'album", http.StatusInternalServerError)
		return
	}
	defer tracksResponse.Body.Close()

	// Analyse de la réponse JSON dans une structure temporaire
	var tracksResponseData struct {
		Tracks []manager.Tracks `json:"data"`
	}
	err = json.NewDecoder(tracksResponse.Body).Decode(&tracksResponseData)
	if err != nil {
		log.Println("Erreur lors de l'analyse de la réponse JSON des pistes de l'album :", err)
		http.Error(w, "Erreur lors de l'analyse de la réponse JSON des pistes de l'album", http.StatusInternalServerError)
		return
	}

	// user du Lmodèle de page pour afficher les données de l'album et les pistes associées
	data := struct {
		Album  manager.Album
		Tracks []manager.Tracks
		Pseudo string
	}{
		Album:  album,
		Tracks: tracksResponseData.Tracks,
		Pseudo: pseudo,
	}
	inittemplate.Temp.ExecuteTemplate(w, "album", data)
}

// fonctionnalité de recherche
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
	Connected(w, r)
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
	// SearchResultsquerVoice(query)
	inittemplate.Temp.ExecuteTemplate(w, "search", result)
}

//	func SearchResultsquerVoice(query string) {
//		Message := fmt.Sprintf(" Résultats des recherches pour: %s", query)
//		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("(New-Object -ComObject SAPI.SpVoice).Speak('%s')", Message))
//		err := cmd.Run()
//		if err != nil {
//			fmt.Println("Error executing PowerShell command:", err)
//		}
//	}
//
// Gestionnaire pour ajouter ou retirer des favoris
func AddHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AddHandler")
	Connected(w, r)
	// Récupérer le pseudo de la session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}
	pseudo := session.Values["pseudo"].(string)

	// Vérifier la méthode de la requête
	switch r.Method {
	case http.MethodPost:
		log.Println("Ajout de la music...")
		// Lire le corps de la requête JSON
		var data manager.Favori
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Erreur de décodage du JSON", http.StatusBadRequest)
			return
		}

		// Lire les données actuelles du fichier Liked.json
		users, err := manager.ReadLikedFile()
		if err != nil {
			http.Error(w, "Erreur de lecture du fichier Liked.json", http.StatusInternalServerError)
			return
		}

		// Trouver l'utilisateur dans la liste des utilisateurs ou créer un nouvel utilisateur
		userIndex := manager.FindUser(users, pseudo)
		if userIndex == -1 {
			// Créer un nouvel utilisateur avec aucun favori
			newUser := manager.User{
				Pseudo:  pseudo,
				Favoris: []manager.Favori{},
			}
			users = append(users, newUser)
			userIndex = len(users) - 1
		}

		// Ajouter le favori à l'utilisateur
		addFavorite(&users[userIndex], data)
		log.Println("Favori ajouté avec succès")
		// Enregistrer les modifications dans le fichier Liked.json
		err = manager.WriteLikedFile(users)
		if err != nil {
			http.Error(w, "Erreur d'écriture du fichier Liked.json", http.StatusInternalServerError)
			return
		}

		// Répondre avec un message de succès
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Favori ajouté avec succès")
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
}

// Fonction pour ajouter un favori à un utilisateur
func addFavorite(user *manager.User, data manager.Favori) {
	// Ajouter le favori à l'utilisateur
	user.Favoris = append(user.Favoris, data)
}

func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveHandler")
	Connected(w, r)
	// Récupérer le pseudo de la session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}
	pseudo := session.Values["pseudo"].(string)

	// Obtenir l'ID de la musique de l'URL
	idParts := strings.Split(r.URL.Path, "/")
	idMusic := idParts[len(idParts)-1]

	// Lire les données actuelles du fichier Liked.json
	users, err := manager.ReadLikedFile()
	if err != nil {
		http.Error(w, "Erreur de lecture du fichier Liked.json", http.StatusInternalServerError)
		return
	}
	userIndex := manager.FindUser(users, pseudo)
	if userIndex == -1 {
		http.Error(w, "Utilisateur non trouvé", http.StatusBadRequest)
		return
	}

	// Supprimer le favori de l'utilisateur en utilisant l'ID de la musique
	removeFavorite(&users[userIndex], idMusic)
	// Enregistrer les modifications dans le fichier Liked.json
	err = manager.WriteLikedFile(users)
	if err != nil {
		http.Error(w, "Erreur d'écriture du fichier Liked.json", http.StatusInternalServerError)
		return
	}

	// Répondre avec un message de succès
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Favori supprimé avec succès")
}

// Fonction pour retirer un favori d'un utilisateur
func removeFavorite(user *manager.User, idMusic string) {
	// Créer une nouvelle liste de favoris pour l'utilisateur
	var newFavorites []manager.Favori
	// Trouver et retirer le favori de l'utilisateur avec l'ID de la musique
	for _, f := range user.Favoris {
		if f.IDMusic != idMusic {
			newFavorites = append(newFavorites, f)
		}
	}
	user.Favoris = newFavorites
}

func FavorisHandler(w http.ResponseWriter, r *http.Request) {
	Connected(w, r)
	// Récupérer le pseudo de la session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}
	pseudo := session.Values["pseudo"].(string)

	// Lire les données actuelles du fichier Liked.json
	users, err := manager.ReadLikedFile()
	if err != nil {
		http.Error(w, "Erreur de lecture du fichier Liked.json", http.StatusInternalServerError)
		return
	}

	// Trouver l'index de l'utilisateur dans la liste des utilisateurs
	userIndex := manager.FindUser(users, pseudo)
	if userIndex == -1 {
		http.Error(w, "Utilisateur non trouvé", http.StatusBadRequest)
		return
	}

	// Récupérer les favoris de l'utilisateur
	favoris := users[userIndex].Favoris

	inittemplate.Temp.ExecuteTemplate(w, "favoris", favoris)
}
