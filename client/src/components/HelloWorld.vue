<template>
  <div class="hello">
    <h1>Hello {{ name }}!!</h1>
    <h1>{{ msg }}</h1>
    <h2>Essential Links</h2>
    <button @click="signOut">Sign out</button>
  </div>
</template>

<script>
import firebase from 'firebase'

export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App',
      name: firebase.auth().currentUser.email
    }
  },
  methods: {
    signOut: function () {
      firebase.auth().signOut().then(() => {
        localStorage.removeItem('jwt')
        this.$router.push('/signin')
      })
    }
  },
  created: function () {
    firebase.auth().currentUser.getIdToken(/* forceRefresh */ true).then(idToken => {
      localStorage.setItem('jwt', idToken);
    }).catch(error => {
      console.log(error)
    })
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
