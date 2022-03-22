# Serveur de base en Golang

Exemple d'un serveur basique de serveur tcp en go, ce serveur est complet, multithread et gère plusieurs clients, il est possible qu'il manque des mutex.

## Compilation

Pour démarrer sur le serveur, se placer sur le répertoire du projet puis taper:

```go
go run .
```

Pour faire un build:

```go
go build
```

## Utilisation

Bien que cet exemple soit pour donner une base avec un code simple, il est utilisable avec telnet pour faire des essais, pour les "commandes" voir dans le fichier server.go.

