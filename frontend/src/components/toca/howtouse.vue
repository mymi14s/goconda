<script setup>
import DataTableExport from './components/DataTableExport.vue'

const customers = [
  {
    id: 1,
    name: 'Amy Elsner',
    country: { name: 'USA' },
    representative: { name: 'Agent A' },
    status: 'new',
    verified: true,
  },
  {
    id: 2,
    name: 'Bernardo Dominic',
    country: { name: 'Canada' },
    representative: { name: 'Agent B' },
    status: 'qualified',
    verified: false,
  },
  // ... more data
]

const columns = [
  { field: 'name', label: 'Name' },
  { field: 'country.name', label: 'Country' },
  { field: 'representative.name', label: 'Agent' },
  { field: 'status', label: 'Status' },
  { field: 'verified', label: 'Verified' },
]
</script>

<template>
  <DataTableExport :data="customers" :columns="columns" :rowsPerPage="10" />
</template>
