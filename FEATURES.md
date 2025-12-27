# 🎉 LEAP SSH Manager - Tier 1 & Tier 2 Özellikleri Eklendi!

## ✅ Eklenen Özellikler

### 📊 **Tier 1 Özellikleri (Temel)**

#### 1. ✏️ **Edit Komutu** (`leap edit`)
- Mevcut bağlantıları düzenleme
- Tüm alanları güncelleme (host, user, port, password, key, tags, jump host, notes)
- Mevcut değerleri default olarak gösterme
- İnteraktif prompt sistemi

**Kullanım:**
```bash
leap edit myserver
```

#### 2. 🗑️ **Delete Komutu** (`leap delete`)
- Tek veya çoklu bağlantı silme
- Onay isteme (güvenlik)
- `--force` flag ile onaysız silme
- Silinen/bulunamayan bağlantıları raporlama

**Kullanım:**
```bash
leap delete oldserver
leap delete server1 server2 server3
leap delete myserver --force
```

**Aliases:** `rm`, `remove`

#### 3. 🧪 **Test/Ping Komutu** (`leap test`)
- TCP port kontrolü
- Latency ölçümü (ms)
- SSH auth testi
- Tek, tüm veya tag bazlı test

**Kullanım:**
```bash
leap test myserver
leap test --all
leap test --tag production
```

**Aliases:** `ping`, `check`

**Çıktı Örneği:**
```
⚡ Connection Health Check
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

myserver (user@host:22)
  ✓ Port 22: Open
  ⏱️  Latency: 45ms
  ✓ SSH Auth: Success
```

#### 4. ⭐ **Favoriler Sistemi** (`leap favorite`)
- Bağlantıları favorilere ekleme/çıkarma
- Favori listesi görüntüleme
- Toggle mekanizması

**Kullanım:**
```bash
leap favorite myserver      # Toggle
leap favorites              # List all
```

**Aliases:** `fav`, `star`, `favs`

---

### 📊 **Tier 2 Özellikleri (Orta)**

#### 5. 📤 **Export/Import** (`leap export`, `leap import`)
- JSON ve YAML formatı desteği
- Dosyaya veya stdout'a export
- Merge modu ile import
- Backup ve paylaşım için ideal

**Kullanım:**
```bash
# Export
leap export backup.json
leap export backup.yaml --format yaml
leap export  # stdout'a yazdır

# Import
leap import backup.json
leap import backup.yaml --merge  # Mevcut olanları güncelle
```

**Çıktı Örneği:**
```
⚡ Import Connections
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Added server1
✓ Added server2
⟳ Updated production
⊘ Skipped staging (already exists)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Added: 2  ⟳ Updated: 1  ⊘ Skipped: 1
```

#### 6. 🖥️ **Remote Exec** (`leap exec`)
- Uzak sunucularda komut çalıştırma
- Tek, tüm veya tag bazlı çalıştırma
- Gerçek zamanlı çıktı
- Jump host desteği

**Kullanım:**
```bash
leap exec myserver "uptime"
leap exec --all "df -h"
leap exec --tag web "systemctl status nginx"
```

#### 7. 📁 **File Transfer** (`leap upload`, `leap download`)
- SCP kullanarak dosya transferi
- Recursive (klasör) desteği
- Jump host desteği
- Progress gösterimi

**Kullanım:**
```bash
# Upload
leap upload myserver ./local.txt /remote/path/
leap upload myserver ./folder/ /remote/ --recursive

# Download
leap download myserver /remote/file.txt ./local/
leap download myserver /remote/folder/ ./ --recursive
```

#### 8. 📝 **Notlar** (`leap notes`)
- Bağlantılara not ekleme
- Not görüntüleme
- Not düzenleme

**Kullanım:**
```bash
leap notes myserver             # Görüntüle
leap notes myserver --edit      # Düzenle
```

---

## 🔧 **Yapısal İyileştirmeler**

### Config Struct Güncellemeleri
Connection struct'ına yeni alanlar eklendi:

