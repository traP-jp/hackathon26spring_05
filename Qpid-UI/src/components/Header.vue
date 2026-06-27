<template>
  <header class="app-header">
    <div class="header-container">
      <!-- 左側：ロゴエリア -->
      <div class="logo">
        <router-link to="/">
          <span class="logo-icon">🏹</span>
          <span class="logo-text">Qpid</span>
        </router-link>
      </div>

      <!-- 右側：アイコン付きメニュー（右端揃え） -->
      <nav class="nav-menu">
        <ul>
          <!-- ホーム画面へ -->
          <li>
            <router-link to="/" class="nav-item" title="ホーム">
              <span class="nav-icon">🏠</span>
              <span class="nav-label">ホーム</span>
            </router-link>
          </li>
          <!-- 自己紹介（Myself）画面へ -->
          <li>
            <router-link to="/myself" class="nav-item" title="自己紹介">
              <span class="nav-icon">👤</span>
              <span class="nav-label">自己紹介</span>
            </router-link>
          </li>
          <!-- いいねされた人（後で実装用：現状は/#へ） -->
          <li>
            <router-link to="/likes" class="nav-item" title="いいねされた人">
              <span class="nav-icon">❤️</span>
              <span class="nav-label">いいね</span>
            </router-link>
          </li>
        </ul>
      </nav>
    </div>
  </header>
</template>

<style scoped>
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

.header-container {
  display: flex;
  justify-content: space-between; /* ロゴを左、メニューを右に引き離す */
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
  font-size: 1.6rem;
}
.logo-text {
  font-size: 1.5rem;
  font-weight: 900;
  letter-spacing: 0.5px;
  color: #ff4a7d; /* キューピッドらしい華やかなピンク */
}

/* 右側メニュー（右端揃え） */
.nav-menu ul {
  display: flex;
  list-style: none;
  gap: 24px; /* アイコン同士の間隔 */
  margin: 0;
  padding: 0;
}

/* ナビゲーションアイテム単体のスタイル（縦並びにしてアプリ感を出す） */
.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #666666;
  gap: 4px;
  font-size: 0.8rem;
  font-weight: 500;
  transition: all 0.2s ease;
  padding: 4px 8px;
}

.nav-icon {
  font-size: 1.4rem; /* アイコンをやや大きめに */
}

/* ホバー時と、現在開いているアクティブなページのスタイル */
.nav-item:hover,
.router-link-active {
  color: #ff4a7d !important;
  transform: translateY(-2px); /* 軽く上に浮き上がらせるエフェクト */
}

/* スマホ対応：画面が小さくなったらラベルを消してアイコンだけにする */
@media (max-width: 480px) {
  .nav-label {
    display: none;
  }
  .nav-menu ul {
    gap: 16px;
  }
}
</style>