const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.account.token,
  userId: state => state.account.userId,
  userName: state => state.account.userName,
  roles: state => state.account.roles,
  permission_routes: state => state.permission.routes
}
export default getters
