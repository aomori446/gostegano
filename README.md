# gostegano 

画像にテキストメッセージを埋め込んだり、埋め込まれたメッセージを復号するための軽量なGo製ツールです。

## 機能

- ✅ 画像にメッセージを埋め込む（ステガノグラフィ）
- ✅ 埋め込まれたメッセージを抽出する
- ✅ ファイルまたはURLから画像を読み込む
- ✅ 使用後にソースファイルを自動削除（オプション）

## インストール

```go
go install github.com/aomori446/gostegano/example@latest
```
### or

```bash
git clone https://github.com/aomori446/gostegano.git
cd gostegano/example
go build -o example
```

## 使用方法

### メッセージを画像に埋め込む

```bash
./example -encode -source input.png -message "ひみつのメッセージ" -target output.png
```
- encode：エンコードモードを有効化
- source：元の画像（ローカルファイルまたはURL）
- message：埋め込むメッセージ
- target：出力ファイル名（.png 形式必須）

### 画像からメッセージを抽出する

```bash
./example -decode -source output.png
```

- decode：デコードモードを有効化
- source：メッセージが埋め込まれたPNG画像

### オプション：使用後に元ファイルを削除
```bash
./example -encode -source input.png -message "ひみつのメッセージ" -target output.png -rf
./example -decode -source output.png -rf
```

## ライセンス
[MIT LICENSE](https://github.com/aomori446/gostegano/blob/main/LICENSE)
