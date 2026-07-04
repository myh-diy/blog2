import { createRouter, createWebHistory } from 'vue-router'
import PostDetailView from '../views/PostDetailView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/post/:slug',
      name: 'post',
      component: PostDetailView,
    },
  ],
})

export default router
