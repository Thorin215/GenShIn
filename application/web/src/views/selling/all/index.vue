<template>
  <div class="container">
    <el-alert type="success">
      <p>账户ID: {{ userId }}</p>
    </el-alert>

    <!-- 数据集创建部分 -->
    <div class="dataset-container">
      <el-input
        v-model="newDatasetName"
        placeholder="请输入要创建的数据集名称"
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
        v-model="metadata.subTasks"
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

      <el-button
        type="primary"
        @click="createDataset"
        class="dataset-button"
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
import { mapGetters } from 'vuex'
import { uploadSet } from '@/api/upload'

export default {
  name: 'AllSelling',
  data() {
    return {
      loading: false,
      newDatasetName: '',      // 用户输入的新数据集名称
      createdDatasetName: '',  // 返回的创建的数据集名称
      creationTime: '',
     // lastModified: '',
      metadata: {
        tasks: [],
        modalities: [],
        formats: [],
        subTasks: [],
        languages: [],
        libraries: [],
        tags: [],
        license: ''
      },
      availableTasks: ['任务1', '任务2', '任务3'],  // 示例数据
      availableModalities: ['模态1', '模态2', '模态3'],  // 示例数据
      availableFormats: ['格式1', '格式2', '格式3'],  // 示例数据
      availableSubTasks: ['子任务1', '子任务2', '子任务3'],  // 示例数据
      availableLanguages: ['语言1', '语言2', '语言3'],  // 示例数据
      availableLibraries: ['库1', '库2', '库3'],  // 示例数据
      availableTags: ['标签1', '标签2', '标签3'],  // 示例数据
      availableLicenses: ['许可证1', '许可证2', '许可证3']  // 示例数据
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
    createDataset() {
      if (!this.newDatasetName) {
        this.$message({
          type: 'warning',
          message: '数据集名称不能为空!'
        })
        return
      }

      const currentTime = new Date().toISOString(); // 获取当前时间的 ISO 字符串

      // 处理 metadata，将选择的数据转换为数组
      const metadata = {
        ...this.metadata
      };

      this.loading = true;
      uploadSet({
        Name: this.newDatasetName,  // 修改为数据集名称
        Owner: this.userId,
        CreationTime: currentTime,
        metadata: metadata
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
  max-width: 500px; /* 控制容器宽度 */
}

.dataset-select {
  width: 100%;
  margin: 10px 0;
}

.dataset-button {
  margin: 10px 0;
}

.dataset-info {
  margin-top: 20px;
  color: #409EFF;
  font-size: 16px;
}
</style>
