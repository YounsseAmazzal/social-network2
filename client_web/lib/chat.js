import { isRecord, readBoolean, readNumber, readString } from "@/lib/api";

export const MESSAGE_MAX_LENGTH = 1000;

export function exceedsMessageLimit(text) {
  return new TextEncoder().encode(text ?? "").length > MESSAGE_MAX_LENGTH;
}
export function normalizeChatMessage(value) {
  if (!isRecord(value)) return null;
  // Accept both the REST shape (message/sent_time) and the websocket shape
  // the backend broadcasts (group_message/send_date)
  const text = readString(value, [
    "text",
    "Text",
    "message",
    "content",
    "group_message",
  ]);
  if (!text) return null;
  return {
    id: readString(value, ["id", "ID"]),
    senderId: readString(value, ["senderID", "sender_id", "SenderID"]),
    receiverId: readString(value, [
      "receiverID",
      "receiver_id",
      "ReceiverID",
      "userchatwith_id",
    ]),
    text,
    createdAt: readString(value, [
      "createDate",
      "created_at",
      "sent_at",
      "sent_time",
      "SentAt",
      "send_date",
    ]),
  };
}
export function normalizeGroupConversationMessage(value) {
  if (!isRecord(value)) return null;
  const text = readString(value, ["message", "text", "group_message", "content"]);
  const groupId = readString(value, ["group_id", "grouupid", "GroupID"]);
  if (!text || !groupId) return null;
  return {
    id: readString(value, ["id", "ID"]),
    senderId: readString(value, ["sender_id", "senderID", "SenderID"]),
    groupId,
    text,
    createdAt: readString(value, [
      "sent_at",
      "sentAt",
      "send_date",
      "CreateDate",
    ]),
    senderNickname: readString(
      value,
      ["senderNickname", "userFullname", "fullname", "fullName"],
      "Member",
    ),
    avatarUrl: readString(value, ["avatarURL", "user_avatar", "avatar"]),
  };
}

export function normalizeSocketNotification(value) {
  if (!isRecord(value)) return null;
  const type = readString(value, ["type", "category", "notification_category"]);
  if (!type) return null;
  const entityId = readString(value, [
    "entity_id",
    "entityID",
    "EntityID",
    "target_id",
  ]);
  const message = readString(value, [
    "message",
    "content",
    "Content",
    "display_message",
  ]);
  return {
    id: readString(value, ["id", "ID"]),
    userId: readString(value, ["user_id", "userID", "UserID", "recipient_id"]),
    actorId: readString(value, ["actor_id", "actorID", "ActorID", "sender_id"]),
    groupId: entityId,
    type,
    entityId,
    entityType: readString(value, [
      "entity_type",
      "entityType",
      "EntityType",
      "target_type",
    ]),
    content: message,
    message,
    isRead: readBoolean(value, ["is_read", "isRead"]),
    createdAt: readString(value, ["created_at", "createdAt", "CreatedAt"]),
  };
}
export function normalizeTypingStatus(value) {
  if (!isRecord(value)) return null;
  const data = isRecord(value.data) ? value.data : value;
  const from = readString(data, ["from", "From"]);
  const to = readString(data, ["to", "To"]);
  if (!from || !to) return null;
  return {
    from,
    to,
    isTyping: readBoolean(data, ["isTyping", "IsTyping"]),
  };
}
export function normalizeUserStatus(value) {
  if (!isRecord(value)) return null;
  const userId = readString(value, ["userID", "userId", "user_id", "UserID"]);
  if (!userId) return null;
  return {
    userId,
    isOnline: readBoolean(value, ["online", "isOnline", "Online"]),
  };
}
export function normalizeChatUsers(value) {
  if (!Array.isArray(value)) return [];
  return value.flatMap((item) => {
    if (!isRecord(item)) return [];
    const id = readString(item, ["id", "ID"]);
    if (!id) return [];
    return [
      {
        id,
        nickname: readString(
          item,
          ["nickname", "Nickname", "username"],
          "Unknown user",
        ),
        avatarUrl: readString(item, ["avatar_url", "avatarURL", "avatar"]),
        isConnected: readBoolean(item, ["is_online", "is_online"]),
      },
    ];
  });
}
export function normalizeGroupMessages(value) {
  if (!Array.isArray(value)) return [];
  return value.flatMap((item) => {
    const message = normalizeGroupConversationMessage(item);
    return message ? [message] : [];
  });
}
export function normalizeVoteStatus(value) {
  if (typeof value === "number")
    return value === 1 ? "like" : value === -1 ? "dislike" : null;
  if (value === "like" || value === "dislike") return value;
  return null;
}
export function normalizeVoteResponse(value) {
  if (!isRecord(value)) return null;
  const likes = readNumber(value, ["likes"]);
  const dislikes = readNumber(value, ["dislikes"]);
  if (!("likes" in value) || !("dislikes" in value)) return null;
  return {
    likes,
    dislikes,
    userVote: normalizeVoteStatus(value.userVote),
  };
}
export function normalizeCommentCount(value) {
  return readNumber(isRecord(value) ? value : {}, [
    "totalComments",
    "comments",
  ]);
}
