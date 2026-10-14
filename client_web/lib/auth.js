import { isRecord, readBoolean, readNumber, readString } from "@/lib/api";
export function normalizeAuthUser(value) {
  if (!isRecord(value)) {
    throw new Error("The server returned an invalid user.");
  }
  const nickname = readString(value, ["nickname"], "").trim();
  const username = readString(value, ["username"], nickname);
  return {
    id: readString(value, ["id", "ID"]),
    nickname,
    username,
    firstname: readString(value, ["first_name", "firstName", "Firstname"]),
    lastname: readString(value, ["last_name", "lastName", "Lastname"]),
    age: readNumber(value, ["age", "Age"]),
    gender: readString(value, ["gender", "Gender"]),
    email: readString(value, ["email", "Email"]),
    avatar: readString(value, [
      "avatarURL",
      "avatar",
      "Avatar",
    ]),
    aboutMe: readString(value, ["aboutMe", "about_me"]),
    isPrivate: readBoolean(value, ["isPrivate", "is_private"]),
  };
}
