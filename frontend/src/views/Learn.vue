<template>
  <div class="learn-container">
    <div class="header">
      <h1>词迹</h1>
      <div class="user-info">
        <span class="username">{{ userStore.user.username }}</span>
        <span class="role" :class="userStore.user.role">{{ userStore.user.role }}</span>
        <button @click="goToStats" class="stats-button">学习状况</button>
        <button v-if="userStore.user.role === 'admin'" @click="goToMonitor" class="monitor-button">监控系统</button>
        <button v-if="userStore.user.role === 'admin'" @click="goToAdmin" class="admin-button">管理后台</button>
        <button v-if="userStore.user.role !== 'admin'" @click="goToUserAdmin" class="user-admin-button">我的单词</button>
        <button @click="handleLogout" class="logout-button">登出</button>
      </div>
    </div>

    <div class="nav-tabs">
      <button 
        :class="['tab-button', { active: activeTab === 'all' }]"
        @click="switchTab('all')"
      >
        全部学习
      </button>
      <button 
        :class="['tab-button', { active: activeTab === 'cet4' }]"
        @click="switchTab('cet4')"
      >
        四级词汇
      </button>
      <button 
        :class="['tab-button', { active: activeTab === 'custom' }]"
        @click="switchTab('custom')"
      >
        我的单词
      </button>
    </div>

    <!-- 搜索栏 -->
    <div class="search-section">
      <div class="search-box">
        <svg class="search-icon" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <path d="M21 21l-4.35-4.35"/>
        </svg>
        <input
          v-model="searchQuery"
          @input="handleSearchInput"
          placeholder="搜索单词或中文释义..."
          class="search-input"
        />
        <button v-if="searchQuery" @click="clearSearch" class="search-clear">✕</button>
      </div>

      <!-- 搜索结果 -->
      <div v-if="searchResults.length > 0" class="search-results">
        <div class="search-results-header">
          搜索结果 ({{ searchResults.length }})
        </div>
        <div
          v-for="result in searchResults"
          :key="result.source + '-' + result.id"
          class="search-result-item"
          @click="selectSearchResult(result)"
        >
          <div class="result-word">
            {{ result.word }}
            <span class="result-source" :class="result.source">{{ result.source === 'system' ? '词库' : '自定义' }}</span>
          </div>
          <div class="result-translation">{{ result.translation }}</div>
          <div v-if="result.example_sentence" class="result-sentence" v-html="result.example_sentence"></div>
        </div>
      </div>

      <!-- 无结果提示 -->
      <div v-if="searchPerformed && searchQuery && searchResults.length === 0" class="search-no-results">
        未找到匹配的单词
      </div>
    </div>

    <div class="learning-content">
      <!-- 进度条 -->
      <div class="progress-section">
        <div class="progress-info">
          <span>本轮学习进度: {{ userStore.completed }} / {{ userStore.total }}</span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: userStore.progress * 100 + '%' }"></div>
        </div>
      </div>

      <!-- 自定义词库：添加单词 -->
      <div v-if="activeTab === 'custom'" class="custom-section">
        <button @click="showAddWord = !showAddWord" class="add-word-toggle">
          {{ showAddWord ? '收起' : '+ 添加单词' }}
        </button>
        <div v-if="showAddWord" class="add-word-form">
          <input v-model="customForm.word" placeholder="单词" class="form-input" />
          <input v-model="customForm.translation" placeholder="中文释义" class="form-input" />
          <input v-model="customForm.example_sentence" placeholder="例句（可选）" class="form-input" />
          <input v-model="customForm.example_sentence_translation" placeholder="例句翻译（可选）" class="form-input" />
          <div class="form-actions">
            <button @click="addCustomWord" class="save-btn" :disabled="!customForm.word || !customForm.translation">保存</button>
          </div>
          <div v-if="addWordMsg" class="form-msg" :class="{ success: addWordSuccess, error: !addWordSuccess }">{{ addWordMsg }}</div>
        </div>
      </div>

      <!-- 学习完成弹窗 -->
      <div v-if="showCompletionModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
          <div class="emoji">🎉</div>
          <h3>恭喜，您已完成本轮学习</h3>
          <p>您已完成本轮所有单词的学习</p>
          <button @click="restartLearning" class="restart-button">进行新的一轮</button>
        </div>
      </div>

      <!-- 单词卡片 -->
      <div v-if="userStore.currentWord" class="word-card">
        <div class="word">{{ userStore.currentWord?.word }}</div>
        
        <!-- 当点击不认识时显示中文翻译、短语和短语翻译 -->
        <div v-if="userStore.showMeaning" class="meaning-section">
          <div class="meaning">{{ userStore.currentWord?.translation }}</div>
        </div>
        <div v-if="userStore.showMeaning && userStore.currentWord?.example_sentence" class="sentence-section">
          <div class="sentence" v-html="userStore.currentWord?.example_sentence"></div>
          <div class="sentence-translation">{{ userStore.currentWord?.example_sentence_translation }}</div>
        </div>

        <!-- 按钮组 -->
        <div class="button-group">
          <button 
            v-if="!userStore.showMeaning" 
            @click="handleKnown" 
            class="known-button"
          >
            ✅ 认识
          </button>
          <button 
            v-if="!userStore.showMeaning" 
            @click="handleUnknown" 
            class="unknown-button"
          >
            ❓ 不认识
          </button>
          <button 
            v-if="userStore.showMeaning" 
            @click="handleNext" 
            class="next-button"
          >
            下一个
          </button>
        </div>
      </div>

      <!-- 没有单词时显示提示 -->
      <div v-if="!userStore.currentWord && userStore.total === 0" class="placeholder-card">
        <div class="emoji">📚</div>
        <h3>{{ activeTab === 'custom' ? '暂无自定义单词' : '该词库暂无单词' }}</h3>
        <p>{{ activeTab === 'custom' ? '点击上方"添加单词"开始创建你的专属词库~' : '请联系管理员添加单词~' }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { useUserStore } from '../stores/user'
import apiClient from '../api/axios'

export default {
  name: 'Learn',
  data() {
    return {
      activeTab: 'all',
      showCompletionModal: false,
      showAddWord: false,
      addWordMsg: '',
      addWordSuccess: false,
      customForm: {
        word: '',
        translation: '',
        example_sentence: '',
        example_sentence_translation: ''
      },
      searchQuery: '',
      searchResults: [],
      searchPerformed: false,
      searchTimer: null,
      roundComplete: false
    }
  },
  watch: {
    'userStore.progress': {
      handler(newVal) {
        // 不再自动弹出完成弹窗，由 roundComplete 标志控制
      },
      immediate: true
    }
  },
  computed: {
    userStore() {
      return useUserStore()
    }
  },
  mounted() {
    this.loadAllWords()
  },
  methods: {
    async switchTab(level) {
      this.activeTab = level
      this.showCompletionModal = false
      this.roundComplete = false
      await this.loadWords(level)
    },
    async loadWords(level) {
      await this.userStore.getNextWord(level)
      await this.userStore.getProgress(level)
      this.userStore.showMeaning = false
    },
    async loadAllWords() {
      await this.loadWords(this.activeTab)
    },
    async handleKnown() {
      if (this.userStore.currentWord) {
        if (this.activeTab === 'custom') {
          await this.userStore.submitCustomAnswer(this.userStore.currentWord.id, true)
        } else {
          await this.userStore.submitAnswer(this.userStore.currentWord.id, true)
        }
        // 先检查进度，判断是否是最后一个单词
        await this.userStore.getProgress(this.activeTab)
        if (this.userStore.progress >= 1) {
          // 最后一个单词：显示释义，等待用户点击"下一个"才弹出完成弹窗
          this.userStore.showMeaning = true
          this.roundComplete = true
        } else {
          await this.userStore.getNextWord(this.activeTab)
          this.userStore.showMeaning = false
        }
      }
    },
    async handleUnknown() {
      if (this.userStore.currentWord) {
        this.userStore.showMeaning = true
        if (this.activeTab === 'custom') {
          await apiClient.post('/learn/custom/answer', {
            word_id: this.userStore.currentWord.id,
            known: false
          })
        } else {
          await apiClient.post('/learn/answer', {
            word_id: this.userStore.currentWord.id,
            known: false
          })
        }
        await this.userStore.getProgress(this.activeTab)
        // 检测是否是最后一个单词
        if (this.userStore.progress >= 1) {
          this.roundComplete = true
        }
      }
    },
    async handleNext() {
      if (this.roundComplete) {
        // 本轮已完成，弹出完成弹窗
        this.showCompletionModal = true
        this.roundComplete = false
        return
      }
      await this.userStore.getNextWord(this.activeTab)
      await this.userStore.getProgress(this.activeTab)
      this.userStore.showMeaning = false
    },
    handleSearchInput() {
      if (this.searchTimer) clearTimeout(this.searchTimer)
      if (!this.searchQuery.trim()) {
        this.searchResults = []
        this.searchPerformed = false
        return
      }
      this.searchTimer = setTimeout(() => {
        this.performSearch()
      }, 300)
    },
    async performSearch() {
      if (!this.searchQuery.trim()) return
      try {
        const response = await apiClient.get('/learn/search', {
          params: { q: this.searchQuery.trim() }
        })
        this.searchResults = response.data.results || []
        this.searchPerformed = true
      } catch (error) {
        console.error('搜索失败:', error)
      }
    },
    clearSearch() {
      this.searchQuery = ''
      this.searchResults = []
      this.searchPerformed = false
      if (this.searchTimer) clearTimeout(this.searchTimer)
    },
    selectSearchResult(result) {
      this.userStore.currentWord = result
      this.userStore.showMeaning = false
      this.searchResults = []
      this.searchQuery = ''
      this.searchPerformed = false
      this.roundComplete = false
    },
    handleLogout() {
      this.userStore.logout()
      this.$router.push('/login')
    },
    goToStats() {
      this.$router.push('/stats')
    },
    goToMonitor() {
      this.$router.push('/monitor')
    },
    goToAdmin() {
      this.$router.push('/admin')
    },
    goToUserAdmin() {
      this.$router.push('/user-admin')
    },
    async addCustomWord() {
      if (!this.customForm.word || !this.customForm.translation) return
      try {
        await apiClient.post('/custom/words', this.customForm)
        this.addWordSuccess = true
        this.addWordMsg = '添加成功！'
        this.customForm = { word: '', translation: '', example_sentence: '', example_sentence_translation: '' }
        this.showAddWord = false
        // 重新加载学习单词
        await this.loadWords('custom')
      } catch (error) {
        this.addWordSuccess = false
        this.addWordMsg = '添加失败：' + (error.response?.data?.error || error.message)
      }
      setTimeout(() => { this.addWordMsg = '' }, 3000)
    },
    async restartLearning() {
      this.showCompletionModal = false
      this.roundComplete = false
      await this.userStore.getNextWord(this.activeTab)
      await this.userStore.getProgress(this.activeTab)
    },
    closeModal() {
      this.showCompletionModal = false
    }
  }
}
</script>

<style scoped>
.learn-container {
  padding: 20px;
  min-height: 100vh;
  background-color: var(--bg-color);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h1 {
  color: var(--primary-color);
  font-size: 2em;
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  font-weight: 600;
  color: var(--text-color);
}

.role {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
}

.role.admin {
  background-color: #e3f2fd;
  color: #1976d2;
}

.role.user {
  background-color: #e8f5e8;
  color: #388e3c;
}

.logout-button {
  background-color: #f44336;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-button:hover {
  background-color: #d32f2f;
  transform: translateY(-2px);
}

.stats-button {
  background-color: #00bcd4;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.stats-button:hover {
  background-color: #0097a7;
  transform: translateY(-2px);
}

.monitor-button {
  background-color: #9c27b0;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.monitor-button:hover {
  background-color: #7b1fa2;
  transform: translateY(-2px);
}

.admin-button {
  background-color: #ffc107;
  color: #333;
  border: none;
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-button:hover {
  background-color: #ffb300;
  transform: translateY(-2px);
}

.user-admin-button {
  background-color: #673ab7;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-admin-button:hover {
  background-color: #5e35b1;
  transform: translateY(-2px);
}

.nav-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
}

.tab-button {
  padding: 12px 24px;
  border: 2px solid var(--border-color);
  border-radius: 24px;
  background-color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tab-button:hover {
  border-color: var(--primary-color);
  transform: translateY(-2px);
}

.tab-button.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.learning-content {
  max-width: 800px;
  margin: 0 auto;
}

.progress-section {
  margin-bottom: 30px;
}

.progress-info {
  margin-bottom: 10px;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
}

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

.word-card {
  background-color: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border: 2px solid var(--border-color);
  text-align: center;
}

.word {
  font-size: 3em;
  font-weight: bold;
  color: #000000;
  margin-bottom: 30px;
}

.meaning-section {
  margin-bottom: 20px;
  padding: 15px 20px;
  background-color: #fff3e0;
  border-radius: 12px;
  border-left: 4px solid #ff9800;
}

.meaning {
  font-size: 1.8em;
  color: #e65100;
  font-weight: 600;
  text-align: center;
}

.sentence-section {
  margin-bottom: 30px;
  text-align: left;
}

.sentence {
  font-size: 1.2em;
  font-style: italic;
  color: #333;
  margin-bottom: 10px;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 12px;
  border-left: 4px solid #2196f3;
}

.sentence .highlight {
  color: var(--primary-color);
  font-weight: bold;
  background-color: #fff3e0;
  padding: 2px 6px;
  border-radius: 4px;
}

.sentence-translation {
  font-size: 1.1em;
  color: #2e7d32;
  padding: 15px;
  background-color: #e8f5e8;
  border-radius: 12px;
  border-left: 4px solid #4caf50;
}

.button-group {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-top: 30px;
}

.known-button,
.unknown-button,
.next-button {
  padding: 12px 32px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--border-color);
}

.known-button {
  background-color: #4caf50;
  color: white;
}

.known-button:hover {
  background-color: #388e3c;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.unknown-button {
  background-color: #ff9800;
  color: white;
}

.unknown-button:hover {
  background-color: #f57c00;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.next-button {
  background-color: var(--primary-color);
  color: white;
}

.next-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.cet4-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 500px;
}

.placeholder-card {
  background-color: white;
  border-radius: 16px;
  padding: 60px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border: 2px solid var(--border-color);
  text-align: center;
}

.emoji {
  font-size: 4em;
  margin-bottom: 20px;
}

.placeholder-card h3 {
  color: var(--text-color);
  margin-bottom: 10px;
}

.placeholder-card p {
  color: #666;
  font-size: 16px;
}

.completion-card {
  background-color: white;
  border-radius: 16px;
  padding: 60px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border: 2px solid var(--border-color);
  text-align: center;
}

.completion-card .emoji {
  font-size: 4em;
  margin-bottom: 20px;
}

.completion-card h3 {
  color: var(--primary-color);
  margin-bottom: 10px;
  font-size: 1.5em;
}

.completion-card p {
  color: #666;
  font-size: 16px;
  margin-bottom: 30px;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-content {
  background-color: white;
  border-radius: 20px;
  padding: 50px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
  text-align: center;
  max-width: 400px;
  width: 90%;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-content .emoji {
  font-size: 5em;
  margin-bottom: 20px;
  animation: bounce 1s ease infinite;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.modal-content h3 {
  color: var(--primary-color);
  margin-bottom: 15px;
  font-size: 1.8em;
}

.modal-content p {
  color: #666;
  font-size: 16px;
  margin-bottom: 30px;
}

.restart-button {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 12px 32px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.restart-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.custom-section {
  margin-bottom: 20px;
}

.add-word-toggle {
  width: 100%;
  padding: 12px;
  border: 2px dashed var(--border-color);
  border-radius: 12px;
  background-color: white;
  color: var(--primary-color);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.add-word-toggle:hover {
  border-color: var(--primary-color);
  background-color: #fffaf0;
  transform: translateY(-1px);
}

.add-word-form {
  margin-top: 12px;
  padding: 20px;
  background-color: white;
  border-radius: 12px;
  border: 2px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-input {
  padding: 10px 14px;
  border: 2px solid var(--border-color);
  border-radius: 10px;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(245, 124, 0, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.save-btn {
  padding: 8px 24px;
  border: none;
  border-radius: 20px;
  background-color: var(--primary-color);
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.save-btn:hover:not(:disabled) {
  background-color: var(--button-hover);
  transform: translateY(-1px);
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-msg {
  text-align: center;
  font-size: 13px;
  font-weight: 500;
}

.form-msg.success {
  color: #4caf50;
}

.form-msg.error {
  color: #f44336;
}

/* 搜索栏样式 */
.search-section {
  max-width: 800px;
  margin: 0 auto 20px;
}

.search-box {
  display: flex;
  align-items: center;
  background-color: white;
  border: 2px solid var(--border-color);
  border-radius: 24px;
  padding: 8px 16px;
  transition: all 0.3s ease;
}

.search-box:focus-within {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(245, 124, 0, 0.1);
}

.search-icon {
  color: #999;
  flex-shrink: 0;
  margin-right: 10px;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 15px;
  padding: 6px 0;
  background: transparent;
  color: var(--text-color);
}

.search-clear {
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  font-size: 16px;
  padding: 4px 8px;
  border-radius: 50%;
  transition: all 0.2s;
}

.search-clear:hover {
  background-color: #f0f0f0;
  color: #333;
}

.search-results {
  margin-top: 12px;
  background-color: white;
  border-radius: 16px;
  border: 2px solid var(--border-color);
  max-height: 400px;
  overflow-y: auto;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.search-results-header {
  padding: 12px 20px;
  font-weight: 600;
  font-size: 14px;
  color: #666;
  border-bottom: 1px solid var(--border-color);
  background-color: #fafafa;
  border-radius: 14px 14px 0 0;
}

.search-result-item {
  padding: 14px 20px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background-color: #fffaf0;
}

.result-word {
  font-size: 1.2em;
  font-weight: 700;
  color: var(--text-color);
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.result-source {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 8px;
}

.result-source.system {
  background-color: #e3f2fd;
  color: #1565c0;
}

.result-source.custom {
  background-color: #f3e5f5;
  color: #7b1fa2;
}

.result-translation {
  font-size: 15px;
  color: #e65100;
  font-weight: 500;
  margin-bottom: 4px;
}

.result-sentence {
  font-size: 13px;
  color: #555;
  font-style: italic;
}

.result-sentence .highlight {
  color: var(--primary-color);
  font-weight: bold;
}

.search-no-results {
  margin-top: 12px;
  padding: 30px;
  text-align: center;
  color: #999;
  background-color: white;
  border-radius: 16px;
  border: 2px solid var(--border-color);
}
</style>
