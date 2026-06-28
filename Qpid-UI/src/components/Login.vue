<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-title">
        <img class="logo-text" src="https://img.icons8.com/?size=100&id=5IeEIBOnAlpA&format=png&color=000000" alt="from Icons8" />
        <h1>Qpid</h1>
      </div>
      
      <button class="start-button" :disabled="isSubmitting" @click="handleLogin">
        {{ isSubmitting ? 'Starting...' : 'Start' }}
      </button>
      
      <p class="terms">
        上のボタンを押すと利用規約に同意したものとみなします
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'

const route = useRoute()
const router = useRouter()
const isSubmitting = ref(false)

const handleLogin = async () => {
  if (isSubmitting.value) return

  isSubmitting.value = true
  try {
    const response = await fetch(`/api/signup`,{
      method: "POST",
      headers:{
        "content-type":"application/json"
      },
      body: JSON.stringify({ agreed: true })
    });

    if(!response.ok){
      console.log("Error : Not OK")
      const errorText = await response.text()
      console.log("バックエンドから返ってきた生の文字:", errorText)
      toast.error("登録に失敗しました")
      return
    }

    const userData = await response.json();
    console.log("APIから取得したデータ:", userData)
    await router.push((route.query.redirect as string | undefined) ?? '/')
    
  }catch(error){
    console.log("Error : ",error)
    toast.error("通信エラーが発生しました")
  } finally {
    isSubmitting.value = false
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
.login-box {
  text-align: center;
}
.icon {
  width: 100px; /* アイコンサイズ */
  height: 100px;
}
.start-button {
  padding: 15px 40px;
  font-size: 1.2rem;
  cursor: pointer;
  margin-top: 20px;
  border: 1px solid #339af0;
  background: #4dabf7;
  position: absolute;
  top: 70%;
  transform: translateX(-50%);
}
.terms {
  margin-top: 50px;
  font-size: 0.8rem;
  color: #666;
  position: absolute;
  top: 80%;
  transform: translateX(-50%);
}
.login-title{
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 50px;
  margin-bottom: 50px;
  position: absolute;
  top: 30%;
  transform: translateX(-50%);
}
.login-title h1 {
  margin: 0;
  font-size: 7rem;
  font-weight: 900;
  letter-spacing: 0.5px;
  color: #ff4a7d;
  
}

</style>
