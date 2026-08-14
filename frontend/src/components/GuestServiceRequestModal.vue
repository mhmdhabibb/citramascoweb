<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { serviceRequestService } from '@/services/admin/serviceRequestService'
import { roomService } from '@/services/roomService'
import { useToastStore } from '@/stores/toastStore'
import type { Room, ServiceRequestCategory } from '@/types'

const props = defineProps<{
  isOpen: boolean
  defaultRoomId?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const toastStore = useToastStore()
const loading = ref(false)
const rooms = ref<Room[]>([])
const isSubmitted = ref(false)

const form = ref({
  room_id: '',
  guest_name: '',
  guest_phone: '',
  category: 'incident_broken_item' as ServiceRequestCategory,
  title: '',
  description: '',
})

const categoryOptions = [
  { value: 'incident_broken_item', label: '🍷 Gelas / Barang Pecah (Insiden)', defaultTitle: 'Gelas/Barang Pecah di Kamar' },
  { value: 'extra_cleaning', label: '🧹 Pembersihan Ekstra / Tumpahan Air', defaultTitle: 'Permintaan Pembersihan Tambahan' },
  { value: 'amenities_request', label: '🪥 Tambahan Handuk / Amenities', defaultTitle: 'Permintaan Handuk / Perlengkapan' },
  { value: 'maintenance_repair', label: '🛠️ Kerusakan Fasilitas (AC / Lampu / Air)', defaultTitle: 'Laporan Kerusakan Fasilitas' },
  { value: 'other', label: '❓ Bantuan Lainnya', defaultTitle: 'Permintaan Bantuan Khusus' },
]

onMounted(async () => {
  try {
    const data = await roomService.getAll()
    rooms.value = data || []
    if (props.defaultRoomId) {
      form.value.room_id = props.defaultRoomId
    } else if (rooms.value.length > 0 && rooms.value[0]) {
      form.value.room_id = rooms.value[0].id
    }
  } catch (error) {
    console.error('Failed to fetch rooms', error)
  }
})

const onCategoryChange = () => {
  const selected = categoryOptions.find(c => c.value === form.value.category)
  if (selected && !form.value.title) {
    form.value.title = selected.defaultTitle
  }
}

const handleSubmit = async () => {
  if (!form.value.room_id || !form.value.guest_name || !form.value.title || !form.value.description) {
    toastStore.error('Mohon lengkapi semua kolom yang wajib diisi.')
    return
  }

  try {
    loading.value = true
    await serviceRequestService.create({
      room_id: form.value.room_id,
      guest_name: form.value.guest_name,
      guest_phone: form.value.guest_phone,
      category: form.value.category,
      title: form.value.title,
      description: form.value.description,
    })

    isSubmitted.value = true
    toastStore.success('Laporan berhasil dikirim ke Resepsionis!')
    emit('success')
  } catch (error: any) {
    toastStore.error(error.message || 'Gagal mengirim laporan')
  } finally {
    loading.value = false
  }
}

const resetAndClose = () => {
  isSubmitted.value = false
  const fallbackRoomId = rooms.value.length > 0 && rooms.value[0] ? rooms.value[0].id : ''
  form.value = {
    room_id: props.defaultRoomId || fallbackRoomId,
    guest_name: '',
    guest_phone: '',
    category: 'incident_broken_item',
    title: '',
    description: '',
  }
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="resetAndClose">
    <div class="modal-card">
      <div class="modal-header">
        <div class="header-icon-wrap">
          <span class="icon">🛎️</span>
        </div>
        <div>
          <h2>Layanan Kamar & Bantuan Tamu</h2>
          <p>Hubungi Resepsionis untuk permintaan pembersihan, insiden, atau perlengkapan kamar</p>
        </div>
        <button class="close-btn" @click="resetAndClose">✕</button>
      </div>

      <!-- Sukses State -->
      <div v-if="isSubmitted" class="success-body">
        <div class="success-icon">✓</div>
        <h3>Laporan Anda Telah Diterima!</h3>
        <p>
          Resepsionis kami telah menerima laporan Anda dan segera menugaskan tim <strong>Housekeeping</strong> ke kamar Anda.
        </p>
        <button class="btn btn-primary" @click="resetAndClose">Tutup Jendela</button>
      </div>

      <!-- Form Body -->
      <form v-else @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-grid-2">
          <div class="form-group">
            <label>Pilih Unit Kamar <span class="req">*</span></label>
            <select v-model="form.room_id" class="input-field" required>
              <option v-for="r in rooms" :key="r.id" :value="r.id">
                {{ r.name }} ({{ r.code }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>Nama Tamu <span class="req">*</span></label>
            <input
              v-model="form.guest_name"
              type="text"
              class="input-field"
              placeholder="Contoh: Pak Budi"
              required
            />
          </div>
        </div>

        <div class="form-grid-2">
          <div class="form-group">
            <label>Jenis Bantuan / Insiden <span class="req">*</span></label>
            <select v-model="form.category" @change="onCategoryChange" class="input-field" required>
              <option v-for="opt in categoryOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>No. Kontak / WhatsApp (Opsional)</label>
            <input
              v-model="form.guest_phone"
              type="tel"
              class="input-field"
              placeholder="0812xxxx"
            />
          </div>
        </div>

        <div class="form-group">
          <label>Judul Laporan <span class="req">*</span></label>
          <input
            v-model="form.title"
            type="text"
            class="input-field"
            placeholder="Contoh: Gelas pecah di lantai dekat meja"
            required
          />
        </div>

        <div class="form-group">
          <label>Rincian Masalah / Permintaan <span class="req">*</span></label>
          <textarea
            v-model="form.description"
            rows="3"
            class="input-field"
            placeholder="Jelaskan kebutuhan Anda agar tim kami dapat membawa peralatan yang tepat..."
            required
          ></textarea>
        </div>

        <div class="notice-box">
          <span>💡</span>
          <span>Laporan Anda akan otomatis diteruskan ke tim Frontdesk & Housekeeping secara real-time.</span>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="resetAndClose" :disabled="loading">
            Batal
          </button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? 'Mengirim...' : 'Kirim ke Resepsionis' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.7);
  backdrop-filter: blur(6px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-card {
  background: white;
  width: 100%;
  max-width: 600px;
  border-radius: 1.25rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  overflow: hidden;
  animation: scaleUp 0.2s ease-out;
}

@keyframes scaleUp {
  from { transform: scale(0.95); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

.modal-header {
  padding: 1.5rem;
  background: linear-gradient(135deg, #1e293b, #0f172a);
  color: white;
  display: flex;
  align-items: center;
  gap: 1rem;
  position: relative;
}

.header-icon-wrap {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
  flex-shrink: 0;
}

.modal-header h2 {
  font-size: 1.15rem;
  font-weight: 700;
  margin: 0;
}

.modal-header p {
  font-size: 0.78rem;
  color: #94a3b8;
  margin: 2px 0 0 0;
}

.close-btn {
  position: absolute;
  right: 1.25rem;
  top: 1.25rem;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.modal-body {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

@media (max-width: 600px) {
  .form-grid-2 {
    grid-template-columns: 1fr;
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 700;
  color: #334155;
}

.req {
  color: #ef4444;
}

.input-field {
  padding: 0.65rem 0.85rem;
  border: 1.5px solid #cbd5e1;
  border-radius: 0.6rem;
  font-size: 0.85rem;
  color: #1e293b;
  outline: none;
  transition: border-color 0.2s;
}

.input-field:focus {
  border-color: #6366f1;
}

.notice-box {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 0.6rem;
  padding: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.6rem;
  font-size: 0.75rem;
  color: #64748b;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.btn {
  padding: 0.65rem 1.25rem;
  border-radius: 0.6rem;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
}

.btn-secondary:hover {
  background: #e2e8f0;
}

.btn-primary {
  background: #6366f1;
  color: white;
}

.btn-primary:hover {
  background: #4f46e5;
}

.success-body {
  padding: 3rem 2rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.success-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #ecfdf5;
  color: #059669;
  border: 2px solid #a7f3d0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  font-weight: bold;
}

.success-body h3 {
  font-size: 1.25rem;
  color: #0f172a;
  margin: 0;
}

.success-body p {
  font-size: 0.85rem;
  color: #64748b;
  max-width: 420px;
  line-height: 1.5;
  margin: 0;
}
</style>
