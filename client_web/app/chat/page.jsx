"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import {
  DEFAULT_AVATAR_PATH,
  buildBackendAssetUrl,
  fetchfromApi,
  getErrorMessage,
} from "@/lib/api";
import { normalizeChatUsers } from "@/lib/chat";
import { fetchAllJoinnedGroups } from "@/lib/groups";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import { useUnreadConversations } from "@/providers/UnreadConversationProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/chat-directory.module.css";

export default function ConversationDirectoryRoute() {
  const { user } = useAuthSession();
  const { showToast } = useToastMessages();
  const { unreadByUser, unreadByGroup, markThreadRead, markGroupRead } =
    useUnreadConversations();
  const [activeTab, setActiveTab] = useState("people");
  const [chatUsers, setChatUsers] = useState([]);
  const [joinedGroups, setJoinedGroups] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const userId = user?.id ?? "";

  useEffect(() => {
    if (!userId) return;
    let isCurrentRequest = true;
    setIsLoading(true);
    void Promise.all([
      fetchfromApi("/UsersForChat"),
      fetchAllJoinnedGroups().catch(() => []),
    ])
      .then(([chatResponse, groups]) => {
        if (!isCurrentRequest) return;
        setChatUsers(normalizeChatUsers(chatResponse.Userstochat));
        setJoinedGroups(groups);
      })
      .catch((error) => {
        if (isCurrentRequest) {
          showToast(getErrorMessage(error, "Failed to load chats."), "error");
        }
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [userId]);

  return (
    <main className={uiStyles.uiPage}>
      <header className={styles.chatHeader}>
        <div>
          <h1>Chats</h1>
          <p>Continue a conversation with people or groups.</p>
        </div>
      </header>

      <div className={styles.chatTabs} role="tablist" aria-label="Chat lists">
        <button
          type="button"
          role="tab"
          aria-selected={activeTab === "people"}
          className={activeTab === "people" ? styles.chatActiveTab : ""}
          onClick={() => setActiveTab("people")}
        >
          People
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={activeTab === "groups"}
          className={activeTab === "groups" ? styles.chatActiveTab : ""}
          onClick={() => setActiveTab("groups")}
        >
          Groups
        </button>
      </div>

      {isLoading ? (
        <p className={uiStyles.uiEmptyState} role="status">
          Loading chats…
        </p>
      ) : activeTab === "people" ? (
        chatUsers.length > 0 ? (
          <div className={styles.chatDirectoryGrid}>
            {chatUsers.map((chatUser) => {
              const unreadCount = unreadByUser[chatUser.id] ?? 0;
              return (
                <Link
                  key={chatUser.id}
                  href={`/chatprivite?Userchat_id=${chatUser.id}`}
                  className={styles.chatPersonCard}
                  onClick={() => void markThreadRead(chatUser.id)}
                >
                  <span className={styles.chatAvatarContainer}>
                    <img
                      src={buildBackendAssetUrl(
                        chatUser.avatarUrl,
                        DEFAULT_AVATAR_PATH,
                      )}
                      alt=""
                      width={48}
                      height={48}
                      className={styles.chatAvatar}
                    />
                    <span
                      className={`${styles.chatOnlineDot} ${chatUser.isConnected ? styles.chatOnline : ""}`}
                      aria-label={chatUser.isConnected ? "Online" : "Offline"}
                    />
                  </span>
                  <span className={styles.chatCardText}>
                    <strong>{chatUser.nickname}</strong>
                    <small>Start a conversation</small>
                  </span>
                  {unreadCount > 0 && (
                    <span className={styles.chatUnreadBadge}>{unreadCount}</span>
                  )}
                </Link>
              );
            })}
          </div>
        ) : (
          <p className={uiStyles.uiEmptyState}>
            You have no connections yet.{" "}
            <Link href="/followers">Find people</Link>
          </p>
        )
      ) : joinedGroups.length > 0 ? (
        <div className={styles.chatDirectoryGrid}>
          {joinedGroups.map((group) => {
            const unreadCount = unreadByGroup[group.id] ?? 0;
            return (
              <Link
                key={group.id}
                href={`/gouupDetails?Grouup_id=${group.id}&tab=chat`}
                className={styles.chatPersonCard}
                onClick={() => void markGroupRead(group.id)}
              >
                <span className={styles.chatGroupAvatar} aria-hidden="true">
                  {group.title.slice(0, 1).toUpperCase()}
                </span>
                <span className={styles.chatCardText}>
                  <strong>{group.title}</strong>
                  <small>{group.description || "Group conversation"}</small>
                </span>
                {unreadCount > 0 && (
                  <span className={styles.chatUnreadBadge}>{unreadCount}</span>
                )}
              </Link>
            );
          })}
        </div>
      ) : (
        <p className={uiStyles.uiEmptyState}>
          Join a group to start a conversation.{" "}
          <Link href="/groups">Browse groups</Link>
        </p>
      )}
    </main>
  );
}
