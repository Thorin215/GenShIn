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
        <h4>{{ dataset.name }}</h4>
        <p>所有者: {{ dataset.owner }}</p>
        <el-button type="text" @click="viewLogs(dataset)">查看修改日志</el-button>
        <el-button type="text" @click="viewMetadata(dataset)">查看详细信息</el-button>
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
import { getDatasetMetadata } from '@/api/upload';
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

.dataset-card {
  margin-bottom: 20px;
  cursor: pointer;
}

.dataset-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr); /* 创建5列 */
  gap: 20px; /* 设置间隙 */
  margin-top: 50px;
  margin-left: 100px;
  width: 80%;
}
</style>
