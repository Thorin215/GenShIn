import request from '@/utils/request'

export function uploadSentence(data) {
  return request({
    url: 'http://localhost:8888/api/v1/uploadSentence',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// uploadSet 函数用于将用户ID和数据集ID上传到服务器
export function uploadSet(data) {
  return request({
    url: 'http://localhost:8888/api/v1/uploadSet',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function updateVersion(data) {
  return request({
    url: 'http://localhost:8888/api/v1/updateVersion',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function getDatasetMetadata(data) {
  return request({
    url: 'http://localhost:8888/api/v1/getDatasetMetadata',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}