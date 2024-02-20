# GROUPIE_TRACKER

<img alt="Go" src="https://img.shields.io/badge/go-%2300ADD8.svg?&style=for-the-badge&logo=go&logoColor=white"/><img alt="JavaScript" src="https://img.shields.io/badge/javascript-%23323330.svg?&style=for-the-badge&logo=javascript&logoColor=%23F7DF1E"/>
<img alt="HTML5" src="https://img.shields.io/badge/html5-%23E34F26.svg?&style=for-the-badge&logo=html5&logoColor=white"/>
<img alt="CSS3" src="https://img.shields.io/badge/css3-%231572B6.svg?&style=for-the-badge&logo=css3&logoColor=white"/>

Le projet **Groupie Tracker** consistait à recevoir une API donnée et manipuler les données qu'elle contenait. L'objectif étant de réaliser un site web et d'y afficher un certain nombre de donnée venant de l'API.

# Objectifs

* Afficher les artistes :
    * Nom du groupe
    * Image
    * Création du groupe
    * Date du premier album
    * Les noms des membres

    ![exemple](https://cdn.discordapp.com/attachments/826340732117712916/839429860706615306/Capture_decran_2021-05-05_a_11.14.18.png)

* Recherche agit sur (Rafraichissement dynamique) :
    * Le nom du groupe
    * La date de création
    * Le nom des membres
    * La date du premier album

    ![exemple](https://cdn.discordapp.com/attachments/826340732117712916/839433027012657158/Capture_decran_2021-05-05_a_11.26.54.png)

* Filtrage (Rafraichissement dynamique) : 
    * Avec un système de pagination en fonction du nombre de résultat afficher (10 max/pages & interaction avec le filtre)

        * La date de création du groupe (Select input)
        * Le nombre de membre dans le groupe (CheckBox)
        * La date du premier album (Range filter)
        
![exemple](https://cdn.discordapp.com/attachments/826340732117712916/839431747250683914/Capture_decran_2021-05-05_a_11.21.42.png)


* Géolocalisation

    * Google Maps API
        * Maps JavaScript API
        * Géocoding API

    *Le service peut-être amené à être interrompu (Version d'essai)*

* Séparation en page de nagivation

* Gestion des erreurs

    * Erreur 404
    * Erreur 500

![error](https://cdn.discordapp.com/attachments/826340732117712916/839427272390082560/Capture_decran_2021-05-05_a_11.04.03.png)

# Pour commencer

Tout d'abord, pour tester notre site web, veillez à télécharger l'intégralité du dossier **GROUPIE_TRACKER**. Une fois cela fait, rendez-vous sur *Visual Studio Code 2*. 

## Lancement du serveur

Pour lancer le serveur, il vous suffit d'utiliser la commande suivante dans votre terminal : ``go run main.go``. 

*⚠️ Pour les utilisateurs Mac possédant une erreur lors de l'éxécution de cette commande, suivez les instructions suivantes :*

- ``env GO111MODULE=off go build main.go``
- ``./main``

Une fois le serveur lancé, des logs apparaissent dans votre terminal en couleur. Vous pouvez cliquer directement sur votre terminal pour ouvrir la page web ou [ici](http://localhost:8000).
![log](https://cdn.discordapp.com/attachments/826340732117712916/839426124756025344/Capture_decran_2021-05-05_a_10.58.49.png)

*Des erreurs peuvent être présente sous cette forme :*
![error](https://cdn.discordapp.com/attachments/826340732117712916/839427272390082560/Capture_decran_2021-05-05_a_11.04.03.png)

# Versions

* 1.0 - Release 05/05/21

# Auteurs

* Elouan DUMONT - [@ByMSRT](https://github.com/ByMSRT)
* Kévin GUYODO - [@kevinguyodo](https://github.com/kevinguyodo)
* Mathis VERON - [@mveron13](https://github.com/mveron13)
* Tao BOURMAUD - [@taobourmaud](https://github.com/taobourmaud)

# Développement avec 

<img alt="Go" src="https://img.shields.io/badge/go-%2300ADD8.svg?&style=for-the-badge&logo=go&logoColor=white"/><img alt="JavaScript" src="https://img.shields.io/badge/javascript-%23323330.svg?&style=for-the-badge&logo=javascript&logoColor=%23F7DF1E"/>
<img alt="HTML5" src="https://img.shields.io/badge/html5-%23E34F26.svg?&style=for-the-badge&logo=html5&logoColor=white"/>
<img alt="CSS3" src="https://img.shields.io/badge/css3-%231572B6.svg?&style=for-the-badge&logo=css3&logoColor=white"/>
