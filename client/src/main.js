// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import firebase from 'firebase'

Vue.config.productionTip = false

let app
const config = {
  apiKey: 'AIzaSyBpMPTk-JMpR_KG-T8ElsoFvOajwEmBKE8',
  authDomain: 'okra-205912.firebaseapp.com',
  databaseURL: 'https://okra-205912.firebaseio.com',
  projectId: 'okra-205912',
  storageBucket: 'okra-205912.appspot.com',
  messagingSenderId: '634081150354'
}
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
