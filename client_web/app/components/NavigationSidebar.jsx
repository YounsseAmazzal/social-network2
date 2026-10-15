"use client";
import Image from "next/image";
import Link from "next/link";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useUnreadConversations } from "@/providers/UnreadConversationProvider";
import styles from "@/app/styles/navigation-sidebar.module.css";
import shellStyles from "@/app/styles/application-shell.module.css";

const navigationItems = [
  { href: "/", label: "Feed", icon: "/icons/homepage.svg" },
  { href: "/groups", label: "Groups", icon: "/icons/groupsicon.svg" },
];

export default function NavigationSidebar() {
  const { user } = useAuthSession();
  const { totalUnread, totalGroupUnread } = useUnreadConversations();

  if (!user) return null;

  return (
    <aside
      className={`${styles.sidebar} ${shellStyles.shellPrimarySidebar}`}
      aria-label="Primary sidebar"
    >
      <nav className={styles.sidebarNavigation}>
        {navigationItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={styles.sidebarNavigationLink}
          >
            <Image src={item.icon} alt={item.label} width={18} height={18} />
            <span>{item.label}</span>
            {item.href === "/groups" && totalGroupUnread > 0 && (
              <span className={styles.sidebarNavigationBadge}>{totalGroupUnread}</span>
            )}
          </Link>
        ))}
        <Link href="/chat" className={styles.sidebarNavigationLink}>
          <Image src="/icons/chats.svg" alt="Chats" width={18} height={18} />
          <span>Chats</span>
          {totalUnread > 0 && (
            <span className={styles.sidebarNavigationBadge}>{totalUnread}</span>
          )}
        </Link>
      </nav>
    </aside>
  );
}
