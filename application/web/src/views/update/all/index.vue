<template>
  <div class="container">
    <!-- 数据集创建表单 -->
    <el-form ref="form" class="dataset-form">
      <el-form-item label="数据集名称">
        <el-input v-model="name" placeholder="请输入数据集名称" />
      </el-form-item>
      <el-form-item label="版本说明">
        <el-input v-model="change_log" placeholder="请输入版本说明" />
      </el-form-item>
      <el-form-item label="行数">
        <el-input type="number" v-model.number="rows" placeholder="请输入行数" />
      </el-form-item>
      <el-form-item label="文件上传">
        <el-upload
          class="upload-demo"
          :show-file-list="true"
          :file-list="fileList"
          :auto-upload="false"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          action=""
          multiple
          >
          <el-button size="large" type="primary">上传文件</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="large" @click="submitForm">提交</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { addDatasetVersion } from '@/api/dataset'
import { uploadFile } from '@/api/file'

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
      owner: '', // 初始为空
      fileListx: [],            // 存储文件的哈希值与文件名
      fileList: [],            // 上传的文件列表
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
    handleFileChange(file, fileList) {
      this.fileList = fileList;

      // 定义支持的文件类型，包括图片、文本、音频、代码文件、压缩包和CSV文件
      const allowedTypes = [
        'image', // 图片
        'text',  // 文本文件
        'audio', // 音频文件
        'application/x-zip-compressed', // ZIP 压缩包
        'application/zip',              // ZIP 压缩包
        'application/x-rar-compressed', // RAR 压缩包
        'text/csv',                     // CSV 文件
        'application/csv'               // CSV 文件
      ];

      // 定义支持的文件扩展名
      const codeExtensions = ['.py', '.js', '.java', '.c', '.cpp', '.csv'];

      // 检查文件的 MIME 类型和扩展名
      const isFile = allowedTypes.some(type => file.raw.type.includes(type)) ||
                    codeExtensions.some(ext => file.name.endsWith(ext));
      if (!isFile) {
        this.$message.error('只允许上传图片、文本、音频文件、代码文件、压缩包或CSV文件。');
        file.status = 'error';
        return false;
      }            
      this.$message.success('文件类型检查通过!');


      const form = new FormData();
      form.append('file', file.raw);
      uploadFile(form).then(response => {
        this.fileListx.push({ hash: response, filename: file.name });
        this.$message.success(`文件上传成功! 文件哈希: ${response}`); 
        file.status = 'success';
      }).catch(error => {
        this.$message.error('文件上传失败!');
        file.status = 'error';
      });

      return file.status === 'success';
    },

    handleFileRemove(file, fileList) {
      console.log(file, fileList);
      this.fileList = fileList.concat(file); // disable file removal
    },

    submitForm() {
      if (!this.name) {
        this.$message({
          type: 'warning',
          message: '数据集名称不能为空!'
        });
        return;
      }

      // const filesArray = this.files.split(',').map(file => file.trim());

      // 获取当前时间并转换为 ISO 8601 格式
      const isoDateString = new Date().toISOString().substring(0, 19) + 'Z';

      const dataToSubmit = { 
        owner: this.owner,
        name: this.name,
        version: {
          creation_time: isoDateString,  // 格式化后的时间字符串
          change_log: this.change_log,
          rows: this.rows,
          files: this.fileListx,
        },
      };

      addDatasetVersion(dataToSubmit)
        .then(response => {
          this.$message({
            type: 'success',
            message: '版本新增成功!'
          });
        })
        .catch(error => {
          this.$message({
            type: 'error',
            message: '版本新增失败!'
          });
        });
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
.upload-demo i {
  font-size: 28px;
  color: #409EFF;
}
</style>
