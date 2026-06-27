<script setup lang="ts">
import { ref } from 'vue'
import { toast } from 'vue3-toastify'

// ユーザーデータの型定義
interface UserProfile {
  id: string
  name: string
  faculty: string
  origin: string
  like_category: string
  like_thing: string
  dislike_category: string
  dislike_thing: string
  tool: string
  hobbies: string[]
  status: string
  bio: string
}

// 初期データ（ここにバックエンドから取得した値をそのまま入れても動きます）
const editForm = ref<UserProfile>({
  id:'n3',
  name: 'εИ',
  faculty: '情報理工学院 情報工学系 B2',
  origin: '高知県',
  like_category: '食べ物',
  like_thing: 'ラーメン',
  dislike_category: '言語',
  dislike_thing: 'TEX',
  tool: 'Python',
  hobbies: ['勉強', 'くねくね', '料理'],
  status: 'オートマトンおじさん',
  bio: 'Pythonはいいぞ！\n最近サウンドを始めました'
})

// 新しい趣味タグの入力用
const newHobbyInput = ref('')

// 変更を保存する（APIを叩く場合はここで行います）
const saveProfile = () => {
  console.log('保存されるデータ:', editForm.value)
  toast.success('プロフィールを保存しました！ ✨')
}

// 趣味タグの追加
const addHobbyTag = () => {
  const tag = newHobbyInput.value.trim()
  if (!tag) return
  if (editForm.value.hobbies.includes(tag)) {
    toast.warning('既に存在するタグです')
    return
  }
  editForm.value.hobbies.push(tag)
  newHobbyInput.value = ''
}

// 趣味タグの削除
const removeHobbyTag = (index: number) => {
  editForm.value.hobbies.splice(index, 1)
}
</script>

<template>
  <div class="myself-screen">
    <div class="profile-card-frame">
      
      <div class="column-left">
        <div class="name-section">
          <label class="label">名前：</label>
          <div class="name-id-wrap">
            <span class="user-name-text">{{ editForm.name }}</span>
            <span class="user-id-text"> (@{{ editForm.id }})</span>
          </div>
        </div>

        <div class="avatar-box">
          <img :src="`https://q.trap.jp/api/v3/public/icon/${editForm.id}`" alt="avatar" class="avatar-img" />
        </div>

        <hr class="divider-line" />

        <div class="info-list">
          <div class="info-row">
            <span class="label">学部/系：</span>
            <input v-model="editForm.faculty" type="text center-text" class="edit-input center-text" />
          </div>

          <div class="info-row">
            <span class="label">出身：</span>
            <input v-model="editForm.origin" type="text center-text" class="edit-input center-text" />
          </div>

          <div class="info-row like-group">
            <span class="label">好きな〇〇：</span>
            <div class="like-inputs">
              <input v-model="editForm.like_category" type="text" class="edit-input center-text" placeholder="例)食べ物" />
              <input v-model="editForm.like_thing" type="text" class="edit-input center-text" placeholder="例)ラーメン" />
            </div>
          </div>

          <div class="info-row like-group">
            <span class="label">嫌いな〇〇：</span>
            <div class="like-inputs">
              <input v-model="editForm.dislike_category" type="text" class="edit-input center-text" placeholder="例)アルゴリズム" />
              <input v-model="editForm.dislike_thing" type="text" class="edit-input center-text" placeholder="例)スターリンソート" />
            </div>
          </div>
        </div>
      </div>

      <div class="vertical-border"></div>

      <div class="column-right">
        <div class="info-list">
          
          <div class="info-row">
            <span class="label">好きな創作ツール：</span>
            <input 
              v-model="editForm.tool"
              type="text" class="edit-input center-text"
              placeholder="Python,アイビスペイント..." 
            />
          </div>

          <div class="info-row hobby-group">
            <span class="label">趣味タグ：</span>
            <div class="hobby-content">
              <div class="hobby-input-wrap">
                <input 
                  v-model="newHobbyInput" 
                  type="text" 
                  placeholder="タグを入力してEnter/ボタン" 
                  class="edit-input" 
                  @keydown.enter.prevent="addHobbyTag"
                />
                <button class="btn-add" @click="addHobbyTag">＋</button>
              </div>

              <div class="tags-container">
                <div v-for="(tag, idx) in editForm.hobbies" :key="idx" class="tag-item editable">
                  <span class="btn-remove-tag" @click="removeHobbyTag(idx)">×</span>
                  #{{ tag }}
                </div>
              </div>
            </div>
          </div>

          <div class="info-row bio-row">
            <span class="label">普段の様子：</span>
            <!-- <input v-model="editForm.status" type="text" class="edit-input" /> -->
            <textarea v-model="editForm.status" rows="5" class="edit-textarea"></textarea>
          </div>

          <div class="info-row bio-row">
            <span class="label">自由記述欄：</span>
            <textarea v-model="editForm.bio" rows="5" class="edit-textarea"></textarea>
          </div>

        </div>
      </div>

    </div>

    <div class="action-footer">
      <button class="btn-action btn-save" @click="saveProfile">保存する</button>
    </div>
  </div>
