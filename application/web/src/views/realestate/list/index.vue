<template>
  <div class="container">
    <el-alert type="success">
      <p>账户ID: {{ userId }}</p>
    </el-alert>
    <div class="dataset-grid">
  <el-card
    v-for="dataset in datasets"
    :key="dataset.name"
    class="dataset-card"
  >
    <h4 class="dataset-name">{{ dataset.name }}</h4>
    <p class="dataset-owner">所有者: {{ dataset.owner }}</p>
    <div class="dataset-actions">
      <el-button type="primary" icon="el-icon-download" @click="DownloadDatasets(dataset.name, userId)">下载</el-button>
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
          <template slot-scope="scope">
            <div v-for="file in scope.row.files" :key="file">{{ file || '无文件' }}</div>
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
import { GetAllDataSet } from '@/api/datasets';
import { getDatasetMetadata, DownloadDatasets } from '@/api/upload';
export default {
  name: 'DataSetsTable',
  data() {
    return {
      datasets: [],
      logs: [],
      metadata: {},
      dialogVisible: false,
      metadataDialogVisible: false,
      selectedDataset: null,
    };
  },
  computed: {
    ...mapGetters(['userId', 'userName', 'roles']),
  },
  mounted() {
    this.fetchDataSets();
  },
  methods: {
    async fetchDataSets() {
      try {
        const response = await GetAllDataSet();
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
        const response2 = await getDatasetMetadata( { owner: dataset.owner, name: dataset.name });
        console.log(response2);
        this.metadata = response2; // 获取元数据
        this.metadataDialogVisible = true;
      } catch (error) {
        console.error('Error fetching metadata:', error);
        this.$message.error('获取元数据失败');
      }
    },
    async DownloadDatasets(dataSetName) {
    const userId = this.userId; // 从 Vuex 获取用户 ID
    const timestamp = new Date().getTime(); // 获取当前时间戳
    try {
      // 调用后端 API 来处理下载请求
      const response = await DownloadDataset({
        name: dataSetName,
        userId: userId,
        timestamp: timestamp
      });
      // 处理响应，例如重定向到下载链接或显示下载信息
      console.log('Download response:', response);
      if (response && response.data && response.data.downloadUrl) {
        // 假设后端返回一个包含下载链接的响应
        window.location.href = response.data.downloadUrl;
      } else {
        this.$message.error('下载请求失败');
      }
    } catch (error) {
      console.error('Error initiating download:', error);
      this.$message.error('下载请求失败');
    }
  },
    closeDialog() {
      this.dialogVisible = false;
      this.logs = []; // 清空日志数据
    },
    closeMetadataDialog() {
      this.metadataDialogVisible = false;
      this.metadata = {}; // 清空元数据
    },
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
