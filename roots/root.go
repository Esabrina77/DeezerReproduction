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
	http.HandleFunc("/aboutUs", controller.AboutUsHandler)
	http.HandleFunc("/artist/", controller.ArtistHandler)
	http.HandleFunc("/search", controller.SearchHandler)
	http.HandleFunc("/radio", controller.RadiosHandler)
	http.HandleFunc("/editorial", controller.EditorialsHandler)
	http.HandleFunc("/album/", controller.AlbumHandler)
	http.HandleFunc("/connexion", controller.ConnexionHandler)
	http.HandleFunc("/inscription", controller.InscriptionHandler)
	http.HandleFunc("/treatmentI", controller.TreatInscriptionHandler)
	http.HandleFunc("/treatmentC", controller.TreatConnexionHandler)
	http.HandleFunc("/404", controller.NotFoundHandler)
	http.HandleFunc("/add-remove", controller.AddHandler)
	http.HandleFunc("/remove/", controller.RemoveHandler)
	http.HandleFunc("/favoris", controller.FavorisHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		initTemplate.Temp.ExecuteTemplate(w, "404", nil)
	})
	http.HandleFunc("/error500", controller.Handle500)
	http.HandleFunc("/error301", controller.Handle301)
	http.HandleFunc("/error403", controller.Handle403)
	http.HandleFunc("/error503", controller.Handle503)
	http.HandleFunc("/error400", controller.Handle400)
	http.HandleFunc("/error505", controller.Handle505)

	  port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Fallback au port 8080 si la variable d'environnement PORT n'est pas définie
    }

    if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
        log.Fatal(err)
    }
}
