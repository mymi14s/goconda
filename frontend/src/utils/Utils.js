import { unref } from "vue"
import alertify from 'alertifyjs';
import { utils as XLSXUtils, writeFile } from 'xlsx'
import { saveAs } from 'file-saver'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'

export class Validator {
  constructor() {}

  isValidEmail(email) {
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
  }

  error(text) {
    alertify.alert('Error', text).set('modal', true)
  }

  success(text) {
    alertify.alert('Success', text).set('modal', true)
  }
}

export class DocExport {
  constructor() {}


  // Export to CSV
  toCSV = (formatData, dataList, exportFields, filename) => {

    const headers = exportFields.join(',')
    const rows = unref(dataList).map(formatData)
    const csvData = [
      headers,
      ...rows.map((row) =>
        exportFields.map((field) => `"${row[field]}"`).join(',')
      ),
    ].join('\n')

    const blob = new Blob([csvData], { type: 'text/csv;charset=utf-8;' })
    saveAs(blob, `${filename}.csv`)
  }

  // Export to XLSX
  toXLSX (formatData, dataList, exportFields, filename) {
    const rows = unref(dataList).map(formatData)
    const worksheet = XLSXUtils.json_to_sheet(rows)
    const workbook = XLSXUtils.book_new()
    XLSXUtils.book_append_sheet(workbook, worksheet, 'Customers')
    writeFile(workbook, filename+'.xlsx')
  }

  // Export to JSON
  toJSON (formatData, dataList, exportFields, filename) {
    const jsonData = JSON.stringify(
      unref(dataList).map(formatData),
      null,
      2
    )
    const blob = new Blob([jsonData], { type: 'application/json' })
    saveAs(blob, filename+'.json')
  }

  // Export to PDF
  toPDF(formatData, dataList, exportFields, filename) {
    const doc = new jsPDF()
    autoTable(doc, {
      head: [exportFields.map((f) => f.toUpperCase())],
      body: unref(dataList).map((c) => {
        const f = formatData(c)
        return exportFields.map((field) => f[field])
      }),
      styles: { fontSize: 8 },
    })
    doc.save(filename+'.pdf')
  }
}

export default { Validator, DocExport}