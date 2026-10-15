const TIME_INTERVALS = [
  { unit: "year", seconds: 31536000 },
  { unit: "month", seconds: 2592000 },
  { unit: "week", seconds: 604800 },
  { unit: "day", seconds: 86400 },
  { unit: "hour", seconds: 3600 },
  { unit: "minute", seconds: 60 },
  { unit: "second", seconds: 1 },
];

const relativeTimeFormatter = new Intl.RelativeTimeFormat("en", {
  numeric: "auto", ///ghaafi icovete "yesteday" directtly l 1 day ago
});


export function describeRelativeTimestamp(value) {
  if (!value) return "";
  const date = new Date(value);
  
  if (Number.isNaN(date.getTime())) return "Invalid date";

  const elapsedSeconds = Math.floor((Date.now() - date.getTime()) / 1000);

  if (elapsedSeconds < 5) return "just now";

  for (const interval of TIME_INTERVALS) {
    if (elapsedSeconds >= interval.seconds) {
      const count = Math.floor(elapsedSeconds / interval.seconds);
      return relativeTimeFormatter.format(-count, interval.unit);
    }
  }

  return "just now";
}
export function formatPostTimestamp(value) {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime()) || date.getFullYear() <= 1) return "";

  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

/** Formats the calendar date and clock time shown for a group event. */
export function formatGroupEventSchedule(value) {
  if (!value) return "Date pending";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("en", {
    weekday: "short",
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

/** Formats only the clock time displayed beside a chat message. */
export function formatChatTimestamp(value) {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";

  return new Intl.DateTimeFormat("en", {
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}
