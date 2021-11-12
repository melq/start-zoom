StartZoom
===
登録したZoomをブラウザで開いてくれるWindows上で動作するツールです

<h1>使用方法</h1>

・`go build main.go` で実行ファイルを生成できます

・ダウンロードしたフォルダ内にある`myzoom.bat`をすでにpathの通っているディレクトリに置いてください
(これを行わないと、タスクスケジューラに会議を登録する機能が使用できません)

<h1>使用上の注意</h1>

・同名の会議は登録できません

・操作を間違えてしまった場合、戻るの表示がなければアプリを一度終了して抜けてください

・同時刻に開くZoomが2つ以上あるときは、一番最後に登録したものが開かれます

・現在進行中のZoomの終了時刻前でも、次のZoomの開始時刻(設定してある「開始時刻の余裕」も含める)を優先して起動します

・Windowsのタスクスケジューラを利用して、会議登録時に起動を予約しておくことができます
