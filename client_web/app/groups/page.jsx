"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { getErrorMessage } from "@/lib/api";
import {
  createGroup,
  createGroupJoinRequest,
  fetchAllJoinnedGroups,
  fetchAllgroupsSuggested,
} from "@/lib/groups";
import { useUnreadConversations } from "@/providers/UnreadConversationProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/group-directory.module.css";
export default function GroupDirectoryRoute() {
  const { showToast } = useToastMessages();
  const { unreadByGroup, totalGroupUnread } = useUnreadConversations();
  const [activeTab, setActiveTab] = useState("discover");
  const [isCreateFormOpen, setIsCreateFormOpen] = useState(false);
  const [joinedGroups, setJoinedGroups] = useState([]);
  const [suggestedGroups, setSuggestedGroups] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [joiningGroupId, setJoiningGroupId] = useState(null);
  const [formValues, setFormValues] = useState({ title: "", description: "" });
  useEffect(() => {
    let isCurrentRequest = true;
    setIsLoading(true);
    void Promise.all([fetchAllJoinnedGroups(), fetchAllgroupsSuggested()])
        // //console.log("The suggeested suggeested suggeested suggeested",rs)
    .then(([joined, suggested]) => {
      if (!isCurrentRequest) return;
      setJoinedGroups(joined);
        //console.log("The joined joined joined joined",joined)
        //console.log("The suggeested suggeested suggeested suggeested",suggested)
        setSuggestedGroups(suggested);
      })
      .catch((error) => {
        if (isCurrentRequest) {
          showToast(getErrorMessage(error, "Failed to load groups."), "error");
        }
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, []);
  async function handleCreateGroup(event) {
    event.preventDefault();
    if (!formValues.title.trim() || !formValues.description.trim()) {
      showToast("Group title and description are required.", "error");
      return;
    }
    setIsCreating(true);
    try {
      const newGroup = await createGroup({
        group_title: formValues.title.trim(),
        description: formValues.description.trim(),
      });

            const formattedGroup = {
        id: newGroup.id,
        title: newGroup.title || formValues.title,
        description: newGroup.description || formValues.description,
      };
      // //console.log("HHHHHHHHHHHHHHHHHHHH",newGroup.id)
      setJoinedGroups((currentGroups) => [formattedGroup, ...currentGroups]);
      setFormValues({ title: "", description: "" });
      setIsCreateFormOpen(false);
      setActiveTab("joined");

      showToast("Group created successfully.", "success");
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to create the group."), "error");
    } finally {
      setIsCreating(false);
    }
  }
  async function JoinGrouprequestHandleer(group) {
    setJoiningGroupId(group.id);
    try {
      await createGroupJoinRequest(group.id);
      setSuggestedGroups((currentGroups) =>
        currentGroups.map((currentGroup) =>
          currentGroup.id === group.id
            ? { ...currentGroup, requestId: 1 }
            : currentGroup,
        ),
      );
      showToast("Join request sent.", "success");
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to send the request."), "error");
    } finally {
      setJoiningGroupId(null);
    }
  }
  const visibleGroups = activeTab === "discover" ? suggestedGroups : joinedGroups;
    //console.log("visibleGroupsvisibleGroupsvisibleGroupsvisibleGroups0",joinedGroups)
  return (
    <main className={uiStyles.uiPage}>
      <header className={styles.groupHeader}>
        <div>
          <h1>Groups</h1>
          <p>Find a community or create one of your own.</p>
        </div>
        <button
          type="button"
          className={uiStyles.uiPrimaryButton}
          aria-expanded={isCreateFormOpen}
          onClick={() => setIsCreateFormOpen((isOpen) => !isOpen)}
        >
          {isCreateFormOpen ? "Close" : "Create group"}
        </button>
      </header>

      {isCreateFormOpen && (
        <form className={styles.groupCreateForm} onSubmit={handleCreateGroup}>
          <h2>Create a new group</h2>
          <label htmlFor="group-title">Title</label>
          <input
            id="group-title"
            value={formValues.title}
            maxLength={255}
            onChange={(event) =>
              setFormValues((current) => ({
                ...current,
                title: event.target.value,
              }))
            }
            required
          />
          <label htmlFor="group-description">Description</label>
          <textarea
            id="group-description"
            value={formValues.description}
            maxLength={1000}
            rows={4}
            onChange={(event) =>
              setFormValues((current) => ({
                ...current,
                description: event.target.value,
              }))
            }
            required
          />
          <button
            type="submit"
            className={uiStyles.uiPrimaryButton}
            disabled={isCreating}
          >
            {isCreating ? "Creating…" : "Create group"}
          </button>
        </form>
      )}

      <div className={styles.groupTabs} role="tablist" aria-label="Group lists">
        <button
          type="button"
          role="tab"
          aria-selected={activeTab === "discover"}
          className={activeTab === "discover" ? styles.groupActiveTab : ""}
          onClick={() => setActiveTab("discover")}
        >
          Discover
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={activeTab === "joined"}
          className={activeTab === "joined" ? styles.groupActiveTab : ""}
          onClick={() => setActiveTab("joined")}
        >
          My groups
          {totalGroupUnread > 0 && (
            <span className={styles.groupTabBadge}>{totalGroupUnread}</span>
          )}
        </button>
      </div>

      {isLoading ? (
        <p className={uiStyles.uiEmptyState} role="status">
          Loading groups…
        </p>
      ) : visibleGroups.length > 0 ? (
        <div className={styles.groupGrid}>
          {visibleGroups.map((group) => {
            const unreadCount = unreadByGroup[group.id] ?? 0;
            return (
              <article key={group.id||group.ID||""} className={styles.groupCard}>
                <div className={styles.groupTitleRow}>
                  <h2>{group.title}</h2>
                  {activeTab === "joined" && unreadCount > 0 && (
                    <span className={styles.groupUnreadBadge}>{unreadCount}</span>
                  )}
                </div>
                <p>{group.description || "No description yet."}</p>
  {activeTab === "joined" ? (
  <Link
    href={`/gouupDetails?Grouup_id=${group.id || group.ID || ""}`}
    className={uiStyles.uiSecondaryButton}
  >
    See Group
  </Link>
) : (group.request_status === "pending" || group.request_status === "followpending") ? (
  <button
    type="button"
    className={uiStyles.uiGhostButton}
    disabled
  >
    Request pending
  </button>
) : (
  <button
    type="button"
    className={uiStyles.uiSecondaryButton}
    disabled={joiningGroupId === (group.id || group.ID)}
    onClick={() => void JoinGrouprequestHandleer(group)}
  >
    {joiningGroupId === (group.id || group.ID)
      ? "Sending…"
      : "Request to join"}
  </button>
)}
              </article>
            );
          })}
        </div>
      ) : (
        <p className={uiStyles.uiEmptyState}>
          {activeTab === "discover"
            ? "No suggested groups are available."
            : "You have not joined any groups yet."}
        </p>
      )}
    </main>
  );
}