</template>

<style scoped>
.myself-screen {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 10px;
  box-sizing: border-box;
}

/* 外枠のカード (デザインの再現) */
.profile-card-frame {
  width: 100%;
  max-width: 1000px;
  border: 1px solid #707070;
  border-radius: 40px;
  background: #ffffff;
  display: flex;
  padding: 40px;
  box-sizing: border-box;
  min-height: 520px;
}

/* 左右カラム */
.column-left, .column-right {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.column-left {
  padding-right: 40px;
}
.column-right {
  padding-left: 40px;
  justify-content: center;
}

.vertical-border {
  width: 1px;
  background-color: #707070;
}

.name-section {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}
.name-input {
  max-width: 150px;
}

.avatar-box {
  width: 130px;
  height: 130px;
  border: 1px solid #707070;
  background-color: #e9ecef;
  margin: 0 auto;
}
.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.divider-line {
  width: 100%;
  border: none;
  border-top: 1px solid #707070;
  margin: 24px 0;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-row {
  display: flex;
  align-items: center;
  font-size: 1rem;
  color: #333;
}

.label {
  font-weight: 500;
  color: #495057;
  white-space: nowrap;
}

/* 好きな〇〇 */
.like-group {
  align-items: flex-start;
}
.like-inputs {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
  max-width: 200px;
}

/* 趣味タグ */
.hobby-group {
  align-items: flex-start;
}
.hobby-content {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}
.hobby-input-wrap {
  display: flex;
  gap: 8px;
}
.btn-add {
  border: 1px solid #707070;
  background: #fff;
  border-radius: 4px;
  padding: 4px 16px;
  cursor: pointer;
}
.btn-add:hover {
  background: #f1f3f5;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag-item {
  border: 1px solid #707070;
  border-radius: 12px;
  padding: 3px 12px;
  font-size: 0.9rem;
  background: #fff5f5;
  border-color: #ffc9c9;
  display: flex;
  align-items: center;
  gap: 6px;
}
.btn-remove-tag {
  cursor: pointer;
  color: #fa5252;
  font-weight: bold;
}

/* 自由記述欄 */
.bio-row {
  align-items: flex-start;
  flex-direction: column;
  gap: 6px;
}

/* 入力パーツ */
.edit-input {
  border: 1px solid #ced4da;
  border-radius: 4px;
  padding: 6px 10px;
  font-size: 0.95rem;
  width: 100%;
  box-sizing: border-box;
}
.edit-input.center-text {
  text-align: center;
  border-radius: 12px;
  border: 1px solid #707070;
}
.edit-textarea {
  width: 100%;
  border: 1px solid #ced4da;
  border-radius: 4px;
  padding: 10px;
  font-size: 0.95rem;
  resize: vertical;
  box-sizing: border-box;
}

/* 保存ボタン */
.action-footer {
  margin-top: 30px;
}
.btn-action {
  border: 1px solid #339af0;
  background: #4dabf7;
  color: white;
  padding: 8px 48px;
  border-radius: 4px;
  font-size: 1.05rem;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-action:hover {
  background: #339af0;
}
</style>