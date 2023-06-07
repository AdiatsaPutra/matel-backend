import Vue from 'vue'
import Router from 'vue-router'
import { normalizeURL, decode } from 'ufo'
import { interopDefault } from './utils'
import scrollBehavior from './router.scrollBehavior.js'

const _4d7bc0ea = () => interopDefault(import('../pages/cms.vue' /* webpackChunkName: "pages/cms" */))
const _27ad38b8 = () => interopDefault(import('../pages/leasing.vue' /* webpackChunkName: "pages/leasing" */))
const _65e83e6c = () => interopDefault(import('../pages/login.vue' /* webpackChunkName: "pages/login" */))
const _126f29be = () => interopDefault(import('../pages/pengguna.vue' /* webpackChunkName: "pages/pengguna" */))
const _199cfd5a = () => interopDefault(import('../pages/upload-data/index.vue' /* webpackChunkName: "pages/upload-data/index" */))
const _03643db3 = () => interopDefault(import('../pages/index.vue' /* webpackChunkName: "pages/index" */))

const emptyFn = () => {}

Vue.use(Router)

export const routerOptions = {
  mode: 'history',
  base: '/',
  linkActiveClass: 'nuxt-link-active',
  linkExactActiveClass: 'nuxt-link-exact-active',
  scrollBehavior,

  routes: [{
    path: "/cms",
    component: _4d7bc0ea,
    name: "cms"
  }, {
    path: "/leasing",
    component: _27ad38b8,
    name: "leasing"
  }, {
    path: "/login",
    component: _65e83e6c,
    name: "login"
  }, {
    path: "/pengguna",
    component: _126f29be,
    name: "pengguna"
  }, {
    path: "/upload-data",
    component: _199cfd5a,
    name: "upload-data"
  }, {
    path: "/",
    component: _03643db3,
    name: "index"
  }],

  fallback: false
}

export function createRouter (ssrContext, config) {
  const base = (config._app && config._app.basePath) || routerOptions.base
  const router = new Router({ ...routerOptions, base  })

  // TODO: remove in Nuxt 3
  const originalPush = router.push
  router.push = function push (location, onComplete = emptyFn, onAbort) {
    return originalPush.call(this, location, onComplete, onAbort)
  }

  const resolve = router.resolve.bind(router)
  router.resolve = (to, current, append) => {
    if (typeof to === 'string') {
      to = normalizeURL(to)
    }
    return resolve(to, current, append)
  }

  return router
}
