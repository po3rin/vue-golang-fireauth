# vue-golang-fireauth

[![Go Report Card](https://goreportcard.com/badge/github.com/po3rin/vue-golang-fireauth)](https://goreportcard.com/report/github.com/po3rin/vue-golang-fireauth)

## Introduction

<img width="855" alt="スクリーンショット 2018-06-18 17.52.34.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/47db3a74-6c0c-1c50-682b-f7d6b2158ba9.png">


これはQiita記事にも記載されています。
https://qiita.com/po3rin/items/d3e016d01162e9d9de80


こんにちは。僕です。最近は Vue.js で SPA、Go言語で APIサーバーを実装しています。ここで意識するのはフロントエンドとバックエンドの分離です。分離させたい理由は二つ。

* アプリ化の際や、今後、他のサービスと連携させるためにも、バックエンドはAPIサーバーに徹したい。
* Vue.jsで開発中にオートリロードしたい(Vue.jsをビルドして出来たファイルをGo言語でserveする形にすると、毎回ビルドするのが面倒)

しかし、これをやろうとすると認証どうするの問題が付いてきます。APIサーバ、フロントの両方をセキュアにしようとするとそれだけで時間取られます。そこで使うのが認証基盤である Firebase Authentication です。これを使うと爆速で認証が作れます。そう、爆速です。Vue.js + Firebase の 解説記事はめちゃくちゃありますが、そこで返ってくるJWT を使ったAPIサーバーでの認証はあまり記事になかったので、今回はそこも含めたセキュアなWEBアプリケーションの実装をハンズオン形式で紹介します。タイトル通り、フロントエンドは Vue.js、バックエンドは Go言語、認証基盤として Firebase を使います。

御託はいいからコード見せろという方はこちらへどうぞ
https://github.com/po3rin/vue-golang-fireauth

## 今回使う技術の概要

### Vue.js

<img src="https://qiita-image-store.s3.amazonaws.com/0/186028/f854d45f-ac35-8092-ef62-cff61b4bb3ea.png" width=20%>

Vue.jsは現在、非常に人気のあるJavaScriptフレームワークです。GitHubでのstarはjQuery、Reactを抑え、現在最も注目されているJSフレームワークです。導入のし易さや、軽量感、学習コストが少ないことを売りにしている。

Vue.jsの入門記事を前に作ったので、初めて触る人はこちらも参考にしてみてください。
[VueCLIからVue.js入門①【VueCLIで出てくるファイルを概要図で理解】](https://qiita.com/po3rin/items/3968f825f3c86f9c4e21)

### Go言語

<img src="https://qiita-image-store.s3.amazonaws.com/0/186028/c5eebff9-91cc-fce6-a46f-2c697f5fc8d1.png" width=20%>

フロントエンド畑で育った僕が、なぜサーバーサイドにGo言語を選ぶかというと、主に下記の理由

* 構文 & 型システムがシンプル。動的言語しかやってこなかった僕でもすんなり受け入れてくれた。
* 並列処理が書きやすい。JavaScriptで苦戦した非同期処理が難なく書ける。
* 標準パッケージのサポートが強力。コードのフォーマットも自動でやってくれる。

とにかく開発者に無駄なことをさせないという気迫がすごい。ゆえに書いてて、読んでて気持ち良い。

### Firebase

<img src="https://qiita-image-store.s3.amazonaws.com/0/186028/dc9dd216-2e8e-f2f1-6fe0-3dba3f82da03.png" width=16%>

FirebaseはMBaas(mobile backend as a service)と呼ばれるものの一種で，オンラインでサインアップするだけでサーバのセットアップやメンテナンスに煩わされることなく使い始めることができ，リアルタイム同期型データベースやユーザ認証等の様々な機能が使えます。趣味で開発する規模であればだいたい無料で使えます。

### JWT

<img src="https://qiita-image-store.s3.amazonaws.com/0/186028/b3fe9aef-86d7-55a6-0b2b-f45abc8dcb80.jpeg" width=40%>

JWT(Json Web Token)と書いてジョットと読むらしい。このトークンを使ってクライアント＆サーバー間で認証できる。特徴としては、

* 発行者が鍵を使ってJSONを署名し、トークンとして使う。
* JSONなので任意の情報を含めることができる(ユーザー名など)。
* 発行者は鍵を使ってトークンの検証を行う為、改竄を検知できる。

上記の特徴から、僕が特にメリットに感じた点は「ステートレス」という点。サーバ側にセッションストアを持たなくても検証ができる。

JWTは以下のフォーマットに従った文字列です。

```
{base64エンコードしたhead１er}.{base64エンコードしたclaims}.{署名}
```

実際にこちらのサイトでJWTのエンコード、デコードが試せる。
https://jwt.io/

JWTをGo言語で解説した記事を前に作ったので、JWTの解説についてはこちらを参考にしてください。
[Go言語で理解するJWT認証 実装ハンズオン](https://qiita.com/po3rin/items/740445d21487dfcb5d9f)

## 実装するアーキテクチャ

サインアップはフロント側で行います。サインアップしたアカウントでサインインすると Firebase から JWT が返却されます。このJWTを使ってAPIサーバーにアクセスします。図にすると下記のような形になります。

<img width="859" alt="スクリーンショット 2018-06-18 18.00.38.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/69c4a410-0ff6-8cda-e1ce-32860424d162.png">


## Vue.jsでSPAを作る

### 必要なインストール
Vue.jsの開発環境を整えるためにvue-cliを使います。そのためnode.jsをインストールしておいてください。インストールされているか下記で確認しましょう。

```bash
$ node -v
v9.5.0

$ npm -v
v5.6.0
```

### Vueアプリケーションの雛形作成

vue-cliをインストールしましょう

```bash
$ npm install -g vue-cli
```
これでVueアプリケーションの雛形が作成できます。早速任意の場所で下記を実行。test-vueの部分はプロジェクト名&ディレクトリ名になります。

```bash
$ vue init webpack test-vue
```

いろいろ聞かれますが、お好みで設定を変えれます。全部EnterでもOKです。実行が終わったら下記を実行し、雛形を確認してみましょう。

```bash
$ cd test-vue
$ npm run dev
```

最終的にコマンドに出ているURLにブラウザからアクセスしてみてください。Vueアプリケーションの土台ができています。

### まずは認証なしSPAを作る

src/HelloWorld を認証後のマイページにすることを想定していきます。そのためにSigninページ、Signupページを追加で作ります。

まずはAPIを叩くために、Promise ベースの HTTPクライアントである axios を導入します。

```bash

$ npm install axios --save
```

これを使って、APIサーバーからデータを受け取って表示するマイページを作成します。
では早速、src/components/HelloWorld.vueを編集します。

```vue

<template>
  <div class='hello'>
    <h1>{{ msg }}</h1>
    <h2>Essential Links</h2>
    <button @click="apiPublic">public</button>
    <button @click="apiPrivate">private</button>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App'
    }
  },
  methods: {
    apiPublic: async function () {
      let res = await axios.get('http://localhost:8000/public')
      this.msg = res.data
    },
    apiPrivate: async function () {
      let res = await axios.get('http://localhost:8000/private')
      this.msg = res.data
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
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
button {
  margin: 10px 0;
  padding: 10px;
}
</style>
```

下のようになります

<img width="790" alt="スクリーンショット 2018-06-19 18.11.14.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/3818ad0b-732c-639d-df44-36035c75808e.png">

今はAPIを叩くボタンを押しても、まだAPIサーバーを作ってないのでエラーが出るはずです。

そして src/components/Signup.vue と src/components/Signin.vue を作成します。
まずは src/components/Signup.vue から

```vue

<template>
  <div class="signup">
    <h2>Sign up</h2>
    <input type="text" placeholder="Username" v-model="email">
    <input type="password" placeholder="Password" v-model="password">
    <button>Register</button>
    <p>Do you have an account?
      <router-link to="/signin">sign in now!!</router-link>
    </p>
  </div>
</template>

<script>
export default {
  name: 'Signup',
  data () {
    return {
      email: '',
      password: ''
    }
  }
}
</script>

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
.signup {
  margin-top: 20px;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center
}
input {
  margin: 10px 0;
  padding: 10px;
}
button {
  margin: 10px 0;
  padding: 10px;
}
</style>
```

下のようになっているはず。signInページへ遷移するリンクはつけていないので　'localhost:8080/#/signup'　のようにURLで直接みてみましょう。

<img width="739" alt="スクリーンショット 2018-06-19 18.05.05.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/02b9adc9-ce96-0cbb-f488-f6f8caf4a5c0.png">


まだサインアップはできません。後程つけていきます。次は src/components/Signin.vue です。

```vue

<template>
    <div class="signin">
        <h2>Sign in</h2>
        <input type="text" placeholder="email" v-model="email">
        <input type="password" placeholder="Password" v-model="password">
        <button>Signin</button>
        <p>You don't have an account?
            <router-link to="/signup">create account now!!</router-link>
        </p>
    </div>
</template>

<script>
export default {
  name: 'Signin',
  data: function () {
    return {
      email: '',
      password: ''
    }
  }
}
</script>

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
.signin {
  margin-top: 20px;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center
}
input {
  margin: 10px 0;
  padding: 10px;
}
button {
  margin: 10px 0;
  padding: 10px;
}
</style>
```

そしてこうなります。'localhost:8080/#/signin'　のようにURLで直接みてみましょう。

<img width="792" alt="スクリーンショット 2018-06-19 18.04.40.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/14c833c8-acca-f485-23c8-decbf2c3c6a4.png">

実際に動くサインイン機能は後程実装します。

最後に、作成したページへのルーティングを設定します。src/router/index.js を編集します。

```js
import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import Signup from '@/components/Signup'
import Signin from '@/components/Signin'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '*',
      redirect: 'signin'
    },
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/signup',
      name: 'Signup',
      component: Signup
    },
    {
      path: '/signin',
      name: 'Signin',
      component: Signin
    }
  ]
})
```

これでサインアップ用、サインイン用、マイページの3ページができました。
では早速APIサーバーの実装に入りましょう。

## Vue.jsから叩くAPIサーバーをGo言語で実装する

<img width="815" alt="スクリーンショット 2018-06-18 18.44.12.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/74fa4d1e-4c95-1c64-43cf-a218609c97a3.png">

### 簡易なHTTPサーバー実装

認証なしの簡易APIサーバーを作ります。今回は軽量なウェブツールキット gorilla/mux を使いましょう。
まずはmain.goを作成します。

```go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func public(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello public!\n"))
}

func private(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello private!\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/public", public)
	r.HandleFunc("/private", private)

	log.Fatal(http.ListenAndServe(":8000", r))
}
```

実際にAPIが叩けるか terminal 等から確認しましょう。

```bash
$ go run main.go

$ curl localhost:8080/public
hello public!

$ curl localhost:8080/private
hello private!
```

いいですね！簡単にHTTPサーバーをかけました。ただし、Vue.js からはまだ叩けません。CORSのエラーが出ます。サーバー側でCORSの設定をする必要があります。

### CORSの設定

main関数の中を編集します。私はlocalhost:8080を設定していますが、各自、Vue.jsが立ち上がっているIPアドレス + ポート番号を設定してください。更に今回は HTTP Repuest の Header に Authorizationヘッダーをつけるので、これも許可します。関数の最後の行でCORSを設定してます。CORSを設定する handlers は github.com/gorilla/handlers を使っているのでこちらもimportしておきましょう。

```go
// ...
import (
  // ...
  "github.com/gorilla/handlers"
)

// ...

func main() {
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization"})

	r := mux.NewRouter()
	r.HandleFunc("/public", public)
	r.HandleFunc("/private", private)

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)))
}
```

これで簡易的ですがAPIサーバーが出来ました。Vue.jsからAPIを叩いて画面の表示が変わることを確認しましょう。

## Firebase Authentication を使って認証機能を追加する

<img width="887" alt="スクリーンショット 2018-06-18 18.45.55.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/a263f340-8145-db0f-b350-ac0d8489d959.png">

では本題の認証をつけていきましょう。まずは上の形を目指します。

### Firebase を設定
https://firebase.google.com から初めます。
コンソールに入り、プロジェクトを新規作成します。

ウェブアプリにFirebaseを追加する をクリックします。
プロジェクトでFirebaseを使うために必要な設定項目が表示されるのでこれをコピーしておきます。

そして今回の実装例としてメールアドレスでの認証を行うので、Authenticationの メール/パスワード を有効にします。もちろんGoogleアカウント認証や、GitHub認証もここから設定できます。これでメールアドレスとパスワードによる認証の準備が出来ました。

Vue.js で Firebase を使うためのモジュールもここで install しておきましょう。

```bash
$ npm install firebase --save
```
src/main.jsに先ほどの設定を組み込みます。下記を自分の設定に書き換えてください。

```js
// ...
import firebase from 'firebase'

Vue.config.productionTip = false

const config = {
  apiKey: 'YOUR_KEY',
  authDomain: 'YOUR_DOMAIN.firebaseapp.com',
  databaseURL: 'YOUR_DOMAIN.firebaseio.com',
  projectId: 'YOUR_ID',
  storageBucket: 'YOUR_BUCKET_ID.appspot.com',
  messagingSenderId: 'YOUR_SENDER_ID'
}
firebase.initializeApp(config)

// ... Vueインスタンス作成　省略
```

これでFirebaseをVue.jsで使う準備が出来ました。

### サインアップ機能

早速 Signup 機能をつけていきましょう。Signup.vue を書き換えます。buttonタグにイベントと、SignUpメソッドを追加しています。

```vue

<template>
  <div class="signup">
    <h2>Sign up</h2>
    <input type="text" placeholder="Username" v-model="email">
    <input type="password" placeholder="Password" v-model="password">
    <button @click="signUp">Register</button>
    <p>Do you have an account?
      <router-link to="/signin">sign in now!!</router-link>
    </p>
  </div>
</template>

<script>
import firebase from 'firebase'
export default {
  name: 'Signup',
  data () {
    return {
      email: '',
      password: ''
    }
  },
  methods: {
    signUp: function () {
      firebase.auth().createUserWithEmailAndPassword(this.email, this.password).then(user => {
        console.log('Create account: ', user.email)
      }).catch(error => {
        console.log(error.message)
      })
    }
  }
}
</script>

<!-- css省略 -->
```

入力のあったメアドとパスワードを createUserWithEmailAndPassword()でFirebaseに送ってユーザーを作成しています。あとはこのアカウントでログインできる処理を追加します。

### サインイン機能

先ほどサインアップしたアカウントでログインできるようにします。Signin.vue を書き換えます。buttonタグにイベントと、signInメソッドを追加します。

```vue

<template>
    <div class="signin">
        <h2>Sign in</h2>
        <input type="text" placeholder="email" v-model="email">
        <input type="password" placeholder="Password" v-model="password">
        <button @click="signIn">Signin</button>
        <p>You don't have an account?
            <router-link to="/signup">create account now!!</router-link>
        </p>
    </div>
</template>

<script>
import firebase from 'firebase'
export default {
  name: 'Signin',
  data: function () {
    return {
      email: '',
      password: ''
    }
  },
  methods: {
    signIn: function () {
      firebase.auth().signInWithEmailAndPassword(this.email, this.password).then(res => {
        localStorage.setItem('jwt', res.user.qa)
        this.$router.push('/')
      }, err => {
        alert(err.message)
      })
    }
  }
}
</script>

<!-- css省略 -->
```

ここでのポイントはログイン時に返ってきたユーザー情報から、サーバーでの認証に使うJWT(res.user.qa)をローカルストレージに保管している点です。これを使って後ほどAPIサーバーの認証を突破します。ログインした後は "/" に遷移させます。

### サインアウト機能

サインアウト機能も簡単です。firebase.auth().signOut()を実行するだけです。この時に /signin にリダイレクトするようにしておきます。また、JWTもローカルストレージから削除します。また、せっかくのマイページなので、ユーザーのメアドを表示するようにしておきましょう。firebase.auth().currentUser.email で取得できます。

```html
<template>
  <div class='hello'>
    <h1>Hello {{ name }}!!</h1>
    <h1>{{ msg }}</h1>
    <h2>Essential Links</h2>
    <button @click="signOut">Sign out</button>
    <button @click="apiPublic">public</button>
    <button @click="apiPrivate">private</button>
  </div>
</template>

<script>
import axios from 'axios'
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
    },
    apiPublic: async function () {
      let res = await axios.get('http://localhost:8000/public')
      this.msg = res.data
    },
    apiPrivate: async function () {
      let res = await axios.get('http://localhost:8000/private')
      this.msg = res.data
    }
  }
}
</script>

<!-- css省略 -->
```

### vue-router で認証済みか確認

このままではログインしてなくてもURL直打ちでマイページに入れてしまいます。なので、認証が必要なルーターにログイン済みか確認するコードを追加する必要があります。src/router/index.js を編集しましょう。

```js

// ...
import firebase from 'firebase'

// ...

let router = new Router({
  routes: [
    {
      path: '*',
      redirect: 'signin'
    },
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld,
      meta: { requiresAuth: true }
    },
    {
      path: '/signup',
      name: 'Signup',
      component: Signup
    },
    {
      path: '/signin',
      name: 'Signin',
      component: Signin
    }
  ]
})

// router.beforeEach()を追加
router.beforeEach((to, from, next) => {
  let currentUser = firebase.auth().currentUser
  let requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  if (requiresAuth && !currentUser) next('signin')
  else if (!requiresAuth && currentUser) next()
  else next()
})

export default router

```

大事なのは '/' のみに設定した meta: { requiresAuth: true } です。これで、このrouteに認証が必要かを判断します。
firebase.auth().currentUser で現在ログインしているユーザーを返します。これでログインしていない場合は、/signin にリダイレクトされるようになりました。

### Vue.jsのライフサイクルに合うようにFirebaseを初期化する

Vue.jsでFirebase Authentication を使う際の最大のポイントはここになると思います。ここまでの実装ではログインした後でも、ブラウザを更新したら /signin　にリダイレクトされてしまいます。

理由としては、先ほど実装した beforeEach が Firebase が初期化される前に実行されるので、アプリケーションの最初のロード時に　firebase.auth().currentUser が null を返してしまう為です。

これを回避するために onAuthStateChanged を使います。これはユーザーの認証状況が変更されたら実行されるオブザーバーです。現在のユーザを取得したときにAuthオブジェクトが初期化などの中間状態にならないようにすることができます。これでVueインスタンス作成のコードを包むことで、Firebase の初期化後に Vueインスタンスが作成されるようになります。src/main.js を編集しましょう！

```js
// ...

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

// ...
```

これでログイン後にブラウザ更新しても /signin にリダイレクトされなくなりました。

## APIサーバーをJWT認証でセキュアにする

<img width="859" alt="スクリーンショット 2018-06-18 18.00.38.png" src="https://qiita-image-store.s3.amazonaws.com/0/186028/69c4a410-0ff6-8cda-e1ce-32860424d162.png">


いよいよ最終段階。最初に見せた上のような形まで持って行きます。
APIサーバーは現状、誰でも叩けるようになっています。アプリケーションとしては Vue.js でサインインに成功した人だけが API を叩けるようにしたいところです。そこで使うのが JWT です。先ほど、ローカルストレージに保存した JWT を使ってサインイン済みのユーザーか検証します。

### Firebase Admin SDK Go セットアップの準備

Go言語で Firebase を使うための設定をしていきます。まずは　Firebase Admin SDK Go　を使えるように設定していきます。まずは必要なパッケージを読み込み

```bash
$ go get -u firebase.google.com/go
$ go get -u google.golang.org/api/option
```

サービスアカウントの認証情報が含まれる JSON ファイル をGo言語で読み込んでFirebaseのセットアップが完了します！
このJSONファイルは Firebase Console または Google Cloud Consoleで取得可能です。僕は Cloud Platform Console で取得しました。

【Firebase Consoleの場合】
- プロジェクトの設定ページの [サービス アカウント] タブに移動し、[サービス アカウント] タブの [Firebase Admin SDK] セクション下部にある [新しい秘密鍵を生成] ボタンをクリックします。

【Cloud Platform Consolの場合】
- [IAM と管理] > [サービス アカウント] にアクセスします。そして、新しい秘密鍵を生成し、ローカルに JSON ファイルを保存します

さらに詳しいセットアップのやり方は公式ドキュメント「サーバーに Firebase Admin SDK を追加する」へ
https://firebase.google.com/docs/admin/setup?authuser=0

### JWT を Go言語 + Firebase で検証

JWTを検証するミドルウェアを作成します。このミドルウェアでハンドラーをラップしてあげれば、ラップした全てのAPIに検証機能がつきます。
先ほど作成した鍵ファイルへのパスは環境変数で読み込み、下記のコードでセットアップします。Go言語での環境変数は os.Getenv("環境変数名") で読み込めます。

```go

import (
    // ...

    firebase "firebase.google.com/go"
    "google.golang.org/api/option"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // Firebase SDK のセットアップ
		opt := option.WithCredentialsFile(os.Getenv("CREDENTIALS"))
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("error: %v\n", err)
                        os.Exit(1)
		}
		auth, err := app.Auth(context.Background())
		if err != nil {
			fmt.Printf("error: %v\n", err)
                        os.Exit(1)
		}

        // クライアントから送られてきた JWT 取得
		authHeader := r.Header.Get("Authorization")
                idToken := strings.Replace(authHeader, "Bearer ", "", 1)

                // JWT の検証
		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
                        // JWT が無効なら Handler に進まず別処理
                        fmt.Printf("error verifying ID token: %v\n", err)
                        w.Write([]byte("error verifying ID token\n"))
                        return
		}
		log.Printf("Verified ID token: %v\n", token)
		next.ServeHTTP(w, r)
	}
}

// ...

```

verifyIDToken()で JWT の検証を Firebase で行なっています。エラーハンドリング等は分かりやすさの為に簡易化していますが、これでミドルウェアは完成です。こいつでprivateハンドラーをラップします。

```go
func main() {
    // ...

    r.HandleFunc("/private", authMiddleware(private))

    // ...
}

```

これで有効なJWTをもつ人だけが /private からデータをもらうことができるようになりました。 実際にAPIをterminalから叩いて見ましょう

```
$ go run main.go

$ curl localhost:8080/public
hello public!

$ curl localhost:8080/privte
error verifying ID token
```

JWT をサーバーに送ってないので /private だけ求めていた'hello private!'が返ってきませんでした。Vue.js側で JWT をHeaderに入れてHTTPリクエストするメソッドを実装しましょう。HelloWorld.vue の apiPrivate関数を編集します。

```js
// ...
apiPrivate: async function () {
    let res = await axios.get('http://localhost:8000/private', {
    headers: {'Authorization': `Bearer ${localStorage.getItem('jwt')}`}
    })
    this.msg = res.data
}
// ...
```

ここでは axiosで Authorization Headerをつけています。ローカルストレージに保管してあるJWTをサーバーに送って検証してもらいます。ここまできたらprivateボタンを推して、データが返ってくるか確認して見ましょう。

これで Vue.js + Go言語 + Firebase で認証付きWEBアプリケーションが完成しました！

## まとめ

Vue.js + Go言語 + Firebase を使って簡単に認証付きWEBアプリケーションができました。クライアントとサーバーが分離しているので、他のマイクロサービスと連携させるのも楽だと思います。ここからユーザーごとにMySQL等でデータを持たせたい場合は、JWT の中の sub を uidとして保存し、ユーザーを識別する形も取れるようです。これからはこれをベースにガンガン開発できます。 もし詰まった箇所があればコメント欄で教えてください！

## 参考記事

[Firebase Authentication 日本語ドキュメント](https://firebase.google.com/docs/auth/?hl=ja)
[Vue 2 + Firebase: How to build a Vue app with Firebase authentication system in 15 minutes](https://medium.com/@anas.mammeri/vue-2-firebase-how-to-build-a-vue-app-with-firebase-authentication-system-in-15-minutes-fdce6f289c3c)
[Vue.js + Firebase を使って爆速でユーザ認証を実装する](https://qiita.com/sin_tanaka/items/ea149a33bd9e4b388241)
[遂にFirebase Admin SDK Goが登場！](https://qiita.com/koki_cheese/items/2d111b2b074bfa697776)
[Goで始めるMiddleware](https://qiita.com/tnakata/items/ea962f1cdad21c2f68aa)
