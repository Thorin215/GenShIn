import request from '@/utils/request'

// 显示所有数据集
export function GetAllDataSet(){
  return request({
    url: 'http://localhost:8888/api/v1/getalldataset',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
  })
}
// 查询销售(可查询所有，也可根据发起销售人查询)(发起的)
export function queryChangeLog(data) {
  return request({
    url: 'http://localhost:8888/api/v1/getChangeLog',
    method: 'get',
    headers: {
      'Content-Type': 'application/json'
    },
    params:{
    data_set_id:data
    }
  })
}

// 根据参与销售人、买家(买家AccountId)查询销售(参与的)
export function querySellingListByBuyer(data) {
  return request({
    url: 'http://localhost:8888/api/v1/querySellingListByBuyer',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// 买家购买
export function createSellingByBuy(data) {
  return request({
    url: 'http://localhost:8888/api/v1/createSellingByBuy',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// 更新销售状态（买家确认、买卖家取消）Status取值为 完成"done"、取消"cancelled" 当处于销售中状态，卖家要取消时，buyer为""空
export function updateSelling(data) {
  return request({
    url: 'http://localhost:8888/api/v1/updateSelling',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// 发起销售
export function createSelling(data) {
  return request({
    url: 'http://localhost:8888/api/v1/createSelling',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// export function uploadSentence(data) {
//   return request({
//     url: '/uploadSentence',
//     method: 'post',
//     data
//   })
// }