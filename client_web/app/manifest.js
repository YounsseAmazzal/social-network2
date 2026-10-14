export default function createWebAppManifest() {
  return {
    name: "Social Network",
    short_name: "Social",
    description: "Connect, share posts, and join group conversations.",
    start_url: "/",
    display: "standalone",
    background_color: "#f6f6fb",
    theme_color: "#5b3df5",
    icons: [
      {
        src: "",
        sizes: "512x512",
        type: "image/jpeg",
      },
    ],
  };
}
