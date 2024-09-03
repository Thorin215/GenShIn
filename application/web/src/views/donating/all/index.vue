<template>
  <div class="container">
    <!-- <el-alert
      type="success"
      title="账户信息"
    >
      <p>账户ID: {{ userId }}</p>
      <p>用户名: {{ userName }}</p>
       <p>余额: ￥{{ balance }} 元</p> 
    </el-alert> -->
    
    <!-- 新添加的表单部分 -->
    <el-form :model="datasetVersion" ref="form" class="dataset-form">
      <el-form-item label="创建时间">
        <el-input v-model="datasetVersion.creationTime" placeholder="请输入创建时间" />
      </el-form-item>
      <el-form-item label="版本说明">
        <el-input v-model="datasetVersion.changeLog" placeholder="请输入版本说明" />
      </el-form-item>
      <el-form-item label="行数">
        <el-input type="number" v-model.number="datasetVersion.rows" placeholder="请输入行数" />
      </el-form-item>
      <el-form-item label="文件哈希列表">
        <el-input v-model="datasetVersion.files" placeholder="用逗号分隔的文件哈希" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="large" @click="submitForm">提交</el-button>
      </el-form-item>
    </el-form>

    <!-- <div v-if="donatingList.length === 0" class="no-data">
      <el-alert
        title="查询不到数据"
        type="warning"
      />
    </div> -->
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { queryDonatingList, updateDonating } from '@/api/donating'

export default {
  name: 'AllDonating',
  data() {
    return {
      loading: true,
      donatingList: [],
      datasetVersion: {
        creationTime: '',
        changeLog: '',
        rows: 0,
        files: ''
      }
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
    queryDonatingList().then(response => {
      if (response !== null) {
        this.donatingList = response
      }
      this.loading = false
    }).catch(_ => {
      this.loading = false
    })
  },
  methods: {
    updateDonating(item, type) {
      let tip = ''
      if (type === 'done') {
        tip = '确认接受捐赠'
      } else {
        tip = '取消捐赠操作'
      }
      this.$confirm('是否要' + tip + '?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'success'
      }).then(() => {
        this.loading = true
        updateDonating({
          donor: item.donor,
          grantee: item.grantee,
          objectOfDonating: item.objectOfDonating,
          status: type
        }).then(response => {
          this.loading = false
          if (response !== null) {
            this.$message({
              type: 'success',
              message: tip + '操作成功!'
            })
          } else {
            this.$message({
              type: 'error',
              message: tip + '操作失败!'
            })
          }
          setTimeout(() => {
            window.location.reload()
          }, 1000)
        }).catch(_ => {
          this.loading = false
        })
      }).catch(() => {
        this.loading = false
        this.$message({
          type: 'info',
          message: '已取消' + tip
        })
      })
    },
    submitForm() {
      const filesArray = this.datasetVersion.files.split(',').map(file => file.trim())
      const dataToSubmit = { ...this.datasetVersion, files: filesArray }

      // 处理提交逻辑，例如发送数据到后端
      console.log('提交的数据:', dataToSubmit)
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


