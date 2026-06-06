<template>
  <div class="admin-container">
    <div class="header">
      <h1>词迹 - 单词管理</h1>
      <div class="user-info">
        <span class="username">{{ userStore.user.username }}</span>
        <span class="role admin">{{ userStore.user.role }}</span>
        <button @click="handleLogout" class="logout-button">登出</button>
      </div>
    </div>

    <div class="admin-nav">
      <button @click="$router.push('/admin/words')" class="nav-button active">单词管理</button>
      <button @click="$router.push('/monitor')" class="nav-button">后台监控</button>
    </div>

    <div class="words-container">
      <div class="action-bar">
        <div class="search-bar">
          <input type="text" v-model="keyword" placeholder="搜索单词或释义..." @keyup.enter="search" class="search-input">
          <button @click="search" class="search-button">搜索</button>
        </div>
        <button @click="showAddModal = true" class="add-button">添加单词</button>
      </div>

      <div class="words-list">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>单词</th>
              <th>中文释义</th>
              <th>例句</th>
              <th>级别</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="words.length === 0">
              <td colspan="6" class="empty-tip">暂无数据</td>
            </tr>
            <tr v-for="word in words" :key="word.id">
              <td>{{ word.id }}</td>
              <td>{{ word.word }}</td>
              <td>{{ word.translation }}</td>
              <td v-html="word.example_sentence"></td>
              <td>{{ word.level === 'cet4' ? '四级' : '全部' }}</td>
              <td>
                <button @click="editWord(word)" class="edit-button">编辑</button>
                <button @click="deleteWord(word.id)" class="delete-button">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="pagination">
        <button :disabled="currentPage <= 1" @click="goToPage(currentPage - 1)" class="page-btn">上一页</button>
        <button v-for="p in pageRange" :key="p" @click="goToPage(p)"
          :class="['page-btn', { active: p === currentPage }]">{{ p }}</button>
        <button :disabled="currentPage >= totalPages" @click="goToPage(currentPage + 1)" class="page-btn">下一页</button>
        <span class="page-info">共 {{ totalWords }} 条，{{ totalPages }} 页</span>
      </div>
    </div>

    <!-- 添加/编辑单词模态框 -->
    <div v-if="showAddModal || showEditModal" class="modal">
      <div class="modal-content">
        <h3>{{ showAddModal ? '添加单词' : '编辑单词' }}</h3>

        <div class="form-group">
          <label for="word">单词</label>
          <input type="text" id="word" v-model="formData.word" placeholder="请输入单词">
        </div>

        <div class="form-group">
          <label for="translation">中文释义</label>
          <input type="text" id="translation" v-model="formData.translation" placeholder="请输入中文释义">
        </div>

        <div class="form-group">
          <label for="example_sentence">例句</label>
          <input type="text" id="example_sentence" v-model="formData.example_sentence" placeholder="请输入例句，单词会自动标红">
        </div>

        <div class="form-group">
          <label for="level">级别</label>
          <select id="level" v-model="formData.level">
            <option value="all">全部</option>
            <option value="cet4">四级</option>
          </select>
        </div>

        <div class="modal-buttons">
          <button @click="showAddModal = false; showEditModal = false" class="cancel-button">取消</button>
          <button @click="saveWord" class="save-button">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useUserStore } from '../stores/user'
import axios from 'axios'

