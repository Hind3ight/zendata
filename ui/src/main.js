import Vue from 'vue'
import App from './App.vue'

// import { Button, message } from 'ant-design-vue';

import VueI18n from 'vue-i18n'
import ConfigProvider from "ant-design-vue/lib/config-provider";
import Button from "ant-design-vue/lib/button";
import 'ant-design-vue/lib/button/style';

import Menu from "ant-design-vue/lib/menu";
import SubMenu from "ant-design-vue/lib/menu";
import MenuItem from "ant-design-vue/lib/menu";
import 'ant-design-vue/lib/menu/style';

import Table from "ant-design-vue/lib/table";
import 'ant-design-vue/lib/table/style';

import Tag from "ant-design-vue/lib/tag";
import Divider from "ant-design-vue/lib/divider";
import Icon from "ant-design-vue/lib/icon";

import zhCN from './assets/lang/zh-CN'
import router from "./router"

Vue.config.productionTip = false

Vue.component(Button)
Vue.use(VueI18n)
Vue.use(ConfigProvider)
Vue.use(Menu)
Vue.use(SubMenu)
Vue.use(MenuItem)
Vue.use(Button)
Vue.use(Table)
Vue.use(Tag)
Vue.use(Icon)
Vue.use(Divider)

const i18n = new VueI18n({
  locale: 'zh-CN',
  messages: {
    'zh-CN': { ...zhCN }
  },
});

new Vue({
  router,
  i18n,
  render: h => h(App),
}).$mount('#app')
