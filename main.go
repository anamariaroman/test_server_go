package main

import (
	//Packages
	"fmt"      // para imprimir en consola
	"log"      // registro
	"net/http" // net- Creacion del servidor

	// http- proporciona implementaciones de servidor y cliente HTTP

	"github.com/gorilla/mux" //Es para el servidor, manejara todos los Endpoints.
	// Implementa un enrutador de solicitudes y un despachador para hacer coincidir
	// las solicitudes entrantes con su controlador respectivo

	"github.com/gorilla/websocket" // Es para los websockets, cliente-servidor, que nos permiten
	// establecer una comunicación bidireccional (es decir, que los datos pueden fluir del cliente al servidor y viceversa)
	//y dúplex (es decir, que la comunicación en ambas direcciones puede suceder de manera simultánea
)

// Estructura para mantener y crear un mensaje (Saludo) en formato de string, que se va a recibir del Cliente o del front-end
type Message struct {
	Greeting string `json:"greeting"`
}

// Variables
// Upgrader: Es el actualizador de paquetes de websocket , actualiza la conexión del servidor HTTP al protocolo WebSocket
var (
	wsUpgrader = websocket.Upgrader{

		//ReadBufferSize y WriteBufferSize especifican el tamaño de estos buffers

		ReadBufferSize:  1024, //MB
		WriteBufferSize: 1024,

		//El buffer es un espacio temporal de memoria física el cual se usa para almacenar información
		//mientras se envía de un lado a otro
	}

	// Puntero para la conexion websocket
	wsConn *websocket.Conn
)

// Endpoint HTTP

// Enrutador Respuestas y Objeto de Peticiones
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Cros-Origin resource sharing (Error de curso)
	// Chequea el origen en el localhost
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		// De lo contrario, verifica http.Request
		// Asegura que esta bien acceder
		return true
	}
	var err error
	// Actualizar que las conexiones permanezcan conectadas
	wsConn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	// Aplazar la cierre de la conexion
	defer wsConn.Close()

	// event loop(Bucle de eventos)
	for {
		// Esperar que entre un mensaje
		var msg Message

		// conexion ws para leer el mensaje en formato JSON
		// el puntero obtiene la direccion del mensaje y se almacena en la variable msg
		// y luego volvera a esperar el proximo mensaje
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error reading JSON: %s\n", err.Error())
			// cierra el ciclo
			break
		}

		// Se imprime en consola el mensaje recibido
		fmt.Printf("Message Received: %s\n", msg.Greeting)
		SendMessage("Hello, Client!")
	}
}

func SendMessage(msg string) {
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Printf("error sending message: %s\n", err.Error())
	}
}

func main() {

	// Enrutador: Definicion de rutas
	router := mux.NewRouter()

	// HandleFunc registra una nueva ruta con un comparador para la ruta URL.
	router.HandleFunc("/socket", WsEndpoint)

	// Log.Fatal informara si hay un problema en la ejecucion en la salida del sistema operativo
	// http.ListenAndServe, abre el puerto del servidor y bloquea para siempre la espera de los clientes
	// ListenAndServe escucha en la dirección de red TCP
	log.Fatal(http.ListenAndServe(":9100", router))

}
