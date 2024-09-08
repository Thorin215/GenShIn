import request from '@/utils/request'

export function queryAllUsers() {
  return request({
    url: 'http://localhost:8888/api/v1/queryAllUsers',
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'post'
  })
}

export function queryUser(data) {
  return request({
    url: 'http://localhost:8888/api/v1/queryUser',
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'post',
    data
  })
}
