"use client";
import { usePathname } from "next/navigation";
import SiteHeader from "@/app/components/SiteHeader";
import NavigationSidebar from "@/app/components/NavigationSidebar";
import { NetworkSuggestions } from "@/app/components/NetworkSuggestions";
import styles from "@/app/styles/application-shell.module.css";
export default function ApplicationShell({ children }) {
  const pathname = usePathname();
  const isAuthenticationPage =
    pathname === "/Signin" || pathname === "/Signup";
  if (isAuthenticationPage) {
    return <main className={styles.shellAuthPage}>{children}</main>;
  }
  return (
    <div className={styles.appShell}>
      <SiteHeader />
      <div className={styles.shellContentGrid}>
        <NetworkSuggestions />
        <NavigationSidebar />
        <main className={styles.shellRouteContent}>{children}</main>
      </div>
    </div>
  );
}
