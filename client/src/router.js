const { createRouter, createWebHashHistory } = VueRouter

const routes = [
  { path: '/', redirect: '/feed' },
  { path: '/login', component: () => import('./views/Login.vue.js') },
  { path: '/register', component: () => import('./views/Register.vue.js') },
  { path: '/feed', component: () => import('./views/Feed.vue.js') },
  { path: '/market', component: () => import('./views/Market.vue.js') },
  { path: '/chat', component: () => import('./views/Chat.vue.js') },
  { path: '/chat-room', component: () => import('./views/ChatRoom.vue.js') },
  { path: '/communities', component: () => import('./views/Communities.vue.js') },
  { path: '/profile', component: () => import('./views/Profile.vue.js') },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const publicPaths = ['/login', '/register']
  if (!token && !publicPaths.includes(to.path)) {
    next('/login')
  } else {
    next()
  }
})
