<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
/* 全局样式变量 */
:root {
  /* 暖色调配色 */
  --bg-color: #FAFAE8; /* RGB: 250,250,232 */
  --border-color: #F3F3CA; /* RGB: 243,243,202 */
  --primary-color: #FFD966;
  --button-hover: #FFC107;
  --text-color: #5D3A1A; /* 深棕色 */
  --progress-color: #F4A261;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  line-height: 1.5;
  font-weight: 400;
  background-color: var(--bg-color);
  color: var(--text-color);
  min-width: 320px;
  min-height: 100vh;
}

#app {
  width: 100%;
  min-height: 100vh;
}

/* 按钮通用样式 */
button {
  border-radius: 24px;
  border: 2px solid var(--border-color);
  padding: 0.8em 1.6em;
  font-size: 1em;
  font-weight: 600;
  font-family: inherit;
  background-color: var(--primary-color);
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

button:hover {
  background-color: var(--button-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

button:focus,
button:focus-visible {
  outline: 4px auto rgba(255, 217, 102, 0.5);
}

/* 卡片通用样式 */
.card {
  background-color: white;
  border-radius: 16px;
  border: 2px solid var(--border-color);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

/* 进度条样式 */
.progress-bar {
  width: 100%;
  height: 20px;
  background-color: var(--border-color);
  border-radius: 10px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background-color: var(--progress-color);
  transition: width 0.5s ease;
}

/* 高亮文本 */
.highlight {
  color: var(--primary-color);
  font-weight: bold;
}
</style>
