"use client";
import Link from "next/link";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/status-page.module.css";
export default function RouteErrorBoundary({ reset }) {
  return (
    <main className={styles.statusPage}>
      <section className={styles.statusCard}>
        <div className={styles.statusIcon} aria-hidden="true">
          😮‍💨
        </div>
        <h1>There is something  wrong</h1>
        <p>
          An Unexpected Error . Retry Again or return to the feed.
        </p>
        <div className={styles.statusActions}>
          <button
            type="button"
            className={uiStyles.uiPrimaryButton}
            onClick={reset}
          >
            Try again
          </button>
          <Link href="/" className={uiStyles.uiSecondaryButton}>
            Back to feed
          </Link>
        </div>
      </section>
    </main>
  );
}
