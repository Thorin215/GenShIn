<template>
  <div class="container">

    <div class="dataset-grid">
      <el-container>
            <div class="dataset-query">
              <el-select v-model="taskvalue" multiple filterable boolean placeholder="请选择任务" @change = "fetchDataSets">
                <el-option
                    v-for="task in optionsTasks"
                    :key="task.value"
                    :label="task.label"
                    :value="task.value">
                </el-option>
              </el-select>
            </div>
        <div class="dataset-query">
          <el-select v-model="modalityvalue" multiple filterable boolean placeholder="请选择数据模态" @change = "fetchDataSets">
            <el-option
                v-for="modality in optionsModalities"
                :key="modality.value"
                :label="modality.label"
                :value="modality.value">
            </el-option>
          </el-select>
        </div>
        <div class="dataset-query">
          <el-select v-model="formatvalue" multiple filterable boolean placeholder="请选择文件格式" @change = "fetchDataSets">
            <el-option
                v-for="format in optionsFormats"
                :key="format.value"
                :label="format.label"
                :value="format.value">
            </el-option>
          </el-select>
        </div>
        <div class="dataset-query">
          <el-select v-model="subtaskvalue" multiple filterable boolean placeholder="请选择子任务" @change = "fetchDataSets">
            <el-option
                v-for="subtask in optionsSubTasks"
                :key="subtask.value"
                :label="subtask.label"
                :value="subtask.value">
            </el-option>
          </el-select>
        </div>
        <div class="dataset-query">
          <el-select v-model="languagevalue" multiple filterable boolean placeholder="请选择语言" @change = "fetchDataSets">
            <el-option
                v-for="language in optionsLanguages"
                :key="language.value"
                :label="language.label"
                :value="language.value">
            </el-option>
          </el-select>
        </div>
        <div class="dataset-query">
          <el-select v-model="libraryvalue" multiple filterable boolean placeholder="请选择适用库" @change = "fetchDataSets">
            <el-option
                v-for="library in optionsLibraries"
                :key="library.value"
                :label="library.label"
                :value="library.value">
            </el-option>
          </el-select>
        </div>
        <div class="dataset-query">
          <el-select v-model="tagvalue" multiple filterable boolean placeholder="请选择标签" @change = "fetchDataSets">
            <el-option
                v-for="tag in optionsTags"
                :key="tag.value"
                :label="tag.label"
                :value="tag.value">
            </el-option>
          </el-select>
        </div>
        <div class="dataset-query">
          <el-select v-model="licensevalue"  filterable boolean placeholder="请选择许可证" @change = "fetchDataSets">
            <el-option
                v-for="license in optionsLicenses"
                :key="license.value"
                :label="license.label"
                :value="license.value">
            </el-option>
          </el-select>
        </div>
  
  
  
      </el-container>
      </div>

    <div class="dataset-grid">
      <el-card
        v-for="dataset in datasets"
        :key="dataset.name"
        class="dataset-card"
      >
        <h4 class="dataset-name">{{ dataset.name }}</h4>
        <p class="dataset-owner">所有者: {{ dataset.owner }}</p>
        <p class="dataset-downloads">下载次数 : {{ dataset.downloads }}</p>
        <div class="dataset-actions">
          <!-- Removed download button -->
          <el-button type="success" icon="el-icon-edit" @click="viewLogs(dataset)">查看修改日志</el-button>
          <el-button type="info" icon="el-icon-info" @click="viewMetadata(dataset)">查看详细信息</el-button>
          <el-button type="danger" icon="el-icon-delete" @click="deleteDataset(dataset)">删除</el-button>
        </div>
      </el-card>
    </div>

    <!-- 修改日志对话框 -->
    <el-dialog title="修改日志" :visible.sync="dialogVisible" width="60%" @close="closeDialog">
      <el-table :data="logs" style="width: 100%">
        <el-table-column prop="creation_time" label="时间戳"></el-table-column>
        <el-table-column prop="change_log" label="变更日志"></el-table-column>
        <el-table-column prop="files" label="文件">
          <template v-slot="scope">
            <el-button type="primary" icon="el-icon-download" @click="downloadFiles(selectedDataset.name, scope.row.files)">下载文件</el-button>
          </template>
        </el-table-column>
      </el-table>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">返回</el-button>
      </span>
    </el-dialog>

 <!-- 详细信息对话框 -->
 <el-dialog title="详细信息" :visible.sync="metadataDialogVisible" width="60%" @close="closeMetadataDialog">
      <el-form :model="metadata" label-width="120px">
        <el-form-item label="任务">
          <div v-for="task in metadata.tasks" :key="task">{{ task || '无任务' }}</div>
        </el-form-item>
        <el-form-item label="数据模态">
          <div v-for="modality in metadata.modalities" :key="modality">{{ modality || '无模态' }}</div>
        </el-form-item>
        <el-form-item label="文件格式">
          <div v-for="format in metadata.formats" :key="format">{{ format || '无格式' }}</div>
        </el-form-item>
        <el-form-item label="子任务">
          <div v-for="subtask in metadata.sub_tasks" :key="subtask">{{ subtask || '无子任务' }}</div>
        </el-form-item>
        <el-form-item label="语言">
          <div v-for="language in metadata.languages" :key="language">{{ language || '无语言' }}</div>
        </el-form-item>
        <el-form-item label="适用库">
          <div v-for="library in metadata.libraries" :key="library">{{ library || '无库' }}</div>
        </el-form-item>
        <el-form-item label="标签">
          <div v-for="tag in metadata.tags" :key="tag">{{ tag || '无标签' }}</div>
        </el-form-item>
        <el-form-item label="许可证">
          <div>{{ metadata.license || '无许可证' }}</div>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeMetadataDialog">返回</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { queryAllDatasets, queryDatasetMetadata, deleteDataset, queryAllVersions } from '@/api/dataset';
