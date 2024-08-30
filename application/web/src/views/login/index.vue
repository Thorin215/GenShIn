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
          v-for="item in accountList"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        >
          <span style="float: left">{{ item.name }}</span>
          <span style="float: right; color: #8492a6; font-size: 13px">{{ item.id }}</span>
        </el-option>
      </el-select>

      <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleLogin">立即进入</el-button>

      <div class="tips">
        <span style="margin-right:20px;">Tips: 选择不同用户角色进行数据共享测试</span>
      </div>

    </el-form>
  </div>
</template>

<script>
import { queryAccountList } from '@/api/account'

export default {
  name: 'Login',
  data() {
    return {
      loading: false,
      redirect: undefined,
      accountList: [],
      value: ''
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
    queryAccountList().then(response => {
      if (response) {
        this.accountList = response
        this.$message('数据加载成功')
      }
    })
  },
  methods: {
    handleLogin() {
      if (this.value) {
        this.loading = true
        this.$store.dispatch('account/login', this.value).then(() => {
          this.$router.push({ path: this.redirect || '/' })
          this.loading = false
        }).catch(() => {
          this.loading = false
        })
      } else {
        this.$message('请选择用户角色')
      }
    },
    selectGet(accountId) {
      this.value = accountId
    }
  }
}
</script>

<style lang="scss" scoped>
$bg:#2d3a4b;
$dark_gray:#889aa4;
$light_gray:#eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background: url('https://c7k8t9m10.github.io/medias/featureimages/8.jpg');
  width:100%;
  height:100%;
  position:fixed;
  background-size:100% 100%;
  overflow: hidden;
  .example {
  font-family: 'alarm_clockregular', sans-serif;
  .subtitle {
      font-size: 30px;
      color: #000000;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
   }
  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }
  .login-select{
   padding: 20px 0px 30px 0px;
   min-height: 100%;
   width: 100%;
   background-color: transparent;
   overflow: hidden;
   text-align: center;
  }
  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
}
</style>
