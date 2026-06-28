<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

// 1. ダミーのユーザーデータ（バックエンドと接続するまでの繋ぎ）
interface UserProfile {
  username: string
  displayname: string
  affiliations:string[]
  major: string
  hometown: string
  favoriteTopic:{
    topic:string,
    value:string
  }
  dislikedTopic:{
    topic:string,
    value:string
  }
  technologies: string[] 
  tags: string[]         
  status: string
  bio: string
}

const dummyUsers: UserProfile[] = [
  {
    username: 'n3',
    displayname: 'εИ',
    affiliations: ['algo','game', 'sysad'],
    major: '情報理工学院 情報工学系 B2',
    hometown: '高知県',
    favoriteTopic:{
      topic:'もの',
      value:'genMira'
    },
    dislikedTopic:{
      topic:'もの',
      value:'Unity'
    },
    technologies: ['Python'],
    tags: ['勉学','くねくね','料理'],
    status: 'オートマトン大好き',
    bio: 'Pythonはいいぞ！\n最近サウンドを始めました'
  },
  // {
  //   id: "Suima",
  //   name: '睡魔',
  //   department: 'all',
  //   faculty: '生命理工学院 B2',
  //   origin: '東京都',
  //   like_category: '飲み物',
  //   like_thing: 'Monster',
  //   dislike_category: '言葉',
  //   dislike_thing: 'およー',
  //   tool: 'Tex',
  //   hobby: 'Tex,',
  //   status: 'TeXおじさん',
  //   bio: 'TeXをやりましょう'
  // }
]

const currentUserIndex = ref(0)
const currentUser = ref<UserProfile | null>(null)
const users = ref<UserProfile[]>([])
// 2. ジェスチャー・操作の管理用変数
let startX = 0
let isDragging = false
const swipeOffset = ref(0) // 視覚的なアニメーション用

const notify = (name: string|undefined, action: string) => {
  toast(`${name} さんに 【${action}】 をしました！`, {
    autoClose: 1000,
    "position": "bottom-left",
  });
}

// アクション処理（バックエンドにデータを送る場合はここで行う）
const handleAction = async(action: 'Like' | 'Nope') => {
  //toast.success(`${currentUser.value?.name} さんに 【${action}】 をしました！`)
  if (!currentUser.value) return;
  await sendAction(action, currentUser.value.username);
  notify(currentUser.value?.displayname,action);

  
  // 次のユーザーへ（データがなくなったらnull）
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

  const threshold = 100 // スワイプを確定させるしきい値（ピクセル）
  if (swipeOffset.value > threshold) {
    handleAction('Like') // 右に大きくスワイプ
  } else if (swipeOffset.value < -threshold) {
    handleAction('Nope') // 左に大きくスワイプ
  } else {
    swipeOffset.value = 0 // しきい値を超えなければ中央に戻す
  }
}

// 4. PCのキーボード（矢印キー）イベントハンドラ
const handleKeyDown = (e: KeyboardEvent) => {
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
    
    console.log("[getReccomend]ユーザー取得成功")
    // 最初のユーザーをセット
    if (users.value.length > 0) {
      currentUser.value = users.value[0]??null;
      currentUserIndex.value = 0;
      console.log("[getReccomend]user info",users.value[0])
    }
  } catch (error) {
    console.error("ユーザー詳細取得エラー:", error);
  }
};

