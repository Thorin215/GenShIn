<template>
  <div class="upload-container">
    <el-form ref="uploadForm" :model="uploadForm" label-width="100px" class="upload-form" @submit.native.prevent="submitForm">
      <el-form-item label="句子" prop="sentence">
        <el-input v-model="uploadForm.sentence" placeholder="输入句子"></el-input>
      </el-form-item>

      <el-form-item label="标签" prop="label">
        <el-select v-model="uploadForm.label" placeholder="选择标签">
          <el-option :value="true" label="真"></el-option>
          <el-option :value="false" label="假"></el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="数据集编号" prop="dataset_id">
        <el-input v-model="uploadForm.dataset_id" placeholder="输入数据集编号"></el-input>
      </el-form-item>

      <el-form-item label="句子编号" prop="sentence_id">
        <el-input v-model="uploadForm.sentence_id" placeholder="输入句子编号"></el-input>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="submitForm">上传数据</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { uploadSentence } from '@/api/upload' // 确保这个路径是正确的

export default {
  name: 'UploadForm',
  data() {
    return {
      uploadForm: {
        sentence: '',
        label: null,
        dataset_id: '',
        sentence_id: ''
      },
      loading: false
    }
  },
  methods: {
    async submitForm() {
      this.$refs.uploadForm.validate(async (valid) => {
        if (valid) {
          this.loading = true
          try {
            const response = await uploadSentence(this.uploadForm)
            if (response.data) {
              this.$message.success('数据上传成功!')
            } else {
              this.$message.error('数据上传失败!')
            }
          } catch (error) {
            this.$message.error('请求失败!')
          } finally {
            this.loading = false
          }
        } else {
          this.$message.warning('请填写正确的信息!')
        }
      })
    },
    resetForm() {
      this.$refs.uploadForm.resetFields()
    }
  }
}
</script>

<style scoped>
.upload-container {
  max-width: 600px;
  margin: 0 auto;
}

.upload-form {
  padding: 20px;
}
</style>
