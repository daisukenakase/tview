# tview

画像をターミナル上に表示するコマンドラインツールです。画像を端末サイズに合わせて縮小し、既定では背景色のカラーセルで描画します。

## 対応環境

- Go 1.26以降
- Windows、Linux、macOS
- JPEG、PNG、WebP、GIF、BMP、TIFF（アニメーション GIF は先頭フレームを表示）
- カラー対応ターミナル（非対応の場合はグレースケールASCIIへ自動切り替え）

## ビルド

### Windows（PowerShell）

```powershell
go build -trimpath -ldflags="-s -w" -o tview.exe .
```

### Linux / macOS

```sh
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o tview .
```

## 使い方

```powershell
.	view.exe .\image.jpg
```

カラーの文字で表示する場合:

```powershell
.	view.exe -color-ascii .\image.jpg
```

表示幅の上限を指定する場合:

```powershell
.	view.exe -width 100 .\image.jpg
```

画像は端末の幅と高さに合わせて縮小されます。カラー非対応の端末や標準出力がリダイレクトされた場合は、グレースケールASCIIで表示します。標準出力が端末でない場合、幅は既定値の80になります。

## テスト

```sh
go test ./...
```