const sendAction = async (action: 'Like' | 'Nope', username: string) => {
  const endpoint = action === 'Like' ? '/api/me/likes' : '/api/me/nopes';
  
  try {
    console.log({username})
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify({ username: username })
    });

    if (response.status === 204) {
      console.log(`${action} 成功: ${username}`);
      // 成功時の処理（必要に応じてトースト通知など）
    } else if (response.status === 409) {
      toast.warn("既にアクション済みです");
    } else {
      throw new Error(`アクション失敗: ${response.status}`);
    }
  } catch (error) {
    console.error("通信エラー:", error);
    toast.error("アクションの送信に失敗しました");
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
  console.log("Matching Start...")
  // getReccomend()
  users.value=dummyUsers
  currentUser.value = dummyUsers[0] ?? null
  }
)
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
  <div class="matching-screen">
    <div 
      v-if="currentUser" 
      class="flexible-stage"
      @mousedown="touchStart"
      @mousemove="touchMove"
      @mouseup="touchEnd"
      @mouseleave="touchEnd"
      @touchstart="touchStart"
      @touchmove="touchMove"
      @touchend="touchEnd"
      :class="{ 'is-dragging': isDragging }"
      :style="{ 
        transform: `translateX(${swipeOffset}px)`, 
        transition: isDragging ? 'none' : 'transform 0.3s ease' 
      }"
    >  
      <div class="absolute-item pos-department">
        <span class="label">所属:</span> {{ currentUser.affiliations.join(', ') }}
      </div>

      <div class="absolute-item pos-origin">
        <span class="label">出身:</span> {{ currentUser.hometown }}
      </div>
      
      <div class="absolute-item pos-faculty">
        <span class="label">学部/系:</span> {{ currentUser.major }}
      </div>
      
      <div class="absolute-item pos-like">
        <span class="label">好きな{{currentUser.favoriteTopic?.topic}}:</span> {{ currentUser.favoriteTopic?.value }}
      </div>
      
      <div class="absolute-item pos-dislike">
        <span class="label">嫌いな{{currentUser.dislikedTopic?.topic}}:</span> {{ currentUser.dislikedTopic?.value }}
      </div>

      <div 
        class="card-center"
      >
        <div class="avatar-box">
          <img :src="`https://q.trap.jp/api/v3/public/icon/${currentUser.username}`" alt="avatar" class="avatar-img" draggable="false" />
        </div>
        <div class="user-name">{{ currentUser.displayname }} (@{{currentUser.username}})</div>
      </div>

      <div class="absolute-item pos-tool">
        <span class="label">好きな創作ツール:</span> {{ currentUser.technologies?.join(', ') }}
      </div>
      
      <div class="absolute-item pos-hobby">
        <span class="label">趣味:</span> {{ currentUser.tags?.join(', ') }}
      </div>
      
      <div class="absolute-item pos-status">
        <span class="label">普段の様子:</span>
        <p class="bio-text">{{ currentUser.status }}</p>

      </div>
      
      <div class="absolute-item pos-bio">
        <span class="label">自由記述欄:</span>
        <p class="bio-text">{{ currentUser.bio }}</p>
      </div>

    </div>

    <div v-else class="no-more-users">
      <h2>今日の条件に合う人は全員チェックしました！</h2>
      <p>新しい友達が増えるのをお楽しみに！</p>
    </div>
  </div>
</template>

<style scoped>
.matching-screen {
  width: 100%;
  min-height: calc(100vh - 70px);
  display: flex;
  justify-content: center;
  align-items: center;
  user-select: none; 
  overflow: hidden;
}

.matching-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  max-width: 1100px;
  padding: 20px 0;
  gap: 40px;
}

.flexible-stage {
  position: relative; 
  width: 100%;
  height: 600px; 
  background: transparent;
}

/* 全ての絶対配置要素の共通ルール */
.absolute-item {
  position: absolute;
  font-size: 1.1rem;
  color: #333;
}

.label {
  font-weight: bold;
  color: #666;
}

/* ==========================================
   各項目の位置をパーセントで指定
   ========================================== */

/* 例：所属は上から20%、左から5%の位置 */
.pos-department {
  top:  0%;
  left: 14%;
}

/* 学部/系は上から35%、左から5% */
.pos-faculty {
  top: 40%;
  left: 0%;
}

.pos-origin {
  top: 18%;
  left: 5%;
}

.pos-like {
  top: 62%;
  left: 7%;
}

.pos-dislike {
  top: 80%;
  left: 14%;
}

/* --- 右側エリア --- */
.pos-tool {
  top: 0%;
  left: 75%;
}

.pos-hobby {
  top: 22%;
  left: 80%;
}

.pos-status {
  top: 45%;
  left: 70%;
  width: 30%;
}

.pos-bio {
  top: 70%;
  left: 60%;
  width: 30%;
}

.bio-text {
  white-space: pre-wrap;
  word-break: break-all;
  background: #f1f3f5;
  padding: 12px;
  border-radius: 8px;
  font-size: 0.95rem;
  margin-top: 6px;
  /*長すぎる自己紹介をマッチング画面で省略*/
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 5;
  line-clamp: 5;
  height:100px;
  overflow: auto;
}

/* --- 中央のアバター（画面のど真ん中に固定） --- */
.card-center {
  position: absolute;
  top: 40%;
  left: 50%;
  /* translate(-50%, -50%) で要素の「中心」をど真ん中に合わせています */
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: grab;
  z-index: 10;
}

.avatar-box {
  width: 220px;
  height: 220px;
  background-color: #eee;
  border: 1px solid #ccc;
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
  border-radius: 4px;
}
.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.user-name {
  margin-top: 25px;
  font-size: 1.4rem;
  font-weight: bold;
}

.no-more-users {
  text-align: center;
  color: #ff4a7d;
}

/* 終了画面 */
.no-more-users {
  text-align: center;
  color: #666;
}
.no-more-users h2 {
  color: #ff4a7d;
  margin-bottom: 12px;
}
</style>