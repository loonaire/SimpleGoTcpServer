package server

import (
	"log"
	"net"
	"strconv"
	"time"
)

type TcpServer struct {
	ip         string
	port       string
	clientList map[uint]Client
	/* Ici il est possible de personnalisé les données du client en ajoutant des éléments */
}

type Client struct {
	conn         net.Conn
	connectionId uint // clé de la map clientList qui sert à "identifier" le client, ainsi il a un identifiant unique
	clientList   *map[uint]Client
}

func NewTcpServer(ip string, port string) TcpServer {
	server := TcpServer{ip: ip, port: port, clientList: make(map[uint]Client)}
	return server
}

func (server *TcpServer) StartServer() {
	var connectionId uint = 0 // nombre de clients qui se sont connecté
	log.Println("Démarrage du serveur...")
	listener, err := net.Listen("tcp", server.ip+":"+server.port)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
	defer listener.Close()
	log.Println("Serveur démarré sur le port: ", server.port, "\nEn attente de connexion...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Erreur lors de la connexion d'un client: ", err)
		}
		client := &Client{conn: conn, clientList: &server.clientList, connectionId: connectionId}
		server.clientList[connectionId] = *client
		connectionId++
		go client.handleClient()
	}
}

func (client *Client) handleClient() {
	log.Println("Nouveau client connecté!")
	client.conn.Write([]byte("Bienvenue sur le serveur!\n"))

	for {

		// gère le timeout de la connexion, le temps avant deconnexion si la connexion ne reçoit ni envoi aucun paquet
		//client.conn.SetDeadline(time.Now().Add(120 * time.Second)) // si on a reçu ni envoyer aucun paquet vers le client
		client.conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // si on a pas reçu de paquet du client, on le deconnecte (ne compte pas si on lui envoi un paquet)
		packet := make([]byte, 200)
		packetlen, err := client.conn.Read(packet)
		if err != nil {

			// gère l'inactivité du client
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				client.send("Deconnexion pour inactivité\n")
				log.Println("Deconnexion pour inactivité", err)
				break
			} else {
				log.Fatal("Erreur lors de la lecture du paquet ", err)
			}
		}
		/* A partir d'ici il est possible de supprimer tout le contenu de la boucle for et de repartir sur quelque chose de propre */
		log.Println("Paquet de ", packetlen, " reçu, contenu:")
		log.Printf("%q\n", packet[:packetlen])

		client.send("Paquet reçu: " + string(packet[:packetlen]))

		switch string(packet[:packetlen-2]) {
		case "whoami":
			client.send(client.conn.RemoteAddr().String() + " id du client: " + strconv.FormatInt(int64(client.connectionId), 10) + "\n")
		case "sendall":
			client.sendToAllClients("Envoi a tout le monde\n")
		case "getclients":
			for _, elt := range *client.clientList {
				client.send(string(elt.conn.RemoteAddr().String()) + "\n")
			}
		}

		if string(packet[:packetlen-2]) == "bye" {
			break
		}

	}
	// ferme la connexion et supprime les informations du client
	client.conn.Close()
	delete(*client.clientList, client.connectionId)
}

func (client *Client) send(packetData string) {
	log.Printf("Envoi au client %s du paquet: %q\n", client.conn.RemoteAddr(), packetData)
	client.conn.Write([]byte(packetData))
}

func (client *Client) sendToAllClients(packetData string) {
	for _, element := range *client.clientList {
		element.send(packetData)
	}
}
