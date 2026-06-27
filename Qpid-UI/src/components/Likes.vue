<script setup lang="ts">
import { ref, computed } from 'vue'

//注意：色とかはまだ何も決めていません。全部仮です

// タブの切り替え状態を管理 ('liked' = Likeした人, 'likedBy' = Likeされた人)
const activeTab = ref<'liked' | 'likedBy'>('liked')

// Likeした人一覧のダミーデータ
const dummyLikedUsers = [
  { 
    id: 1, 
    name: 'n3', 
    bio: '明日晴れますよ(断定）\nうにお願いします。（注文）' 
  },
  { 
    id: 2, 
    name: 'Jiro', 
    bio: 'バックエンドエンジニアです。\nGo言語が好きです。' 
  },
  { 
    id: 3, 
    name: 'Saburo', 
    bio: 'デザイナー兼フロントエンド。\nUI/UXにこだわりがあります!' 
  },
  { 
    id: 4, 
    name: 'Shiro', 
    bio: 'AIに興味があります。\nよろしくお願いします。' 
  },
  { 
    id: 5, 
    name: 'Goro', 
    bio: 'これは自己紹介文のサンプルです、パイソンはいいぞよりも長くこのままだとタブ一覧からはみ出すかもしれないので、先ほどのCSSでしっかり「...」になるかテストするための長い文章です。' 
  },
  { 
    id: 6, 
    name: 'Rokuro', 
    bio: 'プログラミング初心者です！\n楽しく開発したいです。' 
  }
]

// Likeされた人一覧のダミーデータ（こちらも4人に設定）
const dummyLikedByUsers = [
  { 
    id: 101, 
    name: 'Hanako', 
    bio: 'React派ですがVueも触ってみてます!\n仲良くしてください。' 
  },
  { 
    id: 102, 
    name: 'Keiko', 
    bio: '趣味はカフェ巡りです。\n休日はもくもく会によく行きます。' 
  },
  { 
    id: 103, 
    name: 'Mari', 
    bio: 'TypeScript最高!\n型がないと不安になります。' 
  },
  { 
    id: 104, 
    name: 'Yumi', 
    bio: 'インフラエンジニア。\nAWSメインで触ってます。' 
  }
]

// 現在選択されているタブに応じて、表示する配列を切り替える
const displayUsers = computed(() => {
  return activeTab.value === 'liked' ? dummyLikedUsers : dummyLikedByUsers
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

    <div class="cards-grid">
      <div v-for="user in displayUsers" :key="user.id" class="card">
        <div class="icon-placeholder"></div>
        <p class="user-name">{{ user.name }}</p>
        <p class="user-bio">{{ user.bio }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ページ全体を画面の横幅いっぱいまで広げる */
.likes-page {
  padding: 40px 5%; /* 左右に少し余白を取りつつ、画面幅いっぱいに広げる */
  max-width: 100%; /* ★ここを100%にすることで右端まで伸びます */
  margin: 0 auto;
}

/* タブのスタイル（変更なし） */
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
  grid-template-columns: repeat(4, 1fr); /* 画面幅をきっちり4等分して広げる */
  gap: 40px; /* カード同士の間隔 */
}

/* カードの比率を維持したまま拡大 */
.card {
  background-color: #f0f0f0;
  border: 1px solid #aaa;
  border-radius: 20px;
  padding: 10%; /* ★固定のpxではなく%にすることで、カード拡大に合わせて内側の余白も広がる */
  display: flex;
  flex-direction: column;
  aspect-ratio: 2 / 3; /* ★ここで縦横の比率を固定（横2：縦3） */
}

/* アイコンの比率も維持 */
.icon-placeholder {
  width: 50%; /* ★カード幅の50%の大きさに設定 */
  aspect-ratio: 1 / 1; /* ★常に正方形を維持する */
  border: 1px solid #555;
  background-color: #fff;
  margin-bottom: 20px;
}

.user-name {
  font-size: 1.8rem; /* カードが大きくなるのに合わせて文字も少し大きく */
  margin: 0 0 15px 0;
  color: #333;
}

.user-bio {
  font-size: 1.3rem; /* こちらも少し大きく */
  white-space: pre-wrap;
  margin: 0;
  color: #555;
  line-height: 1.6;
  
  /* 紹介文が省略のための設定 */
  display: -webkit-box;
  -webkit-line-clamp: 4; /* カードが縦に長くなった分、4行まで表示させる */
  line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>