import {
  login
} from '@/api/account'
import {
  getToken,
  setToken,
  removeToken
} from '@/utils/auth'
import {
  resetRouter
} from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    userId: '',
    userName: '',
    roles: []
  }
}

const state = getDefaultState()

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_ACCOUNTID: (state, userId) => {
    state.userId = userId
  },
  SET_USERNAME: (state, userName) => {
    state.userName = userName
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  }
}

const actions = {
  login({
    commit
  }, userId) {
    return new Promise((resolve, reject) => {
      login({
        args: [{
          id: userId
        }]
      }).then(response => {
        commit('SET_TOKEN', response[0].userId)
        setToken(response[0].userId)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },
  getInfo({
    commit,
    state
  }) {
    return new Promise((resolve, reject) => {
      login({
        args: [{
          userId: state.token
        }]
      }).then(response => {
        var roles
        if (response[0].userName === 'Admin') {
          roles = ['admin']
        } else {
          roles = ['admin'] //editor
        }
        commit('SET_ROLES', roles)
        commit('SET_ACCOUNTID', response[0].userId)
        commit('SET_USERNAME', response[0].userName)
        resolve(roles)
      }).catch(error => {
        reject(error)
      })
    })
  },
  logout({
    commit
  }) {
    return new Promise(resolve => {
      removeToken()
      resetRouter()
      commit('RESET_STATE')
      resolve()
    })
  },

  resetToken({
    commit
  }) {
    return new Promise(resolve => {
      removeToken()
      commit('RESET_STATE')
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
// import {
//   login
// } from '@/api/account'
// import {
//   getToken,
//   setToken,
//   removeToken
// } from '@/utils/auth'
// import {
//   resetRouter
// } from '@/router'

// const getDefaultState = () => {
//   return {
//     token: getToken(),
//     userId: '',
//     userName: '',
//     roles: []
//   }
// }

// const state = getDefaultState()

// const mutations = {
//   RESET_STATE: (state) => {
//     Object.assign(state, getDefaultState())
//   },
//   SET_TOKEN: (state, token) => {
//     state.token = token
//   },
//   SET_ACCOUNTID: (state, userId) => {
//     state.userId = userId
//   },
//   SET_USERNAME: (state, userName) => {
//     state.userName = userName
//   },
//   SET_ROLES: (state, roles) => {
//     state.roles = roles
//   }
// }

// const actions = {
//   login({
//     commit,
//     dispatch
//   }, userId) {
//     return new Promise((resolve, reject) => {
//       login({
//         args: [{
//           userId: userId
//         }]
//       }).then(response => {
//         commit('SET_TOKEN', response[0].userId)
//         setToken(response[0].userId)
//         // 在登录成功后获取用户信息
//         dispatch('getInfo').then(roles => {
//           if (roles.includes('admin')) {
//             // 如果是管理员，跳转到管理员界面
//             resolve('/admin')
//           } else {
//             // 否则跳转到默认界面
//             resolve('/')
//           }
//         }).catch(error => {
//           reject(error)
//         })
//       }).catch(error => {
//         reject(error)
//       })
//     })
//   },
//   getInfo({
//     commit,
//     state
//   }) {
//     return new Promise((resolve, reject) => {
//       login({
//         args: [{
//           userId: state.token
//         }]
//       }).then(response => {
//         var roles
//         if (response[0].userName === '管理员') {
//           roles = ['admin']
//         } else {
//           roles = ['editor']
//         }
//         commit('SET_ROLES', roles)
//         commit('SET_ACCOUNTID', response[0].userId)
//         commit('SET_USERNAME', response[0].userName)
//         resolve(roles)
//       }).catch(error => {
//         reject(error)
//       })
//     })
//   },
//   logout({
//     commit
//   }) {
//     return new Promise(resolve => {
//       removeToken()
//       resetRouter()
//       commit('RESET_STATE')
//       resolve()
//     })
//   },

//   resetToken({
//     commit
//   }) {
//     return new Promise(resolve => {
//       removeToken()
//       commit('RESET_STATE')
//       resolve()
//     })
//   }
// }

// export default {
//   namespaced: true,
//   state,
//   mutations,
//   actions
// }