export default {
  name: 'AdminWords',
  data() {
    return {
      words: [],
      currentPage: 1,
      pageSize: 20,
      totalPages: 0,
      totalWords: 0,
      keyword: '',
      showAddModal: false,
      showEditModal: false,
      formData: {
        id: '',
        word: '',
        translation: '',
        example_sentence: '',
        level: 'all'
      }
    }
  },
  computed: {
    userStore() {
      return useUserStore()
    },
    pageRange() {
      const range = []
      const total = this.totalPages
      let start = Math.max(1, this.currentPage - 2)
      let end = Math.min(total, start + 4)
      if (end - start < 4) {
        start = Math.max(1, end - 4)
      }
      for (let i = start; i <= end; i++) {
        range.push(i)
      }
      return range
    }
  },
  mounted() {
    this.loadWords()
  },
  methods: {
    async loadWords() {
      try {
        const params = { page: this.currentPage, page_size: this.pageSize }
        if (this.keyword.trim()) {
          params.keyword = this.keyword.trim()
        }
        const response = await axios.get('http://localhost:8080/api/admin/words', {
          headers: {
            Authorization: `Bearer ${this.userStore.token}`
          },
          params
        })
        const res = response.data
        this.words = res.data || []
        this.totalWords = res.total || 0
        this.totalPages = res.total_pages || 0
      } catch (error) {
        const msg = error.response?.data?.error || '请求失败'
        alert('获取单词失败: ' + msg)
      }
    },
    goToPage(page) {
      if (page < 1 || page > this.totalPages) return
      this.currentPage = page
      this.loadWords()
    },
    search() {
      this.currentPage = 1
      this.loadWords()
    },
    editWord(word) {
      this.formData = { ...word }
      this.showEditModal = true
      this.showAddModal = false
    },
    async saveWord() {
      try {
        if (this.showAddModal) {
          await axios.post('http://localhost:8080/api/admin/words', this.formData, {
            headers: {
              Authorization: `Bearer ${this.userStore.token}`
            }
          })
        } else if (this.showEditModal) {
          await axios.put(`http://localhost:8080/api/admin/words/${this.formData.id}`, this.formData, {
            headers: {
              Authorization: `Bearer ${this.userStore.token}`
            }
          })
        }
        this.showAddModal = false
        this.showEditModal = false
        this.currentPage = 1
        this.keyword = ''
        this.loadWords()
      } catch (error) {
        const msg = error.response?.data?.error || '保存失败'
        alert('保存失败: ' + msg)
      }
    },
    async deleteWord(id) {
      if (confirm('确定要删除这个单词吗？')) {
        try {
          await axios.delete(`http://localhost:8080/api/admin/words/${id}`, {
            headers: {
              Authorization: `Bearer ${this.userStore.token}`
            }
          })
          this.loadWords()
        } catch (error) {
          const msg = error.response?.data?.error || '删除失败'
          alert('删除失败: ' + msg)
        }
      }
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

.admin-nav {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
}

.nav-button {
  padding: 12px 24px;
  border: 2px solid var(--border-color);
  border-radius: 24px;
  background-color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-button:hover {
  border-color: var(--primary-color);
  transform: translateY(-2px);
}

.nav-button.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.words-container {
  background-color: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border: 2px solid var(--border-color);
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 12px;
  flex-wrap: wrap;
}

.search-bar {
  display: flex;
  gap: 8px;
  flex: 1;
  max-width: 400px;
}

.search-input {
  flex: 1;
  padding: 10px 16px;
  border: 2px solid var(--border-color);
  border-radius: 24px;
  font-size: 14px;
  outline: none;
  transition: all 0.3s ease;
}

.search-input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(245, 124, 0, 0.1);
}

.search-button {
  padding: 10px 20px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-1px);
}

.add-button {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.add-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.words-list {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

th {
  background-color: var(--bg-color);
  font-weight: 600;
  color: var(--text-color);
}

tr:hover {
  background-color: rgba(245, 124, 0, 0.05);
}

.empty-tip {
  text-align: center;
  padding: 40px 0;
  color: #999;
  font-size: 16px;
}

.edit-button, .delete-button {
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  margin-right: 5px;
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

/* 分页样式 */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 24px;
  flex-wrap: wrap;
}

.page-btn {
  min-width: 36px;
  padding: 8px 14px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  background-color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-color);
}

.page-btn:hover:not(:disabled):not(.active) {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.page-btn.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
  font-weight: 600;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  margin-left: 12px;
  font-size: 13px;
  color: #999;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);
  border: 2px solid var(--border-color);
  max-width: 500px;
  width: 100%;
}

.modal-content h3 {
  color: var(--text-color);
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: var(--text-color);
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 12px;
  border: 2px solid var(--border-color);
  border-radius: 12px;
  font-size: 16px;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(245, 124, 0, 0.1);
}

.modal-buttons {
  display: flex;
  gap: 15px;
  justify-content: flex-end;
  margin-top: 30px;
}

.cancel-button,
.save-button {
  padding: 10px 20px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--border-color);
}

.cancel-button {
  background-color: white;
  color: var(--text-color);
}

.cancel-button:hover {
  background-color: #f5f5f5;
  transform: translateY(-2px);
}

.save-button {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.save-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}
</style>
