import Link from "next/link";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/status-page.module.css";
export default function NotFoundRoute() {
  return (
    <main className={styles.statusPage}>
      <section className={styles.statusCard}>
        <div className={styles.statusCode}>404</div>
        <h1>Page not found</h1>
        <p>The page you’re looking for doesn’t exist or may have been moved.</p>
        <Link href="/" className={uiStyles.uiPrimaryButton}>
          Back to feed
        </Link>
      </section>
    </main>
  );
}
