import os
import re

replacements = {
    r'"nomor kamar sudah digunakan"': '"room number is already in use"',
    r'"unit kamar tidak ditemukan"': '"room unit not found"',
    r'"kamar tidak ditemukan"': '"room not found"',
    r'"laporan tidak ditemukan"': '"report not found"',
    r'"🚨 Laporan Tamu \(\%s\): \%s"': '"🚨 Guest Report (%s): %s"',
    r'"Tamu \%s di kamar \%s melaporkan: \%s\. Mohon segera ditinjau dan ditugaskan\."': '"Guest %s in room %s reported: %s. Please review and assign immediately."',
    r'"Kamar"': '"Room"',
    r'"Gagal memuat data dari database"': '"Failed to load data from database"',
    r"'Unit Kamar'": "'Room Unit'",
    r"'Dibatalkan'": "'Cancelled'",
    r"'Malam'": "'Nights'",
    r"'Okupansi'": "'Occupancy'",
    r"'Tidak ada kamar di database'": "'No rooms in database'",
    r"'Nama Tamu'": "'Guest Name'",
    r"'Kamar'": "'Room'",
    r"'Durasi'": "'Duration'",
    r"'Total Biaya'": "'Total Cost'",
    r"'Gagal mengekspor dokumen PDF'": "'Failed to export PDF document'",
    r'"Total Unit Kamar:"': '"Total Room Units:"',
    r'"=== RINCIAN PERFORMA PER UNIT KAMAR ==="': '"=== ROOM UNIT PERFORMANCE DETAILS ==="',
    r'"Nama Kamar,Jumlah Reservasi,Dibatalkan,Total Malam,Occupancy Rate \(%\),Total Pendapatan \(IDR\),ADR \(IDR\),RevPAR \(IDR\)"': '"Room Name,Reservations,Cancelled,Total Nights,Occupancy Rate (%),Total Revenue (IDR),ADR (IDR),RevPAR (IDR)"',
    r'"Kode Booking,Nama Tamu,Kamar,Check-In,Check-Out,Malam,Tamu,Channel,Status,Total Harga \(IDR\)"': '"Booking Code,Guest Name,Room,Check-In,Check-Out,Nights,Guests,Channel,Status,Total Price (IDR)"',
    r"'Gagal mengekspor file CSV'": "'Failed to export CSV file'",
    r'Unit Kamar Terdaftar di Database': 'Registered Room Units in Database',
    r'<span class="kpi-unit">tamu</span>': '<span class="kpi-unit">guest(s)</span>',
    r'<span class="kpi-unit">kamar</span>': '<span class="kpi-unit">room(s)</span>',
    r'Tidak ada jadwal kedatangan tamu hari ini di database': 'No guest arrivals scheduled today',
    r'Tidak ada kepulangan tamu hari ini di database': 'No guest departures scheduled today',
    r'Daftar tamu yang akan datang dari database': 'Upcoming guest arrivals',
    r'Belum ada jadwal kedatangan tamu di database': 'No upcoming guest arrivals',
    r'Semua Kamar \(': 'All Rooms (',
    r'Malam kamar terjual': 'Room nights sold',
    r'Tabel rincian performa okupansi, ADR, dan pendapatan tiap kamar dari database': 'Room performance, occupancy, ADR, and revenue details',
    r'Belum ada kamar terdaftar di database': 'No rooms registered in database',
    r'Tingkat hunian kamar harian pada bulan': 'Daily room occupancy rate for month',
    r'Rekap riwayat check-in & check-out harian serta catatan kedatangan tamu\.': 'Daily check-in & check-out history summary and guest arrival notes.',
    r'\(Denah Unit Kamar Real-time\)': '(Live Physical Room Matrix)',
    r'Mengubah Status Fisik Kamar': 'Change Physical Room Status',
    r"'Pilih tipe katalog kamar!'": "'Select room catalog type!'",
    r"'Pilih lantai kamar!'": "'Select room floor!'",
    r'No\. Kamar': 'Room No.',
    r'Catat Permintaan Tamu Baru': 'Record New Guest Request',
    r"'Tamu meminta pembersihan tambahan / ganti sprei kamar\.'": "'Guest requested extra cleaning / bed sheet change.'",
    r"'Silakan pilih kamar terlebih dahulu\.'": "'Please select a room first.'",
    r'Jumlah tamu dewasa melebihi kapasitas kamar': 'Number of adults exceeds room capacity',
    r'Jumlah anak melebihi kapasitas kamar': 'Number of children exceeds room capacity',
}

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    new_content = content
    for pattern, replacement in replacements.items():
        new_content = re.sub(pattern, replacement, new_content, flags=re.IGNORECASE)
        
    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Updated {filepath}")

for root, dirs, files in os.walk('.'):
    if 'node_modules' in dirs: dirs.remove('node_modules')
    if '.git' in dirs: dirs.remove('.git')
    for f in files:
        if f.endswith('.go') or f.endswith('.vue') or f.endswith('.ts'):
            process_file(os.path.join(root, f))

