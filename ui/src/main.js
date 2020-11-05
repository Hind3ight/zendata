import Vue from 'vue'
import App from './App.vue'

import Button from "ant-design-vue/lib/button";
import 'ant-design-vue/lib/button/style';

import router from "./router"

Vue.config.productionTip = false

Vue.component(Button.name, Button);
Vue.component(Button.Group.name, Button.Group);

Vue.use(Button);

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
