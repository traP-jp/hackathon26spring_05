<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

// 1. ダミーのユーザーデータ（バックエンドと接続するまでの繋ぎ）
interface UserProfile {
  username: string       // DBの PRIMARY KEY
  name: string           // name
  major: string          // 学部/系
  hometown: string       // 出身
  like_topic: string     // 好きな〇〇（カテゴリ名）
  like_value: string     // 好きなものの具体的な値
  dislike_topic: string  // 嫌いな〇〇（カテゴリ名）
  dislike_value: string  // 嫌いなものの具体的な値
  tool: string           // 好きな創作ツール
  usual_situation: string // 普段の様子
  bio: string            // 自由記述欄
  tags?: string[]        // DBのtagsテーブルから取得する想定の趣味タグ
}

const dummyUsers: UserProfile[] = [
  {
    username: 'n3',
    name: 'εИ',          // name 
    major: '情報理工学院 情報工学系 B2',
    hometown: '高知県',
    like_topic: '食べ物',
    like_value: 'ラーメン',
    dislike_topic: '言語',
    dislike_value: 'TEX',
    tool: 'Python',
    usual_situation: 'オートマトンおじさん',
    bio: 'Pythonはいいぞ！\n最近サウンドを始めました',
    tags: ['勉学', 'くねくね', '料理']
  },
  {
    username: 'Suima',
    name: '睡麻',        // name が復活しました！
    major: '生命理工学院 B2',
    hometown: '東京都',
    like_topic: '飲み物',
    like_value: 'Monster',
    dislike_topic: '言葉',
    dislike_value: 'およー',
    tool: 'Tex',
    usual_situation: 'TeXおじさん',
    bio: 'TeXをやりましょう',
    tags: ['Tex']
  }
]

const currentUserIndex = ref(0)
const currentUser = ref<UserProfile | null | undefined>(dummyUsers[0])
// 2. ジェスチャー・操作の管理用変数
let startX = 0
let isDragging = false
const swipeOffset = ref(0) // 視覚的なアニメーション用

const notify = (name: string|undefined, action: string) => {
  toast(`${name} さんに 【${action}】 をしました！`, {
    autoClose: 1000,
    "position": "bottom-center",
  });
}

// アクション処理（バックエンドにデータを送る場合はここで行う）
const handleAction = (action: 'Like' | 'Nope') => {
  //toast.success(`${currentUser.value?.name} さんに 【${action}】 をしました！`)
  notify(currentUser.value?.name, action); // 元通り name で通知されるように戻しました

  
  // 次のユーザーへ（データがなくなったらnull）
  currentUserIndex.value++
  if (currentUserIndex.value < dummyUsers.length) {
    currentUser.value = dummyUsers[currentUserIndex.value]|| null
  } else {
    currentUser.value = null
  }
  swipeOffset.value = 0
}

// 3. マウス・スマホのドラッグ/スワイプイベントハンドラ
// ユーザーのご指定（右スワイプ/右矢印 = Like、左スワイプ/左矢印 = Nope）で判定します
// 1. 各イベントの型を明示的に指定（Vue 3 / TypeScript環境）
const touchStart = (e: any) => {
  isDragging = true
  startX = e.touches ? e.touches[0].clientX : e.clientX
}

const touchMove = (e: any) => {
  if (!isDragging) return
  const currentX = e.touches ? e.touches[0].clientX : e.clientX
  swipeOffset.value = currentX - startX
}

const touchEnd = () => {
  if (!isDragging) return
  isDragging = false

  const threshold = 80 
  if (swipeOffset.value > threshold) {
    handleAction('Like') // 右に大きくスワイプ
  } else if (swipeOffset.value < -threshold) {
    handleAction('Nope') // 左に大きくスワイプ
  } else {
    swipeOffset.value = 0 // しきい値を超えなければ中央に戻す
  }
}
</script>

