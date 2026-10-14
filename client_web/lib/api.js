const configuredApiUrl = process.env.NEXT_PUBLIC_API_URL?.trim();

export const API_URL = (configuredApiUrl || "http://localhost:8081").replace(
  /\/$/,
  "",
);
export const DEFAULT_PUBLIC_PIC_PATH = `${API_URL}/socialup/posts/Public_pic.jpg`;
export const DEFAULT_AVATAR_PATH = `${API_URL}/socialup/auth/user_default_avatar.avif`;

export class ApiError extends Error {
  constructor(message, status, payload) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

export function isRecord(value) {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}


export function readString(source, keys, fallback = "") {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === "string") return value;
    if (typeof value === "number") return String(value);
  }
  return fallback;
}


export function readNumber(source, keys, fallback = 0) {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === "number" && Number.isFinite(value)) return value;
    if (typeof value === "string" && value.trim()) {
      const parsedValue = Number(value);
      if (Number.isFinite(parsedValue)) return parsedValue;
    }
  }
  return fallback;
}


export function readBoolean(source, keys, fallback = false) {
  for (const key of keys) {
    const value = source[key];
    if (typeof value === "boolean") return value;
    if (typeof value === "number") return value !== 0;
    if (typeof value === "string") {
      if (value === "true" || value === "1") return true;
      if (value === "false" || value === "0") return false;
    }
  }
  return fallback;
}

export function buildBackendAssetUrl(path, fallback = "") {
  
  const normalizedPath = path?.trim();
  if (!normalizedPath) return fallback;
  if (/^http?:\/\//i.test(normalizedPath)) return normalizedPath;
  
  if (/^\/(icons|img|sound)\//.test(normalizedPath)) return normalizedPath;
  
  const relativePath = normalizedPath.startsWith("/") ? normalizedPath : `/${normalizedPath}`;
  //console.log("THe imaaaage ggggggggggggggggggggggggggggg gggggggggg relativePath",relativePath)
    //console.log("THe imaaaage ggggggggggggggggggggggggggggg gggggggggg path",`${API_URL}${path}`)
  return `${API_URL}${relativePath}`;
}

/** Returns a safe user-facing message for caught unknown errors. */
export function getErrorMessage(error, fallback) {
  return error instanceof Error && error.message.trim()
    ? error.message
    : fallback;
}


export async function fetchfromApi(endpoint, options = {}) {
  const headers = new Headers(options.headers);
  const isFormData = options.body instanceof FormData;

  if (!isFormData && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
    credentials: "include",
  });
  // const responsejson= response
  // //console.log("fygrfyur-------",response);
  
  
  const responseText = await response.text();
  //  //console.log("the response from backend",res) 
  let responseData = null;
  
  if (responseText) {
    try {
      responseData = JSON.parse(responseText);
    } catch {
      responseData = responseText;
    }
  }
  // //console.log("the response from backend -------------",responseData) 
  // //console.log("response fromn the backend",response.json)
  if (!response.ok) {
    const payloadMessage = isRecord(responseData)
    ? readString(responseData, ["ERorrmessages", "message"], response.statusText)
    : response.statusText;
    const errorMessage = Array.isArray(payloadMessage)
    ? "The request could not be completed."
    : payloadMessage || "The request could not be completed.";
    throw new ApiError(errorMessage, response.status, responseData);
  }
  
  return responseData;
}
