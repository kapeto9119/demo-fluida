/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // CSS configuration is simpler now that we use local files
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig; 