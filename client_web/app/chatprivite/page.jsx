"use client";
import { useCallback, useEffect, useRef, useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import {
  DEFAULT_AVATAR_PATH,
  fetchfromApi,
  getErrorMessage,
  buildBackendAssetUrl,
  ApiError,
} from "@/lib/api";
import {
  exceedsMessageLimit,
  normalizeChatMessage,
  normalizeChatUsers,
} from "@/lib/chat";
import { EMOJI_OPTIONS } from "@/lib/emojis";
import { formatChatTimestamp } from "@/lib/time";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useRealtimeSocket } from "@/providers/RealtimeSocketProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/direct-conversation.module.css";
import { useUnreadConversations } from "@/providers/UnreadConversationProvider"; 

const MESSAGE_PAGE_SIZE = 20;
const STICK_TO_BOTTOM_THRESHOLD = 120;
const TOP_TRIGGER_THRESHOLD = 40;


const RELATIVE_UNIT_SECONDS = {
  second: 1,
  minute: 60,
  hour: 60 * 60,
  day: 60 * 60 * 24,
  week: 60 * 60 * 24 * 7,
  month: 60 * 60 * 24 * 30,
  year: 60 * 60 * 24 * 365,
};

function pad(value) {
  return String(value).padStart(2, "0");
}

/** Formats a Date as the backend's UTC "YYYY-MM-DD HH:MM:SS" cursor. */
function toBackendCursor(date) {
  return (
    `${date.getUTCFullYear()}-${pad(date.getUTCMonth() + 1)}-${pad(
      date.getUTCDate(),
    )}` +
    ` ${pad(date.getUTCHours())}:${pad(date.getUTCMinutes())}:${pad(
      date.getUTCSeconds(),
    )}`
  );
}


function toPaginationCursor(createdAt, now = Date.now()) {
  if (!createdAt) return null;

  const parsed = new Date(createdAt);
  if (!Number.isNaN(parsed.getTime())) {
    return toBackendCursor(parsed);
  }

  const match = /(\d+)\s*(second|minute|hour|day|week|month|year)/i.exec(
    createdAt,
  );
  if (!match) return null;

  const amount = Number(match[1]);
  const unit = match[2].toLowerCase();
  const seconds = RELATIVE_UNIT_SECONDS[unit];
  if (!Number.isFinite(amount) || !seconds) return null;

  return toBackendCursor(new Date(now - amount * seconds * 1000));
}



function isFollowRestrictionError(error) {
  if (error instanceof ApiError) {
    return error.status === 403 || /follow/i.test(error.message);
  }
  return /follow/i.test(getErrorMessage(error, ""));
}

/**
 * Websocket refusals arrive as plain text ({ type: "error", data: "..." }).
 * Only the mutual-follow restriction may lock the composer: for it the front
 * end renders its own message ("You cannot send messages until you follow
 * each other."). Every other refusal — message too long, empty, self-chat …
 * — is only surfaced as a toast so the user can correct the message and
 * resend it.
 */
function isFollowRestrictionText(errorText) {
  return /follow/i.test(errorText || "");
}

function flipPrivateHistoryMessage(message, currentUserId, receiverId) {
  if (!message || !currentUserId || !receiverId) return message;
  if (
    message.senderId !== currentUserId &&
    message.senderId !== receiverId
  ) {
    return message;
  }
  const trueSenderId =
    message.senderId === currentUserId ? receiverId : currentUserId;
  return {
    ...message,
    senderId: trueSenderId,
    receiverId: trueSenderId === currentUserId ? receiverId : currentUserId,
  };
}

