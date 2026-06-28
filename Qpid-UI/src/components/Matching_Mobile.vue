<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

// 1. ダミーのユーザーデータ（バックエンドと接続するまでの繋ぎ）
interface UserProfile {
  username: string       // DBの PRIMARY KEY
  name: string           // 班の人が残してほしいと言っていたサークルの人の名前
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



const currentUserIndex = ref(0)
const currentUser = ref<UserProfile | null>(null)
const users = ref<UserProfile[]>([])
// 2. ジェスチャー・操作の管理用変数
let startX = 0
let isDragging = false
const swipeOffset = ref(0) // 視覚的なアニメーション用

// 不透明度の最大値
const MAX_OPACITY = 0.8;
// 不透明度が最大になるスワイプ量 (px)
const OPACITY_THRESHOLD = 150; 

// スワイプ量（swipeOffset）から、Like/Nopeそれぞれの奥のレイヤーの不透明度を計算
const swipeOpacity = computed(() => {
  const offset = swipeOffset.value;
  const ratio = Math.min(Math.abs(offset) / OPACITY_THRESHOLD, 1);
  const opacity = ratio * MAX_OPACITY;

  return {
    // 右にスワイプ（プラス）した時はLIKE（赤いハート側）を明るく
    like: offset > 0 ? opacity : 0,
    // 左にスワイプ（マイナス）した時はNOPE（青いハート側）を明るく
    nope: offset < 0 ? opacity : 0
  };
});

const notify = (name: string|undefined, action: string) => {
  toast(`${name} さんに 【${action}】 をしました！`, {
    autoClose: 1000,
    "position": "bottom-center",
  });
}

// アクション処理
const handleAction = (action: 'Like' | 'Nope') => {
  notify(currentUser.value?.name, action);
  
  currentUserIndex.value++
  if (currentUserIndex.value < users.value.length) {
    const nextUser = users.value[currentUserIndex.value];
    currentUser.value = nextUser !== undefined ? nextUser : null;
  }else {
    currentUser.value = null
  }
  swipeOffset.value = 0
}

// 3. マウス・スマホのドラッグ/スワイプイベントハンドラ
// 1. 各イベントの型を明示的に指定（Vue 3 / TypeScript環境）
const touchStart = (e: any) => {
  // 表示するユーザーがもういない場合は、スワイプ操作を受け付けない
  if (!currentUser.value) return

  isDragging = true
  startX = e.touches ? e.touches[0].clientX : e.clientX
}

const touchMove = (e: any) => {
  if (!isDragging) return
  const currentX = e.touches ? e.touches[0].clientX : e.clientX
  swipeOffset.value = currentX - startX
}

// ユーザーのご指定（右スワイプ/右矢印 = Like、左スワイプ/左矢印 = Nope）で判定します
const touchEnd = () => {
  if (!isDragging) return
  isDragging = false

  const threshold = 80 
  if (swipeOffset.value > threshold) {
    handleAction('Like') // 右スワイプ
  } else if (swipeOffset.value < -threshold) {
    handleAction('Nope') // 左スワイプ
  } else {
    swipeOffset.value = 0
  }
}

// 4. PCのキーボード（矢印キー）イベントハンドラ
const handleKeyDown = (e: KeyboardEvent) => {
  if (!currentUser.value) return // ユーザーがいない場合は何もしない
  if (e.key === 'ArrowRight') {
    handleAction('Like')
  } else if (e.key === 'ArrowLeft') {
    handleAction('Nope')
  }
}

const getReccomend = async() =>{
  try{
    //const response = await fetch(`https://qpid.trap.show/api/me`,{
    const response = await fetch(`/api/suggestions`,{
      method: "GET",
      headers:{
        "content-type":"application/json"
      },
    });

    if(!response.ok){
      console.log("Error : Not OK")
      const errorText = await response.text();
      console.log("バックエンドから返ってきた生の文字:", errorText);
    }
    // const errorText = await response.text();
    // console.log("バックエンドから返ってきた生の文字:", errorText);
    const suggestions = await response.json();    
    console.log("[getReccomend]APIから取得したデータ:", suggestions)
    await getReccomendUser(suggestions.map((s: any) => s.username));
    
  }catch(error){
    console.log("Error : ",error)
    toast.error("通信エラーが発生しました")
  }
}

const getReccomendUser = async (userIDs: Array<string>) => {
  try {
    const userPromises = userIDs.map(async (id) => {
      const res = await fetch(`/api/users/${id}`);
      if (!res.ok) return null;
      return res.json();
    });

    const results = await Promise.all(userPromises);
    
    // 取得できたユーザーのみを格納 (nullを除外)
    users.value = results.filter((u) => u !== null);
    
    // 最初のユーザーをセット
    if (users.value.length > 0) {
      currentUser.value = users.value[0]??null;
      currentUserIndex.value = 0;
    }
    console.log("[getReccomend]ユーザー取得成功")
  } catch (error) {
    console.error("ユーザー詳細取得エラー:", error);
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
  console.log("Matching Start...")
  getReccomend()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
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
      class="swipe-overlay-layer"
      :style="{ transition: isDragging ? 'none' : 'opacity 0.3s ease' }"
    >
      <div 
        class="overlay-nope"
        :style="{ opacity: swipeOpacity.nope }"
      ></div>

      <div 
        class="overlay-like"
        :style="{ opacity: swipeOpacity.like }"
      ></div>
    </div>

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
            <div class="tag-badges-container">
              <span v-for="tag in currentUser.tags" :key="tag" class="value font-badge tool-badge">{{ tag }}</span>
            </div>
          </div>

          <div class="info-item">
            <span class="label">好きな{{ currentUser.like_topic }}:</span> 
            <span class="value">{{ currentUser.like_value }}</span>
          </div>

          <div class="info-item">
            <span class="label">嫌いな{{ currentUser.dislike_topic }}:</span> 
            <span class="value">{{ currentUser.dislike_value }}</span>
          </div>

          <div class="info-item block-item">
            <span class="label">普段の様子:</span> 
            <p class="usual-text">{{ currentUser.usual_situation }}</p>
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
  height: calc(100vh - 60px);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0;
  background-color: #f8f9fa;
  box-sizing: border-box;
  user-select: none;
  overflow: hidden;
  cursor: grab;
  position: relative;
}

.matching-mobile-screen:active {
  cursor: grabbing;
}

/* ✨ 奥のレイヤー全体のスタイル ✨ */
.swipe-overlay-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none; /* ドラッグ操作を邪魔しない */
  display: flex;
}

/* 左右それぞれの画像表示エリアの共通設定 */
.overlay-nope,
.overlay-like {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0; 
  background-repeat: no-repeat;
  background-position: center;
  background-size: 80% auto; /* ハートを表示するのは後ろのレイヤー全体 */
}

/* NOPE（左側：青いハート） */
.overlay-nope {
  background-color: rgba(30, 136, 229, 0.08); /* ほんのり青背景 */
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%231E88E5'><path d='M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.5 3 21.9 5.42 21.9 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z'/></svg>");
}

/* LIKE（右側：赤いハート） */
.overlay-like {
  background-color: rgba(255, 74, 125, 0.08); /* ほんのり赤背景 */
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%23FF4A7D'><path d='M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.5 3 21.9 5.42 21.9 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z'/></svg>");
}

.mobile-card-container {
  width: 100%;
  height: 100%;
  background: #ffffff;
  border: none;
  border-radius: 0;
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  z-index: 1; /* 奥のレイヤーより手前に配置 */
  pointer-events: auto;
}

.card-scroll-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  scrollbar-width: thin;
  -webkit-overflow-scrolling: touch;
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

.tag-badges-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: flex-end;
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

.usual-text {
  width: 100%;
  white-space: pre-wrap;
  word-break: break-all;
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

.bio-text {
  width: 100%;
  white-space: pre-wrap;
  word-break: break-all;
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
  z-index: 2;
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