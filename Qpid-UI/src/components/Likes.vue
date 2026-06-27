<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// タブの切り替え状態を管理 ('liked' = Likeした人, 'likedBy' = Likeされた人)
const activeTab = ref<'liked' | 'likedBy'>('liked')

// バックエンドの UserSummary 型の定義
interface UserSummary {

  username: string    // DBの PRIMARY KEY
  name: string        // サークルの人の名前（消さずに残しました！）

  bio?: string
}

// 本物のデータを入れるための空っぽの箱
const likedUsers = ref<UserSummary[]>([])
const likedByUsers = ref<UserSummary[]>([])

// 画面が開いた瞬間に、ブラウザ標準の機能（fetch）でバックエンドからデータを取ってくる
onMounted(async () => {
  try {
    // 1. LIKEした人をバックエンドから取得
    const resLikes = await fetch('http://localhost:8080/api/me/likes')
    if (resLikes.ok) {
      likedUsers.value = await resLikes.json()
    }

    // 2. 自分をLIKEした人をバックエンドから取得
    const resLikedBy = await fetch('http://localhost:8080/api/me/liked-by')
    if (resLikedBy.ok) {
      likedByUsers.value = await resLikedBy.json()
    }

  } catch (error) {
    console.error('APIの取得に失敗しました。サーバーが起動していないか、ログインしていない可能性があります:', error)
  }
})

// 現在選択されているタブに応じて、表示する配列を切り替える
const displayUsers = computed(() => {
  return activeTab.value === 'liked' ? likedUsers.value : likedByUsers.value
})
</script>

<template>
  <div class="likes-page">
    <div class="tab-container">
      <button
        class="tab-button"
        :class="{ active: activeTab === 'liked' }"
        @click="activeTab = 'liked'"
      >
        Likeした人一覧
      </button>
      <button
        class="tab-button"
        :class="{ active: activeTab === 'likedBy' }"
        @click="activeTab = 'likedBy'"
      >
        Likeされた人一覧
      </button>
    </div>

    <div v-if="displayUsers.length === 0" class="no-users">
      まだユーザーがいません。
    </div>

    <div v-else class="cards-grid">
      <div v-for="user in displayUsers" :key="user.username" class="card">
        <div class="icon-placeholder"></div>
        <p class="user-name">{{ user.name || user.username }}</p>
        <p class="user-bio">{{ user.bio || '自己紹介文はまだありません。' }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ページ全体を画面の横幅いっぱいまで広げる */
.likes-page {
  padding: 40px 5%;
  max-width: 100%;
  margin: 0 auto;
}

/* データがないときのメッセージ */
.no-users {
  font-size: 1.5rem;
  color: #777;
  text-align: center;
  margin-top: 50px;
}

/* タブのスタイル */
.tab-container {
  display: flex;
  margin-bottom: 40px;
  border: 1px solid #777;
  border-radius: 6px;
  width: fit-content;
  overflow: hidden;
}

.tab-button {
  padding: 15px 40px;
  font-size: 20px;
  color: #555;
  border: none;
  background-color: #f8f9fa;
  cursor: pointer;
  outline: none;
}

.tab-button.active {
  background-color: #e0e0e0;
}

.tab-button:first-child {
  border-right: 1px solid #777;
}

/* 4列で画面いっぱいに均等に広がるようにする */
.cards-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 40px;
}

/* カードの比率を維持したまま拡大 */
.card {
  background-color: #f0f0f0;
  border: 1px solid #aaa;
  border-radius: 20px;
  padding: 10%;
  display: flex;
  flex-direction: column;
  aspect-ratio: 2 / 3;
}

/* アイコンの比率も維持 */
.icon-placeholder {
  width: 50%;
  aspect-ratio: 1 / 1;
  border: 1px solid #555;
  background-color: #fff;
  margin-bottom: 20px;
}

.user-name {
  font-size: 1.8rem;
  margin: 0 0 15px 0;
  color: #333;
}

.user-bio {
  font-size: 1.3rem;
  white-space: pre-wrap;
  margin: 0;
  color: #555;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/*画面の大きさによる表示の仕方の調整*/
@media (max-width: 768px) {
  .likes-page {
    padding: 20px 4%; /* 左右の余白を少し狭くして画面を広く使う */
  }

  .no-users {
    font-size: 1.2rem;
    margin-top: 30px;
  }

  /* タブを横幅いっぱいに広げてスマホでタップしやすくする */
  .tab-container {
    width: 100%;
    margin-bottom: 25px;
  }

  .tab-button {
    flex: 1; /* ボタンを均等に50%ずつの幅にする */
    padding: 12px 10px;
    font-size: 16px; /* 文字サイズをスマホ用に調整 */
    text-align: center;
  }

  /* スマホの縦長画面に合わせて、4列から「2列」に変更 */
  .cards-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px; /* カード同士の間隔を少し詰める */
  }

  /* スマホだと縦長になりすぎるのを防ぐため、比率を 2/3 から 3/4 に微調整 */
  .card {
    border-radius: 12px;
    padding: 12px;
    aspect-ratio: 3 / 4;
  }

  .icon-placeholder {
    width: 45%;
    margin-bottom: 10px;
  }

  .user-name {
    font-size: 1.2rem; /* スマホで見やすい文字サイズに */
    margin-bottom: 8px;
  }

  .user-bio {
    font-size: 0.9rem; /* 自己紹介文を小さくして、はみ出しを防ぐ */
    line-height: 1.4;
    -webkit-line-clamp: 3; /* 行数を4行から3行に制限してスッキリさせる */
    line-clamp: 3;
  }
}

/* 画面がめちゃくちゃ小さいスマホ（iPhone SEなど）用の微調整 */
@media (max-width: 400px) {
  .cards-grid {
    grid-template-columns: 1fr; /* 1列にして1枚ずつドカンと見せる */
  }
  .card {
    aspect-ratio: auto; /* 1列のときは高さを自動にして文字が漏れないように */
  }
}
</style>