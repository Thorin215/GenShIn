<template>
  <div class="login-container">
    <el-form ref="loginForm" class="login-form" auto-complete="on" label-position="left">
      <div class="example">
        <h1 class='subtitle'>GenShIn</h1>
      </div>
      <div class="title-container">
        <h3 class="title">基于区块链的AI训练数据共享系统</h3>
      </div>
      <el-select v-model="value" placeholder="请选择用户角色" class="login-select" @change="selectGet">
        <el-option
          v-for="item in userList"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        >
          <span style="float: left">{{ item.name }}</span>
          <span style="float: right; color: #8492a6; font-size: 13px">{{ item.id }}</span>
        </el-option>
      </el-select>

      <!-- 按钮容器 -->
      <div class="button-container">
        <el-button :loading="loading" type="primary" class="action-button" @click.native.prevent="handleLogin">登录</el-button>
        <el-button type="success" class="action-button" @click.native.prevent="showRegisterDialog">注册</el-button>
      </div>

      <!-- 注册弹窗 -->
      <el-dialog title="用户注册" :visible.sync="registerDialogVisible" width="30%" center>
        <el-form ref="registerForm" :model="registerForm">
          <el-form-item label="用户名">
            <el-input v-model="registerForm.name"></el-input>
          </el-form-item>
          <el-form-item label="用户ID">
            <el-input v-model="registerForm.id"></el-input>
          </el-form-item>
          <el-form-item label="密码">
            <el-input type="password" v-model="registerForm.password"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button @click="registerDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleRegister">注册</el-button>
        </div>
      </el-dialog>
    </el-form>
  </div>
</template>


<script>
import { queryAllUsers, createUser} from '@/api/user'

export default {
  name: 'Login',
  data() {
    return {
      loading: false,
      redirect: undefined,
      userList: [],
      value: '',
      registerDialogVisible: false, // 控制注册弹窗显示
      registerForm: {
        name: '',
        id: '',
        password: ''
      }
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  created() {
    queryAllUsers().then(response => {
      if (response) {
        this.userList = response
        this.$message.success('数据加载成功')
      }
    })
  },

  methods: {
    handleLogin() {
      if (this.value) {
        this.loading = true
        this.$store.dispatch('user/login', this.value).then(path => {
          this.$router.push({ path: path })
          this.loading = false
        }).catch(() => {
          this.loading = false
        })
      } else {
        this.$message('请选择用户角色')
      }
    },
    selectGet(userId) {
      this.value = userId
    },
    showRegisterDialog() {
      this.registerDialogVisible = true
    },
    handleRegister() {
      if (this.registerForm.name && this.registerForm.id && this.registerForm.password) {
        createUser(this.registerForm).then(() => {
          this.$message.success('注册成功');
          this.registerDialogVisible = false;
        }).catch(() => {
          this.$message.error('注册失败');
        });
      } else {
        this.$message.error('请填写所有字段');
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100%;
  width: 100%;
  height:100%;
  position: fixed;
  background-color: #f2f2f2; // 设置为淡灰色
  background-size: 100% 100%;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .login-select {
    padding: 20px 0px 30px 0px;
    width: 100%;
    background-color: transparent;
    text-align: center;
  }

  .button-container {
    display: flex;            // 使按钮横向排列
    justify-content: space-between;  // 两个按钮之间留空
    gap: 10px;                // 按钮之间的间距
    margin-bottom: 30px;
  }

  .action-button {
    flex: 1;                  // 使按钮均匀分布，占据相同宽度
  }
}
</style>

