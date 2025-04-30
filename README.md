# gostegano 

画像にテキストメッセージを埋め込んだり、埋め込まれたメッセージを復号するための軽量なGo製ツールです。

## 機能

- ✅ 画像にメッセージを埋め込む（ステガノグラフィ）
- ✅ 埋め込まれたメッセージを抽出する
- ✅ ファイルまたはURLから画像を読み込む
- ✅ 使用後にソースファイルを自動削除（オプション）

## インストール

```go
go install github.com/aomori446/gostegano/cmd/gostegano@latest
```
### OR

```bash
git clone https://github.com/aomori446/gostegano.git
cd gostegano/cmd/gostegano
go build -o gostegano
```

## 使用方法

### メッセージを画像に埋め込む

```bash
./gostegano -en -from input.png -msg "ひみつのメッセージ" -to output.png
```
- en：エンコードモードを有効化
- from：元の画像（ローカルファイルまたはURL）
- msg：埋め込むメッセージ
- to：出力ファイル名

### 画像からメッセージを抽出する

```bash
./gostegano -de -from output.png
```

- de：デコードモードを有効化
- from：メッセージが埋め込まれたPNG画像

### オプション：使用後に元ファイルを削除
```bash
./cmd -en -from input.png -msg "ひみつのメッセージ" -to output.png -rm
./cmd -de -from output.png -rm
```

## ライセンス
[MIT LICENSE](https://github.com/aomori446/gostegano/blob/main/LICENSE)
