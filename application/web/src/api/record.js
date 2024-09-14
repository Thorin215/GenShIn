import request from '@/utils/request'

export function queryRecordsByUser(data) {
  return request({
    url: 'http://localhost:8888/api/v1/record/by/user',
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'post',
    data
  })
  
}