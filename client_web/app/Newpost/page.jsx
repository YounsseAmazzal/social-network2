"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import {
  DEFAULT_AVATAR_PATH,
  fetchfromApi,
  getErrorMessage,
  isRecord,
  readString,
  buildBackendAssetUrl,
} from "@/lib/api";
import { useAuthSession } from "@/providers/AuthSessionProvider";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/create-post.module.css";
function normalizeFollowers(value) {
  if (!Array.isArray(value)) return [];
  return value.flatMap((item) => {
    if (!isRecord(item)) return [];
    const id = readString(item, ["id", "ID"]);
    if (!id) return [];
    return [
      {
        id,
        firstName: readString(item, ["firstname", "first_name", "Firstname"]),
        lastName: readString(item, ["lastname", "last_name", "Lastname"]),
        nickname: readString(item, ["nickname", "Nickname"]),
        avatarUrl: readString(item, ["avatar", "avatar_url", "avatarURL"]),
      },
    ];
  });
}
export default function CreatePostRoute() {
  const router = useRouter();
  const { user } = useAuthSession();
  const { showToast } = useToastMessages();
  const [errorMessage, setErrorMessage] = useState("");
  const [privacy, setPrivacy] = useState("publicpost");
  const [followers, setFollowers] = useState([]);
  const [sharedWithUserIds, setSharedWithUserIds] = useState([]);
  const [isLoadingFollowers, setIsLoadingFollowers] = useState(false);
  const [imagePreview, setImagePreview] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const userId = user?.id ?? "";
  //console.log("Inisiaaal state",privacy)
  useEffect(() => {
    if (privacy !== "followers_chosen") {
      setFollowers([]);
      setSharedWithUserIds([]);
      return;
    }
    if (!userId) return;
    let isCurrentRequest = true;
    setIsLoadingFollowers(true);
    void fetchfromApi(`/Userfollowers?useeer_id=${userId}`)
      .then((response) => {
        if (isCurrentRequest)
          //console.log("THeeeeee Response of Followers",response.data)
          setFollowers(normalizeFollowers(response?.data|| []));
          //console.log("T%%%%%%%%%%%%%%%%he followers ",followers)
      })
      .catch(() => showToast("Failed to load followers.", "error"))
      .finally(() => {
        if (isCurrentRequest) setIsLoadingFollowers(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [privacy, userId]);
  useEffect(() => {
    if (!imagePreview) return;
    return () => URL.revokeObjectURL(imagePreview);
  }, [imagePreview]);
  function toggleFollower(userIdToShare, isSelected) {
    setSharedWithUserIds((currentIds) =>
      isSelected
        ? [...new Set([...currentIds, userIdToShare])]
        : currentIds.filter((id) => id !== userIdToShare),
    );
  }
  function Imagechanginghandler(event) {
    const image = event.target.files?.[0];
    setImagePreview(image ? URL.createObjectURL(image) : null);
  }
  function removeImage() {
    setImagePreview(null);
    const input = document.getElementById("shared_image");
    if (input instanceof HTMLInputElement) input.value = "";
  }
  async function Submithandling(event) {
    event.preventDefault();
    setErrorMessage("");
    const formData = new FormData(event.currentTarget);
    formData.append("privacy", privacy);
    const title = String(formData.get("PostTitle") ?? "").trim();
    const description = String(formData.get("PostContent") ?? "").trim();
    if (!title || !description) {
      setErrorMessage("Title and description are required.");
      return;
    }
    if (title.length > 255) {
      setErrorMessage("Title must be 255 characters or fewer.");
      return;
    }
    if (description.length > 500) {
      setErrorMessage("Description must be 500 characters or fewer.");
      return;
    }
    //console.log("Theeeeee priveeeeeeeeeery 0",privacy)
    if (privacy === "followers_chosen") {
      sharedWithUserIds.forEach((id) => formData.append("followers_chosen", id));
    }
    setIsSubmitting(true);
    try {
      await fetchfromApi("/CreatePost", { method: "POST", body: formData });
      showToast("Post created successfully.", "success");
      router.push("/");
    } catch (error) {
      showToast(getErrorMessage(error, "Could not create the post."), "error");
    } finally {
      setIsSubmitting(false);
    }
    // //console.log("The list of the followers",followers)
  }
  return (
    <main className={`${uiStyles.uiPage} ${uiStyles.uiNarrowPage}`}>
      <div className={styles.newPostCard}>
        <header className={styles.newPostHeader}>
          <h1>Create a post</h1>
          <p>Share an update with your network.</p>
        </header>

        {errorMessage && (
          <p className={uiStyles.uiErrorMessage} role="alert">
            {errorMessage}
          </p>
        )}

        <form className={styles.newPostForm} onSubmit={Submithandling}>
          <div className={styles.newPostFormGroup}>
            <label className={styles.newPostLabel} htmlFor="post-title">
              Title
            </label>
            <input
              className={styles.newPostInput}
              id="PostTitle"
              name="PostTitle"
              maxLength={255}
              required
            />
          </div>

          <div className={styles.newPostFormGroup}>
            <label className={styles.newPostLabel} htmlFor="post-description">
              Description
            </label>
            <textarea
              className={styles.newPostTextarea}
              id="PostContent"
              name="PostContent"
              rows={6}
              maxLength={500}
              required
            />
          </div>

          <div className={styles.newPostFormGroup}>
            <label className={styles.newPostLabel} htmlFor="shared_image">
              Image <span>(optional)</span>
            </label>
            <input
              className={styles.newPostFileInput}
              id="shared_image"
              name="shared_image"
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={Imagechanginghandler}
            />
            <label className={styles.newPostFileLabel} htmlFor="shared_image">
              Choose an image
            </label>
            {imagePreview && (
              <div className={styles.newPostPreview}>
                <Image
                  src={imagePreview}
                  alt="Selected post preview"
                  fill
                  sizes="(max-width: 44rem) 100vw, 42rem"
                  unoptimized
                />
                <button
                  type="button"
                  onClick={removeImage}
                  aria-label="Remove image"
                >
                  ×
                </button>
              </div>
            )}
          </div>

          <div className={styles.newPostFormGroup}>
            <label className={styles.newPostLabel} htmlFor="post-privacy">
              Privacy
            </label>
            <select
              className={styles.newPostSelect}
              id="postprivecy"
              value={privacy}
              onChange={(event) => setPrivacy(event.target.value)}
            >
              <option value="publicpost">Public</option>
              <option value="only_followers">Followers only</option>
              <option value="followers_chosen">Specific followers</option>
            </select>
          </div>

          {privacy === "followers_chosen" && (
            <fieldset className={styles.newPostFollowersFieldset}>
              <legend>Share with specific followers</legend>
              {isLoadingFollowers ? (
                <p>Loading followers…</p>
              ) : followers.length > 0 ? (
                <div className={styles.newPostFollowersList}>
                  {followers.map((follower) => (
                    <label key={follower.id} className={styles.newPostFollower}>
                      <Image
                        src={buildBackendAssetUrl(
                          follower.avatar_url,
                          DEFAULT_AVATAR_PATH,
                        )}
                        alt=""
                        width={40}
                        height={40}
                        className={styles.newPostFollowerAvatar}
                      />
                      <span>
                        <strong>
                          {follower.firstName} {follower.lastName}
                        </strong>
                        <small>@{follower.nickname}</small>
                      </span>
                      <input
                        type="checkbox"
                        checked={sharedWithUserIds.includes(follower.id)}
                        onChange={(event) =>
                          toggleFollower(follower.id, event.target.checked)
                        }
                      />
                    </label>
                  ))}
                </div>
              ) : (
                <p>No followers are available.</p>
              )}
            </fieldset>
          )}

          <button
            type="submit"
            className={uiStyles.uiPrimaryButton}
            disabled={isSubmitting}
          >
            {isSubmitting ? "Publishing…" : "Publish post"}
          </button>
        </form>
      </div>
    </main>
  );
}
