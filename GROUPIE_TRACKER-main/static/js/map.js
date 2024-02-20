let tabMarker = [];
let marker;

// Fonction qui va initialiser la map, et qui va pointer en premier lieu sur Ynov Nantes
// Cette fonction va aussi permettre de faire pointer le marqueur sur un autre lieu grâce à la valeur entré dans l'input
function initMap() {
    const newMarker = document.getElementById('search').value
        // Transforme les noms de villes en coordonnées GPS
    const geocoder = new google.maps.Geocoder();
    let initialPlace = { lat: 47.205389240386936, lng: -1.539540743860519 }
    let newPlace = document.getElementById("submit")
        // Définit une map qui pointe sur Ynov Nantes
    const map = new google.maps.Map(document.getElementById("map"), {
        zoom: 8,
        center: initialPlace,
    });


    marker = new google.maps.Marker({
        position: initialPlace,
        title: newMarker,
        map: map,
        clickable: true,
    })

    tabMarker.push(marker)

    // Evènement qui va pointer sur un lieu si il y a un click
    newPlace.addEventListener("click", () => {
        marker.setMap(null)
        geocodeAddress(geocoder, map);
    });

}

// Cette fonction va permettre de traduire un nom de ville en coordonnées GPS, permettant de géolocaliser un lieu.
function geocodeAddress(geocoder, resultsMap) {
    const search = document.getElementById('search')
    geocoder.geocode({ address: search.value }, (results, status) => {
        // Si L'API google fonctionne, on va changer les coordonnées GPS 
        if (status === "OK") {
            resultsMap.setCenter(results[0].geometry.location);
            // On va recréer une map avec les nouvelles coordonnées GPS
            marker = new google.maps.Marker({
                map: resultsMap,
                position: results[0].geometry.location,
            });
            // Fonction qui créer une fenêtre d'information
            infoWindow(search, results)
                // Sinon message d'erreur 
        } else {
            alert("Error type :" + status);
        }
    });

}

// Fonction permettant de créer une petite fenêtre au-dessus du marqueur, lorsqu'on click sur le marqueur.
function infoWindow(search, results) {
    search.value = results[0].formatted_address
    marker.info = new google.maps.InfoWindow({
        content: '<b>' + search.value + '</b>'
    });

    google.maps.event.addListener(marker, 'click', function() {
        this.info.open(map, marker);
    });

}