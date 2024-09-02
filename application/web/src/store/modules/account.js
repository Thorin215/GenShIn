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
        commit('SET_TOKEN', response[0].id)
        setToken(response[0].id)
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
          id: state.token
          // id: userId
        }]
      }).then(response => {
        var roles
        if (response[0].name === 'Admin') {
          roles = ['admin']
        } else {
          roles = ['admin'] //editor
        }
        commit('SET_ROLES', roles)
        commit('SET_ACCOUNTID', response[0].id)
        commit('SET_USERNAME', response[0].name)
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

// import { login } from '@/api/account';
// import { getToken, setToken, removeToken } from '@/utils/auth';
// import { resetRouter } from '@/router';

// const getDefaultState = () => {
//   return {
//     token: getToken(),
//     userId: '',
//     userName: '',
//     roles: []
//   };
// };

// const state = getDefaultState();

// const mutations = {
//   RESET_STATE: (state) => {
//     Object.assign(state, getDefaultState());
//   },
//   SET_TOKEN: (state, token) => {
//     state.token = token;
//   },
//   SET_ACCOUNTID: (state, userId) => {
//     state.userId = userId;
//   },
//   SET_USERNAME: (state, userName) => {
//     state.userName = userName;
//   },
//   SET_ROLES: (state, roles) => {
//     state.roles = roles;
//   }
// };

// const actions = {
//   login({ commit }, userId) {
//     return new Promise((resolve, reject) => {
//       login({
//         args: [{ id: userId }]
//       }).then(response => {
//         const data = response.data; // Access the data array
//         const user = data.find(user => user.id === userId); // Find the specific user

//         if (user) {
//           commit('SET_TOKEN', user.id);
//           setToken(user.id);
//           resolve();
//         } else {
//           reject('User not found');
//         }
//       }).catch(error => {
//         reject(error);
//       });
//     });
//   },

//   getInfo({ commit, state }) {
//     return new Promise((resolve, reject) => {
//       login({
//         args: [{ id: state.token }]
//       }).then(response => {
//         const data = response.data; // Access the data array
//         const user = data.find(user => user.id === state.token); // Find the specific user

//         if (user) {
//           var roles;
//           if (user.name === 'Admin') {
//             roles = ['admin'];
//           } else {
//             roles = ['editor']; // Adjust the role based on the user
//           }
//           commit('SET_ROLES', roles);
//           commit('SET_ACCOUNTID', user.id);
//           commit('SET_USERNAME', user.name);
//           resolve(roles);
//         } else {
//           reject('User not found');
//         }
//       }).catch(error => {
//         reject(error);
//       });
//     });
//   },

//   logout({ commit }) {
//     return new Promise(resolve => {
//       removeToken(); // Must remove token first
//       resetRouter(); // Then reset the router
//       commit('RESET_STATE'); // Reset the state after logout
//       resolve();
//     });
//   },

//   resetToken({ commit }) {
//     return new Promise(resolve => {
//       removeToken();
//       commit('RESET_STATE');
//       resolve();
//     });
//   }
// };

// export default {
//   namespaced: true,
//   state,
//   mutations,
//   actions
// };

