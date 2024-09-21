import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/* Layout */
import Layout from '@/layout'

export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dataset',
    children: [
      {
        path: 'dataset',
        name: 'Dataset',
        component: () => import('@/views/dataset/list/index'),
        meta: {
          title: '数据信息',
          icon: 'xj'
        }
      }]
    },
    {
      path: '/',
      component: Layout,
      redirect: '/upload',
      children: [
        {
        path: 'upload',
        name: 'Upload',
        component: () => import('@/views/upload/all/index'),
        meta: {
          title: '新建数据集',
          icon: 'new2'
        }
      }
      ]
    },
    {
      path: '/',
      component: Layout,
      redirect: '/update',
      children: [{
        path: 'update',
        name: 'Update',
        component: () => import('@/views/update/all/index'),
        meta: {
          title: '更新数据集',
          icon: 'pen'
        }
      }
      ]
    },
    {
      path: '/',
      component: Layout,
      redirect: '/record',
      children: [{
        path: 'record',
        name: 'Record',
        component: () => import('@/views/record/all/index'),
        meta: {
          title: '我的记录',
          icon: 'new3'
        }
      }]
    }
]

export const asyncRoutes = [
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

const createRouter = () => new Router({
  base: '/',
  // mode: 'history', // require service support
  scrollBehavior: () => ({
    y: 0
  }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