```go
type Connection struct {
    // ... mevcut alanlar
    LastUsed     time.Time  // Son kullanım zamanı
    Favorite     bool       // Favori mi?
    Notes        string     // Kullanıcı notları
    UsageCount   int        // Kullanım sayısı
    CreatedAt    time.Time  // Oluşturulma zamanı
}
```

### Yeni Helper Fonksiyonlar
```go
cfg.UpdateLastUsed(name)        // LastUsed ve UsageCount güncelle
cfg.DeleteConnection(name)      // Bağlantı sil
cfg.ToggleFavorite(name)        // Favori toggle
cfg.SetNotes(name, notes)       // Not ekle/güncelle
```

---

## 📋 **Komut Listesi (Güncel)**

| Komut | Aliases | Açıklama |
|-------|---------|----------|
| `leap add` | - | Yeni bağlantı ekle |
| `leap list` | - | Bağlantıları listele |
| `leap connect` | - | Bağlan |
| `leap edit` | - | Bağlantı düzenle |
| `leap delete` | `rm`, `remove` | Bağlantı sil |
| `leap test` | `ping`, `check` | Bağlantı testi |
| `leap favorite` | `fav`, `star` | Favori toggle |
| `leap favorites` | `favs` | Favorileri listele |
| `leap notes` | - | Not görüntüle/düzenle |
| `leap exec` | - | Uzaktan komut çalıştır |
| `leap upload` | - | Dosya yükle |
| `leap download` | - | Dosya indir |
| `leap export` | - | Config export |
| `leap import` | - | Config import |
| `leap tunnel` | - | SSH tunnel |

---

## 🎨 **Tasarım Özellikleri**

Tüm yeni komutlar Laravel CLI tarzında modern tasarıma sahip:

- ⚡ Yeşil header'lar
- ✓/❌ Durum göstergeleri
- 🎨 Renkli çıktılar (ANSI)
- 📊 Unicode çizgiler
- 💡 Yardımcı ipuçları
- 🏷️ İkonlar ve emojiler

---

## 🚀 **Kullanım Örnekleri**

### Senaryo 1: Yeni Sunucu Ekleme ve Test
```bash
leap add production
leap test production
leap favorite production
leap notes production --edit
```

### Senaryo 2: Toplu İşlemler
```bash
leap test --all
leap exec --tag web "systemctl status nginx"
leap export backup-$(date +%Y%m%d).json
```

### Senaryo 3: Dosya Transferi
```bash
leap upload production ./deploy.sh /opt/scripts/
leap exec production "bash /opt/scripts/deploy.sh"
leap download production /var/log/app.log ./logs/
```

### Senaryo 4: Backup ve Paylaşım
```bash
# Backup
leap export ~/backups/leap-$(date +%Y%m%d).json

# Yeni makinede restore
leap import ~/backups/leap-20231227.json
```

---

## 📈 **İstatistikler**

- **Toplam Komut:** 15
- **Yeni Eklenen:** 8
- **Toplam Dosya:** 13 Go dosyası
- **Kod Satırı:** ~1500+ satır
- **Özellik Sayısı:** 15+

---

## ✅ **Test Edildi**

- ✅ Derleme başarılı
- ✅ Tüm komutlar help'te görünüyor
- ✅ Laravel tarzı modern çıktılar
- ✅ Hata yönetimi
- ✅ Flag desteği

---

## 🎯 **Sonraki Adımlar (Tier 3)**

İsterseniz şunları da ekleyebiliriz:

1. **Gruplar** - Bağlantıları gruplama
2. **İstatistikler** - Kullanım analizi
3. **Cloud Sync** - AWS/DO entegrasyonu
4. **Profiller** - Çoklu config dosyası
5. **Alias Sistemi** - Özel komutlar
6. **Tema Desteği** - Renk temaları

---

**Tüm Tier 1 ve Tier 2 özellikleri başarıyla eklendi! 🎉**
