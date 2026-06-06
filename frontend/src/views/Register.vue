<template>
  <div class="register-container">
    <div class="register-card">
      <h1>词迹</h1>
      <h2>用户注册</h2>
      
      <div class="form-group">
        <label for="username">用户名</label>
        <input 
          type="text" 
          id="username" 
          v-model="form.username" 
          placeholder="请输入用户名（3-12个字符）"
        >
      </div>
      
      <div class="form-group">
        <label for="email">邮箱</label>
        <input 
          type="email" 
          id="email" 
          v-model="form.email" 
          placeholder="请输入邮箱地址"
        >
      </div>
      
      <div class="form-group">
        <label for="password">密码</label>
        <input 
          type="password" 
          id="password" 
          v-model="form.password" 
          placeholder="请输入密码（至少7个字符，不能包含中文）"
        >
      </div>
      
      <div class="form-group">
        <label for="confirmPassword">确认密码</label>
        <input 
          type="password" 
          id="confirmPassword" 
          v-model="form.confirmPassword" 
          placeholder="请再次输入密码"
        >
      </div>
      
      <div class="button-group">
        <button @click="handleRegister" class="primary-button">注册</button>
        <button @click="goToLogin" class="secondary-button">返回登录</button>
      </div>
    </div>
  </div>
</template>

<script>
import { useUserStore } from '../stores/user'

export default {
  name: 'Register',
  data() {
    return {
      form: {
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
      }
    }
  },
  methods: {
    validateForm() {
      // 验证用户名
      if (!this.form.username || this.form.username.length < 3) {
        alert('用户名必须大于2个字符')
        return false
      }
      if (this.form.username.length > 12) {
        alert('用户名不能超过12个字符')
        return false
      }

      // 验证邮箱
      if (!this.form.email) {
        alert('请输入邮箱地址')
        return false
      }
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(this.form.email)) {
        alert('请输入有效的邮箱地址')
        return false
      }

      // 验证密码
      if (!this.form.password || this.form.password.length < 7) {
        alert('密码必须大于6个字符')
        return false
      }

      // 验证密码不能包含中文
      for (let char of this.form.password) {
        if (char >= '\u4e00' && char <= '\u9fff') {
          alert('密码不能包含中文字符')
          return false
        }
      }

      // 验证确认密码
      if (this.form.password !== this.form.confirmPassword) {
        alert('两次输入的密码不一致')
        return false
      }

      return true
    },
    async handleRegister() {
      if (!this.validateForm()) {
        return
      }

      try {
        const userStore = useUserStore()
        await userStore.register({
          username: this.form.username,
          email: this.form.email,
          password: this.form.password
        })
        alert('注册成功，请登录')
        this.$router.push('/login')
      } catch (error) {
        if (error.response && error.response.data && error.response.data.error) {
          alert('注册失败: ' + error.response.data.error)
        } else {
          alert('注册失败: 无法连接到服务器，请稍后重试')
        }
      }
    },
    goToLogin() {
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--bg-color);
}

.register-card {
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
  box-sizing: border-box;
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
</style>
