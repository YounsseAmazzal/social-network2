"use client";
import { createContext, useContext, useEffect, useState } from "react";
import { usePathname } from "next/navigation";
import { fetchfromApi, isRecord, readString } from "@/lib/api";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useRealtimeSocket } from "@/providers/RealtimeSocketProvider";
const UnreadConversationContext = createContext(undefined);

function transformSendersToMap(senders, idKey) {
  if (!Array.isArray(senders)) return {};
  
  return senders.reduce((map, item) => {
    const id = item[idKey];
    const count = item.unread_total;
    if (id && typeof count === "number") {
      map[id] = count;
    }
    return map;
  }, {});
}
function getActiveChatUserId(pathname) {
  if (!pathname.startsWith("/chatprivite")) return "";
  
  const params = new URLSearchParams(window.location.search);
  console.log('the talker id',params.get("Userchat_id"))
  return params.get("Userchat_id") || "";
}
function normalizeCountMap(value) {
  if (!isRecord(value)) return {};
  return Object.entries(value).reduce((counts, [key, count]) => {
    if (typeof count === "number" && Number.isFinite(count)) {
      counts[key] = Math.max(0, count);
    }
    return counts;
  }, {});
}
function parseChannelMessage(value) {
  if (!isRecord(value)) return null;
  const type = readString(value, ["type"]);
  const userId = readString(value, ["userId"]);
  const groupId = readString(value, ["groupId"]);
  if (
    (type === "active_chat" ||
      type === "inactive_chat" ||
      type === "thread_read") &&
    userId
  ) {
    return { type, userId };
  }
  if (type === "group_read" && groupId) return { type, groupId };
  return null;
}
async function fetchUnreadCounts() {
  const [chatResponse, groupResponse] = await Promise.all([
    fetchfromApi("/ChatNothification"),
    fetchfromApi("/GroupChatNothification"),
  ]);
  // console.log("API Chat Response:", chatResponse);
  // console.log(" API Group Response:", groupResponse);
  return {
    users: transformSendersToMap(chatResponse.senders, "sender_id"),
    groups: transformSendersToMap(groupResponse.senders, "group_id"),
  };
}
function removeCount(counts, id) {
  if (!(id in counts)) return counts;
  const nextCounts = { ...counts };
  delete nextCounts[id];
  return nextCounts;
}
export function UnreadConversationProvider({ children }) {
  const pathname = usePathname();
  const { user } = useAuthSession();
  const { latestNotification } = useRealtimeSocket();
  const [unreadByUser, setUnreadByUser] = useState({});
  const [unreadByGroup, setUnreadByGroup] = useState({});
  const [activeGroupConversationId, setActiveGroupConversationId] =
    useState("");
  const [recentlyLeftChatUserId, setRecentlyLeftChatUserId] = useState("");
  const [channel, setChannel] = useState(null);
  const [activeChatsInOtherTabs, setActiveChatsInOtherTabs] = useState(
    new Set(),
  );
  const activeChatUserId = getActiveChatUserId(pathname);
  const userId = user?.id ?? "";
  useEffect(() => {
    if (typeof BroadcastChannel === "undefined") return;
    const notificationChannel = new BroadcastChannel("chat_notifications");
    notificationChannel.onmessage = (event) => {
      const message = parseChannelMessage(event.data);
      if (!message) return;
      if (message.type === "active_chat") {
        setActiveChatsInOtherTabs((currentIds) =>
          new Set(currentIds).add(message.userId),
        );
      } else if (message.type === "inactive_chat") {
        setActiveChatsInOtherTabs((currentIds) => {
          const nextIds = new Set(currentIds);
          nextIds.delete(message.userId);
          return nextIds;
        });
      } else if (message.type === "thread_read") {
        setUnreadByUser((currentCounts) =>
          removeCount(currentCounts, message.userId),
        );
      } else {
        setUnreadByGroup((currentCounts) =>
          removeCount(currentCounts, message.groupId),
        );
      }
    };
    setChannel(notificationChannel);
    return () => {
      notificationChannel.close();
    };
  }, []);
  useEffect(() => {
    if (!channel || !activeChatUserId) return;
    channel.postMessage({ type: "active_chat", userId: activeChatUserId });
    return () => {
      channel.postMessage({ type: "inactive_chat", userId: activeChatUserId });
    };
  }, [activeChatUserId, channel]);
  useEffect(() => {
    if (!userId) {
      setUnreadByUser({});
      setUnreadByGroup({});
      return;
    }
    let isCurrentRequest = true;
    void fetchUnreadCounts()
      .then((counts) => {
        if (!isCurrentRequest) return;
        // console.log("couuunts users ",counts)
        // console.log("couuunts users ",counts.groups)
        setUnreadByUser(counts.users);
        setUnreadByGroup(counts.groups);
      })
      .catch(() => undefined);
    return () => {
      isCurrentRequest = false;
    };
  }, [userId]);
  useEffect(() => {
    if (!recentlyLeftChatUserId) return;
    const timeoutId = window.setTimeout(
      () => setRecentlyLeftChatUserId(""),
      1000,
    );
    return () => window.clearTimeout(timeoutId);
  }, [recentlyLeftChatUserId]);
  useEffect(() => {
    if (!activeChatUserId) return;
    setRecentlyLeftChatUserId("");
    return () => setRecentlyLeftChatUserId(activeChatUserId);
  }, [activeChatUserId]);
  async function restoreCountsAfterFailure() {
    try {
      const counts = await fetchUnreadCounts();
      setUnreadByUser(counts.users);
      setUnreadByGroup(counts.groups);
    } catch {
      // Keep the optimistic state if the server is unavailable.
    }
  }
  async function markThreadRead(actorId) {
    if (!userId || !actorId) return;
    setUnreadByUser(prev => removeCount(prev, actorId));
    channel?.postMessage({ type: "thread_read", userId: actorId });
    try {
      await fetchfromApi(`/ReadchatNothification?Userchat_id=${actorId}`, {
        method: "POST",
        // body: JSON.stringify({}),
      });
    } catch {
      await restoreCountsAfterFailure();
    }
  }
  async function markGroupRead(groupId) {
    if (!userId || !groupId) return;
   setUnreadByGroup(prev => {
        const next = { ...prev };
        delete next[groupId];
        return next;
    });
        channel?.postMessage({ type: "group_read", groupId });
    try {
      await fetchfromApi(`/ReadGroupNothification?Grouup_id=${groupId}`, {
        method: "POST",
        // body: JSON.stringify({}),
      });
    } catch {
      await restoreCountsAfterFailure();
    }
  }
  useEffect(() => {
    if (!userId || !activeChatUserId || !unreadByUser[activeChatUserId]) return;
    void markThreadRead(activeChatUserId);
  }, [activeChatUserId, unreadByUser, userId]);
  useEffect(() => {
    if (!latestNotification) return;
    handleIncomingNotification(latestNotification);
  }, [latestNotification]);
  async function handleIncomingNotification(notification) {
    if (!userId || notification.userId !== userId) return;
    // The backend categories are "direct_message"/"group_chat_message";
    // keep the legacy "message"/"group_message" spellings too.
    const isDirectMessage =
      notification.type === "message" ||
      notification.type === "direct_message";
    const isGroupMessage =
      notification.type === "group_message" ||
      notification.type === "group_chat_message";
    if (isDirectMessage && notification.actorId) {
      if (
        notification.actorId === activeChatUserId ||
        notification.actorId === recentlyLeftChatUserId ||
        activeChatsInOtherTabs.has(notification.actorId)
      ) {
        await markThreadRead(notification.actorId);
        return;
      }
      setUnreadByUser((currentCounts) => ({
        ...currentCounts,
        [notification.actorId]: (currentCounts[notification.actorId] ?? 0) + 1,
      }));
      return;
    }
    if (isGroupMessage && notification.groupId) {
      if (notification.groupId === activeGroupConversationId) {
        await markGroupRead(notification.groupId);
        return;
      }
      setUnreadByGroup((currentCounts) => ({
        ...currentCounts,
        [notification.groupId]: (currentCounts[notification.groupId] ?? 0) + 1,
      }));
    }
  }
  const totalUnread = Object.values(unreadByUser).reduce(
    (total, count) => total + count,
    0,
  );
  const totalGroupUnread = Object.values(unreadByGroup).reduce(
    (total, count) => total + count,
    0,
  );
  return (
    <UnreadConversationContext.Provider
      value={{
        unreadByUser,
        unreadByGroup,
        totalUnread,
        totalGroupUnread,
        markThreadRead,
        markGroupRead,
        setActiveGroupConversation: (groupId) =>
          setActiveGroupConversationId(groupId ?? ""),
      }}
    >
      {children}
    </UnreadConversationContext.Provider>
  );
}
export function useUnreadConversations() {
  const context = useContext(UnreadConversationContext);
  if (!context) {
    throw new Error(
      "useUnreadConversations must be used inside UnreadConversationProvider.",
    );
  }
  return context;
}
