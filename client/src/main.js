// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import firebase from 'firebase'
import config from '../firebase.config'

Vue.config.productionTip = false

let app
firebase.initializeApp(config)
firebase.auth().onAuthStateChanged(user => {
  /* eslint-disable no-new */
  if (!app) {
    new Vue({
      el: '#app',
      router,
      components: { App },
      template: '<App/>'
    })
  }
})
