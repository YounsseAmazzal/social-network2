"use client";
import { useState } from "react";
import Link from "next/link";
import { getErrorMessage } from "@/lib/api";
import { useAuthSession } from "../../providers/AuthSessionProvider.jsx";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import styles from "@/app/styles/auth.module.css";
export default function SignInRoute() {
  const { Signin } = useAuthSession();
  const { showToast } = useToastMessages();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  async function Submithanler(event) {
    event.preventDefault();
    //console.log("the identifier and password",identifier,password)
    setIsSubmitting(true);
    try {
      await Signin({ IdentifierEN: identifier, password });
      showToast("Welcome back!", "success");
    } catch (error) {
      showToast(getErrorMessage(error, "Login failed."), "error");
    } finally {
      setIsSubmitting(false);
    }
  }
  return (
    <section className={styles.authSection}>
      <div className={styles.authCard}>
        <h1 className={styles.authTitle}>Welcome back</h1>
        <form className={styles.authForm} onSubmit={Submithanler}>
          <div className={styles.register}>
            <label className={styles.authLabel} htmlFor="identifier">
              Email or nickname
            </label>
            <input
              className={styles.authInput}
              id="identifier"
              name="identifiant"
              autoComplete="username"
              value={identifier}
              onChange={(event) => setIdentifier(event.target.value)}
              required
            />
          </div>

          <div className={styles.register}>
            <label className={styles.authLabel} htmlFor="password">
              Password
            </label>
            <input
              className={styles.authInput}
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </div>

          <button
            type="submit"
            className={styles.authSubmitButton}
            disabled={isSubmitting}
          >
            {isSubmitting ? "Signing in…" : "Sign in"}
          </button>
        </form>

        <p className={styles.authFooter}>
          Need an account? <Link href="/Signup">Create one</Link>
        </p>
      </div>
    </section>
  );
}
