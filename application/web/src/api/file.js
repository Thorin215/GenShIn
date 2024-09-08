import request from '@/utils/request'
import Axios from 'axios'
import FileSaver from 'file-saver'

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

export function downloadFile(data) {
    return Axios({
        url: 'http://localhost:8888/api/v1/downloadFile',
        responseType: 'blob',
        method: 'post',
        headers: {
            'Content-Type': 'application/json'
        },
        data,
    }).then((response) => {
        FileSaver.saveAs(response.data, data.filename)
        return response
    })
}

export function downloadFilesCompressed(data) {
    return Axios({
        url: 'http://localhost:8888/api/v1/downloadFilesCompressed',
        method: 'post',
        responseType: 'blob',
        headers: {
            'Content-Type': 'application/json'
        },
        data,
    }).then((response) => {
        FileSaver.saveAs(response.data, data.zipname + ".zip")
        return response
    })
}
