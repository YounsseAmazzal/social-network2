import { isRecord, readNumber, readString } from "@/lib/api";
import { normalizeVoteStatus } from "@/lib/chat";
import {fetchWithFallback,extractGroupCollection} from "../lib/groups.js"
function normalizeFeedPost(value) {
  if (!isRecord(value)) return null;
  const id = readString(value, ["id", "ID"]);
  if (!id) return null;
  return {
    id,
    title: readString(value, ["title", "Title"], "Untitled"),
    description: readString(value, ["description", "Description"]),
    image: readString(value, ["image","image_url", "Image"]),
    authorId: readString(value, ["createrId", "authorId", "author_id"]),
    createrName: readString(
      value,
      ["createrName", "author_name"],
      "Unknown author",
    ),
    authorAvatar: readString(value, ["authorAvatar", "author_avatar"]),
    createDate: readString(value, ["createDate", "created_at", "createdAt"]),
    likesNum: readNumber(value, ["likes", "likesNum"]),
    dislikesNum: readNumber(value, ["dislikes", "dislikesNum"]),
    Commentsnum: readNumber(value, ["Commentsnum", "totalComments"]),
    voteStatus:
      typeof value.vote_status === "number" ? value.vote_status : null,
  };
}
function normalizeComment(value) {
  if (!isRecord(value)) return null;
  const id = readString(value, ["id", "ID"]);
  if (!id) return null;
  return {
    id,
    content: readString(value, ["content", "Content"]),
    image_url: readString(value, ["image_url", "Image_URL"]),
    creater_id: readString(value, ["creater_id", "creater_Id", "Creater_Id"]),
    creater_name: readString(value, ["creater_name"], "Anonymous"),
    creater_avatar: readString(value, ["creater_avatar", "author_avatar"]),
    likes: readNumber(value, ["likes", "Likes"]),
    created_date: readString(value, ["createDate", "created_date", "createdAt"]),
    lastCreatedAt: readString(value, ["lastCreateDate", "last_create_date"]),
  };
}
export async function fetcFeedPosts(groupId,limit,offset) {
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")

  const response = await fetchWithFallback([`/AllPosts?limit=${limit}&offset=${offset}`], {
    // method: "GET",
    // body: JSON.stringify({ offset: 0, limit: 20 }),
  });
  // //console.log("Grouuups Posts Response",response)
  // //console.log("Grouuups Posts Response",response)
  // //console.log("Grouuups Posts Response",response)
  // //console.log("Grouuups Posts Response",response)
  //console.log("Feeed Posts Response",extractGroupCollection(response).map(normalizeFeedPost))
  return extractGroupCollection(response).map(normalizeFeedPost);
}
function normalizePost(value) {
  if (!isRecord(value)) return null;
  const id = readString(value, ["id", "ID"]);
  if (!id) return null;
  const rawComments = Array.isArray(value.PostComments)
    ? value.PostComments
    : Array.isArray(value.PostComments)
      ? value.PostComments
      : [];
  return {
    id,
    title: readString(value, ["title", "Title"], "Untitled"),
    description: readString(value, ["description", "Description"]),
    image: readString(value, ["image_url", "Image"]),
    authorId: readString(value, ["authorID", "authorId", "author_id"]),
    createrName: readString(
      value,
      ["createrName", "author_name"],
      "Unknown author",
    ),
    authorAvatar: readString(value, ["authorAvatar", "author_avatar"]),
    createDate: readString(value, ["createDate", "created_at", "createdAt"]),
    likesNum: readNumber(value, ["likesNum", "Likes"]),
    dislikesNum: readNumber(value, ["dislikesNum", "Dislikes"]),
    userrating: normalizeVoteStatus(value.userrating ?? value.vote_status),
    PostComments: rawComments.flatMap((comment) => {
      const normalizedComment = normalizeComment(comment);
      return normalizedComment ? [normalizedComment] : [];
    }),
  };
}
export function normalizeFeedPosts(value) {
  if (!Array.isArray(value)) return [];
  return value.flatMap((post) => {
    const normalizedPost = normalizeFeedPost(post);
    return normalizedPost ? [normalizedPost] : [];
  });
}
export function normalizePostResponse(value) {
  if (!isRecord(value)) return null;
  return normalizePost(value.post ?? value);
}
export function getNextVote(currentVote, nextVote) {
  return currentVote === nextVote ? null : nextVote;
}
export function getOptimisticVoteCounts(
  currentVote,
  nextVote,
  likes,
  dislikes,
) {
  const wasLiked = currentVote === "like";
  const wasDisliked = currentVote === "dislike";
  const isRemovingVote = currentVote === nextVote;
  return {
    userVote: getNextVote(currentVote, nextVote),
    likes:
      nextVote === "like"
        ? isRemovingVote
          ? Math.max(0, likes - 1)
          : likes + 1
        : wasLiked
          ? Math.max(0, likes - 1)
          : likes,
    dislikes:
      nextVote === "dislike"
        ? isRemovingVote
          ? Math.max(0, dislikes - 1)
          : dislikes + 1
        : wasDisliked
          ? Math.max(0, dislikes - 1)
          : dislikes,
  };
}
