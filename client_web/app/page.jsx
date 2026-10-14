"use client";
import { useEffect, useState ,Suspense} from "react";
import { useSearchParams } from "next/navigation"; 

import Image from "next/image";
import Link from "next/link";
import { fetchfromApi, buildBackendAssetUrl, DEFAULT_AVATAR_PATH, getErrorMessage } from "@/lib/api";
import { fetcFeedPosts, normalizeFeedPosts } from "@/lib/posts";
import { formatPostTimestamp } from "@/lib/time";
import { useAuthSession } from "../providers/AuthSessionProvider.jsx";
import { useToastMessages } from "@/providers/ToastMessageProvider";
import uiStyles from "@/app/styles/ui.module.css";
import styles from "@/app/styles/feed.module.css";
//  const [currentOffset, setCurrentOffset] = useState(0);
//   const [hasMore, setHasMore] = useState(true);
//   const [isMoreLoading, setIsMoreLoading] = useState(false);
//   const Page_limiit = 10;
export function FeedContent() {
  const searchParams = useSearchParams();
  // convert the limit and the offset to number to avoid problems on calculation
const limit = Number(searchParams.get("limit") || "10");
const offset = Number(searchParams.get("offset") || "0");
  const [currentOffset, setCurrentOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [isMoreLoading, setIsMoreLoading] = useState(false);
  const { user } = useAuthSession();
  const { showToast } = useToastMessages();
  const [posts, setPosts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const userId = user?.id ?? "";
  //console.log("HOOMe Page the id ",userId)
  useEffect(() => {
    if (!userId) {
      setIsLoading(false);
      return;
    }
    let isCurrentRequest = true;
    setIsLoading(true);
    void fetcFeedPosts("/AllPosts",limit,offset)
    .then((response) => {
      //console.log("the resssponse date check ",response)
      if (!isCurrentRequest) return;
          const initialPosts = normalizeFeedPosts(response);
          // console.log("feed posts",response)
          // console.log("feed posts",initialPosts)
        setPosts(initialPosts);
          if (initialPosts.length < limit) {
          setHasMore(false);
        }
      })
      .catch((error) => {
        if (isCurrentRequest) {
          showToast(getErrorMessage(error, "Failed to load posts."), "error");
        }
      })
      .finally(() => {
        if (isCurrentRequest) setIsLoading(false);
      });
    return () => {
      isCurrentRequest = false;
    };
  }, [userId],limit,offset);


   const handleLoadMore = async () => {
    if (isMoreLoading || !hasMore) return;
    
    setIsMoreLoading(true);
    const nextOffset = currentOffset + limit;

    try {
      const response = await fetcFeedPosts("/AllPosts", limit, nextOffset);
      const newPosts = normalizeFeedPosts(response);

      if (newPosts.length === 0) {
        setHasMore(false);
      } else {
        setPosts((prev) => [...prev, ...newPosts]);
        setCurrentOffset(nextOffset);
        
        if (newPosts.length < limit) {
          setHasMore(false);
        }
      }
    } catch (error) {
      showToast(getErrorMessage(error, "Failed to load more posts."), "error");
    } finally {
      setIsMoreLoading(false);
    }
  };
  if (isLoading) {
    return (
      <div className={uiStyles.uiPage}>
        <p className={uiStyles.uiEmptyState} role="status">
          Loading posts…
        </p>
      </div>
    );
  }
  //console.log("The response from server ",posts)
  return (
    <main className={`${uiStyles.uiPage} ${styles.feed}`}>
      <header className={styles.feedHeader}>
        <div>
          <h1>Your feed</h1>
          <p>Recent posts from people and groups you follow.</p>
        </div>
        <Link href="/Newpost" className={uiStyles.uiPrimaryButton}>
          Create post
        </Link>
      </header>


      {posts.length > 0 ? (
  <div className={styles.feedPostList}>
    {posts.map((post) => (
      <article key={post.id} className={styles.feedPostCard}>
        <div className={styles.feedPostHeader}>
          <div>
            <h2>{post.title}</h2>
            <p>By <strong>@{post.createrName || "unknown"}</strong></p>
          </div>
          <time dateTime={post.createDate}>
            {formatPostTimestamp(post.createDate)}
          </time>
        </div>

        {post.image && (
          <img
            src={buildBackendAssetUrl(post.image)}
            alt="Post content"
            width={800}
            height={420}
            className={styles.feedPostImage}
          />
        )}

        <div className={styles.feedPostFooter}>
          <div className={styles.feedPostMetrics}>
            <span>
              <Image src="/icons/like.svg" alt="Likes" width={16} height={16} />{" "}
              {post.likesNum || 0}
            </span>
            <span>
              <Image
                src="/icons/comments.svg"
                alt="Comments"
                width={16}
                height={16}
              />{" "}
              {post.Commentsnum || 0}
            </span>
          </div>
          
          <Link
            href={`/GetPost?Possst_id=${post.id}`}
            className={styles.feedDiscussionLink}
          >
            View Post
          </Link>
        </div>
      </article>
    ))}

         <div className={styles.loadMoreWrapper}>
            {hasMore ? (
              <button 
                className={uiStyles.uiSecondaryButton} 
                onClick={handleLoadMore}
                disabled={isMoreLoading}
              >
                {isMoreLoading ? "Loading..." : "See More Posts"}
              </button>
            ) : (
              <p className={uiStyles.uiEmptyState}>You have seen all the posts!</p>
            )}
          </div>
        </div>
        
      ) : (
        <div className={uiStyles.uiEmptyState}>
          <strong>No posts yet</strong>
          <p>Follow more people or create the first post.</p>
        </div>
      )}

    </main>
  );
}
export default function FeedRoute() {
  return (
    <Suspense fallback={<p>Loading Feed...</p>}>
      <FeedContent />
    </Suspense>
  );
}