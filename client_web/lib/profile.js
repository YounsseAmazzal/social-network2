import {
  fetchfromApi,
  isRecord,
  readBoolean,
  readNumber,
  readString,
} from "../lib/api.js";
function normalizeProfileUser(value) {
  const user = isRecord(value) ? value : {};
  return {
    id: readString(user, ["id", "ID"]),
    firstname: readString(user, ["firstname", "first_name", "Firstname"]),
    lastname: readString(user, ["lastname", "last_name", "Lastname"]),
    nickname: readString(user, ["nickname", "Nickname"]),
    email: readString(user, ["email", "Email"]),
    age: readNumber(user, ["age", "Age"]),
    gender: readString(user, ["gender", "Gender"]),
    aboutMe: readString(user, ["aboutMe", "about_me"]),
    avatarUrl: readString(user, [
      "avatar_url",
      "avatarURL",
      "avatar",
      "AvatarURL",
    ]),
  };
}
function normalizeProfilePost(value) {
  if (!isRecord(value)) return null;
  const id = readString(value, ["id", "ID"]);
  if (!id) return null;
  return {
    id,
    title: readString(value, ["title", "Title"], "Untitled"),
    content: readString(value, ["content", "Content", "description"]),
    likes: readNumber(value, ["likes", "Likes","likes_count"]),
    comments: readNumber(value, ["comments", "Comments","comments_count"]),
    date: readString(value, ["date", "Date", "createDate"]),
    avatar_url: readString(value, ["avatar_url", "Image","image_url"]),
  };
}
function normalizeProfileDetails(value) {
  const profile = isRecord(value) ? value : {};
  //console.log("^^^^^^^^^^isOwner^^^^^^^^^ ",profile.isOwner)
  const user = normalizeProfileUser(profile.user);
  // //console.log("the useeer normazeee",user)
  const posts = Array.isArray(profile.posts)
    ? profile.posts.flatMap((post) => {
        const normalizedPost = normalizeProfilePost(post);
        return normalizedPost ? [normalizedPost] : [];
      })
    : [];
  return {
    user,
    isPrivate: readBoolean(profile, ["isPrivate", "is_private"]),
    myAccount: readBoolean(profile, ["is_Owner", "isOwner"]),
    canSee: readBoolean(profile, ["canSee", "can_See"]),
    followers: readNumber(profile, ["followers", "Followers"]),
    following: readNumber(profile, ["following", "Following"]),
    postsCount: readNumber(profile, ["postsCount", "posts_count"]),
    posts,
  };
}
function normalizePublicUser(value) {
  if (!isRecord(value)) return null;
  const id = readString(value, ["id", "ID"]);
  if (!id) return null;
  const rawRelationship = value.relationship;
  const relationship =
    rawRelationship === "following" ||
    rawRelationship === "follower" ||
    rawRelationship === "mutual" ||
    rawRelationship === "requested"
      ? rawRelationship
      : "none";
  return {
    id,
    firstname: readString(value, ["firstname", "first_name", "Firstname"]),
    lastname: readString(value, ["lastname", "last_name", "Lastname"]),
    nickname: readString(value, ["nickname", "Nickname"]),
    avatar: readString(value, ["avatar", "avatar_url", "avatarURL"]),
    aboutMe: readString(value, ["aboutMe", "about_me"]),
    isPrivate: readBoolean(value, ["isPrivate", "is_private"]),
    relationship,
  };
}
export function normalizeProfileResponse(value) {
  if (!isRecord(value) || !isRecord(value.Seeprofile)) return null;

  const rawData = value.Seeprofile;
  const rawUser = rawData.user || {};

  return {
    isPrivate: readBoolean(rawData, ["isPrivate", "is_private"]),
    isOwner: readBoolean(rawData, ["isOwner", "is_owner"]),
    isFollowing: readBoolean(rawData, ["isFollowing", "is_following"]),
    canSee: readBoolean(rawData, ["canSee", "can_see"]),
    
    relationship: readString(rawUser, ["Rolation", "followStatus", "relationship"], "none"),
    followersCount:readString(rawData, ["followersCount", "followers_Count"]),
    followingCount:readString(rawData, ["following_Count", "followingCount"]),
    postsCount:readString(rawData, ["postsCount", "posts_Count"]),
    profile: normalizeProfileDetails(rawData),
  };
}

export function normalizePublicUsers(value) {
  if (!Array.isArray(value)) return [];
  return value.flatMap((item) => {
    const user = normalizePublicUser(item);
    return user ? [user] : [];
  });
}
/** Loads and normalizes either the signed-in or requested profile. */
export async function fetchProfile(profileId) {
  const endpoint = profileId
    ? `/UserProfile?user_id=${(profileId)}`
    : "/UserProfile";
    console.log("the backend endpoint",endpoint)
    //console.log("the response of the profile", normalizeProfileResponse(await fetchfromApi(endpoint)))
  return normalizeProfileResponse(await fetchfromApi(endpoint));
}
export async function fetchProfileConnections(userId) {
  const [followersResponse, followingResponse] = await Promise.all([
    fetchfromApi(`/Userfollowers?useeer_id=${userId}`),
    fetchfromApi(`/UserfollowinngList?useeer_id=${userId}`),
  ]);
  console.log("the response of User Following",followersResponse)
  return {
    followers: normalizePublicUsers(followersResponse?.data|| []),
    following: normalizePublicUsers(followingResponse?.data|| []),
  };
}
export async function fetchPendingFollowRequests() {
  const response = await fetchfromApi("/followers/Pendiing");
  // console.log("the requuuuets of the folllowwwwing",response)
  if (Array.isArray(response.message)) {
    return normalizePublicUsers(response.message);
  }
  if (isRecord(response.message)) {
    return normalizePublicUsers(response.message.followers);
  }
  return [];
}