function ChatContent() {
  const searchParams = useSearchParams();
  const receiverId = searchParams.get("Userchat_id") || "";

  const { user } = useAuthSession();
  const { latestMessage, typingStatus, userStatus, sendTyping, socket, sendSocketMessage } =
    useRealtimeSocket();
  useToastMessages();
  const { markThreadRead, unreadByUser } = useUnreadConversations();
  const [messages, setMessages] = useState([]);
  const [conversationPartner, setConversationPartner] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isForbidden, setIsForbidden] = useState(false);
  const [canSend, setCanSend] = useState(true);
  const [messageText, setMessageText] = useState("");
  const [isPeerTyping, setIsPeerTyping] = useState(false);
  const [isEmojiPickerOpen, setIsEmojiPickerOpen] = useState(false);
  const [hasSentTypingSignal, setHasSentTypingSignal] = useState(false);
  const [typingDeadline, setTypingDeadline] = useState(0);
  const currentUserId = user?.id ?? "";

  const { showToast } = useToastMessages();

  const [hasMore, setHasMore] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [showScrollButton, setShowScrollButton] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);

  const scrollContainerRef = useRef(null);
  const isFetching = useRef(false);
  const topLatchArmed = useRef(true);

    useEffect(() => {
  if (receiverId && unreadByUser[receiverId] > 0) {
    markThreadRead(receiverId);
  }
  }, [receiverId, messages.length, unreadByUser, markThreadRead]);
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
    if (!receiverId) return;
    let isCurrentRequest = true;
    setIsLoading(true);
    setIsForbidden(false);
    setConversationPartner(null);
    setCanSend(true);
    setHasMore(true);
    setMessages([]);
    setShowScrollButton(false);
    setUnreadCount(0);
    topLatchArmed.current = true;
    isFetching.current = false;

    fetchfromApi(
      `/Usersmessages?Userchat_id=${receiverId}&limit=${MESSAGE_PAGE_SIZE}&last_time=0`,
    )
      .then((response) => {
        if (!isCurrentRequest) return;

        const rawData = response.data || response;
        const rawMessages = rawData.messages || response.messages || [];
        const normalized = Array.isArray(rawMessages)
          ? rawMessages
              .map((m) => normalizeChatMessage(m))
              .filter(Boolean)
              .map((m) => flipPrivateHistoryMessage(m, currentUserId, receiverId))
          : [];

        setMessages(normalized);

        const partnerData =
          response.conversationalist || response.talker || rawData.talker;
        if (partnerData) {
          setConversationPartner(normalizeChatUsers([partnerData])[0] || null);
        }

        setCanSend(response.can_send !== false);
        if (normalized.length < MESSAGE_PAGE_SIZE) setHasMore(false);

        requestAnimationFrame(() => scrollToBottom("auto"));
      })
      .catch((error) => {
        if (!isCurrentRequest) return;
        if (isFollowRestrictionError(error)) {
          setIsForbidden(true);
          showToast(
            getErrorMessage(error, "Follow the user to chat with him."),
            "error",
          );
          return;
        }
        if (error instanceof ApiError && error.status === 404) return;
        showToast(
          getErrorMessage(error, "Failed to load the conversation."),
          "error",
        );
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
          setTimeout(() => {
        isFetching.current = false;
        setIsLoadingMore(false);
      }, 500); 
      });

    return () => {
      isCurrentRequest = false;
    };
  }, [receiverId, currentUserId, scrollToBottom]);

  useEffect(() => {
    if (!latestMessage || "groupId" in latestMessage) return;
      if (latestMessage.type === "error"||latestMessage.Datatype === "error") {
      const errorText = String(latestMessage.data ?? "").trim();
      if (isFollowRestrictionText(errorText)) {
        ///avoid send meesaage on front when there is an error message
        setCanSend(false);
        return;
      }
      showToast?.(errorText || "The message could not be sent.", "error");
    return;
  }
    const isRelated =
      latestMessage.senderId === receiverId ||
      latestMessage.receiverId === receiverId;
    if (!isRelated) return;

    const wasNearBottom = getIsNearBottom();
    let added = false;
    setMessages((prev) => {
      if (
        latestMessage.id &&
        prev.some((m) => m.id === latestMessage.id)
      ) {
        return prev;
      }

      const pendingIndex = prev.findIndex(
        (m) =>
          m.pending &&
          m.senderId === latestMessage.senderId &&
          m.text === latestMessage.text,
      );
      if (pendingIndex !== -1) {
        const next = [...prev];
        next[pendingIndex] = latestMessage;
        added = false;
        return next;
      }

      added = true;
      return [...prev, latestMessage];
    });

    if (added && wasNearBottom) {
      requestAnimationFrame(() => scrollToBottom("smooth"));
    } else if (added) {
      setShowScrollButton(true);
      setUnreadCount((count) => count + 1);
    }
  }, [latestMessage, receiverId, getIsNearBottom, scrollToBottom]);

  useEffect(() => {
    setIsPeerTyping(
      Boolean(
        typingStatus &&
          typingStatus.from === receiverId &&
          typingStatus.to === currentUserId &&
          typingStatus.isTyping,
      ),
    );
  }, [currentUserId, receiverId, typingStatus]);

  useEffect(() => {
    if (!userStatus || !conversationPartner) return;
    if (userStatus.userId !== conversationPartner.id) return;
    setConversationPartner((currentPartner) =>
      currentPartner
        ? { ...currentPartner, isConnected: userStatus.isOnline }
        : currentPartner,
    );
  }, [conversationPartner?.id, userStatus]);

  useEffect(() => {
    if (isPeerTyping && getIsNearBottom()) scrollToBottom("smooth");
  }, [isPeerTyping, getIsNearBottom, scrollToBottom]);

  useEffect(() => {
    if (!hasSentTypingSignal || typingDeadline === 0) return;
    const timeoutId = window.setTimeout(
      () => {
        sendTyping(currentUserId, receiverId, false);
        setHasSentTypingSignal(false);
        setTypingDeadline(0);
      },
      Math.max(0, typingDeadline - Date.now()),
    );
    return () => window.clearTimeout(timeoutId);
  }, [currentUserId, hasSentTypingSignal, receiverId, typingDeadline, sendTyping]);

  useEffect(() => {
    if (!isEmojiPickerOpen) return;
    function closeEmojiPicker(event) {
      const target = event.target;
      if (target instanceof Element && target.closest("#direct-emoji-picker"))
        return;
      setIsEmojiPickerOpen(false);
    }
    document.addEventListener("mousedown", closeEmojiPicker);
    return () => document.removeEventListener("mousedown", closeEmojiPicker);
  }, [isEmojiPickerOpen]);

  useEffect(() => {
    return () => {
      if (currentUserId && receiverId) {
        sendTyping(currentUserId, receiverId, false);
      }
    };
  }, [currentUserId, receiverId, sendTyping]);

  const handleInputChange = (value) => {
    setMessageText(value);
    if (!value.trim()) {
      if (hasSentTypingSignal) sendTyping(currentUserId, receiverId, false);
      setHasSentTypingSignal(false);
      return;
    }
    if (!hasSentTypingSignal) {
      sendTyping(currentUserId, receiverId, true);
      setHasSentTypingSignal(true);
    }
  };

