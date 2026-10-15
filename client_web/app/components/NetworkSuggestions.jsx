import Image from "next/image";
import Link from "next/link";
import styles from "@/app/styles/navigation-sidebar.module.css";
import shellStyles from "@/app/styles/application-shell.module.css";
export function NetworkSuggestions() {
  return (
    <aside
      className={`${styles.sidebar} ${shellStyles.shellNetworkSidebar}`}
      aria-label="Profile and network suggestions"
    >
      <nav className={styles.sidebarNavigation} aria-label="Profile and network">
        <Link href="/Newpost" className={styles.sidebarNavigationLink}>
          <Image src="/icons/create-post.svg" alt="Create post" width={18} height={18} />
          <span>Create post</span>
        </Link>
        <Link href="/UserProfile" className={styles.sidebarNavigationLink}>
          <Image src="/icons/profile-icon.svg" alt="My profile" width={18} height={18} />
          <span>My profile</span>
        </Link>
        <Link href="/followers" className={styles.sidebarNavigationLink}>
          <Image src="/icons/groupsicon.svg" alt="Groups" width={18} height={18} />
          <span>See your network</span>
        </Link>
      </nav>
    </aside>
  );
}
