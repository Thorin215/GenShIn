import {
  queryUser
} from '@/api/user'
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
  SET_USERID: (state, userId) => {
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
      queryUser({
        id: userId
      }).then(response => {
        commit('SET_TOKEN', response.id)
        setToken(response.id)
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
      queryUser({
        id: state.token
      }).then(response => {
        var roles
        if (response.id === 'admin') {
          roles = ['admin']
        } else {
          roles = ['editor']
        }
        commit('SET_ROLES', roles)
        commit('SET_USERID', response.id)
        commit('SET_USERNAME', response.name)
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
