<template>
  <div class="container">
    <div class="dataset-grid">
      <el-card
        v-for="dataset in datasets"
        :key="dataset.name"
        class="dataset-card"
      >
        <h4 class="dataset-name">{{ dataset.name }}</h4>
        <p class="dataset-owner">所有者: {{ dataset.owner }}</p>
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
            <el-button type="primary" icon="el-icon-download" @click="downloadFiles(scope.row.files)">下载文件</el-button>
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


    <!-- 文件详情对话框 -->
    <el-dialog title="文件详情" :visible.sync="fileDialogVisible" width="60%" @close="closeFileDialog">
      <div v-if="selectedFiles.length">
        <el-list>
          <el-list-item v-for="file in selectedFiles" :key="file">
            <span>{{ file || '无文件' }}</span>
          </el-list-item>
        </el-list>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeFileDialog">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { GetAllDataSet } from '@/api/datasets';
import { getDatasetMetadata, downloadDataset } from '@/api/upload';

export default {
  name: 'DataSetsTable',
  data() {
    return {
      datasets: [],
      logs: [],
      metadata: {},
      dialogVisible: false,
      fileDialogVisible: false,
      selectedFiles: [],
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
        const response2 = await getDatasetMetadata({ owner: dataset.owner, name: dataset.name });
        console.log(response2);
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
    async downloadFiles(files) {
      try {
        this.$message('加载');
        //log.console('files', files);
        const response3 = await downloadDataset({
          files: files,
          name: this.selectedDataset.name,
          owner: this.selectedDataset.owner, 
        });
        log.console('response3', response3);
        this.$message('加载完成');

        // 从响应中提取文件内容和文件名
        const { files: fileData } = response3.data;
        if (fileData && fileData.length > 0) {
          const file = fileData[0];
          const { filename, content } = file;

          // Base64 解码
          const byteCharacters = atob(content);
          const byteNumbers = new Array(byteCharacters.length);
          for (let i = 0; i < byteCharacters.length; i++) {
            byteNumbers[i] = byteCharacters.charCodeAt(i);
          }
          const byteArray = new Uint8Array(byteNumbers);

          // 创建 Blob 对象
          const blob = new Blob([byteArray], { type: 'application/zip' });

          // 创建 URL 对象来表示文件的下载地址
          const url = window.URL.createObjectURL(blob);

          // 创建一个 a 元素，并设置 href 为文件 URL
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', filename); // 使用从响应中获取的文件名

          // 触发下载
          document.body.appendChild(link);
          link.click();

          // 清理 URL 对象
          window.URL.revokeObjectURL(url);
        } else {
          this.$message.error('未找到要下载的文件');
        }
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
