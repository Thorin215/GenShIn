import request from '@/utils/request'

// 新建房地产(管理员)
export function createRealEstate(data) {
  return request({
    url: 'http://localhost:8888/api/v1/createRealEstate',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// 获取房地产信息(空json{}可以查询所有，指定proprietor可以查询指定业主名下房产)
export function queryRealEstateList(data) {
  return request({
    url: 'http://localhost:8888/api/v1/queryRealEstateList',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}
