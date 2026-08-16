import { createRouter, createWebHistory } from 'vue-router'
import { setupGuards } from './guards'
import Layout from '@/pages/Layout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: Layout,
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: () => import('@/pages/Dashboard.vue') },
        { path: 'warehouses', name: 'warehouses', component: () => import('@/pages/WarehouseManage.vue') },
        { path: 'products', name: 'products', component: () => import('@/pages/ProductManage.vue') },
        { path: 'orders', name: 'orders', component: () => import('@/pages/OrderManage.vue') },
        { path: 'inventory-check', name: 'inventory-check', component: () => import('@/pages/InventoryCheck.vue') },
        { path: 'profile', name: 'profile', component: () => import('@/pages/Profile.vue'), meta: { requiresAuth: true } },
        {
          path: 'operation-logs',
          name: 'operation-logs',
          component: () => import('@/pages/OperationLogs.vue'),
          meta: { requiresAuth: true, roles: ['admin'] },
        },
      ],
    },
    { path: '/login', name: 'login', component: () => import('@/pages/Login.vue') },
  ],
})

setupGuards(router)
export default router
