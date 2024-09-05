import request from '@/utils/request'
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

// uploadSet 函数用于将用户ID和数据集ID上传到服务器
// export function uploadFile(data) {
//   return request({
//     url: 'http://localhost:8888/api/v1/uploadFile',
//     method: 'post',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     data
//   })
// }

export function uploadFile(formData) {
  return request({
    url: 'http://localhost:8888/api/v1/uploadFile',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data: formData
  })
}

export function downloadDataset(data) {
  return request({
    url: 'http://localhost:8888/api/v1/downloadDataset',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data,
    responseType: 'json'  // 设置响应类型为 json
  });
}

