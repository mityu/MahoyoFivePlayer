# MahoyoFivePlayer

## 何コレ
これは、「魔法使いの夜」というゲームのBGMの一つである「Five」という曲を、ゲーム内でされているようにループ再生させるためのコマンドライン上で動作する音楽プレーヤーです。

## 必要なもの

「魔法使いの夜 ORIGINAL SOUNDTRACK」に収録されている `Five.mp3`

## ビルド
このリポジトリをクローンし、以下のように`make`コマンドを実行してください。ビルドに成功すると、このリポジトリのディレクトリ直下に`Five`という実行ファイルが作成されます。

```sh
$ git clone https://github.com/mityu/MahoyoFivePlayer
$ cd MahoyoFivePlayer
$ make FIVEAUDIO=/path/to/Five.mp3
```

（なお、このリポジトリ直下に `Five.mp3` を配置する場合は、`FIVEAUDIO` を設定する必要はありません。）

## 使えるコマンド
```
e[xit]|q[uit]    Quit this player
p[ause]          Pause music (valid only when playing)
r[esume]         Cancel pausing (valid only when pausing)
s[tatus]         Show status
h[elp]           Show this help
```
