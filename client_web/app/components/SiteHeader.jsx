"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  getErrorMessage,
  buildBackendAssetUrl,
  DEFAULT_AVATAR_PATH,
} from "@/lib/api";
import { InvitaiongroupResponse } from "@/lib/groups";
import { getNotificationRedirectUrl } from "@/lib/notifications";
import { describeRelativeTimestamp } from "@/lib/time";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useNotificationCenter } from "@/providers/NotificationCenterProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import styles from "@/app/styles/site-header.module.css";
export default function SiteHeader() {
  const { user, logout } = useAuthSession();
  const {
    notifications,
    unreadCount,
    isLoading,
    refresh,
    markAllRead,
    markRead,
  } = useNotificationCenter();
  const { showToast } = useToastMessages();
  const router = useRouter();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);
  const [areNotificationsOpen, setAreNotificationsOpen] = useState(false);
  const [respondingInvitation, setRespondingInvitation] = useState(null);
  useEffect(() => {
    function closeMenusOnOutsideClick(event) {
      const target = event.target;
      if (!(target instanceof Element)) return;
      if (target.closest("[data-header-dropdown]")) return;
      setIsProfileMenuOpen(false);
      setAreNotificationsOpen(false);
    }
    document.addEventListener("mousedown", closeMenusOnOutsideClick);
    return () =>
      document.removeEventListener("mousedown", closeMenusOnOutsideClick);
  }, []);
  function closeMobileMenu() {
    setIsMobileMenuOpen(false);
  }
  async function handleLogout() {
    setIsProfileMenuOpen(false);
    await logout();
  }
 async function handleInvitationResponse(notificationId, groupId, decision) {
    const actionKey = `${notificationId}-${decision}`;
    setRespondingInvitation(actionKey);
    try {
      await InvitaiongroupResponse(groupId, decision);
      
      await markRead(notificationId); 
      await refresh();

      if (decision === "acceptinvitasion") {
          router.push(`/gouupDetails?Grouup_id=${groupId}`);
      }
      
      showToast(`Invitation ${decision}ed successfully`, "success");
    } catch (error) {
      showToast(getErrorMessage(error, "Action failed"), "error");
    } finally {
      setRespondingInvitation(null);
    }
}
  async function markAllNotificationsRead() {
    try {
      await markAllRead();
    } catch (error) {
      showToast(
        getErrorMessage(error, "Failed to update notifications."),
        "error",
      );
    }
  }
  async function markNotificationRead(notificationId) {
    try {
      await markRead(notificationId);
    } catch (error) {
      showToast(
        getErrorMessage(error, "Failed to update the notification."),
        "error",
      );
    }
  }
  function toggleNotifications() {
    //console.log("toggele all those nothifications")
    const nextOpenState = !areNotificationsOpen;
    setAreNotificationsOpen(nextOpenState);
    setIsProfileMenuOpen(false);
    if (nextOpenState) {
      void refresh().catch((error) => {
        showToast(
          getErrorMessage(error, "Failed to load notifications."),
          "error",
        );
      });
    }
  }

