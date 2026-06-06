<template>
  <div class="admin-container">
    <div class="header">
      <h1>用户中心</h1>
      <div class="user-info">
        <span class="username">{{ userStore.user.username }}</span>
        <span class="role" :class="userStore.user.role">{{ userStore.user.role }}</span>
        <button @click="goToLearn" class="learn-button">开始学习</button>
        <button @click="handleLogout" class="logout-button">登出</button>
      </div>
    </div>

    <div class="content">
      <div class="form-section">
        <h2>{{ editingWord ? '编辑单词' : '添加自定义单词' }}</h2>
        <form @submit.prevent="handleSubmit" class="word-form">
          <div class="form-group">
            <label for="word">单词</label>
            <input type="text" id="word" v-model="formData.word" placeholder="请输入单词" required>
          </div>
          <div class="form-group">
            <label for="translation">单词翻译</label>
            <input type="text" id="translation" v-model="formData.translation" placeholder="请输入单词翻译" required>
          </div>
          <div class="form-group">
            <label for="exampleSentence">短语/例句</label>
            <textarea id="exampleSentence" v-model="formData.example_sentence" placeholder="请输入短语或例句"></textarea>
          </div>
          <div class="form-group">
            <label for="exampleSentenceTranslation">短语/例句翻译</label>
            <textarea id="exampleSentenceTranslation" v-model="formData.example_sentence_translation" placeholder="请输入短语或例句翻译"></textarea>
          </div>
          <div class="button-group">
            <button type="submit" class="submit-button">{{ editingWord ? '保存' : '添加' }}</button>
            <button type="button" v-if="editingWord" @click="cancelEdit" class="cancel-button">取消</button>
          </div>
        </form>
      </div>

      <div class="list-section">
        <h2>我的单词库</h2>
        <div class="word-list">
          <div v-if="words.length === 0" class="empty-state">
            暂无自定义单词，添加一些单词开始学习吧！
          </div>
          <div v-for="word in words" :key="word.id" class="word-item">
            <div class="word-info">
              <div class="word-header">
                <span class="word-text">{{ word.word }}</span>
                <span class="word-level custom">自定义</span>
              </div>
              <div class="word-detail">
                <p><strong>翻译：</strong>{{ word.translation }}</p>
                <p v-if="word.example_sentence"><strong>例句：</strong>{{ word.example_sentence }}</p>
                <p v-if="word.example_sentence_translation"><strong>例句翻译：</strong>{{ word.example_sentence_translation }}</p>
              </div>
            </div>
            <div class="word-actions">
              <button @click="editWord(word)" class="edit-button">编辑</button>
              <button @click="deleteWord(word)" class="delete-button">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useUserStore } from '../stores/user'
import apiClient from '../api/axios'

export default {
  name: 'UserAdmin',
  data() {
    return {
      words: [],
      editingWord: null,
      formData: {
        word: '',
        translation: '',
        example_sentence: '',
        example_sentence_translation: ''
      }
    }
  },
  computed: {
    userStore() {
      return useUserStore()
    }
  },
  mounted() {
    this.loadWords()
  },
  methods: {
    async loadWords() {
      try {
        const response = await apiClient.get('/custom/words')
        this.words = response.data
      } catch (error) {
        console.error('获取单词列表失败:', error)
      }
    },
    async handleSubmit() {
      try {
        if (this.editingWord) {
          await apiClient.put(`/custom/words/${this.editingWord.id}`, this.formData)
        } else {
          await apiClient.post('/custom/words', this.formData)
        }
        this.resetForm()
        await this.loadWords()
      } catch (error) {
        console.error('操作失败:', error)
        alert('操作失败，请稍后重试')
      }
    },
    editWord(word) {
      this.editingWord = word
      this.formData = {
        word: word.word,
        translation: word.translation,
        example_sentence: word.example_sentence || '',
        example_sentence_translation: word.example_sentence_translation || ''
      }
    },
    cancelEdit() {
      this.resetForm()
    },
    async deleteWord(word) {
      if (!confirm(`确定要删除单词"${word.word}"吗？`)) {
        return
      }
      try {
        await apiClient.delete(`/custom/words/${word.id}`)
        await this.loadWords()
      } catch (error) {
        console.error('删除失败:', error)
        alert('删除失败，请稍后重试')
      }
    },
    resetForm() {
      this.editingWord = null
      this.formData = {
        word: '',
        translation: '',
        example_sentence: '',
        example_sentence_translation: ''
      }
    },
    goToLearn() {
      this.$router.push('/learn')
    },
    handleLogout() {
      this.userStore.logout()
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.admin-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header h1 {
  color: var(--primary-color);
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  font-weight: 600;
}

.role {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.role.admin {
  background-color: #ffd54f;
  color: #f57f17;
}

.role.user {
  background-color: #e3f2fd;
  color: #1565c0;
}

.learn-button, .logout-button {
  padding: 8px 20px;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.learn-button {
  background-color: #4caf50;
  color: white;
}

.learn-button:hover {
  background-color: #43a047;
  transform: translateY(-1px);
}

.logout-button {
  background-color: #f44336;
  color: white;
}

.logout-button:hover {
  background-color: #e53935;
  transform: translateY(-1px);
}

.content {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 30px;
}

.form-section, .list-section {
  background-color: white;
  padding: 30px;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.form-section h2, .list-section h2 {
  color: var(--primary-color);
  margin-top: 0;
  margin-bottom: 25px;
}

.word-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.form-group input,
.form-group textarea {
  padding: 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.button-group {
  display: flex;
  gap: 12px;
  margin-top: 10px;
}

.submit-button, .cancel-button {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.submit-button {
  background-color: var(--primary-color);
  color: white;
}

.submit-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-1px);
}

.cancel-button {
  background-color: #f5f5f5;
  color: #666;
}

.cancel-button:hover {
  background-color: #e0e0e0;
}

.word-list {
  max-height: calc(100vh - 300px);
  overflow-y: auto;
}

.empty-state {
  text-align: center;
  color: #999;
  padding: 40px;
}

.word-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  margin-bottom: 15px;
  transition: all 0.3s ease;
}

.word-item:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  border-color: var(--primary-color);
}

.word-info {
  flex: 1;
  margin-right: 20px;
}

.word-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.word-text {
  font-size: 1.4em;
  font-weight: bold;
  color: #333;
}

.word-level {
  padding: 4px 10px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 600;
}

.word-level.custom {
  background-color: #e8f5e9;
  color: #2e7d32;
}

.word-detail p {
  margin: 6px 0;
  color: #555;
  font-size: 14px;
}

.word-detail strong {
  color: #333;
}

.word-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.edit-button, .delete-button {
  padding: 8px 18px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.edit-button {
  background-color: #2196f3;
  color: white;
}

.edit-button:hover {
  background-color: #1976d2;
  transform: translateY(-1px);
}

.delete-button {
  background-color: #f44336;
  color: white;
}

.delete-button:hover {
  background-color: #d32f2f;
  transform: translateY(-1px);
}
</style>