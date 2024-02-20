const matchList = document.getElementById('card-data'); // Permet d'indiquer ou se mettrons les cartes des artistes

// Permet de créer les différents chemins de l'API venant de notre serveur
const api = "/api/"

const artist = "artists"
const location_ = "locations"
const date = "dates"
const relation = "relation"


// Fonction asynchrone qui s'exécute ligne par ligne 
async function data(url) {
    // Permet de récupérer les données de l'API en JSON
    const res_artist = await fetch(url);
    const dataArtist = await res_artist.json();
    // Appelle de function permettant la création de la carte avec son contenu mais aussi la création de la pop-up avec son contenu
    show(dataArtist);
}

data(api + artist);


var arr; // Notre tableau final avec les cartes des groupes
var mini; // Valeur du select de la date minimum
var maxi; // Valeur du select de la date maximum
var tab; // Un tableau dans lequel on va mettre toutes les cartes des artistes à aller chercher à la fin pour push dans arr
var ArrOfMembers; // Un tableau dans lequel est inscrit toutes les valeurs des checkbox pour savoir quel nombre de membre l'utilisateur veut.

function show(dataArtist) {
    if (dataArtist === undefined) { // dataArtist est undefined après utilisation d'un filtre donc on le redéfini s'il est undefined
        data(api + artist)
    }
    arr = []
    tab = []
    mini = document.getElementById("Selectmini").value
    maxi = document.getElementById("Selectmaxi").value
    tab = dataArtist.map(match => `
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
                </div>
            <button class="btn" type="button">Voir plus ...</button>
            </div>
    </div>
    `);
    SelectionArtist(dataArtist)
    let nbrPage = Math.ceil(arr.length / 10)
    pageNumber(nbrPage)
    matchList.innerHTML = arr.slice(0, 10).join(''); // On join notre tableau de 10 résultats par page pour avoir un affichage sans "," et correct
}
// Fonction pour initialiser la pop-up
function setupPopup() {
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
                popup.classList.toggle("open");
            }
        })
        // Création d'un événement pour l'ouverture/fermeture de la pop-up
    popupCloseBtn.addEventListener("click", () => { popup.classList.toggle("open"); });

    popup.addEventListener("click", function(event) {
        if (event.target == popup) {
            popup.classList.toggle("open");
        }
    })
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

var min = document.getElementById("mini")
var max = document.getElementById("maxi")

// Création menu déroulant, permettant le tri des artiste(s)/groupes selon la dates de création 
function minOrMax(param) {
    let select = document.createElement("select")
    select.id = "Select" + param.id
    param.appendChild(select)
    select.setAttribute("onchange", "show()")
    if (param == min) {
        minimum(select)
    }
    if (param == max) {
        maximum(select)
    }
}

function minimum(param) {
    for (let i = 1958; i <= 2021; i++) {
        var optionMin = document.createElement("option")
        optionMin.value = i
        optionMin.text = i
        param.appendChild(optionMin)
    }
}

function maximum(param) {
    for (let l = 2021; l >= 1958; l--) {
        var optionMax = document.createElement("option")
        optionMax.value = l
        optionMax.text = l
        param.appendChild(optionMax)
    }
}

minOrMax(min)
minOrMax(max)
    // Fonction permettant de filtrer le nombre de membres par groupes
function selectionArtistWithMembers() {
    let index
    for (index = 1; index <= 9; index++) {
        if (index <= 8) {
            //La valeur de chaque bouton seront donné grâce à index
            selectMembers(index, index)
        } else {
            // Appel de fonction selectMembers avec la valeur 0, qui changera avec la fonction pushOnArr
            selectMembers("AllMembers", 0)
        }
    }
}

// Fonction permettant d'initialisé les checkbox, leurs valeurs, leurs id ...
function selectMembers(name, value) {
    let MemberSort = document.getElementById("MemberSort")
    let check = document.createElement("input")
    let label = document.createElement("label")
    label.textContent = " " + name + " :"
    check.id = "Checkbox" + name
    check.value = value
    if (check.value == 0) {
        check.checked = true
    }
    check.setAttribute("type", "checkbox")
    check.setAttribute("onclick", "pushOnArr()")
    check.setAttribute("onclick", "show()")
    MemberSort.appendChild(label)
    MemberSort.appendChild(check)
}

selectionArtistWithMembers()

function SelectionArtist(dataArtist) {
    pushOnArr()
    for (let i = mini; i <= maxi; i++) { // i correspond à toutes les années entre le 1er select et le second
        for (let l = 0; l < dataArtist.length; l++) { // l correspond à chaque carte d'artistes
            if (dataArtist[l].creationDate == i) { // Si la date de creation du groupe de tel artiste = une année entre le mini et le maxi
                for (let ArrIndex = 0; ArrIndex < ArrOfMembers.length; ArrIndex++) { // ArrIndex est l'index du tableau des membres de chaque groupe
                    if (dataArtist[l].members.length == ArrOfMembers[ArrIndex]) { // Si le nombre de membre du groupe = un des nombres du filtre choisis
                        for (let z = 1958; z <= output.innerHTML; z++) { // z correspond aux années de 1958 à la valeur du range
                            const Album = dataArtist[l].firstAlbum.split('-')
                            const YearOfAlbum = Album[2]
                            if (YearOfAlbum == z) { // Si la date du 1er albulm = une date entre 1958 et le range 
                                // if (arr.length < 10){
                                arr.push(tab[l]) // On push dans arr(tableau vide) la carte de tel artiste si tout les filtres sont checked
                                    // }
                            }
                        }
                    }
                }
            }
        }
    }
}

// Fonction qui renvoyé la valeur des Checkbox à la fonction ci-dessus
function pushOnArr() {
    ArrOfMembers = [];
    for (let NbrDeCheckbox = 1; NbrDeCheckbox <= 8; NbrDeCheckbox++) {
        if (document.getElementById("Checkbox" + NbrDeCheckbox).checked === true) {
            ArrOfMembers.push(document.getElementById("Checkbox" + NbrDeCheckbox).value)
        } else {
            if (document.getElementById("CheckboxAllMembers").checked === true) {
                ArrOfMembers.push(document.getElementById("Checkbox" + NbrDeCheckbox).value)
            }
        }
    }
    return ArrOfMembers
}





var range = document.getElementById("input-range")
var output = document.getElementById("dateOutput")
output.innerHTML = range.value

range.oninput = function() {
    output.innerHTML = this.value
}

let contener = document.getElementById("PageChoosing")

// Fonction qui récupère le nombre de page à générer par rapport aux nombres d'artistes affichés
function pageNumber(param) {
    contener.innerHTML = ""
    for (let index = 1; index <= param; index++) {
        let page = document.createElement("button")
        page.setAttribute("onclick", "pagination(" + index + ")")
        page.id = "page" + index
        page.innerText = "Page n° " + index
        contener.appendChild(page)
    }
}

var arrBis = []

// Fonction pour la pagination
function pagination(param) {
    matchList.innerHTML = ""
    if (param == 1) {
        arrBis = arr.slice(0, 10)
    } else if (param == 2) {
        arrBis = arr.slice(10, 20)
    } else if (param == 3) {
        arrBis = arr.slice(20, 30)
    } else if (param == 4) {
        arrBis = arr.slice(30, 40)
    } else if (param == 5) {
        arrBis = arr.slice(40, 50)
    } else {
        arrBis = arr.slice(50)
    }
    matchList.innerHTML = arrBis.join('');
}

setupPopup()