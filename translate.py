import os
import re

replacements = {
    r'"Lantai berhasil dibuat"': '"Floor created successfully"',
    r'"Lantai berhasil diperbarui"': '"Floor updated successfully"',
    r'"Lantai berhasil dihapus"': '"Floor deleted successfully"',
    r'"Proses check-in berhasil, status kamar diubah menjadi Occupied"': '"Check-in successful, room status changed to Occupied"',
    r'"Proses check-out berhasil, status kamar kembali Available"': '"Check-out successful, room status back to Available"',
    r'"Unit nomor kamar berhasil ditambahkan"': '"Room unit added successfully"',
    r'"Unit nomor kamar berhasil diperbarui"': '"Room unit updated successfully"',
    r'"Status kamar berhasil diubah"': '"Room status changed successfully"',
    r'"Unit nomor kamar berhasil dihapus"': '"Room unit deleted successfully"',
    r'"Tamu Check In"': '"Guest Check-in"',
    r'"Tamu Check Out"': '"Guest Check-out"',
    r'"Berhasil mengambil data dashboard"': '"Dashboard data retrieved successfully"',
    r'"Laporan layanan kamar berhasil dikirim ke Resepsionis."': '"Room service report successfully sent to Receptionist."',
    r'"Tugas berhasil diteruskan ke tim Housekeeping."': '"Task successfully forwarded to Housekeeping team."',
    r'"Tugas berhasil diselesaikan oleh Housekeeping."': '"Task successfully completed by Housekeeping."',
    r'"Laporan layanan kamar berhasil dibatalkan."': '"Room service report successfully cancelled."',
    r'"Data lantai tidak ditemukan"': '"Floor data not found"',
    r'"Gagal membuat lantai"': '"Failed to create floor"',
    r'"Reservasi kamar Anda di CM Living telah berhasil dibatalkan sesuai permohonan."': '"Your room reservation at CM Living has been successfully cancelled as requested."',
    r"'Penawaran berhasil diperbarui!'": "'Promotion updated successfully!'",
    r"'Penawaran berhasil ditambahkan!'": "'Promotion added successfully!'",
    r"'Penawaran berhasil dihapus!'": "'Promotion deleted successfully!'",
    r"'Dokumen PDF data performa & transaksi berhasil diunduh!'": "'Performance & transaction PDF document downloaded successfully!'",
    r"'File CSV data performa & transaksi berhasil diunduh!'": "'Performance & transaction CSV file downloaded successfully!'",
    r"'Unit kamar berhasil diperbarui!'": "'Room unit updated successfully!'",
    r"'Unit kamar berhasil ditambahkan!'": "'Room unit added successfully!'",
    r"'Unit kamar berhasil dihapus!'": "'Room unit deleted successfully!'",
    r'"Data berhasil diambil"': '"Data retrieved successfully"',
    r'"Data tidak ditemukan"': '"Data not found"',
    r'"Gagal"': '"Failed"',
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

