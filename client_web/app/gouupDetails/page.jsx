"use client";
import { useEffect, useState, Suspense } from "react";
import Image from "next/image";
import Link from "next/link";
import { useParams, useSearchParams } from "next/navigation";
import {
  DEFAULT_AVATAR_PATH,
  getErrorMessage,
  buildBackendAssetUrl, API_URL
} from "../../lib/api.js";
import {
  groupEventCreation,
  PosstGroupcreate,
  fetchGroupAllInfo,
  fetchallEventsGroup,
  fetchallgroupmemberes,
  fetchGroupPending,
  fetchallpostsgroup,
  FetchFolllowersToinvite,
  respondToGroupRequest,
  InviteUsertogroup,
  participateOnEventVote,
} from "../../lib/groups.js";
import {
  formatGroupEventSchedule,
  describeRelativeTimestamp,
} from "@/lib/time";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useUnreadConversations } from "@/providers/UnreadConversationProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import GroupConversation from "@/app/groups/[id]/GroupConversation";
import uiStyles from "@/app/styles/ui.module.css";
import stylesmember from "@/app/styles/post-discussion.module.css";
import styles from "@/app/styles/group-workspace.module.css";
const groupTabs = [
  { id: "feed", label: "Feed" },
  { id: "events", label: "Events" },
  { id: "members", label: "Members" },
  { id: "chat", label: "Chat" },
];
function getGroupMemberDisplayName(member) {
  return (
    member.nickname ||
    `${member.firstName} ${member.lastName}`.trim() ||
    "Unknown member"
  );
}
function GroupWorkspaceContent() {
  // const params = useParams();
  const { showToast } = useToastMessages();
  const { user } = useAuthSession();
  const { unreadByGroup, markGroupRead } = useUnreadConversations();
  const searchParams = useSearchParams();
  const groupId = searchParams.get("Grouup_id") || "";
  const limit = searchParams.get("limit") || "10";
  const offset = searchParams.get("offset") || "0";

  const currentUserId = user?.id ?? "";
  const [activeTab, setActiveTab] = useState("feed");
  const [groupDetails, setGroupDetails] = useState(null);
  const [events, setEvents] = useState([]);
  const [members, setMembers] = useState([]);
  const [pendingMembers, setPendingMembers] = useState([]);
  const [posts, setPosts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [pageError, setPageError] = useState("");
  const [isPosting, setIsPosting] = useState(false);
  const [isCreatingEvent, setIsCreatingEvent] = useState(false);
  const [votingEventId, setVotingEventId] = useState(null);
  const [processingMemberId, setProcessingMemberId] = useState(null);
  const [isInviteDialogOpen, setIsInviteDialogOpen] = useState(false);
  const [invitableUsers, setInvitableUsers] = useState([]);
  const [isLoadingInvitableUsers, setIsLoadingInvitableUsers] = useState(false);
  const [sendingInviteUserId, setSendingInviteUserId] = useState(null);
  const [invitedUserIds, setInvitedUserIds] = useState(new Set());


  const [currentOffset, setCurrentOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [isMoreLoading, setIsMoreLoading] = useState(false);
  const Page_limiit = 10;


// event offset
const [eventsOffset, setEventsOffset] = useState(0);
const [hasMoreEvents, setHasMoreEvents] = useState(true);
const [isMoreEventsLoading, setIsMoreEventsLoading] = useState(false);
const Events_limit = 5;
  const [postForm, setPostForm] = useState({
    GroupPostTitle: "",
    GroupPostContent: "",
    Group_post_image: null,
  });
  const [eventForm, setEventForm] = useState({
    title: "",
    description: "",
    event_date: "",
  });
  const isOwner = Boolean(
    groupDetails?.creator?.id &&
    currentUserId &&
    currentUserId === groupDetails.creator.id
  );

  const availableTabs = groupTabs.filter(
    (tab) => tab.id !== "pending" || isOwner,
  );
  useEffect(() => {
    if (!groupId) return;
    let isCurrentRequest = true;
    async function loadGroup() {
      setIsLoading(true);
      setPageError("");
      try {
        const [details, memberList, postList] = await Promise.all([
          fetchGroupAllInfo(groupId),
          // fetchallEventsGroup(groupId,Events_limit,0).catch(() => []),
          fetchallgroupmemberes(groupId).catch(() => []),
          fetchallpostsgroup(groupId, limit, offset).catch(() => []),
        ]);
        if (!isCurrentRequest) return;
        // //console.log("the events lisst !!!!!!!!!!!!!!!!t",eventList)
        setGroupDetails(details);
        // setEvents(eventList);
        // //console.log("Theeeeeeeeeee members of grouup", memberList)
        // //console.log("Theeeeeeeeeee groupId of grouup",groupId)
        // //console.log("Theeeeeeeeeee limit of grouup",limit)
        // //console.log("Theeeeeeeeeee postList of grouup",postList)

        setEventsOffset(0);
      // if (eventList.length < Events_limit) setHasMoreEvents(false);
      //  else setHasMoreEvents(true);
        setMembers(memberList);
        setPosts(postList);
        setCurrentOffset(0);
        setHasMore(true);


        /// for event
             setEvents([]); 
        setEventsOffset(0);
        setHasMoreEvents(true);
      } catch (error) {
        if (isCurrentRequest) {
          setPageError(getErrorMessage(error, "Failed to load this group."));
        }
      } finally {
        if (isCurrentRequest) setIsLoading(false);
      }
    }
    void loadGroup();
    return () => {
      isCurrentRequest = false;
    };
  }, [groupId]);

  // when the user change the tab and the groupid 
useEffect(() => {
    if (activeTab === "events" && events.length === 0) {
        setIsMoreEventsLoading(true);
        fetchallEventsGroup(groupId, Events_limit, 0).then(data => {
            setEvents(data);
            setHasMoreEvents(data.length === Events_limit);
            setIsMoreEventsLoading(false);
        });
    }
}, [activeTab, groupId]); 
  useEffect(() => {
    const requestedTab = searchParams.get("tab");
    if (!requestedTab) return;
    const isAllowed =
      groupTabs.some((tab) => tab.id === requestedTab) ||
      (requestedTab === "pending" && isOwner);
    setActiveTab(isAllowed ? requestedTab : "feed");
  }, [isOwner, searchParams]);
  useEffect(() => {
    if (activeTab !== "chat" || !groupId || !unreadByGroup[groupId]) return;
    void markGroupRead(groupId);
  }, [activeTab, groupId, unreadByGroup]);
  useEffect(() => {
    if (activeTab !== "pending" || !isOwner || !groupId) return;
    let isCurrentRequest = true;
    void fetchGroupPending(groupId)
      .then((pending) => {
        if (isCurrentRequest) setPendingMembers(pending);
      })
      .catch((error) => {
        if (isCurrentRequest) {
          showToast(
            getErrorMessage(error, "Failed to load pending requests."),
            "error",
          );
        }
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [activeTab, groupId, isOwner]);
  useEffect(() => {
    if (!isInviteDialogOpen) return;
    function closeOnEscape(event) {
      if (event.key === "Escape") setIsInviteDialogOpen(false);
    }
    document.addEventListener("keydown", closeOnEscape);
    return () => document.removeEventListener("keydown", closeOnEscape);
  }, [isInviteDialogOpen]);
  async function handleCreatePost(event) {
    event.preventDefault();
    if (!groupId) return;
    if (!postForm.GroupPostTitle.trim() || !postForm.GroupPostContent.trim()) {
      showToast("Post title and content are required.", "error");
      return;
    }
    setIsPosting(true);
    try {
      const formData = new FormData();
      formData.append("GroupPostTitle", postForm.GroupPostTitle.trim());
      formData.append("GroupPostContent", postForm.GroupPostContent.trim());
      if (postForm.Group_post_image) formData.append("Group_post_image", postForm.Group_post_image);
      await PosstGroupcreate(groupId, formData);
      setPostForm({ GroupPostTitle: "", GroupPostContent: "", Group_post_image: null });
      setPosts(await fetchallpostsgroup(groupId, limit, offset));
      showToast("Group post created.", "success");
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to create the post."), "error");
    } finally {
      setIsPosting(false);
    }
  }
  async function handleSeeMore() {
    const nextOffset = currentOffset + Page_limiit;
    setIsMoreLoading(true);

    try {
      const newPosts = await fetchallpostsgroup(groupId, Page_limiit, nextOffset);
      if (newPosts.length === 0) {
        setHasMore(false);
        showToast("No more posts to show.", "info");
      } else {
        // combine between the old and the neww posts 
        setPosts((prevPosts) => [...prevPosts, ...newPosts]);
        setCurrentOffset(nextOffset);
        // if he return a valuses changed then the usall       
        if (newPosts.length < Page_limiit) setHasMore(false);
      }
    } catch (error) {
      showToast("Failed to load more posts.", "error");
    } finally {
      setIsMoreLoading(false);
    }
  }
  async function handleCreateEvent(event) {
    event.preventDefault();
    if (!groupId) return;
    const eventDate = new Date(eventForm.event_date);
    if (
      !eventForm.title.trim() ||
      !eventForm.description.trim() ||
      Number.isNaN(eventDate.getTime())
    ) {
      showToast("Fill in a valid event title, description, and date.", "error");
      return;
    }
    setIsCreatingEvent(true);
    try {
      await groupEventCreation(groupId, {
        title: eventForm.title.trim(),
        description: eventForm.description.trim(),
        event_date: eventDate.toISOString(),
      });
      setEventForm({ title: "", description: "", event_date: "" });
        const refreshedEvents = await fetchallEventsGroup(groupId, Events_limit, 0);
        setEvents(refreshedEvents);
        setEventsOffset(0); 
        setHasMoreEvents(refreshedEvents.length === Events_limit);
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to create the event."), "error");
    } finally {
      setIsCreatingEvent(false);
    }
  }
async function handleVote(groupEvent, vote) {
    if (!groupId) return;
    const nextVote = groupEvent.userVote === vote ? "" : vote;
    
    setVotingEventId(groupEvent.id);
    try {
        await participateOnEventVote(groupId, groupEvent.id, nextVote);
        
        const totalToFetch = Number(eventsOffset) + Number(Events_limit);
        
        const refreshedEvents = await fetchallEventsGroup(groupId, totalToFetch, 0);
        setEvents(refreshedEvents);

        showToast(nextVote ? "Participation registered!" : "Participation removed.", "success");
    } catch (error) {
        showToast("Failed to update your vote.", "error");
    } finally {
        setVotingEventId(null);
    }
}
  async function handleMembershipDecision(memberId, decision) {
    if (!groupId) return;
    setProcessingMemberId(memberId);
    try {
      await respondToGroupRequest(groupId, memberId, decision);
      setPendingMembers((currentMembers) =>
        currentMembers.filter((member) => member.id !== memberId),
      );
      if (decision === "accept") {
        setMembers(await fetchallgroupmemberes(groupId));
      }
      showToast(
        decision === "accept" ? "Member accepted." : "Request rejected.",
        "success",
      );
    } catch (error) {
      showToast(
        getErrorMessage(error, "Failed to update the request."),
        "error",
      );
    } finally {
      setProcessingMemberId(null);
    }
  }
  async function openInviteDialog() {
    if (!groupId) return;
    setIsInviteDialogOpen(true);
    setIsLoadingInvitableUsers(true);
    try {
      setInvitableUsers(await FetchFolllowersToinvite(groupId));
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to load followers."), "error");
    } finally {
      setIsLoadingInvitableUsers(false);
    }
  }
  async function sendInvitation(userId) {
    if (!groupId) return;
    setSendingInviteUserId(userId);
    try {
      await InviteUsertogroup(groupId, userId);
      setInvitedUserIds((currentIds) => new Set(currentIds).add(userId));
      showToast("Invitation sent.", "success");
    } catch (error) {
      showToast(
        getErrorMessage(error, "Failed to send the invitation."),
        "error",
      );
    } finally {
      setSendingInviteUserId(null);
    }
  }
  function stopDialogPropagation(event) {
    event.stopPropagation();
  }
  if (isLoading) {
    return (
      <p className={uiStyles.uiEmptyState} role="status">
        Loading group…
      </p>
    );
  }
  async function handleSeeMoreEvents() {
    const nextOffset = Number(eventsOffset) + Number(Events_limit);
  setIsMoreEventsLoading(true);

  try {
   const newEvents = await fetchallEventsGroup(groupId, Events_limit, nextOffset);
        
        if (newEvents.length === 0) {
            setHasMoreEvents(false);
        } else {
            setEvents((prev) => {
                const existingIds = new Set(prev.map(e => e.id));
                const uniqueNew = newEvents.filter(e => !existingIds.has(e.id));
                return [...prev, ...uniqueNew];
            });
            
            setEventsOffset(nextOffset);
            if (newEvents.length < Events_limit) setHasMoreEvents(false);
        }
  } catch (error) {
    showToast("Failed to load more events.", "error");
  } finally {
    setIsMoreEventsLoading(false);
  }
}
  if (!groupDetails) {
    return (
      <section className={uiStyles.uiEmptyState}>
        <h1>Group unavailable</h1>
        <p>{pageError || "This group could not be opened."}</p>
        <Link href="/groups">Back to groups</Link>
      </section>
    );
  }
  const unreadCount = groupId ? (unreadByGroup[groupId] ?? 0) : 0;
  // //console.log("THe meeeemeerrrrs ",members.length)
  return (
    <main className={uiStyles.uiPage}>
      <header className={styles.groupHeader}>
        <div className={styles.groupCover} aria-hidden="true" />
        <div className={styles.groupHeaderContent}>
          <h1>{groupDetails.group.title}</h1>
          <p>{groupDetails.group.description || "No group description yet."}</p>
          <div className={styles.groupMetadata}>
            <span>{groupDetails.totalMembers} members</span>
            <span>
              Created {describeRelativeTimestamp(groupDetails.group.createdAt)}
            </span>
            <span>
              Owner:{" "}
              <Link href={`/UserProfile/${groupDetails.creator.id}`}>
                {getGroupMemberDisplayName(groupDetails.creator)}
              </Link>
            </span>
            { (
              <button
                type="button"
                className={uiStyles.uiSecondaryButton}
                onClick={() => void openInviteDialog()}
              >
                Invite members
              </button>
            )}
          </div>
        </div>
      </header>

      <div className={styles.groupTabs} role="tablist" aria-label="Group sections">
        {availableTabs.map((tab) => (
          <button
            type="button"
            role="tab"
            aria-selected={activeTab === tab.id}
            key={tab.id}
            className={activeTab === tab.id ? styles.groupActiveTab : ""}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
            {tab.id === "chat" && unreadCount > 0 && activeTab !== "chat" && (
              <span className={uiStyles.uiBadge}>{unreadCount}</span>
            )}
          </button>
        ))}
        {isOwner && (
          <button
            type="button"
            role="tab"
            aria-selected={activeTab === "pending"}
            className={activeTab === "pending" ? styles.groupActiveTab : ""}
            onClick={() => setActiveTab("pending")}
          >
            Requests
          </button>
        )}
      </div>

      {activeTab === "feed" && (
        <div className={styles.groupSectionStack}>
          <form className={styles.groupEditorCard} onSubmit={handleCreatePost}>
            <h2>Create a group post</h2>
            <label className="srOnly" htmlFor="GroupPostTitle">
              Post title
            </label>
            <input
              id="GroupPostTitle"
              value={postForm.GroupPostTitle}
              maxLength={255}
              placeholder="Post title"
              onChange={(event) =>
                setPostForm((current) => ({
                  ...current,
                  GroupPostTitle: event.target.value,
                }))
              }
              required
            />
            <label className="srOnly" htmlFor="GroupPostContent">
              Post content
            </label>
            <textarea
              id="GroupPostContent"
              value={postForm.GroupPostContent}
              maxLength={1000}
              placeholder="What do you want to share?"
              onChange={(event) =>
                setPostForm((current) => ({
                  ...current,
                  GroupPostContent: event.target.value,
                }))
              }
              required
            />
            <label className={styles.groupFileLabel}>
              Add image (optional)
              <input
                type="file"
                accept="image/jpeg,image/png,image/gif"
                onChange={(event) =>
                  setPostForm((current) => ({
                    ...current,
                    Group_post_image: event.target.files?.[0] ?? null,
                  }))
                }
              />
            </label>
            <button
              type="submit"
              className={uiStyles.uiPrimaryButton}
              disabled={isPosting}
            >
              {isPosting ? "Publishing…" : "Publish post"}
            </button>
          </form>

          {posts.length > 0 ? (
            <>

              {posts.map((post) => (
                <article key={post.id} className={styles.groupPostCard}>
                  <div className={styles.groupPostAuthor}>
                    <img
                      src={buildBackendAssetUrl(
                        post.creater.avatarUrl,
                        DEFAULT_AVATAR_PATH,
                      )}
                      alt=""
                      width={40}
                      height={40}
                      className={styles.groupAvatar}
                    />
                    <div>
                      <strong>{getGroupMemberDisplayName(post.creater)}</strong>
                      <small>{(post.createDate)}</small>
                    </div>
                  </div>
                  <h2>{post.title || "Untitled post"}</h2>
                  <p className={styles.groupPostContent}>{post.description}</p>
                  {post.image && (
                    <img
                      src={buildBackendAssetUrl(post.image||DEFAULT_AVATAR_PATH)}
                      alt=""
                      width={800}
                      height={420}
                      className={styles.groupPostImage}
                    />
                  )}
                  <div className={styles.groupPostFooter}>
                    <span>
                      <Image
                        src="/icons/like.svg"
                        alt="Likes"
                        width={16}
                        height={16}
                      />{" "}
                      {post.likes}
                    </span>
                    <span>
                      <Image
                        src="/icons/comments.svg"
                        alt="Comments"
                        width={16}
                        height={16}
                      />{" "}
                      {post.totalComments}
                    </span>
                    <Link href={`/GetPost?Possst_id=${post.id}`}>See Group Post</Link>
                  </div>
                </article>))}
              {hasMore && (
                <button
                  type="button"
                  className={uiStyles.uiSecondaryButton}
                  style={{ width: '100%', marginTop: '1rem' }}
                  disabled={isMoreLoading}
                  onClick={handleSeeMore}
                >
                  {isMoreLoading ? "Loading..." : "See More Posts"}
                </button>
              )}

              {!hasMore && posts.length > 0 && (
                <p className={styles.emptyState}>You've reached the end of the feed.</p>
              )}
            </>
          ) : (
            <p className={uiStyles.uiEmptyState}>No group posts yet.</p>
          )}
        </div>
      )}
     {activeTab === "events" && (
        <div className={styles.groupSectionStack}>
          <form className={styles.groupEditorCard} onSubmit={handleCreateEvent}>
            <h2>Create an event</h2>
            <label className="srOnly" htmlFor="event-title">
              Event title
            </label>
            <input
              id="event-title"
              value={eventForm.title}
              maxLength={255}
              placeholder="Event title"
              onChange={(event) =>
                setEventForm((current) => ({
                  ...current,
                  title: event.target.value,
                }))
              }
              required
            />
            <label className="srOnly" htmlFor="event-description">
              Event description
            </label>
            <textarea
              id="event-description"
              value={eventForm.description}
              maxLength={1000}
              placeholder="Event description"
              onChange={(event) =>
                setEventForm((current) => ({
                  ...current,
                  description: event.target.value,
                }))
              }
              required
            />
            <label className="srOnly" htmlFor="event-date">
              Event date
            </label>
            <input
              id="event-date"
              type="datetime-local"
              value={eventForm.event_date}
              onChange={(event) =>
                setEventForm((current) => ({
                  ...current,
                  event_date: event.target.value,
                }))
              }
              required
            />
            <button
              type="submit"
              className={uiStyles.uiPrimaryButton}
              disabled={isCreatingEvent}
            >
              {isCreatingEvent ? "Creating…" : "Create event"}
            </button>
          </form>

          {events.length > 0 ? (
            <>
            {events.map((groupEvent) => (
              <article key={groupEvent.id} className={styles.groupEventCard}>
                <div className={styles.groupEventHeader}>
                  {/* <div>
                    <p>By {getGroupMemberDisplayName(groupEvent.Creator)}</p>
                  </div> */}
                  <div className={styles.groupPostAuthor}>
                    <img
                      src={buildBackendAssetUrl(groupEvent.Creator.avatarUrl||DEFAULT_AVATAR_PATH)}
                      alt=""
                      width={40}
                      height={40}
                      className={styles.groupAvatar}
                      />
                    <div>
                    <p>By {getGroupMemberDisplayName(groupEvent.Creator)}</p>
                      {/* <strong>{getGroupMemberDisplayName(post.creater)}</strong>
                      <small>{describeRelativeTimestamp(post.createdAt)}</small> */}
                    </div>
                  </div>
                  <time dateTime={groupEvent.event_date}>
                    {formatGroupEventSchedule(groupEvent.event_date)}
                  </time>
                </div>
                  <h2>{groupEvent.title}</h2>
                <p>{groupEvent.description}</p>
                <div className={styles.groupVoteButtons}>
                  <button
                    type="button"
                    className={
                      groupEvent.userVote === "going" ? styles.groupVotedGoing : ""
                    }
                    disabled={votingEventId === groupEvent.id}
                    onClick={() => void handleVote(groupEvent, "going")}
                  >
                    Going ({groupEvent.TotalParticipate})
                  </button>
                  <button
                    type="button"
                    className={
                      groupEvent.userVote === "not going"
                        ? styles.groupVotedNotGoing
                        : ""
                    }
                    disabled={votingEventId === groupEvent.id}
                    onClick={() => void handleVote(groupEvent, "not going")}
                  >
                    Not going ({groupEvent.totalNotparticipate})
                  </button>
                </div>
              </article>
            ))}
            {hasMoreEvents && (
          <button
            type="button"
            className={uiStyles.uiSecondaryButton}
            style={{ width: '100%', marginTop: '1rem' }}
            disabled={isMoreEventsLoading}
            onClick={handleSeeMoreEvents}
          >
            {isMoreEventsLoading ? "Loading..." : "See More Events"}
          </button>
        )}
         </>
          ) : (
            <p className={uiStyles.uiEmptyState}>No events yet.</p>
          )}
        </div>
      )}
      {activeTab === "chat" && groupId && (
        <GroupConversation groupId={groupId} />
      )}

      {activeTab === "members" && (
        <section className={styles.groupMemberSection}>

          {Array.isArray(members) && members.length > 0 ? (
            <div className={styles.groupMemberGrid}>
              {members.map((member) => (
                <Link href={`/UserProfile?user_id=${member.id}`} key={member.id} className={styles.groupMemberCard}>
                  <span>
                    <strong>{member.nickname || "User"}</strong>
                    <small>@{member.nickname || member.id}</small>
                    <img
                      src={buildBackendAssetUrl(member.avatar_url,DEFAULT_AVATAR_PATH)}

                      alt=""
                      width={48}
                      height={48}
                      className={stylesmember.discussionAvatar}
                    >
                    </img>

                  </span>
                </Link>
              ))}
            </div>
          ) : (
            <p className={uiStyles.uiEmptyState}>No members found.</p>
          )}
        </section>
      )}
      
      {activeTab === "pending" && isOwner && (
        <section className={styles.groupMemberSection}>
          {pendingMembers.length > 0 ? (
            <div className={styles.groupPendingList}>
              {pendingMembers.map((member) => (
                <div key={member.id} className={styles.groupPendingMember}>
                  <img
                    src={buildBackendAssetUrl(
                      member.avatarUrl,
                      DEFAULT_AVATAR_PATH,
                    )}
                    alt=""
                    width={44}
                    height={44}
                    className={styles.groupAvatar}
                  />
                  <span>
                    <strong>{getGroupMemberDisplayName(member)}</strong>
                    <small>@{member.nickname || member.id}</small>
                  </span>
                  <div className={styles.groupMemberActions}>
                    <button
                      type="button"
                      className={styles.groupAcceptButton}
                      disabled={processingMemberId === member.id}
                      onClick={() =>
                        void handleMembershipDecision(member.id, "joinaccept")
                      }
                    >
                      Accept
                    </button>
                    <button
                      type="button"
                      className={styles.groupRejectButton}
                      disabled={processingMemberId === member.id}
                      onClick={() =>
                        void handleMembershipDecision(member.id, "joinreject")
                      }
                    >
                      Reject
                    </button>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className={uiStyles.uiEmptyState}>No pending requests.</p>
          )}
        </section>
      )}

      {isInviteDialogOpen && (
        <div
          className={styles.groupDialogBackdrop}
          onClick={() => setIsInviteDialogOpen(false)}
        >
          <div
            className={styles.groupDialog}
            role="dialog"
            aria-modal="true"
            aria-labelledby="invite-dialog-title"
            onClick={stopDialogPropagation}
          >
            <header className={styles.groupDialogHeader}>
              <h2 id="invite-dialog-title">Invite members</h2>
              <button
                type="button"
                onClick={() => setIsInviteDialogOpen(false)}
                aria-label="Close invitation dialog"
              >
                ×
              </button>
            </header>
            <div className={styles.groupInviteList}>
              {isLoadingInvitableUsers ? (
                <p>Loading followers…</p>
              ) : invitableUsers.length === 0 ? (
                <p>There is nobody left to invite.</p>
              ) : (
                invitableUsers.map((invitableUser) => {
                  const hasBeenInvited = invitedUserIds.has(invitableUser.id);
                  return (
                    <div key={invitableUser.id} className={styles.groupInviteRow}>
                      <img
                        src={buildBackendAssetUrl(
                          invitableUser.avatarUrl,
                          DEFAULT_AVATAR_PATH,
                        )}
                        alt=""
                        width={40}
                        height={40}
                        className={styles.groupAvatar}
                      />
                      <span>
                        <strong>
                          {getGroupMemberDisplayName(invitableUser)}
                        </strong>
                        <small>
                          @{invitableUser.nickname || invitableUser.id}
                        </small>
                      </span>
                      <button
                        type="button"
                        disabled={
                          hasBeenInvited ||
                          sendingInviteUserId === invitableUser.id
                        }
                        onClick={() => void sendInvitation(invitableUser.id)}
                      >
                        {hasBeenInvited
                          ? "Invited"
                          : sendingInviteUserId === invitableUser.id
                            ? "Sending…"
                            : "Invite"}
                      </button>
                    </div>
                  );
                })
              )}
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
export default function GroupWorkspaceRoute() {
  return (
    <Suspense fallback={<p className={uiStyles.uiEmptyState}>Initializing Workspace...</p>}>
      <GroupWorkspaceContent />
    </Suspense>
  );
}