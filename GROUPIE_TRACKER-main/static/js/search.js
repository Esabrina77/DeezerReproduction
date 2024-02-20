const search = document.getElementById('search'); // Permet d'avoir l'input de recherche et d'agir dessus
const matchList = document.getElementById('card-data'); // Permet d'indiquer ou se mettrons les cartes des artistes

// Permet de créer les différents chemins de l'API venant de notre serveur
const api = "/api/"

const artist = "artists"
const location_ = "locations"
const date = "dates"
const relation = "relation"


// Fonction asynchrone qui va récuperer les données de l'API de notre serveur et nous renvoie un "JSON Object"
const searchArtists = async searchText => {
    const res_artist = await fetch(api + artist);
    const artistsData = await res_artist.json();
    // Système de recheche avec regex qui agit sur l'API "artist"
    let matches = artistsData.filter(data => {
        const regex = new RegExp(`^${searchText}`, 'gi');

        let allMembers = ""
            // On parcourt le tableau de membres
        for (let index = 0; index < data.members.length; index++) {
            allMembers += data.members[index]
        }
        let resultOfMatches = data.name.match(regex) || (data.creationDate).toString().match(regex) || allMembers.match(regex) || (data.firstAlbum).toString().match(regex)
        return resultOfMatches
    });
    // Si on efface notre recherche, remise à zéro de l'affichage HTML
    if (searchText.length === 0) {
        matches = [];
        matchList.innerHTML = '';
    }
    // On renvoie les résultats
    outputHtml(matches);
}

// Mise en forme des résultats dans les cartes
const outputHtml = (matches) => {
    if (matches.length > 0) {
        const html = matches.map(match => `
        <div class="card" id="card">
            <div class="card-header" id="card-header">
                <img src="${match.image}" alt="">
            </div>
                <div class="card-body" id="card-body">
                    <ul>
                        <li><h4>Nom :</h4><br>${match.name}</li>
                        <br>
                        <li><h4>Date de création :</h4><br>${match.creationDate}</li>
                        <br>
                        <li><h4>Membres :</h4><br>${match.members}</li>
                        <br>
                        <li><h4>Premier album :</h4><br>${match.firstAlbum}</li>
                    </ul>
                    <div class="popup-header-cont">
                        <h3>${match.name}</h3>
                    </div>
                    <div class="read-more-cont">
                        <p class="relation" data-url="${match.relations}">...</p>
                        <button class="btn_map" type="button" onclick=redirectMap() >Accéder à la map</button>
                    </div>
                <button class="btn" type="button">Voir plus ...</button>
                </div>
        </div>
        `).join('');
        // Transfert sur le HTML
        let finalhtml = html;
        matchList.innerHTML = finalhtml;
    }
}

// Mise à jour de la recherche à chaque caractère entré par l'utilisateur
search.addEventListener('input', () => searchArtists(search.value))

const cardData = document.querySelector(".row");
const popup = document.querySelector(".popup-box");
const popupCloseBtn = popup.querySelector(".popup-close-btn")

// Création d'un évènement qui va afficher le contenu d'une API dans la pop-up lors d'un click
cardData.addEventListener("click", async function(event) {
    if (event.target.tagName.toLowerCase() == "button") {
        const item = event.target.parentElement;
        const relation = item.querySelector(".relation");
        const pathPart = relation.dataset.url.split("/");
        let res = await fetch(`/api/relation/${pathPart[pathPart.length-1]}`);
        let data = await res.json();
        elementAPI(data, relation);
        const h3 = item.querySelector(".popup-header-cont").innerHTML;
        const readMoreCont = item.querySelector(".read-more-cont").innerHTML;
        popup.querySelector(".popup-header").innerHTML = h3;
        popup.querySelector(".popup-body").innerHTML = readMoreCont
        popupBox();
    }
})

// Création d'un événement pour l'ouverture/fermeture de la pop-up
popupCloseBtn.addEventListener("click", popupBox);

popup.addEventListener("click", function(event) {
    if (event.target == popup) {
        popupBox();
    }
})

function popupBox() {
    popup.classList.toggle("open");
}

function elementAPI(elementJSON, relation) {
    // Transformer le JSON en string
    let json = JSON.stringify(elementJSON.datesLocations)
        //Analyse de la string créer précédemment
    let parseJSON = JSON.parse(json)
    let result = [];
    let index, resultpush

    // Récupération de chaque clé et valeur du fichier JSON
    for (index in parseJSON) {
        resultpush = index + " : " + parseJSON[index]
        result.push(resultpush)

    }

    relation.innerHTML = result.join(', ')

}

function redirectMap() {
    window.location.replace("/map")
}