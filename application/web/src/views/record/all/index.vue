<template>
  <div>
    <el-table :data="downloadRecords" style="width: 100%">
      <el-table-column prop="dataset_owner" label="Dataset Owner" width="180">
      </el-table-column>
      <el-table-column prop="dataset_name" label="Dataset Name" width="180">
      </el-table-column>
      <el-table-column prop="user" label="User" width="180">
      </el-table-column>
      <el-table-column
          prop="time"
          label="Download Time"
          width="180"
        >
          <template slot-scope="scope">
            {{ new Date(scope.row.time).toISOString() }} <!-- 输出 ISO 8601 格式 -->
          </template>
        </el-table-column>
      <el-table-column label="Files">
        <template slot-scope="scope">
          <el-tag
            v-for="(file, index) in scope.row.files"
            :key="index"
            type="primary"
            style="margin-right: 5px"
          >
            <!-- 将file.filename与后端返回的字段一致 -->
            {{ file.filename }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { queryRecordsByUser } from '@/api/record';

export default {
  data() {
    return {
      downloadRecords: [] // Store fetched download records
    };
  },
  computed: {
    ...mapGetters(['userId', 'userName', 'roles']),
  },
  mounted() {
    this.fetchDownloadRecords(); // Fetch download records when component mounts
  },
  methods: {
    async fetchDownloadRecords() {
      try {
        const response = await queryRecordsByUser({ user: this.userId }); // Adjust API endpoint as needed
        this.downloadRecords = response; // Assuming the data is in the 'data' field of the response
      } catch (error) {
        this.$message.error("Failed to fetch download records");
      }
    }
  }
};
</script>


