<script setup>
import { getCurrentInstance } from 'vue'
import avatar from '@/assets/images/avatars/8.jpg'
const socket = getCurrentInstance().appContext.config.globalProperties.$socket;
const StudioWebManager = getCurrentInstance().appContext.config.globalProperties.$StudioWebManager;


const itemsCount = 42
</script>

<template>
  <CDropdown placement="bottom-end" variant="nav-item">
    <CDropdownToggle class="py-0 pe-0" :caret="false">
      <CAvatar :src="avatar" size="md" />
    </CDropdownToggle>
    <CDropdownMenu class="pt-0">
      <CDropdownHeader
        component="h6"
        class="bg-body-secondary text-body-secondary fw-semibold mb-2 rounded-top"
      >
        Account
      </CDropdownHeader>
      
      <CDropdownDivider />
      <CDropdownItem> <CIcon icon="cil-shield-alt" /> Lock Account </CDropdownItem>
      <CDropdownItem @click="StudioWebManager.logout()"> <CIcon icon="cil-lock-locked" /> Logout </CDropdownItem>
    </CDropdownMenu>
  </CDropdown>
</template>
