<template>
  <div class="container">
    <!-- 数据集创建部分 -->
    <div class="dataset-container">
      <el-input
        v-model="newDatasetName"
        placeholder="请输入要创建的数据集名称"
        class="dataset-input"
      />

      <el-input
        v-model.number="rows"
        type="number"
        placeholder="请输入数据集行数"
        class="dataset-input"
      />

      <el-select
        v-model="metadata.tasks"
        placeholder="请选择任务"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="task in availableTasks"
          :key="task"
          :label="task"
          :value="task"
        />
      </el-select>

      <el-select
        v-model="metadata.modalities"
        placeholder="请选择数据模态"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="modality in availableModalities"
          :key="modality"
          :label="modality"
          :value="modality"
        />
      </el-select>

      <el-select
        v-model="metadata.formats"
        placeholder="请选择文件格式"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="format in availableFormats"
          :key="format"
          :label="format"
          :value="format"
        />
      </el-select>

      <el-select
        v-model="metadata.sub_tasks"
        placeholder="请选择子任务"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="subTask in availableSubTasks"
          :key="subTask"
          :label="subTask"
          :value="subTask"
        />
      </el-select>

      <el-select
        v-model="metadata.languages"
        placeholder="请选择语言"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="language in availableLanguages"
          :key="language"
          :label="language"
          :value="language"
        />
      </el-select>

      <el-select
        v-model="metadata.libraries"
        placeholder="请选择适用库"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="library in availableLibraries"
          :key="library"
          :label="library"
          :value="library"
        />
      </el-select>

      <el-select
        v-model="metadata.tags"
        placeholder="请选择标签"
        class="dataset-select"
        multiple
      >
        <el-option
          v-for="tag in availableTags"
          :key="tag"
          :label="tag"
          :value="tag"
        />
      </el-select>

      <el-select
        v-model="metadata.license"
        placeholder="请选择许可证"
        class="dataset-select"
      >
        <el-option
          v-for="license in availableLicenses"
          :key="license"
          :label="license"
          :value="license"
        />
      </el-select>

      <el-upload
        class="upload-demo"
        :before-upload="beforeUpload"
        :show-file-list="true"
        @success="handleUploadSuccess"
      >
        <el-button type="primary">上传文件</el-button>
      </el-upload>

      <el-button
        type="primary"
        @click="createDataset"
        class="dataset-button"
        :loading="loading"
      >
        创建数据集
      </el-button>
      <p v-if="createdDatasetName" class="dataset-info">
        创建成功！数据集名称: {{ createdDatasetName }}, 账户ID: {{ userId }}
      </p>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { uploadSet, uploadFile } from '@/api/upload';  // Assuming you have a method to create dataset

