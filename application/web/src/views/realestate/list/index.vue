<template>
  <div class="container">
    <el-alert
      type="success"
    >
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
        <el-button type="text" @click="viewLogs(dataset.name)">查看修改日志</el-button>
      </el-card>
    </div>

    <!-- 日志对话框 -->
    <el-dialog title="修改日志" :visible.sync="dialogVisible" width="60%" @close="closeDialog">
      <el-table :data="logs" style="width: 100%">
        <el-table-column prop="LogID" label="日志ID"></el-table-column>
        <el-table-column prop="DataSetID" label="数据集ID"></el-table-column>
        <el-table-column prop="ChangeLog" label="变更日志"></el-table-column>
        <el-table-column prop="TimeStemp" label="时间戳"></el-table-column>
      </el-table>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">返回</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import { GetAllDataSet, queryChangeLog } from '@/api/datasets';

export default {
  name: 'DataSetsTable',
  data() {
    return {
      datasets: [],
      logs: [],
      dialogVisible: false,
    };
  },
  computed: {
    ...mapGetters(['userId', 'userName', 'roles']),
  },
  mounted() {
    this.fetchDataSets();
  },
  methods: {
    headClass() { 
                    return "text-align:center"
                },
    async fetchDataSets() {
      try {
        const response = await GetAllDataSet();
        this.datasets = response;
        console.log(response.data);
        this.$message.success('数据加载成功！');
      } catch (error) {
        console.error('Error fetching datasets:', error);
        this.$message.error('数据加载失败');
      }
    },
    async viewLogs(dataSetName) {
      try {
        const response = await queryChangeLog(dataSetName);
        console.log(dataSetName);
        this.logs = response;
        // 逻辑处理，例如打开一个对话框显示日志
        console.log('Logs for dataset:', this.logs);
        this.dialogVisible = true;
      } catch (error) {
        console.error('Error fetching logs:', error);
      }
    },
    closeDialog() {
      this.dialogVisible = false;
    this.logs = []; // 清空日志数据
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

.realEstate-card {
  width: 280px;
  height: 340px;
  margin: 18px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
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