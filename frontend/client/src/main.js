import Vue from 'vue'
import Axios from './plugins/axios'
import './plugins/vuetify'
import App from './App.vue'
import router from './router'
import store from './store'
import VueLogger from 'vuejs-logger';
import Notifications from 'vue-notification'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'

const isProduction = process.env.NODE_ENV === 'production';

const options = {
  isEnabled: true,
  logLevel : isProduction ? 'error' : 'debug',
  stringifyArguments : false,
  showLogLevel : true,
  showMethodName : true,
  separator: '|',
  showConsoleColors: true
};


Vue.use(Notifications);
Vue.use(VueLogger, options);
Vue.prototype.$http = Axios;

// const token = localStorage.getItem('token');
// if (token) {
//   Vue.prototype.$http.defaults.headers.common['Authorization'] = token
// }

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
