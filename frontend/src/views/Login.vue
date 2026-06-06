<template>
  <div class="login-container">
    <div class="login-card">
      <h1>词迹</h1>
      <h2>英语单词学习</h2>
      
      <div class="form-group">
        <label for="username">用户名</label>
        <input type="text" id="username" v-model="username" placeholder="请输入用户名">
      </div>
      
      <div class="form-group">
        <label for="password">密码</label>
        <input type="password" id="password" v-model="password" placeholder="请输入密码">
      </div>
      
      <div class="button-group">
        <button @click="handleLogin" class="primary-button">登录</button>
      </div>
      
      <div class="register-link">
        <span>还没有账号？</span>
        <a @click="goToRegister" class="link">立即注册</a>
      </div>
    </div>
  </div>
</template>

<script>
import { useUserStore } from '../stores/user'

export default {
  name: 'Login',
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    async handleLogin() {
      try {
        const userStore = useUserStore()
        await userStore.login(this.username, this.password)
        if (userStore.isAdmin) {
          this.$router.push('/admin')
        } else {
          this.$router.push('/learn')
        }
      } catch (error) {
        if (error.response && error.response.data && error.response.data.error) {
          alert('登录失败: ' + error.response.data.error)
        } else {
          alert('登录失败: 无法连接到服务器，请稍后重试')
        }
      }
    },
    goToRegister() {
      this.$router.push('/register')
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--bg-color);
}

.login-card {
  background-color: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border: 2px solid var(--border-color);
  max-width: 400px;
  width: 100%;
  text-align: center;
}

h1 {
  color: var(--primary-color);
  font-size: 2.5em;
  margin-bottom: 10px;
}

h2 {
  color: var(--text-color);
  font-size: 1.2em;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
  text-align: left;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: var(--text-color);
}

input {
  width: 100%;
  padding: 12px;
  border: 2px solid var(--border-color);
  border-radius: 12px;
  font-size: 16px;
  transition: all 0.3s ease;
}

input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(245, 124, 0, 0.1);
}

.button-group {
  display: flex;
  gap: 15px;
  margin-top: 30px;
}

.primary-button {
  flex: 1;
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 12px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.primary-button:hover {
  background-color: var(--button-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.secondary-button {
  flex: 1;
  background-color: white;
  color: var(--primary-color);
  border: 2px solid var(--primary-color);
  padding: 12px;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.secondary-button:hover {
  background-color: var(--primary-color);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.register-link {
  margin-top: 20px;
  font-size: 14px;
  color: #666;
}

.register-link .link {
  color: var(--primary-color);
  cursor: pointer;
  text-decoration: none;
  font-weight: 600;
  margin-left: 5px;
}

.register-link .link:hover {
  text-decoration: underline;
}
</style>
