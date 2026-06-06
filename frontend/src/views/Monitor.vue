<template>
  <div class="monitor-container">
    <div class="header">
      <h1>后台监控系统</h1>
      <div class="user-info">
        <span class="username">{{ userStore.user.username }}</span>
        <span class="role" :class="userStore.user.role">{{ userStore.user.role }}</span>
        <button @click="goToAdmin" class="admin-button">单词管理</button>
        <button @click="goToLearn" class="learn-button">返回学习</button>
        <button @click="handleLogout" class="logout-button">登出</button>
      </div>
    </div>

    <div class="content">
      <!-- 统计卡片 -->
      <div class="stats-section">
        <div class="stat-card">
          <div class="stat-icon">👥</div>
          <div class="stat-content">
            <div class="stat-value">{{ overview.today_visits || 0 }}</div>
            <div class="stat-label">今日访问量</div>
          </div>
        </div>
      </div>

      <!-- 图表区域 -->
      <div class="charts-section">
        <!-- 时段访问量 -->
        <div class="chart-card">
          <h3>各时段访问量</h3>
          <div ref="hourlyChart" class="chart"></div>
        </div>

        <!-- 每日注册量 -->
        <div class="chart-card">
          <h3>近7天注册量</h3>
          <div ref="dailyChart" class="chart"></div>
        </div>
      </div>

      <!-- 错词排行 -->
      <div class="ranking-section">
        <h3>错词排行榜</h3>
        <div class="ranking-list">
          <div v-if="overview.wrong_words?.length === 0" class="empty-state">
            暂无错词数据
          </div>
          <div v-for="(item, index) in overview.wrong_words" :key="index" class="ranking-item">
            <div class="rank-number" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
            <div class="rank-content">
              <div class="word-info">
                <span class="word">{{ item.word }}</span>
                <span class="translation">{{ item.translation }}</span>
              </div>
              <div class="wrong-count">
                <span class="count-label">错误次数：</span>
                <span class="count-value">{{ item.wrong_count }}</span>
              </div>
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
import * as echarts from 'echarts'

export default {
  name: 'Monitor',
  data() {
    return {
      overview: {
        today_visits: 0,
        wrong_words: []
      },
      hourlyVisits: [],
      dailyRegistrations: [],
      hourlyChart: null,
      dailyChart: null
    }
  },
  computed: {
    userStore() {
      return useUserStore()
    }
  },
  mounted() {
    this.loadMonitorData()
  },
  methods: {
    async loadMonitorData() {
      try {
        const [overviewRes, hourlyRes, dailyRes] = await Promise.all([
          apiClient.get('/monitor/overview'),
          apiClient.get('/monitor/hourly-visits'),
          apiClient.get('/monitor/daily-registrations')
        ])

        this.overview = overviewRes.data
        this.hourlyVisits = hourlyRes.data
        this.dailyRegistrations = dailyRes.data

        this.$nextTick(() => {
          this.initHourlyChart()
          this.initDailyChart()
        })
      } catch (error) {
        console.error('加载监控数据失败:', error)
      }
    },
    initHourlyChart() {
      if (this.hourlyChart) {
        this.hourlyChart.dispose()
      }

      const hours = Array.from({ length: 24 }, (_, i) => `${i}:00`)
      
      this.hourlyChart = echarts.init(this.$refs.hourlyChart)
      const option = {
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(255,255,255,0.95)',
          borderColor: '#e0e0e0',
          textStyle: { color: '#333' }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: hours,
          axisLabel: {
            fontSize: 11
          }
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '访问量',
            type: 'line',
            smooth: true,
            data: this.hourlyVisits,
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: 'rgba(255, 152, 0, 0.5)' },
                  { offset: 1, color: 'rgba(255, 152, 0, 0.1)' }
                ]
              }
            },
            lineStyle: {
              color: '#ff9800',
              width: 2
            },
            itemStyle: {
              color: '#ff9800'
            }
          }
        ]
      }
      this.hourlyChart.setOption(option)
    },
    initDailyChart() {
      if (this.dailyChart) {
        this.dailyChart.dispose()
      }

      const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      
      this.dailyChart = echarts.init(this.$refs.dailyChart)
      const option = {
        tooltip: {
          trigger: 'axis',
          backgroundColor: 'rgba(255,255,255,0.95)',
          borderColor: '#e0e0e0',
          textStyle: { color: '#333' }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: days
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '注册量',
            type: 'line',
            smooth: true,
            data: this.dailyRegistrations,
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: 'rgba(76, 175, 80, 0.5)' },
                  { offset: 1, color: 'rgba(76, 175, 80, 0.1)' }
                ]
              }
            },
            lineStyle: {
              color: '#4caf50',
              width: 2
            },
            itemStyle: {
              color: '#4caf50'
            }
          }
        ]
      }
      this.dailyChart.setOption(option)
    },
    goToAdmin() {
      this.$router.push('/admin')
    },
    goToLearn() {
      this.$router.push('/learn')
    },
    handleLogout() {
      this.userStore.logout()
      this.$router.push('/login')
    }
  },
  beforeUnmount() {
    if (this.hourlyChart) {
      this.hourlyChart.dispose()
    }
    if (this.dailyChart) {
      this.dailyChart.dispose()
    }
  }
}
</script>

<style scoped>
.monitor-container {
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

.admin-button, .learn-button, .logout-button {
  padding: 8px 20px;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-button {
  background-color: #ff9800;
  color: white;
}

.admin-button:hover {
  background-color: #f57c00;
  transform: translateY(-1px);
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
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.stats-section {
  display: flex;
  gap: 20px;
}

.stat-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 25px 35px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.stat-icon {
  font-size: 3em;
}

.stat-content {
  color: white;
}

.stat-value {
  font-size: 2.5em;
  font-weight: bold;
  line-height: 1;
}

.stat-label {
  font-size: 1em;
  opacity: 0.9;
  margin-top: 5px;
}

.charts-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
}

.chart-card {
  background-color: white;
  padding: 25px;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.chart-card h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  font-size: 1.2em;
}

.chart {
  width: 100%;
  height: 300px;
}

.ranking-section {
  background-color: white;
  padding: 25px;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.ranking-section h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  font-size: 1.2em;
}

.ranking-list {
  max-height: 500px;
  overflow-y: auto;
}

.empty-state {
  text-align: center;
  color: #999;
  padding: 40px;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.3s ease;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-item:hover {
  background-color: #fafafa;
}

.rank-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.1em;
  flex-shrink: 0;
}

.rank-1 {
  background: linear-gradient(135deg, #ffd700 0%, #ffb800 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(255, 215, 0, 0.4);
}

.rank-2 {
  background: linear-gradient(135deg, #c0c0c0 0%, #a8a8a8 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(192, 192, 192, 0.4);
}

.rank-3 {
  background: linear-gradient(135deg, #cd7f32 0%, #b87333 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(205, 127, 50, 0.4);
}

.rank-4, .rank-5, .rank-6, .rank-7, .rank-8, .rank-9, .rank-10 {
  background-color: #f5f5f5;
  color: #666;
}

.rank-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.word-info {
  display: flex;
  gap: 10px;
  align-items: center;
}

.word {
  font-size: 1.2em;
  font-weight: bold;
  color: #333;
}

.translation {
  color: #666;
}

.wrong-count {
  display: flex;
  align-items: center;
  gap: 5px;
}

.count-label {
  color: #999;
  font-size: 0.9em;
}

.count-value {
  font-size: 1.3em;
  font-weight: bold;
  color: #f44336;
}
</style>
