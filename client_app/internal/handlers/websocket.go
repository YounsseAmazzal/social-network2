package handlers

import (
	"fmt"
	"log"
	"net/http"
	"social/internal/requests"
	"social/internal/services"
	"social/internal/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type WSManager struct {
	/// wheere can i stooore userid and multiple tabs
	Clients  sync.Map
	WsDb     *services.Service
	Upgrader websocket.Upgrader
}

func Neeewsmanger(sevice *services.Service) *WSManager {
	return &WSManager{
		WsDb: sevice,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == "http://localhost:3000"
			},
		},
	}
}
func (wsrepo *WSManager) WebSocket(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/WebSocket", http.MethodGet) {
		return
	}
	cookie, err := req.Cookie("solo_social")
	if err != nil {
		utils.ErrorResponse(res, "sesssion not valiiid ", http.StatusUnauthorized)
		return
	}
	_, err = wsrepo.WsDb.Reepo.Checksession(cookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "Your Session Not valid", http.StatusUnauthorized)
		// log.Printf("Your Session Not valid", err.Error())
		return
	}
	user, err := wsrepo.WsDb.Reepo.GetUserByToken(cookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "Your Session Not validd", http.StatusUnauthorized)
		return
	}

	conn, err := wsrepo.Upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Register the connection fooor the userid
	wsrepo.registerClient(user.ID, conn)

	defer func() {
		wsrepo.unregisterClient(user.ID, conn)
		conn.Close()
	}()

	// message looop
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break // Connection closed
		}
		fmt.Printf("the Payload %s",payload)
		// wsm.routeMessage(userID, payload)
	}
}

// add new conectioon to the user
func (wsm *WSManager) registerClient(userID string, conn *websocket.Conn) {
	// try to bring all the connection of the user
	actual, _ := wsm.Clients.LoadOrStore(userID, make(map[*websocket.Conn]bool))

	connections := actual.(map[*websocket.Conn]bool)

	connections[conn] = true

	if len(connections) == 1 {
		//user online
		// wsm.broadcastStatus(userID, true)
		fmt.Println("the user os noow Online 1 tab ")
	}
}

// delete connection when the user close the taab of lose connection
func (wsm *WSManager) unregisterClient(userID string, conn *websocket.Conn) {
	// bring all the user connections
	val, ok := wsm.Clients.Load(userID)
	if !ok {
		return
	}

	connections := val.(map[*websocket.Conn]bool)
	/// delete the special connection (onne taab)
	delete(connections, conn)
	// check if there is any avaaaliabe conneections
	if len(connections) == 0 {
		// if all connections closed delete from memory
		wsm.Clients.Delete(userID)

		// "Offliiiiiiiine"
		// wsm.broadcastStatus(userID, false)
		fmt.Println("all the tabs are closed use is offline")
	} else {
		fmt.Println("user cloooesed one tabe and has other opeened tabs ")
	}
}
