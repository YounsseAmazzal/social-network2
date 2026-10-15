import {
  ApiError,
  fetchfromApi,
  isRecord,
  readNumber,
  readString,
} from "@/lib/api";
function asPayloadRecord(value) {
  return isRecord(value) ? value : {};
}
function unwrapGroupPayload(response) {
  if (!isRecord(response)) return {};
  return isRecord(response.Groupdetails) ? response.Groupdetails : response;
}
/** Extracts list payloads from the response shapes used by group endpoints. */
export function extractGroupCollection(response) {
  if (Array.isArray(response)) return response;
  
  if (!isRecord(response)) return [];

  if (Array.isArray(response.data)) return response.data;

  const keys = ["Discovergroup", "joineddgroup", "Groupsmembers", "JoinedGroups", "FeedGroups", "members","GroupPosts","GroupEvents","FeedPosts","JoinGrouplist","InvidetUserList"];
  for (const key of keys) {
    if (Array.isArray(response[key])) return response[key];
  }

  return [];
}

export async function fetchWithFallback(endpoints, options = {}) {
  const { allowedFallbackStatuses = [404, 405], ...requestOptions } = options;
  let lastError;
  for (const endpoint of endpoints) {
    try {
      return await fetchfromApi(endpoint, requestOptions);
    } catch (error) {
      lastError = error;
      if (
        !(error instanceof ApiError) ||
        !allowedFallbackStatuses.includes(error.status)
      ) {
        throw error;
      }
    }
  }
  throw lastError instanceof Error
    ? lastError
    : new Error("No group endpoint accepted the request.");
}
function normalizeGroupAuthor(value) {
  const author = asPayloadRecord(value);
  return {
    id: readString(author, ["id", "ID"]),
    nickname: readString(author, ["nickname", "Nickname"]),
    firstName: readString(author, ["firstname", "first_name", "Firstname"]),
    lastName: readString(author, ["lastname", "last_name", "Lastname"]),
    avatarUrl: readString(author, [
      "avatar",
      "avatar_url",
      "avatarURL",
      "Avatar",
      "AvatarURL",
    ]),
  };
}
function normalizeGroup(value) {
  const group = asPayloadRecord(value);
  return {
    id: readString(group, ["id", "ID"]),
    ownerId: readString(group, ["creator_id", "userId", "UserID"]),
    title: readString(group, ["group_title", "Title"]),
    description: readString(group, ["description", "Description"]),
    createdAt: readString(group, ["created_at", "createdAt", "CreateDate"]),
    requestId: readNumber(group, ["request_id", "requestId", "RequestID"]),
      request_status: readString(group, ["request_status", "requeststatus", "Request_status"]),
  };
}
function normalazeGroupEvent(value) {
  const event = asPayloadRecord(value);
  const rawVote = readString(event, ["vote", "Vote", "user_rate"]);
  const userVote =
    rawVote === "going" || rawVote === "not going" ? rawVote : "";
  return {
    id: readString(event, ["id", "ID"]),
    groupId: readString(event, ["group_id", "groupId", "GroupId"]),
    title: readString(event, ["title", "Title"]),
    description: readString(event, [
      "description",
      "Description",
      "descreption",
    ]),
    event_date: readString(event, ["event_date", "eventDate", "EventDate"]),
    created_at: readString(event, ["created_at", "createdAt", "CreateDate"]),
    TotalParticipate: readNumber(event, ["total_participate", "TotalParticipate", "Total_Participate"]),
    totalNotparticipate: readNumber(event, [
      "total_Not_participate",
      "totalNotparticipate",
      "TotalNotParticipate",
    ]),
    userVote,
    Creator: normalizeGroupAuthor(event.Creator ?? event.creator),
  };
}
function normalizeGroupMember(value) {
  const member = value
  return {
    id: readString(member, ["id", "ID"]),
    nickname: readString(member, ["nickname", "Nickname"]),
    firstName: readString(member, ["firstname", "first_name", "Firstname"]),
    lastName: readString(member, ["lastname", "last_name", "Lastname"]),
    avatar_url: readString(member, [
      "avatar",
      "avatar_url",
      "avatarURL",
      "Avatar",
      "AvatarURL",
    ]),
  };
}
function normalizeGroupPost(value) {
  const post = asPayloadRecord(value);
  const authorValue = post.creater ?? post.Author ?? post;
  //console.log("the post creator ",authorValue)
  return {
    id: readString(post, ["id", "ID"]),
    title: readString(post, ["title", "Title"]),
    description: readString(post, ["description", "Description", "content"]),
    createDate: readString(post, ["createDate", "created_at", "createdAt"]),
    totalComments: readNumber(post, [
      "Commentsnum",
      "totalComments",
      "numberOfComments",
      "number_of_comments",
      "comments",
    ]),
    likes: readNumber(post, ["likes", "likesNum", "likesNum"]),
    image: readString(post, ["image_url", "Image", "image", "ImageURL"]),
    creater: normalizeGroupAuthor({
      ...asPayloadRecord(authorValue),
      nickname:
        normalizeGroupAuthor(authorValue).nickname ||
        readString(post, ["createrName", "authorName", "AuthorName"]),
      avatar_url:
        normalizeGroupAuthor(authorValue).avatarUrl ||
        readString(post, ["avatar_url", "AuthorAvatar"]),
    }),
  };
}
// function normalizeVote(value) {
//   return value === "going" || value === "not going" || value === "remove"
//     ? value
//     : "remove";
// }
export async function fetchAllJoinnedGroups() {
  const response = await fetchWithFallback(["/JoinedGroups","/GroupJoined"]);
  //console.log("The response of joinedd groups",response)
  return extractGroupCollection(response).map(normalizeGroup);
}
export async function fetchAllgroupsSuggested() {
  const response = await fetchWithFallback(["/DiscoverGroups"]);
  //console.log("reeeeeeeesponse of groups",extractGroupCollection(response).map(normalizeGroup))
  return extractGroupCollection(response).map(normalizeGroup);
}
export async function createGroup(payload) {
  const response = await fetchWithFallback(["/CreateGroup", "/groupsCreate"], {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return normalizeGroup(unwrapGroupPayload(response));
}
export async function createGroupJoinRequest(groupId) {
  //console.log("the grouuuuup id g",groupId)
  await fetchWithFallback([`/RequesttoJoinGroup?Grouup_id=${groupId}`], {
    method: "POST",
    // body: JSON.stringify({ }),
  });
}
export async function fetchGroupAllInfo(groupId) {
  const response = await fetchWithFallback([`/SeeeTheGroup?Grouup_id=${groupId}`]);
  //console.log("the response of the GROuuup info ",response)
  const payload = unwrapGroupPayload(response.Groupdetails||response);
  //console.log("the pyload of the GROuuup info ",payload)
  return {
    group: normalizeGroup(payload.Groupdetails ?? payload),
    creator: normalizeGroupAuthor(payload.creater ?? payload.Creater),
    totalMembers: readNumber(payload, ["members_count", "Members_Count"]),
  };
}
export async function fetchallEventsGroup(groupId,limit,offset) {
  //console.log("!!!!!!!!!!!!!!!!he entered here!!!!!!!!!!")
  const response = await fetchWithFallback([
    `/GroupEvents?Grouup_id=${groupId}&limit=${limit}&offset=${offset}`,
  ]);
  //console.log("the group events response",response)
  return extractGroupCollection(response).map(normalazeGroupEvent);
}
export async function groupEventCreation(groupId, payload) {
  await fetchWithFallback([`/GroupCreateEvent?Grouup_id=${groupId}`], {
    method: "POST",
    body: JSON.stringify({
      group_id: String(groupId),
      title: payload.title,
      description: payload.description,
      event_date: payload.event_date,
    }),
  });
}
export async function participateOnEventVote(groupId, eventId, vote) {
  // //console.log("!!!!!!!!! poarticipaltion",groupId)
  // //console.log("!!!!!!!!! poarticipaltion",eventId)
  // //console.log("!!!!!!!!! poarticipaltion",vote)

  await fetchWithFallback([`/GroupEventParticipate?Grouup_id=${groupId}&Event_id=${eventId}&vote=${vote}`], {
    method: "POST",
    // body: JSON.stringify({ id: eventId, vote: normalizeVote(vote) }),
  });
}
export async function fetchGroupPending(groupId) {
  //console.log("joining group requests response")
  const response = await fetchWithFallback([`/JoiningGroupRequests?Grouup_id=${groupId}`]);
  //console.log("joining group requests response",response)
  return extractGroupCollection(response, "members").map(normalizeGroupMember);
}
export async function fetchallgroupmemberes(groupId) {
  const response = await fetchWithFallback([
    `/GrouupMembers?Grouup_id=${groupId}`,
  ]);
  const members = extractGroupCollection(response);
  //console.log("!!!!!!!!!!!members",members.map(normalizeGroupMember))
  return members.map(normalizeGroupMember);
  // return extractGroupCollection(response, "members").map(normalizeGroupMember);
}
export async function fetchallpostsgroup(groupId,limit,offset) {
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")
  // //console.log("he entred here ??????????????")

  const response = await fetchWithFallback([`/GetGroupPosts?Grouup_id=${groupId}&limit=${limit}&offset=${offset}`], {
    // method: "GET",
    // body: JSON.stringify({ offset: 0, limit: 20 }),
  });
  // //console.log("Grouuups Posts Response",response)
  // //console.log("Grouuups Posts Response",response)
  // //console.log("Grouuups Posts Response",response)
  // //console.log("Grouuups Posts Response",response)
  console.log("Grouuups Posts Response",extractGroupCollection(response).map(normalizeGroupPost))
  return extractGroupCollection(response).map(normalizeGroupPost);
}
export async function respondToGroupRequest(groupId, userId, decision) {
  await fetchWithFallback([`/RespondtoJoinRequest?Grouup_id=${groupId}&useeer_id=${userId}&type=${decision}`], {
    method: "POST",
    // body: JSON.stringify({
    //   group_id: String(groupId),
    //   user_id: userId,
    //   type: decision,
    // }),
  });
}
export async function PosstGroupcreate(groupId, formData) {
  await fetchWithFallback([`/CreatePostGroup?Grouup_id=${groupId}`], {
    method: "POST",
    body: formData,
  });
}
export async function FetchFolllowersToinvite(groupId) {
  const response = await fetchWithFallback([
    `/InviteUserstogroup?Grouup_id=${groupId}`,
  ]);
  return extractGroupCollection(response, "users").map(normalizeGroupMember);
}
export async function InviteUsertogroup(groupId, userId) {
  await fetchWithFallback([`/InviteGroup?Grouup_id=${groupId}&InvetUser_id=${userId}`], {
    method: "POST",
    // body: JSON.stringify({}),
  });
}
export async function InvitaiongroupResponse(groupId, decision) {
  await fetchWithFallback([`/AcceptGroupInvitation?Grouup_id=${groupId}&type=${decision}`], {
    method: "POST",
    body: JSON.stringify({}),
  });
}
