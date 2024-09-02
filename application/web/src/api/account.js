import request from '@/utils/request'

// // 获取登录界面角色选择列表
// export function queryAccountList() {
//   return request({
//     url: 'http://localhost:8888/api/v1/queryAccountList',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     method: 'post'
//   })
// }

// // 登录
// export function checkAccount() {
//   return request({
//     url: 'http://localhost:8888/api/v1/checkAccount',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     method: 'post'
//   })
// }
// import request from '@/utils/request'

// 获取登录界面角色选择列表
export function queryAccountList() {
  return request({
    url: 'http://localhost:8888/api/v1/queryAccountList',
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'post'
  })
}

// // 登录
// export function checkAccount(accountId) {
//   return request({
//     url: 'http://localhost:8888/api/v1/checkAccount',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     method: 'post',
//     data: { id: accountId }  // Pass accountId as part of the request body
//   })
// }

// // 登录
export function login(data) {
  return request({
    url: 'http://localhost:8888/api/v1/queryAccountList',
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'post',
    data
  })
}

// export function login(data) {
//   return request({
//     url: 'http://localhost:8888/api/v1/checkAccount',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     method: 'post',
//     data
//   })
// }