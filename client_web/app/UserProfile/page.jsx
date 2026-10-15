
"use client";
import { useEffect, useState,Suspense } from "react";
import Image from "next/image";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import {
  DEFAULT_AVATAR_PATH,
  DEFAULT_PUBLIC_PIC_PATH,
  fetchfromApi,
  getErrorMessage,
  buildBackendAssetUrl,
} from "@/lib/api";
import { normalizeAuthUser } from "@/lib/auth";
import {
  fetchPendingFollowRequests,
  fetchProfile,
  fetchProfileConnections,
} from "../../lib/profile.js";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import stylesmember from "@/app/styles/post-discussion.module.css";
import styles from "@/app/styles/profile-route.module.css";
import styles2 from "@/app/styles/group-workspace.module.css";

function getProfileUserDisplayName(user) {
  return (
    `${user.firstname} ${user.lastname}`.trim() ||
    user.nickname ||
    "Unknown user"
  );
}
function ProfileContent() {
  const { user, isCheckingSession, updateUser } = useAuthSession();
  const { showToast } = useToastMessages();
  const searchParams = useSearchParams();
  const [profileResponse, setProfileResponse] = useState(null);
  const [isLoadingProfile, setIsLoadingProfile] = useState(true);
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [pendingRequests, setPendinnngRequests] = useState([]);
  const [isLoadingConnections, setIsLoadingConnections] = useState(false);
  const [activeTab, setActiveTab] = useState("posts");
  const [nickname, setNickname] = useState("");
  const [aboutMe, setAboutMe] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const [avatarFile, setAvatarFile] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [isSaving, setIsSaving] = useState(false);
 const currentUserId = user?.id ?? "";
  const profileId = searchParams.get("user_id") || currentUserId;
  const isOwnProfile = profileId === currentUserId;
  // const isOwnProfile = !profileId || profileId === currentUserId;

  
  const availableTabs = isOwnProfile
    ? ["posts", "followers", "following", "settings","requests"]
    : ["posts", "followers", "following"];
useEffect(() => {
    if (isCheckingSession || !currentUserId) return;
    let isCurrentRequest = true;
    setIsLoadingProfile(true);

    fetchProfile(profileId)
      .then((response) => {
        if (!isCurrentRequest || !response) return;
        
        setProfileResponse(response);

          if (response.isOwner) {
          const u = response.profile.user;
          setNickname(u.nickname || "");
          setAboutMe(u.aboutMe || "");
          setIsPrivate(response.isPrivate);
        }
      })
      .catch((error) => {
        if (isCurrentRequest) {
          showToast(
            getErrorMessage(error, "Failed to load this profile."),
            "error",
          );
        }
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoadingProfile(false);
      });

    return () => { isCurrentRequest = false; };
  }, [currentUserId, isCheckingSession, profileId]);
  useEffect(() => {
    const requestedTab = searchParams.get("tab");
    if (requestedTab && availableTabs.includes(requestedTab)) {
      setActiveTab(requestedTab);
    }
  }, [searchParams]);
  useEffect(() => {
    if (!profileResponse || !currentUserId) return;
    if (activeTab !== "followers" && activeTab !== "following" && activeTab !== "requests") return;
    const targetUserId = profileId || currentUserId;
    let isCurrentRequest = true;
    setIsLoadingConnections(true);
    void fetchProfileConnections(targetUserId)
      .then((connections) => {
        if (!isCurrentRequest) return;
        console.log("the followere " ,connections.following)
        setFollowers(connections.followers);
        setFollowing(connections.following);
      })
      .catch((error) => {
        if (isCurrentRequest) {
          showToast(
            getErrorMessage(error, "Failed to load connections."),
            "error",
          );
        }
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoadingConnections(false);
      });
    if (activeTab === "requests" && isOwnProfile) {
      void fetchPendingFollowRequests()
        .then((requests) => {
          if (isCurrentRequest) setPendinnngRequests(requests);
        })
        .catch((error) => {
          if (isCurrentRequest) {
            showToast(
              getErrorMessage(error, "Failed to load follow requests."),
              "error",
            );
          }
        });
    }
    return () => {
      isCurrentRequest = false;
    };
  }, [activeTab, currentUserId, isOwnProfile, profileId, profileResponse]);
  useEffect(() => {
    if (!avatarPreview) return;
    return () => URL.revokeObjectURL(avatarPreview);
  }, [avatarPreview]);
  async function refreshProfile() {
    const response = await fetchProfile(profileId);
    setProfileResponse(response);
    return response;
  }
  function FollowingStatus(response, isOwnProfile) {
  if (isOwnProfile) return "none";
  
  const rel = response.relationship; 

  if (rel === "followpending" || rel === "pending") return "pending";
  if (rel === "folowAccepted" || rel === "accepted") return "following";
  if (rel === "follower") return "follower"; 
  
  return "none";
}

const handleToggleFollow = () => {
    const rel = profileResponse?.relationship;
    
    if (rel === "folowAccepted" || rel === "followpending") {
        UnfollowHandler();
    } else {
        FollowHandler();
    }
};

  async function FollowHandler() {
    if (!profileId) return;
    try {
      const response = await fetchfromApi(`/Userfollowinng?Userfollowinng_id=${profileId}`, {
        method: "POST",
      });
      showToast(response.message || "Follow request sent.", "success");
      await refreshProfile();
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to follow this user."), "error");
    }
  }
  async function UnfollowHandler() {
    if (!profileId) return;
    try {
      await fetchfromApi(`/UserUnfollow?UserUnfollow_id=${profileId}`, {
        method: "DELETE",
      });
      showToast("Unfollowed.", "success");
      await refreshProfile();
      const connections = await fetchProfileConnections(profileId);
      setFollowers(connections.followers);
      setFollowing(connections.following);
    } catch (error) {
      showToast(
        getErrorMessage(error, "Failed to unfollow this user."),
        "error",
      );
    }
  }
async function handleFollowRequestDecision(followerId, decision) {
  try {
    const endpoint = decision === "accept" 
      ? `/followers/Accepting?FollowAccet_id=${followerId}` 
      : `/followers/Declineinng?FollowDecline_id=${followerId}`;

    await fetchfromApi(endpoint, {
      method: "POST",
    });

    setPendinnngRequests((currentRequests) =>
      currentRequests.filter((request) => request.id !== followerId)
    );

    showToast(
      decision === "accept" ? "Request accepted." : "Request declined.",
      "success"
    );
    
    await refreshProfile();
  } catch (error) {
    showToast(getErrorMessage(error, "Failed to update the request."), "error");
  }
}
  function avataarhandler(event) {
    const file = event.target.files?.[0];
    if (!file) return;
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
  }
  async function handleSaveProfile(event) {
    event.preventDefault();
    setIsSaving(true);
    try {
      let response;
      const payload = { 
        nickname, 
        about_me: aboutMe, 
        is_private: isPrivate 
      };
      
      if (avatarFile) {
        const formData = new FormData();
        formData.append("nickname", nickname);
        formData.append("about_me", aboutMe);
        formData.append("is_private", String(isPrivate)); 
        formData.append("avatar", avatarFile);
        
        response = await fetchfromApi("/ProfileUppdate", {
          method: "PATCH",
          body: formData,
        });
      } else {
        response = await fetchfromApi("/ProfileUppdate", {
          method: "PATCH",
          body: JSON.stringify(payload),
        });
      }
      
      if (response?.user) {
        updateUser(normalizeAuthUser(response.user));
      }
      
      await refreshProfile(); 
      showToast("Profile updated successfully!", "success");
      setAvatarFile(null);
      setAvatarPreview(null);
    } catch (error) {
      showToast(getErrorMessage(error, "Update failed."), "error");
    } finally {
      setIsSaving(false);
    }
  }
  if (isCheckingSession || isLoadingProfile) return <p className={uiStyles.uiEmptyState}>Loading...</p>;
  if (!profileResponse) return <p className={uiStyles.uiEmptyState}>Profile not found.</p>;
  
  // const profileUser = profileResponse;
  //console.log("The profile response1",profileResponse.profile)
  //console.log("The profile response",profileUser)
  //console.log("The profile response",profileUser)
  //console.log("The profile response",profileUser)
  //console.log("The profile response",profileUser)
  //console.log("The profile response",profileUser.canSee)
  //console.log("The profile response",profileUser.canSee)
  //console.log("The profile response",profileUser.canSee)
  
  const followStatus = FollowingStatus(profileResponse, isOwnProfile);
  const isLocked = profileResponse.isPrivate && !profileResponse.isOwner && followStatus !== "following";


const getFollowButtonLabel = () => {
    const rel = profileResponse?.relationship;
    if (rel === "followpending" || rel === "pending") return "Cancel Request";
    if (rel === "folowAccepted" || rel === "accepted") return "Unfollow";
    if (rel === "follower") return "Follow Back";
    return "Follow";
};
// if (!canSeeContent) {

// }
// console.log("the profile responnnnnseeeee",profileResponse)
const avatarSource = avatarPreview || buildBackendAssetUrl(
    profileResponse.profile.user.avatarUrl || user?.image, 
    DEFAULT_AVATAR_PATH
);

  return (
    <main className={`${uiStyles.uiPage} ${styles.profilePage}`}>
   <header className={styles.profileHeader}>
      <img src={avatarSource} alt="" width={104} height={104} className={styles.profileAvatar} />
      <div className={styles.profileSummary}>
        <h1>{profileResponse.profile.user.firstname} {profileResponse.profile.user.lastname}</h1>
        <p>@{profileResponse.profile.user.nickname}</p>
        
        {/* Stats and Bio are always visible in modern social networks */}
        <div className={styles.profileStats}>
          <span><strong>{profileResponse.profile.postsCount}</strong> posts</span>
          <span><strong>{profileResponse.followersCount || 0}</strong> followers</span>
          <span><strong>{profileResponse.followingCount || 0}</strong> following</span>
        </div>
    {profileResponse.profile.user.aboutMe && <p className={styles.profileBio}>{profileResponse.profile.user.aboutMe}</p>}
      {!isOwnProfile && (
          <div className={styles.profileActions}>
            <button
              type="button"
              className={
                followStatus === "none" || followStatus === "follower" 
                  ? uiStyles.uiPrimaryButton 
                  : uiStyles.uiSecondaryButton
              }
              onClick={handleToggleFollow} 
            >
              {getFollowButtonLabel()}
            </button>
          </div>
        )}
      </div>
    </header>
    {isLocked ? (
        <div className={uiStyles.uiEmptyState} style={{ marginTop: '3rem', textAlign: 'center' }}>
          <p style={{ fontSize: '2rem' }}>🔒</p>
          <strong>This account is private</strong>
          <p>Follow this user to see their posts and photos.</p>
                    <div style={{ marginTop: '1rem' }}>
      <button
        type="button"
        className={uiStyles.uiPrimaryButton}
        onClick={handleToggleFollow} 
      >
        {getFollowButtonLabel()}
      </button>
    </div>
        </div>
      ) : (
        <>
    
      <div className={styles.profileTabs} role="tablist" aria-label="Profile sections">
        {availableTabs.map((tab) => (
          <button
            type="button"
            role="tab"
            key={tab}
            aria-selected={activeTab === tab}
            className={activeTab === tab ? styles.profileActiveTab : ""}
            onClick={() => setActiveTab(tab)}
          >
            {tab.charAt(0).toUpperCase() + tab.slice(1)}
          </button>
        ))}
      </div>

      {activeTab === "posts" &&
        (profileResponse.profile.posts.length > 0 ? (
          <div className={styles.profilePostList}>
            {profileResponse.profile.posts.map((post) => (
              <article key={post.id} className={styles.profilePostCard}>
                {post.avatar_url && (
                  <img
                    src={buildBackendAssetUrl(post.avatar_url)}
                    alt=""
                    width={700}
                    height={360}
                    className={styles.profilePostImage}
                  />
                )}
                <div className={styles.profilePostBody}>
                  <h2>{post.title}</h2>
                  <p>{post.description}</p>
                  <footer>
                    <span>
                      <Image
                        src="/icons/like.svg"
                        alt="Likes"
                        width={16}
                        height={16}
                      />{" "}
                      {post.likes || 0}
                    </span>
                    <span>
                      <Image
                        src="/icons/comments.svg"
                        alt="Comments"
                        width={16}
                        height={16}
                      />{" "}
                      {post.comments || 0}
                    </span>
                    <time>{post.date}</time>
                  </footer>
                </div>
              </article>
            ))}
          </div>
        ) : (
          <p className={uiStyles.uiEmptyState}>No posts yet.</p>
        ))}

{activeTab === "followers" && (
  <section className={styles.connectionsSection}>
    {Array.isArray(followers) && followers.length > 0 ? (
      followers.map((f) => (
        <Link 
          href={`/UserProfile?user_id=${f.id}`} 
          key={f.id} 
          className={styles.profilePersonRow}
        >
          <img
            src={buildBackendAssetUrl(f.avatarUrl || f.avatar, DEFAULT_AVATAR_PATH)}
            width={44}
            height={44}
            className={styles.avatarCircle}
            alt={`${f.nickname}'s avatar`}
          />
          <div className={styles.personInfo}>
            <strong>{f.firstname} {f.lastname}</strong>
            <small>@{f.nickname}</small>
          </div>
        </Link>
      )) 
    ) : (
      <p className={uiStyles.uiEmptyState}>No followers yet.</p>
    )}
  </section>
)}

      {activeTab === "following" && (
        <section>
          <h2 className={styles.profileSectionTitle}>Following</h2>
          {isLoadingConnections ? (
            <p className={uiStyles.uiEmptyState}>Loading connections…</p>
          ) : following.length > 0 ? (
            following.map((connection) => (
              <Link
                href={`/UserProfile?user_id=${connection.id}`}
                key={connection.id}
                className={styles.profilePersonRow}
              >
                <img
                  src={buildBackendAssetUrl(
                    connection.avatar,
                    DEFAULT_AVATAR_PATH,
                  )}
                  alt=""
                  width={44}
                  height={44}
                  className={styles.profilePersonAvatar}
                />
                <span>
                  <strong>{getProfileUserDisplayName(connection)}</strong>
                  <small>@{connection.nickname}</small>
                </span>
              </Link>
            ))
          ) : (
            <p className={uiStyles.uiEmptyState}>Not following anyone yet.</p>
          )}
        </section>
      )}
{isOwnProfile && activeTab === "requests" && (
  <section className={styles2.groupMemberSection}>
    {pendingRequests.length > 0 ? (
      <div className={styles2.groupPendingList}>
        {pendingRequests.map((req) => (
          <div key={req.id} className={styles2.groupPendingMember}>
            <img
              src={buildBackendAssetUrl(req.avatarUrl || req.avatar, DEFAULT_AVATAR_PATH)}
              alt=""
              width={44}
              height={44}
              className={styles2.groupAvatar}
            />

            <span>
              <strong>{req.firstname} {req.lastname}</strong>
              <small>@{req.nickname || "user"}</small>
            </span>

            <div className={styles2.groupMemberActions}>
              <button
                type="button"
                className={styles2.groupAcceptButton}
                onClick={() => void handleFollowRequestDecision(req.id, "accept")}
              >
                Accept
              </button>
              <button
                type="button"
                className={styles2.groupRejectButton}
                onClick={() => void handleFollowRequestDecision(req.id, "decline")}
              >
                Reject
              </button>
            </div>
          </div>
        ))}
      </div>
    ) : (
      <p className={uiStyles.uiEmptyState}>No pending follow requests.</p>
    )}
  </section>
)}

      {activeTab === "settings" && isOwnProfile && (
        <form className={styles.profileSettingsCard} onSubmit={handleSaveProfile}>
          <h2>Edit profile</h2>
          <div className={styles.profileAvatarEditor}>
            <img
              src={avatarSource}
              alt=""
              width={76}
              height={76}
              className={styles.profileSettingsAvatar}
              // unoptimized={avatarPreview !== null}
            />
            <label>
              Change avatar
              <input
                type="file"
                accept="image/jpeg,image/png,image/gif"
                onChange={avataarhandler}
              />
            </label>
          </div>
          <label className={styles.profileFormGroup}>
            Nickname
            <input
              value={nickname}
              onChange={(event) => setNickname(event.target.value)}
            />
          </label>
          <label className={styles.profileFormGroup}>
            About me
            <textarea
              value={aboutMe}
              maxLength={1000}
              rows={5}
              onChange={(event) => setAboutMe(event.target.value)}
            />
          </label>
<label className={styles.profilePrivateToggle}>
  <span>
    <strong>Private account</strong>
    <small>
      {isPrivate
        ? "Only approved followers can view your profile."
        : "Anyone can view your profile."}
    </small>
  </span>
  <input
    type="checkbox"
    checked={!!isPrivate} 
    onChange={(event) => setIsPrivate(event.target.checked)} 
  />
</label>
          <button
            type="submit"
            className={uiStyles.uiPrimaryButton}
            disabled={isSaving}
          >
            {isSaving ? "Saving…" : "Save changes"}
          </button>
        </form>
      )}
      </>
      )}
    </main>
  );
}

export default function ProfileRouteView() {
  return (
    <Suspense fallback={<p>Loading Profile...</p>}>
      <ProfileContent />
    </Suspense>
  );
}
