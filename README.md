# hackathon26spring_05

## Members

- genMira
- Arina
- hiko9907
- watapro26B
- tidus
- eoa_te

## コマンド - フロント

ローカル環境で立ち上げる
```
npm run dev
```

立ち上がる先
```
http://localhost:8080/
```

その後中止する
ctrl+C

## Backend

タスクランナーの [Taskfile](https://taskfile.dev) があると便利です.

```sh
# Ubuntu, debian
apt install task
# MacOS
sudo snap install task --classic
```

### Commands

- 初回セットアップ

  ```sh
  task init
  ```

- ローカル環境で起動

  変更を保存すると自動で再起動されます.

  ```sh
  task up
  ```
  
- ローカル環境を停止

  これやらないと手元のパソコンで動きっぱなしになります.

  ```sh
  task down
  ```

- フォーマット

  PR 出す前に行ってください.

  ```sh
  task fmt
  ```

### Environment

バックエンドの環境変数は `.env.local` に設定してください. `task init` で自動生成されます.
