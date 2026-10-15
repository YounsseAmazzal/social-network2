"use client";
import { createContext, useContext,
   useEffect
   , useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { fetchfromApi, isRecord } from "@/lib/api";
import { normalizeAuthUser } from "@/lib/auth";
import styles from "@/app/styles/application-shell.module.css";
const AuthSessionContext = createContext(undefined);
function extractAuthUser(response) {
  if (isRecord(response?.data) && isRecord(response.data.user)) {
    return response.data.user;
  }
  return response.Userinfo ?? response;
}
function buildRegistrationPayload(registration) {
  // //console.log("the sened  form",registration)
  // //console.log("the 3aaaaaksss",!registration.avatar_url)
  if (registration.avatar_url!=null) {
    return JSON.stringify(registration);
  }
  //console.log("reeeeeeeeeee",registration)
  const formData = new FormData();
  Object.entries(registration).forEach(([key, value]) => {
    if (value === null || value === undefined) return;
    formData.append(key, value instanceof File ? value : String(value));
  });
  return formData;
}
export function AuthSessionProvider({ children }) {
  //console.log("the user want to signuuppp")
  const [user, setUser] = useState(null);
  const [isCheckingSession, setIsCheckingSession] = useState(true);
  const router = useRouter();
  const pathname = usePathname();
  useEffect(() => {
    let isCurrentRequest = true;
    async function checkSession() {
      try {
        const response = await fetchfromApi("/Me");
        console.log("api provider Me testing",response)
        if (!isCurrentRequest) return;
        setUser(normalizeAuthUser(extractAuthUser(response)));
      } catch {
        if (isCurrentRequest) setUser(null);
      } finally {
        if (isCurrentRequest) setIsCheckingSession(false);
      }
    }
    void checkSession();
    return () => {
      isCurrentRequest = false;
    };
  }, []);
  //console.log("useeeeeeer infooo",user)
  useEffect(() => {
    if (isCheckingSession) return;
    const isAuthenticationPage =
      pathname === "/Signin" || pathname === "/Signup";
    if (!user && !isAuthenticationPage) {
      router.replace("/Signin");
    } else if (user && isAuthenticationPage) {
      router.replace("/");
    }
  }, [isCheckingSession, pathname, router, user]);
  async function Signin(credentials) {
    //console.log("thaaat is credintial of user",credentials)
    const response = await fetchfromApi("/Signin", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
    //console.log("the response of Signin",response)
    setUser(normalizeAuthUser(extractAuthUser(response)));
    router.replace("/");
  }
  async function Signup(registration) {
    // //console.log("khaaaaaaaasni register",registration)
    const response = await fetchfromApi("/Signup", {
      method: "POST",
      body: buildRegistrationPayload(registration),
    });
    // console.log("the response of SignUp messaaage",response)
    setUser(normalizeAuthUser(extractAuthUser(response)));
    router.replace("/");
  }
  async function logout() {
    try {
      await fetchfromApi("/Logout", { method: "DELETE" });
    } catch {
      // Clear the local session even if the server is unavailable.
    } finally {
      setUser(null);
      router.replace("/Signin");
    }
  }
  function updateUser(updatedUser) {
    if (!isRecord(updatedUser)) return;
    setUser(updatedUser);
  }
  return (
    <AuthSessionContext.Provider
      value={{
        user,
        isCheckingSession,
        Signin,
        Signup,
        logout,
        updateUser,
      }}
    >
      {isCheckingSession ? (
        <div className={styles.shellLoadingScreen} role="status" aria-live="polite">
          <div className={styles.shellLoadingContent}>
            <span className={styles.shellSpinner} aria-hidden="true" />
            <span>Checking your session…</span>
          </div>
        </div>
      ) : (
        children
      )}
    </AuthSessionContext.Provider>
  );
}
export function useAuthSession() {
  const context = useContext(AuthSessionContext);
  if (!context) {
    throw new Error("useAuthSession must be used inside AuthSessionProvider.");
  }
  return context;
}
