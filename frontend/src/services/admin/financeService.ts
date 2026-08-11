import api from '@/services/api'

export interface ChartOfAccount {
  code: string
  name: string
  type: string
  parent_code?: string
  created_at?: string
  updated_at?: string
}

export interface CashTransaction {
  id?: string
  date: string
  description: string
  type: string // 'in' | 'out'
  amount: number
  account_id: string
  account?: ChartOfAccount
  created_at?: string
}

export interface GeneralJournal {
  id?: string
  date: string
  description: string
  account_code: string
  account?: ChartOfAccount
  debit: number
  credit: number
  reservation_id?: string
  created_at?: string
}

export interface Invoice {
  id?: string
  invoice_number?: string
  type: string // "AP" or "AR"
  partner_name: string
  date: string
  due_date: string
  amount: number
  status?: string // "Unpaid", "Paid"
  description: string
  created_at?: string
}

export const financeService = {
  // COA
  async getCOAs() {
    const response = await api.get('/finance/coa')
    return response.data.data as ChartOfAccount[]
  },
  async createCOA(data: ChartOfAccount) {
    const response = await api.post('/finance/coa', data)
    return response.data.data
  },
  async updateCOA(code: string, data: Partial<ChartOfAccount>) {
    const response = await api.put(`/finance/coa/${code}`, data)
    return response.data.data
  },
  async deleteCOA(code: string) {
    const response = await api.delete(`/finance/coa/${code}`)
    return response.data.data
  },
  async importCOA(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await api.post('/finance/coa/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data
  },

  // Cash & Bank
  async getCashTransactions() {
    const response = await api.get('/finance/cash')
    return response.data.data as CashTransaction[]
  },
  async createCashTransaction(data: CashTransaction) {
    const response = await api.post('/finance/cash', data)
    return response.data.data
  },
  async deleteCashTransaction(id: string) {
    const response = await api.delete(`/finance/cash/${id}`)
    return response.data.data
  },

  // General Journal
  async getGeneralJournals() {
    const response = await api.get('/finance/journal')
    return response.data.data as GeneralJournal[]
  },
  async createGeneralJournal(data: GeneralJournal) {
    const response = await api.post('/finance/journal', data)
    return response.data.data
  },

  // Verify Payment
  async verifyPayment(reservation_id: string, status: 'confirmed' | 'rejected') {
    const response = await api.post('/finance/verify-payment', { reservation_id, status })
    return response.data.data
  },

  // Invoices
  async getInvoices() {
    const response = await api.get('/finance/invoices')
    return response.data.data as Invoice[]
  },
  async createInvoice(data: Invoice) {
    const response = await api.post('/finance/invoices', data)
    return response.data.data
  },
  async payInvoice(id: string, account_id: string) {
    const response = await api.post(`/finance/invoices/${id}/pay`, { account_id })
    return response.data
  }
}
