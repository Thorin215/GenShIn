import request from '@/utils/request'

export function queryAllDatasets(){
  return request({
    url: 'http://localhost:8888/api/v1/queryAllDatasets',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
  })
}

export function createDataset(data){
  return request({
    url: 'http://localhost:8888/api/v1/createDataset',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function queryDatasetMetadata(data){
  return request({
    url: 'http://localhost:8888/api/v1/queryDatasetMetadata',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

export function addDatasetVersion(data){
  return request({
    url: 'http://localhost:8888/api/v1/addDatasetVersion',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}