"use client";
import { createContext, useContext, useEffect, useState } from "react";
import { API_URL, isRecord, readString } from "@/lib/api";
import {
  normalizeChatMessage,
  normalizeGroupConversationMessage,
  normalizeSocketNotification,
  normalizeTypingStatus,
  normalizeUserStatus,
} from "@/lib/chat";
import { useAuthSession } from "@/providers/AuthSessionProvider";
const RealtimeSocketContext = createContext(undefined);
/** Longest we wait between reconnection attempts. */
const MAX_RETRY_DELAY_MS = 15000;
function getWebSocketUrl() {
  const protocol = API_URL.startsWith("https://") ? "wss://" : "ws://";
  return `${protocol}${API_URL.replace(/^https?:\/\//, "")}/WebSocket`;
}
function playSound(source) {
  const audio = new Audio(source);
  audio.currentTime = 0;
  void audio.play().catch(() => undefined);
}
export function RealtimeSocketProvider({ children }) {
  const { user } = useAuthSession();
  const [socket, setSocket] = useState(null);
  const [latestMessage, setLatestMessage] = useState(null);
  const [latestNotification, setLatestNotification] = useState(null);
  const [typingStatus, setTypingStatus] = useState(null);
  const [userStatus, setUserStatus] = useState(null);
  const userId = user?.id ?? "";
  useEffect(() => {
    if (!userId) {
      setSocket(null);
      return;
    }
    let disposed = false;
    let retryTimer = null;
    let attempts = 0;
    let webSocket = null;

    function handleMessage(event) {
      try {
        const payload = JSON.parse(event.data);
        // console.log("dataaaatype",event.data)
        // console.log("dataaaatype",event.Data)
        const type = readString(payload, ["type", "Datatype", "wsdatatype"]);
        const data = isRecord(payload.data)
          ? payload.data
          : isRecord(payload.WsData)
            ? payload.WsData
            : isRecord(payload.wsdata)
              ? payload.wsdata
              : payload;
        switch (type) {
          case "private_message":
          case "send_message_group":
          case "group_message": {
            
            const groupId =
              readString(data, ["group_id", "grouupid", "GroupID"]) || "";
            const isPrivate =
              type === "private_message" ||
              (!groupId && "userchatwith_id" in data);
            const message = isPrivate
              ? normalizeChatMessage(data)
              : normalizeGroupConversationMessage(data);

            if (!message) return;
            setLatestMessage(message);
            const isIncoming = message.senderId !== userId;
            if (isIncoming) playSound("/audios/receives.mp3");
            break;
          }

          case "notification_read":
          case "social_notification":
            //wraaaap all the nothifiactions by tge bacckend 
            setLatestNotification(normalizeSocketNotification(data));
            break;

          case "user_status":
            setUserStatus(normalizeUserStatus(data));
            break;

          case "user_typing":
            setTypingStatus(normalizeTypingStatus(data));
            break;

          case "error": {
            /// i refuse  aany action message without check the following 
            //witth 
            const errorText = readString(payload, ["WsData", "wsdata", "data"]);
            setLatestMessage({
              type: "error",
              data: errorText || "The request could not be completed.",
            });
            break;
          }

          default:
            // Unknown event type egnooret and keep the connection alive
            break;
        }
      } catch {
        //keeeep the con aliieve 
      }
    }
    function connect() {
      if (disposed) return;
      webSocket = new WebSocket(getWebSocketUrl());
      setSocket(webSocket);
      webSocket.onopen = () => {
        attempts = 0;
      };
      webSocket.onmessage = handleMessage;
      webSocket.onclose = () => {
        if (disposed) return;
        setSocket((current) => (current === webSocket ? null : current));
        const delay = Math.min(500 * 2 ** attempts, MAX_RETRY_DELAY_MS);
        attempts += 1;
        retryTimer = window.setTimeout(connect, delay);
      };
    }
    connect();
    return () => {
      disposed = true;
      if (retryTimer) window.clearTimeout(retryTimer);
      if (webSocket) {
        webSocket.onclose = null;
        webSocket.onmessage = null;
        if (
          webSocket.readyState === WebSocket.OPEN ||
          webSocket.readyState === WebSocket.CONNECTING
        ) {
          webSocket.close();
        }
      }
    };
  }, [userId]);
  function sendTyping(senderId, recipientId, isTyping) {
    if (socket?.readyState !== WebSocket.OPEN) return;
    socket.send(
      JSON.stringify({
        wsdatatype: "user_typing",
        wsdata: { from: senderId, to: recipientId, isTyping },
      }),
    );
  }
  return (
    <RealtimeSocketContext.Provider
      value={{
        socket,
        latestMessage,
        latestNotification,
        typingStatus,
        userStatus,
        sendTyping,
        playSendSound: () => playSound("/audios/sends.mp3"),
        playReceiveSound: () => playSound("/audios/receives.mp3"),
      }}
    >
      {children}
    </RealtimeSocketContext.Provider>
  );
}
export function useRealtimeSocket() {
  const context = useContext(RealtimeSocketContext);
  if (!context) {
    throw new Error(
      "useRealtimeSocket must be used inside RealtimeSocketProvider.",
    );
  }
  return context;
}
