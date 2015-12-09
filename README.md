# go_shortlink read me

Pour plus d'information lire: conception/DCD.pdf

1. Installation de l'environnement de développement

 - Installer golang depuis la documentation https://golang.org/

 - Installer docker engine et docker compose à partir de http://docs.docker.com/engine/installation/ et http://docs.docker.com/compose/install/. Configurer la variable d'environnement DOCKER_HOST


2. Compiler go_shortlink:

Exécuter le make qui est dans le dossier github.com/bmotto/go_shortlink.

3. Execution

Taper dans une console : go_shortlink

Execution depuis docker : sudo docker-compose up

4. Cas d'utilisation

3.1 Shortlink handler

Pour générer un shortlink il y a plusieurs moyens. Pour un url court de type www.google.com il peut directement être passé dans l'url après /shortlink/url_court accompagné ou non d'un parametre custon permettant de personnalisé le code court. l'utilisateir peut choisir également de mettre les "www" ou de ne pas les mettre, examples:
 - http://127.0.0.1:9999/shortlink/www.facebook.com
 - http://127.0.0.1:9999/shortlink/www.facebook.com&custom=faBo
 - http://127.0.0.1:9999/shortlink/facebook.com
 - http://127.0.0.1:9999/shortlink/facebook.com&custom=faBo

Pour les url un peu plus long contenant un ou plusieurs "/", l'utilisateur doit passer les paramètres dans le body de la requête POST, ces arguments doivent ce trouver au format json, example:
 - {"link":"www.youtube.com/watch?v=RoePjPQP7XE","custom":"millio"}

Par défaut, l'API prend les données transmises dans le json, celui-ci est vide alors l'API prend ceux présents dans l'URL.

3.2 Redirection handler

Si on reprend les exemples si dessus qui on ammené à la création d'un lien court pour l'url facebook.com, il suffit d'utilier faBo pour être redirigé vers facebook.com:
  - http://127.0.0.1:9999/faBo

3.3 Admin handler

En conservant le même example que précédement, pour connaître les statistiques de redirection du code court faBo, il suffit de taper:
  - http://127.0.0.1:9999/admin/faBo
