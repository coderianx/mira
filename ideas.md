# Mira Ideas

## Temel

- `mira update <repo>` — binary'yi tekrar indir, üzerine yaz
- `mira info <repo>` — `mira.json` içeriğini göster

## Paket yönetimi
- Kurulum kaydı (`~/.local/share/mira/state.json`) — `list`, `upgrade-all`, `uninstall <isim>` (list ve uninstall ✅)
- `mira upgrade` — tüm kuruluları güncelle

## Dağıtım
- GitHub Releases desteği (raw yerine API ile son release'i bulma)
- `.tar.gz` / `.zip` desteği (şu an sadece raw binary)
- Branch/tag seçeneği (`--branch` veya `--tag`)
- GitLab / özel registry desteği

## Geliştirme
- `mira init` — bulunulan dizine `mira.json` iskeleti oluşturur
- `mira verify` — kurulu binary'lerin SHA256'sını tekrar kontrol et
- `mira publish` — release oluşturmayı otomatize et

## CLI iyileştirmeleri
- `--dir` flag'i ile kurulum dizini seçme (varsayılan `~/.local/bin`)
- Renkli çıktı
