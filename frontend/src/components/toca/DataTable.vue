<script setup>
import { ref, computed, watch, toRefs } from 'vue'
import { CButton, CFormInput, CCard, CCardBody } from '@coreui/vue'
import { utils as XLSXUtils, writeFile } from 'xlsx'
import { saveAs } from 'file-saver'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'

const props = defineProps({
  data: {
    type: Array,
    required: true,
  },
  columns: {
    type: Array,
    required: true,
    // [{ field: 'name', label: 'Name' }, ...]
  },
  rowsPerPage: {
    type: Number,
    default: 10,
  },
})

const emit = defineEmits(['update:selected'])

const { data, columns, rowsPerPage } = toRefs(props)

const currentPage = ref(1)
const globalSearch = ref('')

// Pagination logic
const totalPages = computed(() =>
  Math.ceil(filteredData.value.length / rowsPerPage.value)
)

const changePage = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
}

// Filtering
const filteredData = computed(() => {
  if (!globalSearch.value) return data.value
  const search = globalSearch.value.toLowerCase()
  return data.value.filter((item) =>
    columns.value.some((col) => {
      const val = String(resolveField(item, col.field)).toLowerCase()
      return val.includes(search)
    })
  )
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * rowsPerPage.value
  return filteredData.value.slice(start, start + rowsPerPage.value)
})

// Utility to resolve nested fields like 'country.name'
function resolveField(obj, path) {
  return path.split('.').reduce((acc, part) => acc && acc[part], obj)
}

// Format customer row for export
function formatRow(row) {
  const formatted = {}
  columns.value.forEach(({ field, label }) => {
    const val = resolveField(row, field)
    formatted[label] = val === true ? 'Yes' : val === false ? 'No' : val
  })
  return formatted
}

// Export helpers
function exportToCSV() {
  const headers = columns.value.map((col) => col.label).join(',')
  const rows = filteredData.value.map(formatRow)
  const csvData = [
    headers,
    ...rows.map((row) =>
      columns.value
        .map((col) => `"${(row[col.label] ?? '').toString().replace(/"/g, '""')}"`)
        .join(',')
    ),
  ].join('\n')

  const blob = new Blob([csvData], { type: 'text/csv;charset=utf-8;' })
  saveAs(blob, 'data.csv')
}

function exportToXLSX() {
  const rows = filteredData.value.map(formatRow)
  const worksheet = XLSXUtils.json_to_sheet(rows)
  const workbook = XLSXUtils.book_new()
  XLSXUtils.book_append_sheet(workbook, worksheet, 'Data')
  writeFile(workbook, 'data.xlsx')
}

function exportToJSON() {
  const jsonData = JSON.stringify(filteredData.value.map(formatRow), null, 2)
  const blob = new Blob([jsonData], { type: 'application/json' })
  saveAs(blob, 'data.json')
}

function exportToPDF() {
  const doc = new jsPDF()
  autoTable(doc, {
    head: [columns.value.map((c) => c.label)],
    body: filteredData.value.map((row) => {
      const formatted = formatRow(row)
      return columns.value.map((col) => formatted[col.label])
    }),
    styles: { fontSize: 8 },
  })
  doc.save('data.pdf')
}
</script>

<template>
  <CCard>
    <CCardBody>
      <div class="d-flex justify-content-between flex-wrap gap-2 mb-3">
        <div class="btn-group" role="group">
          <CButton color="secondary" size="sm" @click="exportToCSV">CSV</CButton>
          <CButton color="secondary" size="sm" @click="exportToXLSX">XLSX</CButton>
          <CButton color="secondary" size="sm" @click="exportToJSON">JSON</CButton>
          <CButton color="secondary" size="sm" @click="exportToPDF">PDF</CButton>
        </div>

        <CFormInput
          v-model="globalSearch"
          placeholder="Global Search"
          class="w-50"
        />
      </div>

      <table class="table table-striped">
        <thead>
          <tr>
            <th v-for="col in columns" :key="col.field">{{ col.label }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, index) in paginatedData" :key="index">
            <td v-for="col in columns" :key="col.field">
              {{ resolveField(row, col.field) }}
            </td>
          </tr>
          <tr v-if="paginatedData.length === 0">
            <td :colspan="columns.length" class="text-center">No data found.</td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <nav aria-label="Page navigation" class="mt-3">
        <ul class="pagination justify-content-center mb-0">
          <li
            class="page-item"
            :class="{ disabled: currentPage === 1 }"
          >
            <button
              class="page-link"
              @click.prevent="changePage(currentPage - 1)"
              :disabled="currentPage === 1"
            >
              Previous
            </button>
          </li>

          <li
            v-for="page in totalPages"
            :key="page"
            class="page-item"
            :class="{ active: currentPage === page }"
          >
            <button
              class="page-link"
              @click.prevent="changePage(page)"
            >
              {{ page }}
            </button>
          </li>

          <li
            class="page-item"
            :class="{ disabled: currentPage === totalPages }"
          >
            <button
              class="page-link"
              @click.prevent="changePage(currentPage + 1)"
              :disabled="currentPage === totalPages"
            >
              Next
            </button>
          </li>
        </ul>
      </nav>
    </CCardBody>
  </CCard>
</template>

<style scoped>
/* Optional custom styles */
</style>