<template>
  <div 
    class="matching-mobile-screen"
    @mousedown="touchStart"
    @mousemove="touchMove"
    @mouseup="touchEnd"
    @mouseleave="touchEnd"
    @touchstart="touchStart"
    @touchmove="touchMove"
    @touchend="touchEnd"
  >
    <div 
      v-if="currentUser" 
      class="mobile-card-container"
      :class="{ 'is-dragging': isDragging }"
      :style="{ 
        transform: `translateX(${swipeOffset}px)`, 
        transition: isDragging ? 'none' : 'transform 0.3s ease' 
      }"
    >  
      <div class="card-scroll-body">
        
        <div class="profile-main">
          <div class="avatar-box">
            <img :src="`https://q.trap.jp/api/v3/public/icon/${currentUser.username}`" alt="avatar" class="avatar-img" draggable="false" />
          </div>
          <h2 class="user-name">{{ currentUser.name }}</h2>
          <span class="user-id">@{{ currentUser.username }}</span>
        </div>

        <div class="profile-details">
          
          <div class="info-item">
            <span class="label">学部/系:</span> 
            <span class="value">{{ currentUser.major }}</span>
          </div>

          <div class="info-item">
            <span class="label">出身:</span> 
            <span class="value">{{ currentUser.hometown }}</span>
          </div>

          <div class="info-item">
            <span class="label">好きな創作ツール:</span> 
            <span class="value font-badge tool-badge">{{ currentUser.tool }}</span>
          </div>

          <div v-if="currentUser.tags && currentUser.tags.length > 0" class="info-item">
            <span class="label">趣味タグ:</span> 
            <span class="value tag-text">{{ currentUser.tags.join('、') }}</span>
          </div>

          <div class="info-item">
            <span class="label">好きな{{ currentUser.like_topic }}:</span> 
            <span class="value">{{ currentUser.like_value }}</span>
          </div>

          <div class="info-item">
            <span class="label">嫌いな{{ currentUser.dislike_topic }}:</span> 
            <span class="value">{{ currentUser.dislike_value }}</span>
          </div>

          <div class="info-item">
            <span class="label">普段の様子:</span> 
            <span class="value italic-text">“ {{ currentUser.usual_situation }} ”</span>
          </div>

          <div class="info-item block-item">
            <span class="label">自由記述欄:</span>
            <p class="bio-text">{{ currentUser.bio }}</p>
          </div>
        </div>

      </div>
    </div>

    <div v-else class="no-more-users">
      <h2>今日の条件に合う人は全員チェックしました！🏹</h2>
      <p>新しい友達が増えるのをお楽しみに！</p>
    </div>
  </div>
</template>

<style scoped>
.matching-mobile-screen {
  width: 100%;
  min-height: calc(100vh - 60px);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 15px;
  background-color: #f8f9fa;
  box-sizing: border-box;
  user-select: none;
  overflow: hidden;
  /* 画面全体のどこをドラッグしてもカーソルが「掴む」マークになるように変更 */
  cursor: grab;
}

.matching-mobile-screen:active {
  cursor: grabbing;
}

.mobile-card-container {
  width: 100%;
  max-width: 360px;
  height: 600px; 
  background: #ffffff;
  border: 1px solid #c8c8c8;
  border-radius: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  /* カード自体のクリックイベントがスクロールの邪魔をしないよう設定 */
  pointer-events: auto;
}

.card-scroll-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  scrollbar-width: thin;
}

.profile-main {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 25px;
}

.avatar-box {
  width: 130px;
  height: 130px;
  border: 1px solid #999;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  background-color: #f0f0f0;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-name {
  font-size: 1.5rem;
  font-weight: bold;
  margin: 12px 0 2px 0;
  color: #222;
}

.user-id {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 5px;
}

.profile-details {
  display: flex;
  flex-direction: column;
  gap: 16px;
  text-align: left;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  border-bottom: 1px dashed #e0e0e0;
  padding-bottom: 8px;
  font-size: 0.95rem;
}

.info-item.block-item {
  flex-direction: column;
  border-bottom: none;
}

.label {
  font-weight: bold;
  color: #666;
  font-size: 0.85rem;
  white-space: nowrap;
  margin-right: 15px;
  padding-top: 2px;
}

.value {
  color: #333;
  text-align: right;
  word-break: break-word;
}

.font-badge {
  background-color: #e3f2fd;
  color: #1e88e5;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
}

.tool-badge {
  background-color: #f1f3f5;
  color: #495057;
  border: 1px solid #dee2e6;
}

.italic-text {
  font-style: italic;
  color: #555;
}

.bio-text {
  width: 100%;
  white-space: pre-wrap;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  padding: 10px;
  border-radius: 8px;
  font-size: 0.9rem;
  margin-top: 6px;
  color: #444;
  line-height: 1.4;
  box-sizing: border-box;
}

.no-more-users {
  text-align: center;
  color: #666;
  padding: 20px;
}
.no-more-users h2 {
  color: #ff4a7d;
  margin-bottom: 12px;
  font-size: 1.5rem;
}
</style>