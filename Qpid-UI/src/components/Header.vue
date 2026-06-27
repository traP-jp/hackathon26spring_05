<template>
  <header class="app-header">
    <div class="header-container">
      <div class="logo">
        <a href="/">
          <span class="logo-icon">🤝</span>
          <span class="logo-text">フレンズLink</span>
        </a>
      </div>

      <nav class="nav-menu">
        <ul>
          <li><a href="#" class="active">ホーム</a></li>
          <li><a href="#">探す</a></li>
          <li><a href="#">メッセージ</a></li>
          <li><a href="#">コミュニティ</a></li>
        </ul>
      </nav>

      <div class="user-actions">
        <button class="icon-btn" aria-label="通知">
          <span class="icon">🔔</span>
        </button>
        <div class="user-profile">
          <img src="https://via.placeholder.com/40" alt="マイページ" class="avatar" />
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
/* ヘッダー全体の固定と背景 */
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 70px;
  background-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  z-index: 1000;
  display: flex;
  align-items: center;
}

/* コンテナで中央寄せと余白を設定 */
.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* ロゴのスタイル */
.logo a {
  display: flex;
  align-items: center;
  text-decoration: none;
  gap: 8px;
}
.logo-icon {
  font-size: 1.5rem;
}
.logo-text {
  font-size: 1.25rem;
  font-weight: bold;
  color: #ff5a60; /* マッチング系アプリでよく使われるアクティブな色 */
}

/* ナビゲーションメニューのスタイル */
.nav-menu ul {
  display: flex;
  list-style: none;
  gap: 30px;
  margin: 0;
  padding: 0;
}
.nav-menu a {
  text-decoration: none;
  color: #555555;
  font-weight: 500;
  transition: color 0.2s ease;
  padding: 8px 4px;
}
.nav-menu a:hover, .nav-menu a.active {
  color: #ff5a60;
  border-bottom: 2px solid #ff5a60;
}

/* ユーザーアクションエリア */
.user-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}
.icon-btn {
  background: none;
  border: none;
  font-size: 1.25rem;
  cursor: pointer;
  padding: 4px;
}
.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #eee;
  cursor: pointer;
}

/* レスポンシブ対応（スマホサイズの場合） */
@media (max-width: 768px) {
  .nav-menu {
    /* スマホでは下部ナビゲーションにするか、ハンバーガーメニューにするのが一般的 */
    /* 今回は簡易的に非表示、または要素を絞る設定に調整してください */
    display: none; 
  }
}
</style>