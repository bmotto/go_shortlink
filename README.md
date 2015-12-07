# go_shortlink read me

Pour plus d'information lire: conception/DCD.pdf

1. Installation de l'environnement

installer golang depuis la doc

installer redis-server et le démarrer


2. ??



3. Cas d'utilisation

3.1 Shortlink handler

Pour générer un shortlink il y a plusieurs moyens. Pour un url court de type www.google.com il peut directement être passé dans l'url après /shortlink/url_court accompagné ou non d'un parametre custon permettant de personnalisé le code court. l'utilisateir peut choisir également de mettre les "www" ou de ne pas les mettre, examples:
 - http://127.0.0.1:9999/shortlink/www.facebook.com
 - http://127.0.0.1:9999/shortlink/www.facebook.com&custom=faBo
 - http://127.0.0.1:9999/shortlink/facebook.com
 - http://127.0.0.1:9999/shortlink/facebook.com&custom=faBo

Pour les url un peu plus long contenant un ou plusieur "/", l'utilisateur doit passer les parametres dans le body de la requête POST, ces arguments doivent ce trouver au format json, example:
 - {"link":"www.youtube.com/watch?v=RoePjPQP7XE","custom":"millio"}

Par défaut, l'API prend les données transmises dans le json, celui-ci est vide alors l'API prend ceux présents dans l'URL.

3.2 Redirection handler

3.3 Admin handler
