export default {
  CONFIG: '/v1/kolide/config',
  FORGOT_PASSWORD: '/v1/kolide/forgot_password',
  HOSTS: '/v1/kolide/hosts',
  INVITES: '/v1/kolide/invites',
  LABEL_HOSTS: (id) => {
    return `/v1/kolide/labels/${id}/hosts`;
  },
  LOGIN: '/v1/kolide/login',
  LOGOUT: '/v1/kolide/logout',
  ME: '/v1/kolide/me',
  QUERIES: '/v1/kolide/queries',
  RESET_PASSWORD: '/v1/kolide/reset_password',
  TARGETS: '/v1/kolide/targets',
  USERS: '/v1/kolide/users',
};
