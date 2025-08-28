<script setup>
import { getCurrentInstance, onMounted, ref, onBeforeMount } from 'vue'
import { useColorModes } from '@coreui/vue'

import { useThemeStore } from '@/stores/theme.js'

import { useRouter } from 'vue-router'

const router = useRouter()

const self = ref('Vue Dev')
const users = ref([])


const { isColorModeSet, setColorMode } = useColorModes(
  'coreui-free-vue-admin-template-theme',
)
const currentTheme = useThemeStore()

const socket = getCurrentInstance().appContext.config.globalProperties.$socket;
const StudioWebManager = getCurrentInstance().appContext.config.globalProperties.$StudioWebManager;


socket.on('connect', (data) => {
  console.log('Received: Connected to the socket server', data);

});

socket.on('disconnect', () => {
  console.log('desiconned to the socket server');
  sendMessageToServer()
});

socket.on('sid-auth', (data) => {
  let response = StudioWebManager.set_sio_sid(data);
  console.log('SIO SID set:', response, data);
});

socket.connect();

onBeforeMount(() => {
  const urlParams = new URLSearchParams(window.location.href.split('?')[1])
  let theme = urlParams.get('theme')

  if (theme !== null && theme.match(/^[A-Za-z0-9\s]+/)) {
    theme = theme.match(/^[A-Za-z0-9\s]+/)[0]
  }

  if (theme) {
    setColorMode(theme)
    return
  }

  if (isColorModeSet()) {
    return
  }

  setColorMode(currentTheme.theme)
})



onMounted(() => {
  try {
    // const is_authenticated = sessionStorage.getItem('is_authenticated');
    // if (!is_authenticated && !window.location.href.includes('/auth/login')){
    //   router.push('/auth/login')
    // }
    // users.value = response.data
  } catch (error) {
    console.error('Failed to fetch users:', error)
  }
})
</script>

<template>
  <router-view />
</template>

<style lang="scss">
// Import Main styles for this application
@use 'styles/style';
// We use those styles to show code examples, you should remove them in your application.
@use 'styles/examples';
</style>
