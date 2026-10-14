"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { getErrorMessage } from "@/lib/api";
import { useAuthSession } from "../../providers/AuthSessionProvider.jsx";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import styles from "@/app/styles/auth.module.css";
const initialRegistration = {
  first_name: "",
  last_name: "",
  nick_name: "",
  date_of_birth: "",
  email: "",
  password: "",
  confirm_password: "",
  gender: "male",
  about_me: "",
  avatar_url: null,
};
export default function CreateAccountRoute() {
  const { Signup } = useAuthSession();
  const { showToast } = useToastMessages();
  const [registration, setRegistration] = useState(initialRegistration);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  useEffect(() => {
    if (!avatarPreview) return;
    return () => URL.revokeObjectURL(avatarPreview);
  }, [avatarPreview]);
  function haandelsubmit(event) {
    const target = event.target;
    if (target instanceof HTMLInputElement && target.type === "file") {
      const avatar = target.files?.[0] ?? null;
      setRegistration((current) => ({ ...current, avatar }));
      setAvatarPreview(avatar ? URL.createObjectURL(avatar) : null);
      return;
    }
    const fieldName = target.name;
    const nextValue = target.value;
    setRegistration((current) => ({ ...current, [fieldName]: nextValue }));
  }
  function removeAvatar() {
    setRegistration((current) => ({ ...current, avatar: null }));
    setAvatarPreview(null);
    const input = document.getElementById("avatar");
    if (input instanceof HTMLInputElement) input.value = "";
  }
  async function Submithanler(event) {
    // //console.log("registration",registration)
    event.preventDefault();
    setFormError("");
    if (registration.password !== registration.confirm_password) {
      setFormError("Passwords do not match.");
      return;
    }
    setIsSubmitting(true);
    try {
      await Signup(registration);
      showToast("Your account was created successfully.", "success");
    } catch (error) {
      showToast(getErrorMessage(error, "Registration failed."), "error");
    } finally {
      setIsSubmitting(false);
    }
  }
  return (
    <section className={styles.authSection}>
      <div className={`${styles.authCard} ${styles.authWideCard}`}>
        <h1 className={styles.authTitle}>Create your account</h1>
        <form className={styles.authForm} onSubmit={Submithanler}>
          <div className={styles.authFormRow}>
            <div className={styles.register}>
              <label className={styles.authLabel} htmlFor="firstname">
                First name
              </label>
              <input
                className={styles.authInput}
                id="firstname"
                name="first_name"
                autoComplete="given-name"
                value={registration.first_name}
                onChange={haandelsubmit}
                required
              />
            </div>
            <div className={styles.register}>
              <label className={styles.authLabel} htmlFor="last_name">
                Last name
              </label>
              <input
                className={styles.authInput}
                id="lastname"
                name="last_name"
                autoComplete="family-name"
                value={registration.last_name}
                onChange={haandelsubmit}
                required
              />
            </div>
          </div>

          <div className={styles.authFormRow}>
            <div className={styles.register}>
              <label className={styles.authLabel} htmlFor="nick_name">
                Nickname <span className={styles.authOptional}>(optional)</span>
              </label>
              <input
                className={styles.authInput}
                id="nickname"
                name="nick_name"
                value={registration.nick_name}
                onChange={haandelsubmit}
              />
            </div>
            <div className={styles.register}>
              <label className={styles.authLabel} htmlFor="date_of_birth">
                Date of birth
              </label>
              <input
                className={styles.authInput}
                id="date"
                name="date_of_birth"
                type="date"
                value={registration.date_of_birth}
                onChange={haandelsubmit}
                required
              />
            </div>
          </div>

          <div className={styles.register}>
            <label className={styles.authLabel} htmlFor="email">
              Email
            </label>
            <input
              className={styles.authInput}
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              value={registration.email}
              onChange={haandelsubmit}
              required
            />
          </div>

          <div className={styles.authFormRow}>
            <div className={styles.register}>
              <label className={styles.authLabel} htmlFor="password">
                Password
              </label>
              <input
                className={styles.authInput}
                id="password"
                name="password"
                type="password"
                autoComplete="new-password"
                value={registration.password}
                onChange={haandelsubmit}
                required
              />
            </div>
            <div className={styles.register}>
              <label className={styles.authLabel} htmlFor="confirm-password">
                Confirm password
              </label>
              <input
                className={styles.authInput}
                id="confirm-password"
                name="confirm_password"
                type="password"
                autoComplete="new-password"
                value={registration.confirm_password}
                onChange={haandelsubmit}
                required
              />
            </div>
          </div>

          <div className={styles.register}>
            <label className={styles.authLabel} htmlFor="gender">
              Gender
            </label>
            <select
              className={styles.authSelect}
              id="gender"
              name="gender"
              value={registration.gender}
              onChange={haandelsubmit}
            >
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>
          </div>

          <div className={styles.register}>
            <label className={styles.authLabel} htmlFor="about_me">
              About <span className={styles.authOptional}>(optional)</span>
            </label>
            <input
              className={styles.authInput}
              id="about"
              name="about_me"
              value={registration.about_me}
              onChange={haandelsubmit}
              maxLength={1000}
            />
          </div>

          <div className={styles.register}>
            <label className={styles.authLabel} htmlFor="avatar">
              Avatar <span className={styles.authOptional}>(optional)</span>
            </label>
            <input
              className={styles.authFileInput}
              id="avatar"
              name="avatar"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={haandelsubmit}
            />
            <label className={styles.authFileLabel} htmlFor="avatar">
              Choose an avatar
            </label>
            {avatarPreview && (
              <div className={styles.authAvatarPreview}>
                <Image
                  src={avatarPreview}
                  alt="Selected avatar preview"
                  fill
                  sizes="104px"
                  unoptimized
                />
                <button
                  type="button"
                  className={styles.authRemoveAvatar}
                  onClick={removeAvatar}
                  aria-label="Remove selected avatar"
                >
                  ×
                </button>
              </div>
            )}
          </div>

          {formError && (
            <p className={styles.authErrorMessage} role="alert">
              {formError}
            </p>
          )}

          <button
            type="submit"
            className={styles.authSubmitButton}
            disabled={isSubmitting}
          >
            {isSubmitting ? "Creating account…" : "Create account"}
          </button>
        </form>

        <p className={styles.authFooter}>
          Already registered? <Link href="/Signin">Sign in</Link>
        </p>
      </div>
    </section>
  );
}
