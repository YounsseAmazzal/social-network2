"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import {
  DEFAULT_AVATAR_PATH,
  fetchfromApi,
  getErrorMessage,
  isRecord,
  readBoolean,
  readString,
  buildBackendAssetUrl,
} from "@/lib/api";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/people-directory.module.css";
const networkTabs = [
  { key: "All", label: "All" },
  { key: "Userfollowers", label: "followers" },
  { key: "Userfollowing", label: "following" },
];
function normalizeRelationship(value) {
  return value === "following" ||
    value === "follower" ||
    value === "mutual" ||
    value === "requested"
    ? value
    : "none";
}
function normalizeNetworkUsers(value) {
  if (!Array.isArray(value)) return [];
  return value.flatMap((item) => {
    if (!isRecord(item)) return [];
    const id = readString(item, ["id", "ID"]);
    if (!id) return [];
    const firstName = readString(item, [
      "firstname",
      "first_name",
      "Firstname",
    ]);
    const lastName = readString(item, ["lastname", "last_name", "Lastname"]);
    const nickname = readString(item, ["nickname", "Nickname"]);
    return [
      {
        id,
        name: `${firstName} ${lastName}`.trim() || nickname || "Unknown user",
        username: nickname ? `@${nickname}` : "",
        avatarUrl: readString(item, ["avatar", "avatar_url", "avatar_URL"]),
        bio: readString(item, ["aboutMe", "about_me"]),
        isPrivate: readBoolean(item, ["isPrivate", "is_private"]),
        relationship: normalizeRelationship(item.Rolation),
      },
    ];
  });
}
async function fetchNetworkUsers(tab) {
  const response = await fetchfromApi(`/FetchNetwork?connect=${tab}`);
  return normalizeNetworkUsers(response.Contactdata);
}
function resolveNextRelationship(relationship, isPrivate) {
  if (relationship === "none") return isPrivate ? "requested" : "following";
  if (
    relationship === "following" ||
    relationship === "mutual" ||
    relationship === "requested"
  ) {
    return "none";
  }
  if (relationship === "follower") return "mutual";
  return "none";
}
export default function PeopleDirectoryRoute() {
  const { showToast } = useToastMessages();
  const [activeTab, setActiveTab] = useState("All");
  const [searchTerm, setSearchTerm] = useState("");
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [processingUserId, setProcessingUserId] = useState(null);
  useEffect(() => {
    let isCurrentRequest = true;
    setIsLoading(true);
    void fetchNetworkUsers(activeTab)
      .then((networkUsers) => {
        if (isCurrentRequest) setUsers(networkUsers);
      })
      .catch((error) => {
        if (!isCurrentRequest) return;
        setUsers([]);
        showToast(getErrorMessage(error, "Failed to load users."), "error");
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [activeTab]);
  async function followinghandleer(user) {
    if (processingUserId) return;
    setProcessingUserId(user.id);
    const previousRelationship = user.relationship;
    //console.log("the Rolation of users",previousRelationship)
    const nextRelationship = resolveNextRelationship(
      previousRelationship,
      user.isPrivate,
    );
    const isRemovingRelationship =
      previousRelationship === "following" ||
      previousRelationship === "mutual" ||
      previousRelationship === "requested";
    setUsers((currentUsers) =>
      currentUsers.map((currentUser) =>
        currentUser.id === user.id
          ? { ...currentUser, relationship: nextRelationship }
          : currentUser,
      ),
    );
    try {
  const endpoint = isRemovingRelationship ? "UserUnfollow" : "Userfollowinng";
  const idKey = isRemovingRelationship ? "UserUnfollow_id" : "Userfollowinng_id";
  const response = await fetchfromApi(
    `/${endpoint}?${idKey}=${user.id}`, 
    { method: isRemovingRelationship ? "DELETE" : "POST" }
  );
      const message = isRemovingRelationship
        ? "Relationship removed."
        : response.isPending
          ? "Follow request sent."
          : "Now following.";
      showToast(message, "success");
      setUsers(await fetchNetworkUsers(activeTab));
    } catch (error) {
      setUsers((currentUsers) =>
        currentUsers.map((currentUser) =>
          currentUser.id === user.id
            ? { ...currentUser, relationship: previousRelationship }
            : currentUser,
        ),
      );
      showToast(getErrorMessage(error, "Action failed."), "error");
    } finally {
      setProcessingUserId(null);
    }
  }
  const normalizedSearch = searchTerm.trim().toLowerCase();
  const visibleUsers = normalizedSearch
    ? users.filter((user) =>
        [user.name, user.username, user.bio].some((value) =>
          value.toLowerCase().includes(normalizedSearch),
        ),
      )
    : users;

  // console.log("the visible useeeeers ",visibleUsers)
  return (
    <main className={uiStyles.uiPage}>
      <header className={styles.peopleHeader}>
        <div>
          <h1>Connect with people</h1>
          <p>Find friends and manage your follow relationships.</p>
        </div>
        <label className="srOnly" htmlFor="network-search">
          Search users
        </label>
        <input
          id="network-search"
          className={styles.peopleSearchInput}
          type="search"
          value={searchTerm}
          onChange={(event) => setSearchTerm(event.target.value)}
          placeholder="Search by name or username"
        />
      </header>

      <div className={styles.peopleTabs} role="tablist" aria-label="Network lists">
        {networkTabs.map((tab) => (
          <button
            type="button"
            role="tab"
            key={tab.key}
            aria-selected={activeTab === tab.key}
            className={activeTab === tab.key ? styles.peopleActiveTab : ""}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {isLoading ? (
        <p className={uiStyles.uiEmptyState} role="status">
          Loading users…
        </p>
      ) : visibleUsers.length > 0 ? (
        <div className={styles.peopleUserList}>
          {visibleUsers.map((user) => {
            const actionLabel =
              user.relationship === "none"
                ? user.isPrivate
                  ? "Request"
                  : "Follow"
                : user.relationship === "requested"
                  ? "Cancel request"
                  : user.relationship === "follower"
                    ? "Follow back"
                    : "Unfollow";
            return (
              <article key={user.id} className={styles.peopleUserCard}>
                <Link href={`/UserProfile?user_id=${user.id}`} className={styles.peopleUserLink}>
                  <img
                    src={buildBackendAssetUrl(
                      user.avatar_Url,
                      DEFAULT_AVATAR_PATH,
                    )}
                    alt=""
                    width={56}
                    height={56}
                    className={styles.peopleAvatar}
                  />
                  <span className={styles.peopleUserDetails}>
                    <strong>
                      {user.name}
                      {user.isPrivate && (
                        <span title="Private account"> 🔐</span>
                      )}
                    </strong>
                    <small>{user.username}</small>
                    {user.bio && <span>{user.bio}</span>}
                  </span>
                </Link>
                <button
                  type="button"
                  className={`${styles.peopleActionButton} ${styles[user.relationship]}`}
                  disabled={processingUserId === user.id}
                  onClick={() => void followinghandleer(user)}
                >
                  {processingUserId === user.id ? "Working…" : actionLabel}
                </button>
              </article>
            );
          })}
        </div>
      ) : (
        <p className={uiStyles.uiEmptyState}>
          {normalizedSearch ? "There is noo ussers match your search." : "No users found."}
        </p>
      )}
    </main>
  );
}