export default {
  name: 'DataSets',
  data() {
    return {
      loading: false,
      newDatasetName: '',      // 用户输入的新数据集名称
      createdDatasetName: '',  // 返回的创建的数据集名称
      rows: 0,                 // 数据集行数
      metadata: {
        tasks: [],
        modalities: [],
        formats: [],
        sub_tasks: [],
        languages: [],
        libraries: [],
        tags: [],
        license: ''
      },
      fileHash: [],            // 存储文件的哈希值
      availableTasks: ['图像识别', '文本摘要', '语音合成'],  // 示例数据
      availableModalities: ['视觉', '文本', '声音'],  // 示例数据
      availableFormats: ['PNG', 'CSV', 'WAV'],  // 示例数据
      availableSubTasks: ['面部识别', '情感分析', '语音到文本'],  // 示例数据
      availableLanguages: ['英语', '中文', '法语'],  // 示例数据
      availableLibraries: ['TensorFlow', 'PyTorch', 'Keras'],  // 示例数据
      availableTags: ['机器学习', '深度学习', '自然语言处理'],  // 示例数据
      availableLicenses: ['MIT License', 'Apache License', 'GNU General Public License (GPL)']  // 示例数据
    }
  },
  computed: {
    ...mapGetters([
      'userId',
      'roles',
      'userName'
    ])
  },
  methods: {
    beforeUpload(file) {
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
      const isFile = allowedTypes.some(type => file.type.includes(type)) ||
                    codeExtensions.some(ext => file.name.endsWith(ext));
      if (!isFile) {
        this.$message.error('只允许上传图片、文本、音频文件、代码文件、压缩包或CSV文件。');
        return false;
      }            
      this.$message.success('文件类型检查通过!');

      const form = new FormData();
      form.append('file', file);
      uploadFile(form).then(response => {
        this.fileHash.push(response.hash);
        this.$message.success(`文件上传成功! 文件哈希: ${response.hash}`); 
      }).catch(error => {
        this.$message.error('文件上传失败!');
      });
      return true;
    },
    handleUploadSuccess(response) {
      // this.fileHash = response.hash;
      //this.$message.success(`文件上传成功! 文件哈希: ${response.hash}`);
    },

  //   customRequest({ file, onSuccess, onError }) {
  //     this.$message.success('开始上传文件!');
  //     const formData = new FormData();
  //     formData.append('file', file);
  //     console.log('Uploading file:', file); // 调试用

  //     // 直接使用 axios 发送请求
  //     axios.post('http://localhost:8888/api/v1/uploadFile', formData, {
  //       headers: {
  //         'Content-Type': 'multipart/form-data'
  //       }
  //     })
  //     .then(response => {
  //       console.log('Upload response:', response); // 调试用
  //       this.$message.success('文件上传成功！');
  //       onSuccess(response.data); // 调用 onSuccess 通知 el-upload 上传成功
  //     })
  //     .catch(error => {
  //       console.error('Upload error:', error); // 调试用
  //       this.$message.error('文件上传失败！');
  //       onError(error); // 调用 onError 通知 el-upload 上传失败
  //     });
  //   }
  // ,
  customRequest (file){
  //  const form = new FormData();
  //  form.append('file', file);
  //   uploadFile(form).then(response => {
  //     this.$message.success('文件上传成功!');
  //     this.fileHash = response.hash;
  //   }).catch(error => {
  //     this.$message.error('文件上传失败!');
  //   });
},
    createDataset() {
      if (!this.newDatasetName) {
        this.$message({
          type: 'warning',
          message: '数据集名称不能为空!'
        });
        return;
      }

      if (!this.fileHash) {
        this.$message({
          type: 'warning',
          message: '请先上传文件!'
        });
        return;
      }

      const currentTime = new Date().toISOString(); // 获取当前时间的 ISO 字符串

      // 处理 metadata，将选择的数据转换为数组
      const metadata = {
        ...this.metadata
      };

      this.loading = true;
      uploadSet({
        Name: this.newDatasetName,  // 数据集名称
        Owner: this.userId,
        CreationTime: currentTime,
        Rows: this.rows,  // 新添加的行数字段
        metadata: metadata,
        Files: this.fileHash  // 添加文件哈希
      }).then(response => {
        this.loading = false;
        if (response) {
          this.createdDatasetName = response.dataset_name;
          this.$message({
            type: 'success',
            message: `数据集创建成功! 创建时间: ${currentTime}`
          });
        } else {
          this.$message({
            type: 'error',
            message: '数据集名称重复! 请重试。'
          });
        }
      }).catch(error => {
        this.loading = false;
        this.$message({
          type: 'error',
          message: `创建数据集时发生错误: ${error.message}`
        });
      });
      this.fileHash = [];
    }
  }
}
</script>

<style>
.container {
  width: 100%;
  text-align: center;
  min-height: 100%;
  overflow: hidden;
}

.dataset-container {
  margin: 20px auto;
  text-align: center;
  max-width: 600px; /* 控制容器宽度 */
  padding: 20px;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  background-color: #fff;
}

.dataset-input {
  margin-bottom: 20px;
}

.dataset-select {
  width: 100%;
  margin: 10px 0;
}

.dataset-button {
  margin: 20px 0;
}

.dataset-info {
  margin-top: 20px;
  color: #409EFF;
  font-size: 16px;
}
</style>
