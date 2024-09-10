<template>
  <div class="container">
    <!-- <div class="search-bar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索数据集名称"
        prefix-icon="el-icon-search"
        @input="filteredDatasets"
      ></el-input>
    </div> -->
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
import { queryAllDatasets, queryDatasetMetadata } from '@/api/dataset';
import { downloadFilesCompressed } from '@/api/file';

export default {
  name: 'DataSetsTable',
  data() {
    return {
      datasets: [],
      logs: [],
      metadata: {},
      dialogVisible: false,
      fileDialogVisible: false,
      metadataDialogVisible: false,
      selectedFiles: [],
      selectedDataset: null,
      searchQuery: '',
    };
  },
  computed: {
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
        this.datasets = response;
        console.log(this.datasets);
        this.$message.success('数据加载成功！');
      } catch (error) {
        console.error('Error fetching datasets:', error);
        this.$message.error('数据加载失败');
      }
    },
    viewLogs(dataset) {
      this.selectedDataset = dataset;
      this.logs = dataset.versions; // 将选中的数据集版本作为日志
      this.dialogVisible = true;
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