const handleLoadMore = async () => {
  if (isFetching.current || !hasMore) return;

  const container = scrollContainerRef.current;
  if (!container) return;

  isFetching.current = true;
  setIsLoadingMore(true);

  const oldestMessage = messages[0];
  const cursor = toPaginationCursor(oldestMessage?.createdAt);
  
  if (!cursor) {
    setHasMore(false);
    isFetching.current = false;
    setIsLoadingMore(false);
    return;
  }

  const previousHeight = container.scrollHeight;
  const previousTop = container.scrollTop;

  try {
    const response = await fetchfromApi(
      `/Usersmessages?Userchat_id=${receiverId}&limit=${MESSAGE_PAGE_SIZE}&last_time=${encodeURIComponent(cursor)}`
    );
    
    const rawMessages = response.messages || response.data?.messages || [];
    const olderMessages = Array.isArray(rawMessages)
      ? rawMessages
          .map((m) => normalizeChatMessage(m))
          .filter(Boolean)
          .map((m) => flipPrivateHistoryMessage(m, currentUserId, receiverId))
      : [];

    const existingIds = new Set(messages.map((m) => m.id));
    const uniqueOlder = olderMessages.filter((m) => !existingIds.has(m.id));

    if (uniqueOlder.length === 0) {
      setHasMore(false);
    } else {
      setMessages((prev) => [...uniqueOlder, ...prev]);
      
      if (olderMessages.length < MESSAGE_PAGE_SIZE) {
        setHasMore(false);
      }

      requestAnimationFrame(() => {
        if (scrollContainerRef.current) {
          const addedHeight = scrollContainerRef.current.scrollHeight - previousHeight;
          scrollContainerRef.current.scrollTop = previousTop + addedHeight;
        }
      });
    }
  } catch (error) {
    if (isFollowRestrictionError(error)) {
      setIsForbidden(true);
      showToast(
        getErrorMessage(error, "Follow the user to chat with him."),
        "error",
      );
    } else if (error instanceof ApiError) {
      showToast(
        getErrorMessage(error, "Failed to load older messages."),
        "error",
      );
    } else {
      console.error("Load more error:", error);
    }
  } finally {
    setTimeout(() => {
      isFetching.current = false;
      setIsLoadingMore(false);
    }, 1000); 
  }
};

  const handleScroll = (event) => {
    const { scrollTop } = event.currentTarget;

    syncScrollButton();

    // history page can be requested only after leaving the trigger zone.
    if (scrollTop > TOP_TRIGGER_THRESHOLD) {
      topLatchArmed.current = true;
      return;
    }

    if (
      scrollTop <= TOP_TRIGGER_THRESHOLD &&
      topLatchArmed.current &&
      hasMore &&
      !isFetching.current &&
      !isForbidden 
    ) {
      void handleLoadMore();
    }
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const trimmedMessage = messageText.trim();
     if (!trimmedMessage || !canSend || socket?.readyState !== WebSocket.OPEN) {
      if (!canSend) showToast("Permissions changed. You can no longer message this user.", "error");
        return;
      }

    // Mirror the backend's length check so the user gets a toast right away
    // instead of waiting for the refusal over the websocket.
    if (exceedsMessageLimit(trimmedMessage)) {
      showToast("The message is too long. Maximum 1000 characters.", "error");
      return;
    }

    const optimisticMessage = {
      id: `pending-${Date.now()}-${Math.random().toString(36).slice(2)}`,
      senderId: currentUserId,
      receiverId,
      text: trimmedMessage,
      createdAt: new Date().toISOString(),
      pending: true,
    };

    const sent = await sendSocketMessage({
      wsdatatype: "private_message",
      wsdata: {
        userchatwith_id: receiverId,
        message: trimmedMessage,
      },
    });
    if (!sent) {
      return;
    }

    // Show the message immediately; the socket echo replaces it afterwards.
    setMessages((prev) => [...prev, optimisticMessage]);

    setMessageText("");
    if (hasSentTypingSignal) sendTyping(currentUserId, receiverId, false);
    setHasSentTypingSignal(false);
    setTypingDeadline(0);

    // Follow our own message down to the bottom.
    requestAnimationFrame(() => scrollToBottom("smooth"));
  };

  if (isLoading) {
    return (
      <p className={uiStyles.uiEmptyState} role="status">
        Loading conversation…
      </p>
    );
  }
  if (isForbidden) {
    return (
      <p className={uiStyles.uiEmptyState}>
        You must follow each other before starting a conversation.{" "}
        <Link href={`/UserProfile?user_id=${receiverId}`}>View profile</Link>
      </p>
    );
  }
  if (!conversationPartner) {
    return <p className={uiStyles.uiEmptyState}>User not found.</p>;
  }
  return (
    <section
      className={styles.directChatPanel}
      aria-label={`Conversation with ${conversationPartner.nickname}`}
    >
      <header className={styles.directChatHeader}>
        <img
          src={buildBackendAssetUrl(
            conversationPartner.avatarUrl,
            DEFAULT_AVATAR_PATH,
          )}
          alt=""
          width={46}
          height={46}
          className={styles.directHeaderAvatar}
        />
        <div>
          <h1>{conversationPartner.nickname}</h1>
          <span
            className={
              conversationPartner.isConnected ? styles.directOnline : styles.directOffline
            }
          >
            {conversationPartner.isConnected ? "Online" : "Offline"}
          </span>
        </div>
      </header>

      <div className={styles.directMessageListWrapper}>
        <div
          className={styles.directMessageList}
          ref={scrollContainerRef}
          onScroll={handleScroll}
        >
          {isLoadingMore && (
            <p className={styles.directLoadingInfo}>Loading history...</p>
          )}
          {!messages || messages.length === 0 ? (
            <p className={styles.directEmptyState}>Say hello! 👋</p>
          ) : (
            <>
              {messages.map((msg, index) => (
                <div
                  key={msg.id || `${msg.senderId}-${msg.createdAt}-${index}`}
                  className={`${styles.directMessage} ${
                    msg.senderId === currentUserId
                      ? styles.directCurrentUserMessage
                      : styles.directPartnerMessage
                  }`}
                >
                  <span>{msg.text || msg.content}</span>
                  <small>{formatChatTimestamp(msg.createdAt)}</small>
                </div>
              ))}
              {isPeerTyping && (
                <p className={styles.directTypingIndicator}>typing…</p>
              )}
            </>
          )}
        </div>

        {showScrollButton && (
          <button
            type="button"
            className={styles.directScrollToBottomButton}
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
              <span className={styles.directScrollToBottomBadge}>{unreadCount}</span>
            )}
          </button>
        )}
      </div>

      <div className={styles.directComposerArea}>
        {!canSend && (
          <p className={styles.directPermissionWarning}>
            You cannot send messages until you follow each other.
          </p>
        )}
        <div id="direct-emoji-picker" className={styles.directEmojiContainer}>
          {isEmojiPickerOpen && (
            <div
              className={styles.directEmojiPicker}
              role="menu"
              aria-label="Emoji picker"
            >
              {EMOJI_OPTIONS.map((emoji) => (
                <button
                  type="button"
                  role="menuitem"
                  key={emoji}
                  onClick={() => handleInputChange(`${messageText}${emoji}`)}
                >
                  {emoji}
                </button>
              ))}
            </div>
          )}
          <form className={styles.directComposer} onSubmit={handleSubmit}>
            <button
              type="button"
              className={styles.directEmojiButton}
              aria-label="Toggle emoji picker"
              aria-expanded={isEmojiPickerOpen}
              onClick={() => setIsEmojiPickerOpen((isOpen) => !isOpen)}
            >
              🙂
            </button>
            <label className="srOnly" htmlFor="direct-message-input">
              Message
            </label>
            <input
              id="direct-message-input"
              value={messageText}
              onChange={(event) => handleInputChange(event.target.value)}
              placeholder={
                canSend ? "Write a message…" : "Messaging unavailable"
              }
              maxLength={1000}
              disabled={!canSend}
              autoComplete="off"
            />
{canSend && (
        <button
          type="submit"
          className={styles.directSendButton}
          disabled={!messageText.trim()}
          aria-label="Send message"
        >
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
          >
            <path d="m22 2-7 20-4-9-9-4Z" />
            <path d="M22 2 11 13" />
          </svg>
        </button>
      )}
          </form>
        </div>
      </div>
    </section>
  );
}

export default function DirectConversationRoute() {
  return (
    <Suspense fallback={<p>Loading...</p>}>
      <ChatContent />
    </Suspense>
  );
}
