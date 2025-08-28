import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import PrimeVue from 'primevue/config';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column'
import InputText from 'primevue/inputtext';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';


import App from './App.vue'
import router from './router'

import CoreuiVue from '@coreui/vue'
import CIcon from '@coreui/icons-vue'
import { iconsSet as icons } from '@/assets/icons'
import DocsComponents from '@/components/DocsComponents'
import DocsExample from '@/components/DocsExample'
import DocsIcons from '@/components/DocsIcons'


// import 'alertifyjs/build/css/alertify.min.css';
// import 'alertifyjs/build/alertify.min.js';
import axios from 'axios';
import alertify from 'alertifyjs';


import StudioWebManager from './utils/StudioWebManager';
import Utils from './utils/Utils';
import socket from './utils/Socket';

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)



const app = createApp(App)
app.use(pinia)
app.use(router)
app.use(CoreuiVue)
app.use(PrimeVue, { unstyled: true })
app.provide('icons', icons)
app.component('CIcon', CIcon)
app.component('DocsComponents', DocsComponents)
app.component('DocsExample', DocsExample)
app.component('DocsIcons', DocsIcons)
app.component('DataTable', DataTable)
app.component('Column', Column)
app.component('InputText', InputText);
app.component('Dropdown', Dropdown);
app.component('InputNumber', InputNumber);


app.config.globalProperties.$StudioWebManager = new StudioWebManager('/');
app.config.globalProperties.$validator = new Utils.Validator();
app.config.globalProperties.$docExport = new Utils.DocExport();
app.config.globalProperties.$socket = socket;

app.mount('#app')
