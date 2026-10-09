package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social/internal/models"
	"social/internal/requests"
	"social/internal/services"
	"social/internal/utils"
	"sync"
	"time"

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
		fmt.Printf("the Payload %s", payload)
		wsrepo.routeMessage(user.ID, payload)
	}
}

// add new conectioon to the user
func (wsrepo *WSManager) registerClient(userID string, conn *websocket.Conn) {
	// try to bring all the connection of the user
	actual, _ := wsrepo.Clients.LoadOrStore(userID, make(map[*websocket.Conn]bool))

	connections := actual.(map[*websocket.Conn]bool)

	connections[conn] = true

	if len(connections) == 1 {
		//user online
		wsrepo.broadcastStatus(userID, true)
		fmt.Println("the user is noow Online 1 tab ")
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

func (wsm *WSManager) broadcastStatus(userID string, isOnline bool) {
	payload := map[string]any{
		"Datatype": "user_status",
		"WsData": map[string]any{
			"user_id": userID,
			"online":  isOnline,
		},
	}
	// maintaain the resources
	message, _ := json.Marshal(payload)

	// loop on all the user on the hub
	wsm.Clients.Range(func(key, value any) bool {
		targetUserID := key.(string)

		if targetUserID == userID {
			return true
		}
		// all tabs
		connections := value.(map[*websocket.Conn]bool)
		// sent to all tabs
		for conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error broadcasting to %s: %v", targetUserID, err)
			}
		}

		return true //
	})
}
func (wsrepo *WSManager) routeMessage(senderID string, payload []byte) {
	fmt.Println("useerr to rouuuter the message")
	var wsInput models.WebSocketInput
	if err := json.Unmarshal(payload, &wsInput); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}
	// fmt.Println("ssssssssssssssssssssssssssssssss",wsInput.Datatype)
	switch wsInput.Datatype {
	case "send_message_group":
		wsrepo.handleGroupMessage(senderID, wsInput.WsData)
	case "user_typing":
		// wsrepo.handleTyping(senderID, wsInput.WsData)
	case "private_message":
		wsrepo.handelprivitechat(senderID, wsInput.WsData)
	case "notification_read":
		// wsrepo.SendNotificationToUser(senderID,)
	}
	// i will add "private_message",
	//  "notification_read"
}

func (wsrepo *WSManager) handleGroupMessage(senderID string, data map[string]any) {
	fmt.Println("User sent a group message")
	groupID, _ := data["group_id"].(string)
	content, _ := data["message"].(string)

	if groupID == "" || content == "" || len(content) > 1000 {
		return
	}
	// fmt.Println("the grouupid",groupID)
	// fmt.Println("the senderID",senderID)
	// fmt.Println("the content",content)
	_, err := wsrepo.WsDb.Reepo.CreateGroupMessage(groupID, senderID, content)
	if err != nil {
		log.Println("DB Error:", err)
		return
	}
	broadcastPayload := map[string]any{
		"Datatype": "send_message_group",
		"WsData": map[string]any{
			"grouupid":      groupID,
			"sender_id":     senderID,
			"group_message": content,
			"send_date":     time.Now().Format(time.RFC3339),
		},
	}
	group, _ := wsrepo.WsDb.Reepo.GetGroupMetadata(groupID)
	for _, memberID := range group.GroupMember {
		fmt.Println("brodcast the message to the group", broadcastPayload)
		wsrepo.SendToUser(memberID, broadcastPayload)
		// khaasni nziddd Send Notification
		///nothiiif
		if memberID != senderID {
			notification := models.Notification{
				RecipientID: memberID,
				SenderID:    senderID,
				Category:    "group_chat_message",
				TargetID:    groupID,
				TargetType:  "group_chat",
				Message:     "New message in " + group.Title,
			}
			_ = wsrepo.WsDb.Reepo.CreateNotification(&notification)
			wsrepo.SendNotificationToUser(memberID, notification)
		}
	}
}

func (wsrepo *WSManager) handelprivitechat(senderID string, data map[string]any) {
	fmt.Println("User sent a privite messages message")
	userchatwith_id, _ := data["userchatwith_id"].(string)
	message, _ := data["message"].(string)
	// fmt.Println("hhhhhhhhhhhhh userchatwith_id ",userchatwith_id)
	if userchatwith_id == "" || message == "" || len(message) > 1000 {
		return
	}
	if userchatwith_id == senderID {
		return
	}
	caanchat, err := wsrepo.WsDb.Reepo.CoulldChat(senderID, userchatwith_id)
	if !caanchat {
		// fmt.Println("$$$$$$$$$$$$$$$$$$$$$$heee can chat$$$$$$$$$$$$$$$$$$$$$$$$$",)
		return
	}
	msgid, err := wsrepo.WsDb.Reepo.CreatePrivateMessage(userchatwith_id, senderID, message)
	if err != nil {
		log.Println("DB Error:", err)
		return
	}
	broadcastPayload := map[string]any{
		"Datatype": "send_message_group",
		"WsData": map[string]any{
			"id":              msgid,
			"userchatwith_id": userchatwith_id,
			"sender_id":       senderID,
			"group_message":   message,
			"send_date":       time.Now().Format(time.RFC3339),
		},
	}
	wsrepo.SendToUser(userchatwith_id, broadcastPayload)
	wsrepo.SendToUser(senderID, broadcastPayload)

	// nothiiiiiif viaa websocket
	notification := models.Notification{
		RecipientID: userchatwith_id,
		SenderID:    senderID,
		Category:    "direct_message",
		TargetID:    senderID,
		TargetType:  "chat",
		Message:     "You received a new message from a friend",
	}
	_ = wsrepo.WsDb.Reepo.CreateNotification(&notification)
	wsrepo.SendNotificationToUser(userchatwith_id, notification)
}
func (wsrepo *WSManager) SendToUser(userID string, payload any) {
	val, ok := wsrepo.Clients.Load(userID)
	if !ok {
		//means  User is offline
		return
	}

	conns := val.(map[*websocket.Conn]bool)
	data, _ := json.Marshal(payload)

	for conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Write error:", err)
		}
	}
}
func (wsm *WSManager) IsUserOnline(userID string) bool {
	// log.Fatal("yaaaaaaaaaaak")
	_, online := wsm.Clients.Load(userID)
	return online
}

func (wsm *WSManager) SendNotificationToUser(userID string, notification models.Notification) {
	payload := map[string]any{
		"wsdatatype": "notification_read",
		"wsdata":     notification,
	}
	wsm.SendToUser(userID, payload)
}

// func (){

// }
