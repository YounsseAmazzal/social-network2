"use client";
import { createContext, useContext, useEffect, useState } from "react";
import { fetchfromApi, isRecord, readBoolean, readString } from "@/lib/api";
import {
  filterBellNotifications,
  Nothifiaction_limit,
  addNotificationToBell,
  isBellNotificationType,
} from "@/lib/notifications";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useRealtimeSocket } from "@/providers/RealtimeSocketProvider";
const NotificationCenterContext = createContext(undefined);
function normalizeNotification(value) {
  if (!isRecord(value)) return null;

  const id = readString(value, ["id", "ID"]);
  
  const type = readString(value, ["category", "notification_category", "type"]);

  if (!id || !type) {
    console.error("Missing ID or Category in notification:", value);
    return null;
  }

  return {
    id,
    userId: readString(value, ["recipient_id"]),
    actorId: readString(value, ["sender_id"]),
    type: type, 
    entityId: readString(value, ["target_id"]),
    entityType: readString(value, ["target_type"]),
    message: readString(value, ["message", "display_message"]),
    isRead: readBoolean(value, ["is_read", "isRead", "read_status"]),
    createdAt: readString(value, ["created_at", "createdAt"]),
  };
}

async function fetchBellNotifications() {
  const response = await fetchfromApi(
    `/golbalenothification?limit=${Nothifiaction_limit}`,
  );
  
  const rawItems = response.userNothification || response.data || [];
  // console.log("the nothificaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",rawItems)
  const normalizedItems = Array.isArray(rawItems)
    ? rawItems.flatMap((item) => {
        const notification = normalizeNotification(item);
        return notification ? [notification] : [];
      })
    : [];
  // console.log("fillterrrd nothificatttion",filterBellNotifications(normalizedItems))
  return filterBellNotifications(normalizedItems);
}
export function NotificationCenterProvider({ children }) {
  const { user } = useAuthSession();
  const { latestNotification } = useRealtimeSocket();
  const [notifications, setNotifications] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const userId = user?.id ?? "";
  useEffect(() => {
    if (!userId) {
      setNotifications([]);
      setIsLoading(false);
      return;
    }
    let isCurrentRequest = true;
    setIsLoading(true);
    void fetchBellNotifications()
      .then((items) => {
        if (isCurrentRequest) setNotifications(items);
      })
      .catch(() => undefined)
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [userId]);
useEffect(() => {
    if (!latestNotification) return;
    // Socket notifications already arrive normalized from liiibb
    const notification = latestNotification;
    const recipientId = notification.userId || notification.recipient_id;
    if (!notification.id || recipientId !== userId) return;
    if (!isBellNotificationType(notification.type)) return;

    setNotifications((current) =>
      addNotificationToBell(current, notification, Nothifiaction_limit)
    );
  }, [latestNotification, userId]);
  async function refresh() {
    if (!userId) {
      setNotifications([]);
      return;
    }
    setIsLoading(true);
    try {
      setNotifications(await fetchBellNotifications());
    } finally {
      setIsLoading(false);
    }
  }
  async function markAllRead() {
    if (!userId) return;
    await fetchfromApi(`/ReadNormaalNothification?all=${true}`, {
      method: "POST",
      body: JSON.stringify({ all: true }),
    });
    setNotifications((currentNotifications) =>
      currentNotifications.map((notification) => ({
        ...notification,
        isRead: true,
      })),
    );
  }
  async function markRead(notificationId) {
    if (!userId || !notificationId) return;
    
    setNotifications((current) =>
      current.map(n => n.id === notificationId ? { ...n, isRead: true } : n)
    );

    try {
      await fetchfromApi(`/ReadNormaalNothification?notification_id=${notificationId}`, {
        method: "POST",
      });
    } catch {
       refresh();
    }
  }
  const unreadCount = notifications.filter(
    (notification) => !notification.isRead,
  ).length;
  return (
    <NotificationCenterContext.Provider
      value={{
        notifications,
        unreadCount,
        isLoading,
        refresh: async () => setNotifications(await fetchBellNotifications()),
        markAllRead,
        markRead,
      }}
    >
      {children}
    </NotificationCenterContext.Provider>
  );
}
export function useNotificationCenter() {
  const context = useContext(NotificationCenterContext);
  if (!context) {
    throw new Error(
      "useNotificationCenter must be used inside NotificationCenterProvider.",
    );
  }
  return context;
}

