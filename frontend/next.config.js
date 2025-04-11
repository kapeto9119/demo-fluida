/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // CSS configuration is simpler now that we use local files
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        // TEMPORARY FIX: Hardcode the production API URL
        destination: `https://serene-radiance-production.up.railway.app/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig; 