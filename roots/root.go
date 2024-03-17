package roots

import (
	"groupieTrack/controller"
	initTemplate "groupieTrack/templates"
	"log"
	"net/http"
)

func InitServe() {

	FileServer := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", FileServer))
	http.HandleFunc("/home", controller.HomeHandler)

	http.HandleFunc("/genre", controller.GenreHandler)
	http.HandleFunc("/allArtists", controller.AllArtistsHandler)

	http.HandleFunc("/artist/", controller.ArtistHandler)
	http.HandleFunc("/search", controller.SearchHandler)
	http.HandleFunc("/radio", controller.RadiosHandler)
	http.HandleFunc("/editorial", controller.EditorialsHandler)
	http.HandleFunc("/album/", controller.AlbumHandler)
	http.HandleFunc("/filtered-radios", controller.FilteredRadiosHandler)
	http.HandleFunc("/connexion", controller.ConnexionHandler)
	http.HandleFunc("/inscription", controller.InscriptionHandler)
	http.HandleFunc("/treatmentI", controller.TreatInscriptionHandler)
	http.HandleFunc("/treatmentC", controller.TreatConnexionHandler)
	http.HandleFunc("/404", controller.NotFoundHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		initTemplate.Temp.ExecuteTemplate(w, "404", nil)
	})
	if err := http.ListenAndServe(controller.Port, nil); err != nil {
		log.Fatal(err)
	}
}