//   const generalNotifications = notifications.filter(n => 
//     n.type !== "direct_message" && 
//     n.type !== "group_chat_message"
// );
  return (
    <header className={styles.header}>
      <div className={styles.headerBrandArea}>
        <Image
          src="/images/solosocial.jpeg"
          alt=""
          width={44}
          height={44}
          className={styles.headerLogo}
          priority
        />
        <span className={styles.headerBrandName}>social</span>
      </div>

      <button
        type="button"
        className={styles.headerMobileMenuButton}
        onClick={() => setIsMobileMenuOpen((isOpen) => !isOpen)}
        aria-expanded={isMobileMenuOpen}
        aria-controls="primary-navigation"
        aria-label={isMobileMenuOpen ? "Close navigation" : "Open navigation"}
      >
        <span aria-hidden="true">{isMobileMenuOpen ? "×" : "☰"}</span>
      </button>

      <nav
        id="primary-navigation"
        className={`${styles.headerNavigation} ${isMobileMenuOpen ? styles.headerNavigationOpen : ""}`}
        aria-label="Primary navigation"
      >
        {user ? (
          <>
            <Link href="/" className={styles.headerNavLink} onClick={closeMobileMenu}>
              Home
            </Link>
            <div className={styles.headerDropdownContainer} data-header-dropdown>
              <button
                type="button"
                className={styles.headerIconButton}
                onClick={toggleNotifications}
                aria-label={`Notifications${unreadCount ? `, ${unreadCount} unread` : ""}`}
                aria-expanded={areNotificationsOpen}
              >
                <Image
                  src="/icons/notifications-icon.svg"
                  alt="Notifications"
                  width={22}
                  height={22}
                />
                {unreadCount > 0 && (
                  <span className={styles.headerNotificationBadge}>
                    {unreadCount}
                  </span>
                )}
              </button>

              {areNotificationsOpen && (
                <section className={styles.headerNotificationDropdown}>
                  <div className={styles.headerDropdownHeader}>
                    <strong>Notifications</strong>
                    {unreadCount > 0 && (
                      <button
                        type="button"
                        className={styles.headerTextButton}
                        onClick={() => void markAllNotificationsRead()}
                      >
                        Mark all read
                      </button>
                    )}
                  </div>

                  <div className={styles.headerNotificationList}>
                    {isLoading ? (
                      <p className={styles.headerDropdownEmpty}>
                        Loading notifications…
                      </p>
                    ) : notifications.length === 0 ? (
                      <p className={styles.headerDropdownEmpty}>
                        No new notifications.
                      </p>
                    ) : (
                      notifications.map((notification) => {
                        const isInvitation =
                          notification.type === "group_invitation";
                        const itemClassName = `${styles.headerNotificationItem} ${notification.isRead ? "" : styles.headerUnreadNotification}`;
                        if (isInvitation) {
                          return (
                            <article
                              key={notification.id}
                              className={itemClassName}
                            >
                              <p>{notification.message}</p>
                              <time dateTime={notification.createdAt}>
                                {describeRelativeTimestamp(
                                  notification.createdAt,
                                )}
                              </time>
                              {!notification.isRead && (
                                <div className={styles.headerInvitationActions}>
                                  <button
                                    type="button"
                                    className={styles.headerAcceptButton}
                                    disabled={respondingInvitation !== null}
                                    onClick={() =>
                                      void handleInvitationResponse(
                                        notification.id,
                                        notification.entityId,
                                        "acceptinvitasion",
                                      )
                                    }
                                  >
                                    {respondingInvitation ===
                                    `${notification.id}-accept`
                                      ? "Accepting…"
                                      : "Accept"}
                                  </button>
                                  <button
                                    type="button"
                                    className={styles.headerRefuseButton}
                                    disabled={respondingInvitation !== null}
                                    onClick={() =>
                                      void handleInvitationResponse(
                                        notification.id,
                                        notification.entityId,
                                        "refuseinvitaion",
                                      )
                                    }
                                  >
                                    {respondingInvitation ===
                                    `${notification.id}-refuse`
                                      ? "Declining…"
                                      : "Decline"}
                                  </button>
                                </div>
                              )}
                            </article>
                          );
                        }
                        return (
                          <button
                            type="button"
                            key={notification.id}
                            className={itemClassName}
                            onClick={() => {
                              void markNotificationRead(notification.id);
                              const href = getNotificationRedirectUrl(notification);
                              if (href) router.push(href);
                            }}
                          >
                            <span>{notification.message}</span>
                            <time dateTime={notification.createdAt}>
                              {describeRelativeTimestamp(
                                notification.createdAt,
                              )}
                            </time>
                          </button>
                        );
                      })
                    )}
                  </div>
                </section>
              )}
            </div>

            <div className={styles.headerDropdownContainer} data-header-dropdown>
              <button
                type="button"
                className={`${styles.headerProfileButton} ${isProfileMenuOpen ? styles.headerProfileButtonActive : ""}`}
                onClick={() => {
                  setIsProfileMenuOpen((isOpen) => !isOpen);
                  setAreNotificationsOpen(false);
                }}
                aria-expanded={isProfileMenuOpen}
                aria-label="Open profile menu"
              >
                <img
                  src={buildBackendAssetUrl(
                    user.avatar ||DEFAULT_AVATAR_PATH,
                  )}
                  alt=""
                  width={42}
                  height={42}
                />
              </button>

              {isProfileMenuOpen && (
                <section className={styles.headerProfileDropdown}>
                  <div className={styles.headerProfileSummary}>
                    <strong>
                      {user.firstname} {user.lastname}
                    </strong>
                    <span>@{user.username || user.nickname}</span>
                  </div>
                  <Link
                    href="/UserProfile"
                    className={styles.headerDropdownLink}
                    onClick={() => {
                      setIsProfileMenuOpen(false);
                      closeMobileMenu();
                    }}
                  >
                    <Image
                      src="/icons/profile-icon.svg"
                      alt="My profile"
                      width={16}
                      height={16}
                    />
                    My profile
                  </Link>
                  <button
                    type="button"
                    className={styles.headerLogoutButton}
                    onClick={() => void handleLogout()}
                  >
                    <Image
                      src="/icons/logout.svg"
                      alt="Log out"
                      width={16}
                      height={16}
                    />
                    Log out
                  </button>
                </section>
              )}
            </div>
          </>
        ) : (
          <>
            <Link
              href="/Signin"
              className={styles.headerNavLink}
              onClick={closeMobileMenu}
            >
              Signin
            </Link>
            <Link
              href="/Signup"
              className={styles.headerNavLink}
              onClick={closeMobileMenu}
            >
              Signup
            </Link>
          </>
        )}
      </nav>
    </header>
  );
}
