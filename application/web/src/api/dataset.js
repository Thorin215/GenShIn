import request from '@/utils/request'

export function queryAllDatasets(){
  return request({
    url: 'http://localhost:8888/api/v1/dataset/all',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
  })
}

export function createDataset(data){
  return request({
    url: 'http://localhost:8888/api/v1/dataset/create',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function queryDatasetMetadata(data){
  return request({
    url: 'http://localhost:8888/api/v1/dataset/metadata',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function addDatasetVersion(data){
  return request({
    url: 'http://localhost:8888/api/v1/dataset/version/create',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function queryAllVersions(data){
  return request({
    url: 'http://localhost:8888/api/v1/dataset/version/all',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}