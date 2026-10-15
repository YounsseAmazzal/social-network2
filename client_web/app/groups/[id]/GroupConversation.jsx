"use client";
import { useCallback, useEffect, useState, useRef } from "react";
import { fetchfromApi, getErrorMessage } from "@/lib/api";
import {
  exceedsMessageLimit,
  normalizeGroupMessages,
} from "@/lib/chat";
import { EMOJI_OPTIONS } from "@/lib/emojis";
import { formatChatTimestamp } from "@/lib/time";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useUnreadConversations } from "@/providers/UnreadConversationProvider";
import { useRealtimeSocket } from "@/providers/RealtimeSocketProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import styles from "@/app/styles/group-conversation.module.css";
// Distance from the top that triggers loading the previous page.
const TOP_TRIGGER_THRESHOLD = 40;
// If the user is within this many pixels of the bottom we consider the
// conversation "pinned" and keep it scrolled down for new messages.
const STICK_TO_BOTTOM_THRESHOLD = 120;
export default function GroupConversation({ groupId }) {
  const { user } = useAuthSession();
  const { socket, latestMessage } = useRealtimeSocket();
  const { showToast } = useToastMessages();

    const [hasMore, setHasMore] = useState(true);
    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const [showScrollButton, setShowScrollButton] = useState(false);
    const [unreadCount, setUnreadCount] = useState(0);
    const scrollContainerRef = useRef(null);
    const isFetching = useRef(false);
    const topLatchArmed = useRef(true);


  const { markGroupRead, setActiveGroupConversation } =useUnreadConversations();
  const [messages, setMessages] = useState([]);
  const [messageText, setMessageText] = useState("");
  const [isEmojiPickerOpen, setIsEmojiPickerOpen] = useState(false);
  const userId = user?.id ?? "";

  const Limit=10

  const scrollToBottom = useCallback((behavior = "auto") => {
    const container = scrollContainerRef.current;
    if (!container) return;
    if (behavior === "smooth") {
      container.scrollTo({ top: container.scrollHeight, behavior: "smooth" });
    } else {
      container.scrollTop = container.scrollHeight;
    }
  }, []);

  const getIsNearBottom = useCallback(() => {
    const container = scrollContainerRef.current;
    if (!container) return true;
    const distance =
      container.scrollHeight - container.scrollTop - container.clientHeight;
    return distance <= STICK_TO_BOTTOM_THRESHOLD;
  }, []);

  const syncScrollButton = useCallback(() => {
    const isNearBottom = getIsNearBottom();
    setShowScrollButton(!isNearBottom);
    if (isNearBottom) setUnreadCount(0);
    return isNearBottom;
  }, [getIsNearBottom]);

  const handleScrollToBottomClick = useCallback(() => {
    scrollToBottom("smooth");
    setUnreadCount(0);
    setShowScrollButton(false);
  }, [scrollToBottom]);

useEffect(() => {
  if (!groupId) return;
  let isCurrentRequest = true;

  setHasMore(true);
  setShowScrollButton(false);
  setUnreadCount(0);
  topLatchArmed.current = true;
  isFetching.current = false;

  fetchfromApi(`/Grouupmessages?Grouup_id=${groupId}&limit=${Limit}&last_time=0`)
    .then((response) => {
      if (!isCurrentRequest) return;
      
      const rawMessages = response.Groupmessages || [];
      const normalized = normalizeGroupMessages(rawMessages).reverse(); 
      
      setMessages(normalized);
      if (normalized.length < Limit) setHasMore(false);
      
      requestAnimationFrame(() => scrollToBottom("auto"));
    })
    .catch((error) => {
      if (isCurrentRequest) {
        showToast(
          getErrorMessage(error, "Failed to load messages."),
          "error",
        );
      }
    });

  return () => { isCurrentRequest = false; };
}, [groupId, scrollToBottom]);


  useEffect(() => {
    setActiveGroupConversation(groupId);
    return () => setActiveGroupConversation(null);
  }, [groupId]);


const handleLoadMore = async () => {
  if (isFetching.current || !hasMore) return;

  const container = scrollContainerRef.current;
  if (!container) return;

  // Nothing to paginate if the messages already fit without scrolling.
  if (container.scrollHeight <= container.clientHeight) return;

  isFetching.current = true;
  topLatchArmed.current = false; // block further loads until we scroll back down
  setIsLoadingMore(true);

  const oldestMessage = messages[0];
  const cursor = oldestMessage?.createdAt;
  if (!cursor) {
    isFetching.current = false;
    setIsLoadingMore(false);
    return;
  }
  const previousHeight = container.scrollHeight;
  const previousTop = container.scrollTop;

  try {
    const response = await fetchfromApi(`/Grouupmessages?Grouup_id=${groupId}&limit=${Limit}&last_time=${encodeURIComponent(cursor)}`);
    const rawMessages = response.Groupmessages || [];
    const newOldMessages = normalizeGroupMessages(rawMessages).reverse();

    const existingIds = new Set(messages.map((m) => m.id));
    const uniqueOlder = newOldMessages.filter(
      (m) => !m.id || !existingIds.has(m.id),
    );

    if (uniqueOlder.length > 0) {
      setMessages((prev) => {
        const prevIds = new Set(prev.map((m) => m.id));
        const toPrepend = uniqueOlder.filter(
          (m) => !m.id || !prevIds.has(m.id),
        );
        return [...toPrepend, ...prev];
      });
    }

    // A short page means we reached the beginning of the conversation.
    // A page that adds nothing new means the cursor did not advance; stop
    // paging so we never loop and hammer the server.
    if (newOldMessages.length < Limit || uniqueOlder.length === 0) {
      setHasMore(false);
    }

    requestAnimationFrame(() => {
      if (scrollContainerRef.current) {
        const addedHeight =
          scrollContainerRef.current.scrollHeight - previousHeight;
        scrollContainerRef.current.scrollTop = previousTop + addedHeight;
      }
    });
  } catch (error) {
    showToast(
      getErrorMessage(error, "Failed to load older messages."),
      "error",
    );
  } finally {
    isFetching.current = false;
    setIsLoadingMore(false);
  }
};

  const handleScroll = (event) => {
    const { scrollTop } = event.currentTarget;

    syncScrollButton();

    if (scrollTop > TOP_TRIGGER_THRESHOLD) {
      topLatchArmed.current = true;
      return;
    }

    if (
      scrollTop <= TOP_TRIGGER_THRESHOLD &&
      topLatchArmed.current &&
      hasMore &&
      !isFetching.current
    ) {
      void handleLoadMore();
    }
  };



  useEffect(() => {
  if (!latestMessage || !latestMessage.groupId) return;
    if (latestMessage.groupId !== groupId) return;
    const wasNearBottom = getIsNearBottom();
    let added = false;
    setMessages((currentMessages) => {
      if (
        latestMessage.id &&
        currentMessages.some((message) => message.id === latestMessage.id)
      ) {
        return currentMessages;
      }
      added = true;
      return [...currentMessages, latestMessage];
    });

    if (added && wasNearBottom) {
      // Follow the conversation when the user is already at the bottom.
      requestAnimationFrame(() => scrollToBottom("smooth"));
    } else if (added) {
      // Otherwise, surface the jump-to-bottom button with an unread count.
      setShowScrollButton(true);
      setUnreadCount((count) => count + 1);
    }

    if (latestMessage.senderId !== userId) void markGroupRead(groupId);
  }, [groupId, latestMessage, userId, getIsNearBottom, scrollToBottom]);
  // useEffect(() => {
  //   document
  //     .getElementById("group-message-end")
  //     ?.scrollIntoView({);
  // }, [messages]);
  useEffect(() => {
    if (!isEmojiPickerOpen) return;
    function closeEmojiPicker(event) {
      const target = event.target;
      if (target instanceof Element && target.closest("#group-emoji-picker"))
        return;
      setIsEmojiPickerOpen(false);
    }
    document.addEventListener("mousedown", closeEmojiPicker);
    return () => document.removeEventListener("mousedown", closeEmojiPicker);
  }, [isEmojiPickerOpen]);
  function handleSubmit(event) {
    event.preventDefault();
    const trimmedMessage = messageText.trim();
    if (!trimmedMessage || socket?.readyState !== WebSocket.OPEN) return;
    if (exceedsMessageLimit(trimmedMessage)) {
      showToast("The message is too long. Maximum 1000 characters.", "error");
      return;
    }
    socket.send(
      JSON.stringify({
        wsdatatype: "send_message_group",
        wsdata: {
          group_id: groupId,
          message: trimmedMessage,
        },
      }),
    );
    setMessageText("");
    setIsEmojiPickerOpen(false);
  }
  return (
    <section className={styles.groupChatPanel} aria-label="Group conversation">
      {/* <div className={styles.groupMessageList} aria-live="polite"> */}
      <div className={styles.groupMessageListWrapper}>
      <div 
        className={styles.groupMessageList} 
        ref={scrollContainerRef} 
        onScroll={handleScroll}
      >
                {isLoadingMore && <p className={styles.loadingInfo}>Loading older messages...</p>}

        {messages.length === 0 ? (
          <p className={styles.groupEmptyState}>No messages yet. Say hello!</p>
        ) : (
          
          
          messages.map((message, index) => {
            const isCurrentUser = message.senderId === userId;
            return (
              <div
              key={
                message.id ||
                `${message.senderId}-${message.createdAt}-${index}`
              }
              className={`${styles.groupMessageRow} ${isCurrentUser ? styles.groupCurrentUserMessage : ""}`}
              >
                {!isCurrentUser && <strong>{message.senderNickname}</strong>}
                <div className={styles.groupMessageBubble}>
                  <span>{message.text}</span>
                  {message.createdAt && (
                    <time dateTime={message.createdAt}>
                      {formatChatTimestamp(message.createdAt)}
                    </time>
                  )}
                </div>
              </div>
            );
          })
        )}
      <div id="group-message-end" />
      </div>

      {showScrollButton && (
        <button
          type="button"
          className={styles.groupScrollToBottomButton}
          onClick={handleScrollToBottomClick}
          aria-label={
            unreadCount > 0
              ? `Scroll to bottom, ${unreadCount} new message${
                  unreadCount > 1 ? "s" : ""
                }`
              : "Scroll to bottom"
          }
        >
          <svg
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
          >
            <path d="M12 5v14" />
            <path d="m19 12-7 7-7-7" />
          </svg>
          {unreadCount > 0 && (
            <span className={styles.groupScrollToBottomBadge}>{unreadCount}</span>
          )}
        </button>
      )}
      </div>

      <div className={styles.groupComposerArea}>
        <div id="group-emoji-picker" className={styles.groupEmojiContainer}>
          {isEmojiPickerOpen && (
            <div
            className={styles.groupEmojiPicker}
            role="menu"
            aria-label="Emoji picker"
            >
              {EMOJI_OPTIONS.map((emoji) => (
                <button
                type="button"
                role="menuitem"
                key={emoji}
                onClick={() =>
                  setMessageText((current) => `${current}${emoji}`)
                }
                >
                  {emoji}
                </button>
              ))}
            </div>
          )}
          <form className={styles.groupComposer} onSubmit={handleSubmit}>
            <button
              type="button"
              className={styles.groupEmojiButton}
              aria-label="Toggle emoji picker"
              aria-expanded={isEmojiPickerOpen}
              onClick={() => setIsEmojiPickerOpen((isOpen) => !isOpen)}
              >
              🙂
            </button>
            <label className="srOnly" htmlFor="group-message-input">
              Group message
            </label>
            <input
              id="group-message-input"
              value={messageText}
              onChange={(event) => setMessageText(event.target.value)}
              placeholder="Write a group message…"
              maxLength={1000}
              autoComplete="off"
              />
            <button
              type="submit"
              className={styles.groupSendButton}
              disabled={
                !messageText.trim() || socket?.readyState !== WebSocket.OPEN
              }
              aria-label="Send message"
              >
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                aria-hidden="true"
              >
                <path d="m22 2-7 20-4-9-9-4Z" />
                <path d="M22 2 11 13" />
              </svg>
            </button>
          </form>
        </div>
      </div>
                {/* </div> */}
    </section>
  );
}
