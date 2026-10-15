import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Emits .next/standalone/server.js so the Docker image only ships the files
  // the production server actually needs.
  output: "standalone",
};

export default nextConfig;
