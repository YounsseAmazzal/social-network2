"use client";
import { useEffect, useState, Suspense } from "react"; 
import Image from "next/image";
import Link from "next/link";
// import { useParams } from "next/navigation";
import {
  DEFAULT_AVATAR_PATH,
  DEFAULT_PUBLIC_PIC_PATH,
  fetchfromApi,
  getErrorMessage,
  buildBackendAssetUrl,
} from "@/lib/api";
// import { normalizeVoteResponse } from "@/lib/chat";
import { getOptimisticVoteCounts, normalizePostResponse } from "@/lib/posts";
import { formatPostTimestamp } from "@/lib/time";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/post-discussion.module.css";
import { useSearchParams } from "next/navigation"; 

function PostDiscussionContent() {
    const searchParams = useSearchParams();
    const postId = searchParams.get("Possst_id") || "";
    //console.log("POSSSST ID FROM PAARAAM",postId)
    const { showToast } = useToastMessages();
    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [commentText, setCommentText] = useState("");
    const [commentImage, setCommentImage] = useState(null);
    const [commentImagePreview, setCommentImagePreview] = useState(null);
    const [isSubmittingComment, setIsSubmittingComment] = useState(false);

    /// offset and limit
    const [commentsOffset, setCommentsOffset] = useState(0);
    const Commments_limit = 10; 
    const [hasMoreComments, setHasMoreComments] = useState(true);
    const [isMoreLoading, setIsMoreLoading] = useState(false);

    useEffect(() => {
      //console.log("DSFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",postId)
      if (!postId) return;
      let isCurrentRequest = true;
      setIsLoading(true);
      // //console.log("@@@@@@@@@@@@@@@@@@@@")
      // //console.log("@@@@@@@@@@@@@@@@@@@@")
      // //console.log("@@@@@@@@@@@@@@@@@@@@")
      // //console.log("@@@@@@@@@@@@@@@@@@@@")
      // //console.log("@@@@@@@@@@@@@@@@@@@@")
      void fetchfromApi(`/GetPost?Possst_id=${postId}&limit=${Commments_limit}&offset=0`)
      .then((response) => {
          //console.log("The post details content",response)
        if (!isCurrentRequest) return;
        const normalizedPost = normalizePostResponse(response.PostDetails || response.PostGroupDetailed);
        //console.log("normaaalazw post ",normalizedPost)
        setPost(normalizedPost);
        const initiialComments = normalizedPost?.PostComments ?? [];
        // setComments(initiialComments);
        setComments(initiialComments);
          if (initiialComments.length < Commments_limit) {
            setHasMoreComments(false);
          }
      })
      .catch((error) => {
        if (!isCurrentRequest) return;
        setPost(null);
        showToast(getErrorMessage(error, "Failed to load this post."), "error");
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [postId]);
  async function handleLoadMoreComments() {
        const nextOffset = commentsOffset + Commments_limit;
        setIsMoreLoading(true);

        try {
          // request to the only the new comments 
          const response= await fetchfromApi(`/GetPost?Possst_id=${postId}&limit=${Commments_limit}&offset=${nextOffset}`);
            const normalizedPost = normalizePostResponse(response.PostDetails || response.PostGroupDetailed);
            const newComments = normalizedPost?.PostComments ?? [];;
              if (newComments.length === 0) {
                setHasMoreComments(false);
            } else {
                setComments(prev => [...prev, ...newComments]);
                setCommentsOffset(nextOffset);
                if (newComments.length < Commments_limit) setHasMoreComments(false);
            }
        } catch (error) {
            showToast(
              getErrorMessage(error, "Failed to load more comments."),
              "error",
            );
        } finally {
            setIsMoreLoading(false);
        }
    }

   async function handlePostVote(voteType) {
    if (!post) return;
    const endpoint = voteType === "like" ? "/Postlike" : "/Postdislike";

    try {
      await fetchfromApi(`${endpoint}?Possst_id=${postId}`, { method: "POST" });
      const response = await fetchfromApi(`/GetPost?Possst_id=${postId}`);
      const rawData = response.data || response;
      setPost(normalizePostResponse(rawData.PostDetails || rawData.PostGroupDetailed||[]));
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to save your vote."), "error");
    }
  }

  function handleCommentImage(event) {
    const image = event.target.files?.[0];
    if (!image) return;
    if (!image.type.startsWith("image/")) {
      showToast("Select a valid image file.", "error");
      return;
    }
    if (image.size > 5 * 1024 * 1024) {
      showToast("Image must be smaller than 5 MB.", "error");
      return;
    }
    setCommentImage(image);
    const reader = new FileReader();
    reader.onload = () => {
      if (typeof reader.result === "string")
        setCommentImagePreview(reader.result);
    };
    reader.readAsDataURL(image);
  }
  function removeCommentImage() {
    setCommentImage(null);
    setCommentImagePreview(null);
    const input = document.getElementById("comment-image");
    if (input instanceof HTMLInputElement) input.value = "";
  }
  async function handleAddComment(event) {
    event.preventDefault();
    if (!commentText.trim() && !commentImage) return;
    setIsSubmittingComment(true);
    const formData = new FormData();
    formData.append("Content", commentText.trim());
    formData.append("postID", postId);
    if (commentImage) formData.append("comment_immmage", commentImage);
    try {
      await fetchfromApi(`/CreateComment?Possst_id=${postId}`, { method: "POST", body: formData });
      const response = await fetchfromApi(`/GetPost?Possst_id=${postId}`);
      const refreshedPost = normalizePostResponse(response.PostDetails||response.PostGroupDetailed||response);
        if (refreshedPost) {
        setComments(refreshedPost.PostComments || []);
        setCommentsOffset(0);
        setHasMoreComments((refreshedPost.PostComments?.length || 0) >= Commments_limit);
      }
      //console.log("THE ressssponse of post commeenets",refreshedPost)
      setComments(refreshedPost?.PostComments ?? []);
      setCommentText("");
      removeCommentImage();
      showToast("Comment posted.", "success");
    } catch (error) {
      showToast(
        getErrorMessage(error, "Failed to post the comment."),
        "error",
      );
    } finally {
      setIsSubmittingComment(false);
    }
  }
  if (isLoading) {
    return (
      <p className={uiStyles.uiEmptyState} role="status">
        Loading post…
      </p>
    );
  }
  if (!post) {
    return (
      <section className={uiStyles.uiEmptyState}>
        <h1>Post not found</h1>
        <p>This post may have been removed or you may not have access to it.</p>
        <Link href="/">Back to feed</Link>
      </section>
    );
  }
  return (
    <main className={`${uiStyles.uiPage} ${uiStyles.uiNarrowPage}`}>
      <article className={styles.discussionPostCard}>
        <header className={styles.discussionPostHeader}>
          <img
            src={buildBackendAssetUrl(post.authorAvatar)}
            alt=""
            width={48}
            height={48}
            className={styles.discussionAvatar}
          />
          <div>
            <h1>{post.createrName || "Unknown author"}</h1>
            <time dateTime={post.createDate}>
              {formatPostTimestamp(post.createDate)}
            </time>
          </div>
        </header>

        <h2 className={styles.discussionPostTitle}>{post.title}</h2>
        <p className={styles.discussionPostContent}>{post.description}</p>
        {post.image && (
          <img
            src={buildBackendAssetUrl(post.image)}
            alt=""
            width={900}
            height={500}
            className={styles.discussionPostImage}
          />
        )}
        <div className={styles.discussionVoteRow}>
          <button
            type="button"
            className={post.userrating === "like" ? styles.discussionLiked : ""}
            aria-pressed={post.userrating === "like"}
            aria-label={`Like this post (${post.likesNum || 0})`}
            onClick={() => void handlePostVote("like")}
          >
            <Image src="/icons/like.svg" alt="Like" width={18} height={18} />{" "}
            {post.likesNum || 0}
          </button>
          <button
            type="button"
            className={post.userVote === "dislike" ? styles.discussionDisliked : ""}
            aria-pressed={post.userVote === "dislike"}
            aria-label={`Dislike this post (${post.dislikesNum || 0})`}
            onClick={() => void handlePostVote("dislike")}
          >
            <Image src="/icons/dislike.svg" alt="Dislike" width={18} height={18} />{" "}
            {post.dislikesNum || 0}
          </button>
        </div>
      </article>

      <section
        className={styles.discussionCommentsSection}
        aria-labelledby="comments-title"
      >
        <h2 id="Content">Comments ({comments.length})</h2>
        <form className={styles.discussionCommentForm} onSubmit={handleAddComment}>
          <label className="srOnly" htmlFor="comment-text">
            Comment
          </label>
          <textarea
            id="Content"
            value={commentText}
            maxLength={1000}
            placeholder="Share your thoughts"
            onChange={(event) => setCommentText(event.target.value)}
            required={!commentImage}
          />
          <div className={styles.discussionCommentActions}>
            <input
              id="comment_immmage"
              className={styles.discussionFileInput}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleCommentImage}
            />
            <label className={styles.discussionFileLabel} htmlFor="comment_immmage">
              Add image
            </label>
            {commentImagePreview && (
              <span className={styles.discussionPreview}>
                <Image
                  src={commentImagePreview}
                  alt="Comment image preview"
                  fill
                  sizes="52px"
                  unoptimized
                />
                <button
                  type="button"
                  onClick={removeCommentImage}
                  aria-label="Remove image"
                >
                  ×
                </button>
              </span>
            )}
            <button
              type="submit"
              className={uiStyles.uiPrimaryButton}
              disabled={
                isSubmittingComment || (!commentText.trim() && !commentImage)
              }
            >
              {isSubmittingComment ? "Posting…" : "Post comment"}
            </button>
          </div>
        </form>

        {comments.length > 0 ? (
          <div className={styles.discussionCommentList}>
            {comments.map((comment) => (
              <article key={comment.id} className={styles.discussionComment}>
                <img
                  src={buildBackendAssetUrl(
                    comment.authorAvatar,
                    DEFAULT_AVATAR_PATH,
                  )}
                  alt=""
                  width={42}
                  height={42}
                  className={styles.discussionAvatar}
                />
                <div className={styles.discussionCommentBody}>
                  <header>
                    <strong>{comment.creater_name || "Anonymous"}</strong>
                    <time dateTime={comment.created_date || comment.createdAt}>{
                comment.created_date }
                    </time>
                  </header>
                  {comment.content && <p>{comment.content}</p>}
                  {comment.image && (
                    <img
                      src={buildBackendAssetUrl(comment.image)}
                      alt=""
                      width={600}
                      height={350}
                      className={styles.discussionCommentImage}
                    />
                  )}
                </div>
              </article>
            ))}
                  {hasMoreComments && (
                            <button 
                                type="button" 
                                className={uiStyles.uiSecondaryButton}
                                onClick={handleLoadMoreComments}
                                disabled={isMoreLoading}
                                style={{ width: '100%', marginTop: '1rem' }}
                            >
                                {isMoreLoading ? "Loading..." : "Load Older Comments"}
                            </button>
                  )}
          </div>
        ) : (
          <p className={uiStyles.uiEmptyState}>No comments yet.</p>
        )}
      </section>
    </main>
  );
}
export default function PostDiscussionRoute() {
  return (
    <Suspense fallback={<p>Loading UI...</p>}>
      <PostDiscussionContent />
    </Suspense>
  );
}