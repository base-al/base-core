module.exports = {
  apps: [
    {
      name: "base",
      script: "./base",
      instances: "1",
      timestamp: "YYYY-MM-DD HH:mm:ss Z",
      log: "/var/log/base/base_api.log",
    },
  ],
};
