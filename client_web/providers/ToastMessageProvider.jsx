"use client";
import { createContext, useContext, useState } from "react";
import styles from "@/app/styles/toast.module.css";
const ToastMessageContext = createContext(undefined);
const toastIcons = {
  success: "✓",
  error: "!",
  info: "i",
};
export function ToastMessageProvider({ children }) {
  const [toasts, setToasts] = useState([]);
  function showToast(message, type = "info") {
    const id = crypto.randomUUID();
    setToasts((currentToasts) => [...currentToasts, { id, message, type }]);
    window.setTimeout(() => {
      setToasts((currentToasts) =>
        currentToasts.filter((toast) => toast.id !== id),
      );
    }, 3500);
  }
  return (
    <ToastMessageContext.Provider value={{ showToast }}>
      {children}
      <div className={styles.toastRegion} aria-live="polite" aria-atomic="true">
        {toasts.map((toast) => (
          <div
            key={toast.id}
            className={`${styles.toast} ${styles[toast.type]}`}
          >
            <span className={styles.toastIcon} aria-hidden="true">
              {toastIcons[toast.type]}
            </span>
            <span>{toast.message}</span>
          </div>
        ))}
      </div>
    </ToastMessageContext.Provider>
  );
}
export function useToastMessages() {
  const context = useContext(ToastMessageContext);
  if (!context) {
    throw new Error(
      "useToastMessages must be used inside ToastMessageProvider.",
    );
  }
  return context;
}
