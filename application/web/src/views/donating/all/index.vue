<template>
  <div class="container">
    <!-- 数据集创建表单 -->
    <el-form :model="datasetForm" ref="form" class="dataset-form">
      <el-form-item label="数据集名称">
        <el-input v-model="name" placeholder="请输入数据集名称" />
      </el-form-item>
      <el-form-item label="版本说明">
        <el-input v-model="change_log" placeholder="请输入版本说明" />
      </el-form-item>
      <el-form-item label="行数">
        <el-input type="number" v-model.number="rows" placeholder="请输入行数" />
      </el-form-item>
      <el-form-item label="文件哈希列表">
        <el-input v-model="files" placeholder="用逗号分隔的文件哈希" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="large" @click="submitForm">提交</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
// import { queryDonatingList, updateDonating } from '@/api/donating'
import { updateVersion } from '@/api/upload'

export default {
  name: 'AllDonating',
  data() {
    return {
      loading: true,
      donatingList: [],
      name: '', // 数据集名称
      change_log: '',
      rows: 0,
      files: '',
      owner: '' // 初始为空
    }
  },
  computed: {
    ...mapGetters([
      'userId',
      'roles',
      'userName'
    ])
  },
  created() {
    this.owner = this.userId // 在 created 钩子中设置 Owner
  },
  methods: {
    submitForm() {
      if (!this.name) {
        this.$message({
          type: 'warning',
          message: '数据集名称不能为空!'
        })
        return
      }

      const filesArray = this.files.split(',').map(file => file.trim())

      // 获取当前时间并转换为 ISO 8601 格式
      const isoDateString = new Date().toISOString().substring(0, 19) + 'Z'

      const dataToSubmit = { 
        name: this.name,
        owner: this.owner,
        creation_time: isoDateString,  // 格式化后的时间字符串
        change_log: this.change_log,
        rows: this.rows,
        files: filesArray
      }


      // Use imported updateVersion API function
      updateVersion(dataToSubmit)
        .then(response => {
          if (response.dataset_name === this.name) {
            this.$message({
              type: 'success',
              message: '版本更新成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: '版本更新失败!'
            })
          }
        })
        .catch(error => {
          this.$message({
            type: 'error',
            message: '网络错误!'
          })
          console.error('提交数据错误:', error)
        })
    }
  }
}
</script>

<style>
.container {
  width: 60%;
  margin: 20px auto;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  background-color: #fff;
}
.el-alert {
  margin-bottom: 20px;
}
.dataset-form {
  max-width: 600px;
  margin: 0 auto;
}
.el-form-item {
  margin-bottom: 20px;
}
.el-button {
  width: 100%;
}
.no-data {
  text-align: center;
}
</style>
