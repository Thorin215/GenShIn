<template>
  <div class="container">
    <el-alert type="success">
      <p>账户ID: {{ accountId }}</p>
      <!-- <p>用户名: {{ userName }}</p>
      <p>余额: ￥{{ balance }} 元</p> -->
    </el-alert>

    <!-- 新增部分 -->
    <div class="dataset-container">
      <el-input
        v-model="newDatasetId"
        placeholder="请输入要创建的数据集编号"
        class="dataset-input"
      />
      <el-button
        type="primary"
        @click="createDataset"
        class="dataset-button"
      >
        创建数据集
      </el-button>
      <p v-if="createdDatasetId" class="dataset-info">
        创建成功！数据集编号: {{ createdDatasetId }}, 账户ID: {{ accountId }}
      </p>
    </div>

    <!-- <div v-if="sellingList.length === 0" style="text-align: center;">
      <el-alert title="查询不到数据" type="warning" />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val, index) in sellingList" :key="index" :span="6" :offset="1">
        <el-card class="all-card"> -->
          <!-- <div slot="header" class="clearfix">
            <span>{{ val.sellingStatus }}</span>
            <el-button
              v-if="roles[0] !== 'admin' && (val.seller === accountId || val.buyer === accountId) && val.sellingStatus !== '完成' && val.sellingStatus !== '已过期' && val.sellingStatus !== '已取消'"
              style="float: right; padding: 3px 0"
              type="text"
              @click="updateSelling(val, 'cancelled')"
            >
              取消
            </el-button>
            <el-button
              v-if="roles[0] !== 'admin' && val.seller === accountId && val.sellingStatus === '交付中'"
              style="float: right; padding: 3px 8px"
              type="text"
              @click="updateSelling(val, 'done')"
            >
              确认收款
            </el-button>
            <el-button
              v-if="roles[0] !== 'admin' && val.sellingStatus === '销售中' && val.seller !== accountId"
              style="float: right; padding: 3px 0"
              type="text"
              @click="createSellingByBuy(val)"
            >
              购买
            </el-button>
          </div> -->
          <!-- <div class="item">
            <el-tag>房产ID: </el-tag>
            <span>{{ val.objectOfSale }}</span>
          </div>
          <div class="item">
            <el-tag type="success">销售者ID: </el-tag>
            <span>{{ val.seller }}</span>
          </div>
          <div class="item">
            <el-tag type="danger">价格: </el-tag>
            <span>￥{{ val.price }} 元</span>
          </div> -->
          <!-- <div class="item">
            <el-tag type="warning">有效期: </el-tag>
            <span>{{ val.salePeriod }} 天</span>
          </div> -->
          <!-- <div class="item">
            <el-tag type="info">创建时间: </el-tag>
            <span>{{ val.createTime }}</span>
          </div> -->
          <!-- <div class="item">
            <el-tag>购买者ID: </el-tag>
            <span v-if="val.buyer === ''">虚位以待</span>
            <span>{{ val.buyer }}</span>
          </div> -->
        <!-- </el-card>
      </el-col>
    </el-row> -->
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { querySellingList, createSellingByBuy, updateSelling } from '@/api/selling'
import { uploadSet } from '@/api/upload'

export default {
  name: 'AllSelling',
  data() {
  return {
    loading: false,
    sellingList: [],
    newDatasetId: '',        // 用户输入的新数据集编号
    createdDatasetId: '',    // 返回的创建的数据集编号
    creationTime: '',        // 返回的创建时间
    lastModified: ''         // 返回的最后修改时间
  }
},
  computed: {
    ...mapGetters([
      'accountId',
      'roles',
      'userName',
      'balance'
    ])
  },
  created() {
    this.fetchSellingList();
  },
  methods: {
  createDataset() {
    console.log('Create Dataset button clicked');
    if (!this.newDatasetId) {
      this.$message({
        type: 'warning',
        message: '数据集编号不能为空!'
      })
      return
    }

    const currentTime = new Date().toISOString(); // 获取当前时间的 ISO 字符串

    this.loading = true;
    uploadSet({
      dataset_id: this.newDatasetId,
      account_id: this.accountId,
      creation_time: currentTime // 传递创建时间
    }).then(response => {
      console.log('Response:', response); // 调试输出
      this.loading = false;
      if (response) {
        // 更新创建的数据集编号和时间
        this.createdDatasetId = response.dataset_id;
        this.creationTime = response.creation_time;
        this.lastModified = response.last_modified;

        this.$message({
          type: 'success',
          message: `数据集创建成功! 创建时间: ${this.creationTime}, 最后修改时间: ${this.lastModified}`
        });
      } else {
        this.$message({
          type: 'error',
          message: '数据集编号重复! 请重试。'
        });
      }
    }).catch(error => {
      console.error('Error:', error); // 调试输出
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

  .tag {
    float: left;
  }

  .item {
    font-size: 14px;
    margin-bottom: 18px;
    color: #999;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
  }

  .clearfix:after {
    clear: both;
  }

  .all-card {
    width: 280px;
    height: 380px;
    margin: 18px;
  }

  /* 新增部分的样式 */
  .dataset-container {
    margin: 20px 0;
    text-align: center;
  }

  .dataset-input {
    width: 300px;
    margin-right: 10px;
  }

  .dataset-button {
    margin-left: 10px;
  }

  .dataset-info {
    margin-top: 20px;
    color: #409EFF; /* El-Button primary color */
    font-size: 16px;
  }
</style>