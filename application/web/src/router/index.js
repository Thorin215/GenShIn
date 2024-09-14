import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [{
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
  children: [{
    path: 'dataset',
    name: 'Dataset',
    component: () => import('@/views/dataset/list/index'),
    meta: {
      title: '数据信息',
      icon: 'xj'
    }
  }]
}
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/upload',
    component: Layout,
    redirect: '/upload/all',
    name: 'Upload',
    alwaysShow: true,
    meta: {
      title: '新建',
      icon: 'new5'
    },
    children: [{
      path: 'all',
      name: 'UploadAll',
      component: () => import('@/views/upload/all/index'),
      meta: {
        title: '新建数据集',
        icon: 'new2'
      }
    }
    ]
  },
  {
    path: '/update',
    component: Layout,
    redirect: '/update/all',
    name: 'Update',
    alwaysShow: true,
    meta: {
      title: '更新',
      icon: 'new4'
    },
    children: [{
      path: 'all',
      name: 'UpdateAll',
      component: () => import('@/views/update/all/index'),
      meta: {
        title: '更新数据集',
        icon: 'pen'
      }
    }
    ]
  },
  {
    path: '/record',
    component: Layout,
    redirect: '/record/all',
    name: 'Record',
    alwaysShow: true,
    meta: {
      title: '记录',
      icon: 'new2'
    },
    children: [{
      path: 'all',
      name: 'RecordAll',
      component: () => import('@/views/record/all/index'),
      meta: {
        title: '我的记录',
        icon: 'new3'
      }
    }
    ]
  },
  // 404 page must be placed at the end !!!
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