import { downloadFilesCompressed } from '@/api/file';

export default {
  name: 'DataSetsTable',
  data() {
    return {
      datasets: [],
      logs: [],
      metadata: {},
      TempMetaData: {},
      dialogVisible: false,
      fileDialogVisible: false,
      metadataDialogVisible: false,
      selectedFiles: [],
      selectedDataset: null,
      searchQuery: '',
      taskvalue: [],
      modalityvalue: [],
      formatvalue: [],
      subtaskvalue: [],
      languagevalue: [],
      libraryvalue: [],
      tagvalue: [],
      licensevalue: '',
      index:0,
      flag: true,
      optionsTasks: [
        { value: '图像识别', label: '图像识别' },
        { value: '文本摘要', label: '文本摘要' },
        { value: '语音合成', label: '语音合成' }
      ],
      optionsModalities: [
        { value: '视觉', label: '视觉' },
        { value: '文本', label: '文本' },
        { value: '声音', label: '声音' }
      ],
      optionsFormats: [
        { value: 'PNG', label: 'PNG' },
        { value: 'CSV', label: 'CSV' },
        { value: 'WAV', label: 'WAV' }
      ],
      optionsSubTasks: [
        { value: '面部识别', label: '面部识别' },
        { value: '情感分析', label: '情感分析' },
        { value: '语音到文本', label: '语音到文本' }
      ],
      optionsLanguages: [
        { value: '英语', label: '英语' },
        { value: '中文', label: '中文' },
        { value: '法语', label: '法语' }
      ],
      optionsLibraries: [
        { value: 'TensorFlow', label: 'TensorFlow' },
        { value: 'PyTorch', label: 'PyTorch' },
        { value: 'Keras', label: 'Keras' }
      ],
      optionsTags: [
        { value: '机器学习', label: '机器学习' },
        { value: '深度学习', label: '深度学习' },
        { value: '自然语言处理', label: '自然语言处理' }
      ],
      optionsLicenses: [
        { value: 'MIT License', label: 'MIT License' },
        { value: 'Apache License', label: 'Apache License' },
        { value: 'GNU General Public License (GPL)', label: 'GNU General Public License (GPL)' },
        { value: '', label: 'No License' }
      ]
    };
  },
  computed: {
    data() {
      return data
    },
    ...mapGetters(['userId', 'userName', 'roles']),
    filteredDatasets() {
      if (!this.searchQuery) return this.datasets;
      return this.datasets.filter(dataset =>
        dataset.name.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
  },
  mounted() {
    this.fetchDataSets();
  },
  methods: {
    async fetchDataSets() {
      try {
        const response = await queryAllDatasets();
                this.datasets = [];
        for(const item of response){
          this.TempMetaData = await queryDatasetMetadata({ owner: item.owner, name: item.name });
          if(
              this.taskvalue.every(it => this.TempMetaData.tasks.includes(it))
              && this.modalityvalue.every(it1 => this.TempMetaData.modalities.includes(it1))
              && this.formatvalue.every(it2 => this.TempMetaData.formats.includes(it2))
              && this.subtaskvalue.every(it3 => this.TempMetaData.sub_tasks.includes(it3))
              && this.languagevalue.every(it4 => this.TempMetaData.languages.includes(it4))
              && this.libraryvalue.every(it5 => this.TempMetaData.libraries.includes(it5))
              && this.tagvalue.every(it6 => this.TempMetaData.tags.includes(it6))
              && this.TempMetaData.license === this.licensevalue
          ){
            this.datasets.push(item);
          }
        }
        // this.datasets = response;
        console.log(this.datasets);
        this.$message.success('数据加载成功！');
      } catch (error) {
        console.error('Error fetching datasets:', error);
        this.$message.error('数据加载失败');
      }
    },
    async deleteDataset(dataset) {
      try {
        if(this.userId === dataset.owner) {
          const reponse2 = await deleteDataset({name: dataset.name, owner: dataset.owner}); // Use your API call to delete the dataset
          if(reponse2 != "success") {
            this.$message.error('删除失败');
            return;
          }
          this.datasets = this.datasets.filter(d => d.id !== dataset.id); // Update UI after deletion
          this.$message.success(`删除 ${dataset.name} 成功`);
        }else {
          this.$message.error('不能删除他人的数据集');
          return;
        }  
      } catch (error) {
        console.error('删除数据集失败:', error);
        this.$message.error(`删除失败: ${error.message || error}`);
      }
    },
    async viewLogs(dataset) {
  this.selectedDataset = dataset;
  try {
    // 使用 queryAllVersions 来获取该数据集的所有版本信息
    const response = await queryAllVersions({ owner: dataset.owner, name: dataset.name });
    this.logs = response; // 将 API 返回的版本信息作为日志
    this.dialogVisible = true; // 打开日志对话框
  } catch (error) {
    console.error('Error fetching versions:', error);
    this.$message.error('获取修改日志失败');
  }
},
    async viewMetadata(dataset) {
      try {
        const response2 = await queryDatasetMetadata({ owner: dataset.owner, name: dataset.name });
        this.metadata = response2; // 获取元数据
        this.metadataDialogVisible = true;
      } catch (error) {
        console.error('Error fetching metadata:', error);
        this.$message.error('获取元数据失败');
      }
    },
    openFileDialog(files) {
      this.selectedFiles = files;
      this.fileDialogVisible = true;
    },
    closeDialog() {
      this.dialogVisible = false;
      this.logs = []; // 清空日志数据
    },
    closeFileDialog() {
      this.fileDialogVisible = false;
      this.selectedFiles = []; // 清空选中的文件
    },
    closeMetadataDialog() {
      this.metadataDialogVisible = false;
      this.metadata = {}; // 清空元数据
    },
    async downloadFiles(zipname, files) {
      try {
        await downloadFilesCompressed({
          files: files,
          zipname: zipname,
          dataset_owner: this.selectedDataset.owner,
          dataset_name: this.selectedDataset.name,
          user: this.userId
        });
      } catch (error) {
        console.error('下载文件失败:', error);
        this.$message.error(`下载文件失败: ${error.message || error}`);
      }
    }
  },
};
</script>
<style scoped>
.container {
  width: 100%;
  text-align: center;
  min-height: 100%;
  overflow: hidden;
}

/* .search-bar {
  margin-bottom: 20px;
} */

.dataset-query {
  gap: 20px; /* 间距 */
  justify-content: center;
}


.dataset-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20px; /* 间距 */
  justify-content: center;
}

.dataset-card {
  width: calc(33% - 20px); /* 每个卡片占据三分之一的宽度，减去间距 */
  box-shadow: 0 4px 8px rgba(0,0,0,0.1); /* 阴影效果 */
  transition: transform 0.3s ease-in-out; /* 平滑过渡效果 */
}

.dataset-card:hover {
  transform: translateY(-5px); /* 鼠标悬停时上移 */
}

.dataset-name {
  margin: 0;
  font-size: 1.25rem;
  color: #333;
}

.dataset-owner {
  font-size: 0.875rem;
  color: #666;
  margin: 10px 0;
}

.dataset-downloads {
  font-size: 0.875rem;
  color: #666;
  margin: 10px 0;
}

.dataset-actions {
  display: flex;
  gap: 10px; /* 按钮之间的间距 */
}

.el-button {
  flex: 1; /* 按钮平分空间 */
  text-align: left; /* 文字左对齐 */
  padding: 10px 15px; /* 增加内边距 */
}

.el-button i {
  margin-right: 8px; /* 图标与文字之间的间距 */
}
</style>
