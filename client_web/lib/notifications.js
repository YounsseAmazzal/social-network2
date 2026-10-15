export const Nothifiaction_limit = 20;

const Nothifiacation_types = new Set([
  "new_follower",        
  "follow_request",
  "follow_accepted",     
  "group_join_request",
  "group_request_accepted",
  "group_request_rejected",
  "group_invitation",
  "direct_message",      
  "group_chat_message",   
]);

export function isBellNotificationType(type) {
  return typeof type === "string" && Nothifiacation_types.has(type);
}

export function filterBellNotifications(notifications) {
  return (notifications ?? []).filter((notification) =>
    isBellNotificationType(notification.type),
  );
}

export function addNotificationToBell(
  notifications,
  notification,
  limit = Nothifiaction_limit,
) {
  if (!isBellNotificationType(notification.type)) return notifications;
  if (notifications.some((item) => item.id === notification.id))
    return notifications;
  return [notification, ...notifications].slice(0, limit);
}

export function getNotificationRedirectUrl(notification) {
  const { entityId, actorId, type } = notification;
  switch (type) {
    case "like":
    case "comment":
      return entityId ? `/post/${entityId}` : null;
    
    case "new_follower":
            return actorId ? `/UserProfile?user_id=${actorId}` : null;
    case "follow_accepted":
      return actorId ? `/UserProfile?user_id=${actorId}` : null;
    
    case "follow_request":
      return "/UserProfile?tab=requests";
      
    case "group_event":
      return entityId ? `/gouupDetails?Grouup_id=${entityId}&tab=events` : "/groups";
      
    case "group_join_request":
      return entityId ? `/gouupDetails?Grouup_id=${entityId}&tab=pending` : "/groups";
      
    case "group_request_accepted":
    case "group_invitation":
    case "group_chat_message":
      return entityId ? `/gouupDetails?Grouup_id=${entityId}` : "/groups";
      
    case "direct_message":
      return actorId ? `/chatprivite?Userchat_id=${actorId}` : null;
      
    default:
      return "/";
  }
}