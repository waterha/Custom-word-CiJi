<template>
  <div class="stats-container">
    <div class="header">
      <button @click="goBack" class="back-button">
        <span class="back-arrow">←</span> 返回学习
      </button>
      <h1>学习状况</h1>
      <div class="header-placeholder"></div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <template v-if="!loading">
      <!-- 今日统计 -->
      <section class="stats-section">
        <h2 class="section-title">
          <span class="title-icon">📊</span> 今日统计
        </h2>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-value today-total">{{ today.total }}</div>
            <div class="stat-label">今日学习</div>
          </div>
          <div class="stat-card">
            <div class="stat-value today-known">{{ today.known }}</div>
            <div class="stat-label">认识</div>
          </div>
          <div class="stat-card">
            <div class="stat-value today-unknown">{{ today.unknown }}</div>
            <div class="stat-label">不认识</div>
          </div>
          <div class="stat-card">
            <div class="stat-value today-accuracy">{{ todayAccuracy }}%</div>
            <div class="stat-label">正确率</div>
          </div>
        </div>
      </section>

      <!-- 累计统计 -->
      <section class="stats-section">
        <h2 class="section-title">
          <span class="title-icon">📈</span> 累计统计
        </h2>
        <div class="overall-grid">
          <div class="overall-item">
            <div class="overall-value">{{ overall.total_words }}</div>
            <div class="overall-label">词库总数</div>
          </div>
          <div class="overall-item">
            <div class="overall-value">{{ overall.learned_words }}</div>
            <div class="overall-label">已学单词</div>
          </div>
          <div class="overall-item">
            <div class="overall-value">{{ overall.known_words }}</div>
            <div class="overall-label">已认识</div>
          </div>
          <div class="overall-item">
            <div class="overall-value">{{ overall.unknown_words }}</div>
            <div class="overall-label">不认识</div>
          </div>
          <div class="overall-item">
            <div class="overall-value">{{ overall.wrong_words }}</div>
            <div class="overall-label">常错词</div>
          </div>
          <div class="overall-item">
            <div class="overall-value">{{ overall.custom_words }}</div>
            <div class="overall-label">自定义词</div>
          </div>
        </div>

        <!-- 进度条 -->
        <div class="progress-block">
          <div class="progress-header">
            <span>学习进度</span>
            <span class="progress-text">{{ (overall.progress * 100).toFixed(1) }}%</span>
          </div>
          <div class="progress-bar-bg">
            <div class="progress-bar-fill" :style="{ width: (overall.progress * 100).toFixed(1) + '%' }"></div>
          </div>
        </div>
        <div class="progress-block">
          <div class="progress-header">
            <span>总体正确率</span>
            <span class="progress-text">{{ (overall.accuracy * 100).toFixed(1) }}%</span>
          </div>
          <div class="progress-bar-bg accuracy-bg">
            <div class="progress-bar-fill accuracy-fill" :style="{ width: (overall.accuracy * 100).toFixed(1) + '%' }"></div>
          </div>
        </div>
      </section>

      <!-- 学习连续天数 -->
      <section class="stats-section">
        <h2 class="section-title">
          <span class="title-icon">🔥</span> 学习连续天数
        </h2>
        <div class="streak-display">
          <span class="streak-number">{{ streakDays }}</span>
          <span class="streak-unit">天</span>
        </div>
      </section>

      <!-- 本周趋势 -->
      <section class="stats-section">
        <h2 class="section-title">
          <span class="title-icon">📅</span> 近7天学习趋势
        </h2>
        <div class="chart-container">
          <div class="bar-chart">
            <div v-for="(day, index) in weekly" :key="index" class="bar-wrapper">
              <div class="bar-value" v-if="day.count > 0">{{ day.count }}</div>
              <div
                class="bar"
                :style="{ height: barHeight(day.count) }"
                :class="{ 'bar-today': index === 6 }"
              ></div>
              <div class="bar-label">{{ formatDayLabel(day.date) }}</div>
            </div>
          </div>
        </div>
      </section>

      <!-- 常错单词 -->
      <section class="stats-section" v-if="topWrong.length > 0">
        <h2 class="section-title">
          <span class="title-icon">⚠️</span> 常错单词 Top 5
        </h2>
        <div class="wrong-list">
          <div v-for="(item, index) in topWrong" :key="index" class="wrong-item">
            <span class="wrong-rank">{{ index + 1 }}</span>
            <div class="wrong-info">
              <div class="wrong-word">{{ item.word }}</div>
              <div class="wrong-translation">{{ item.translation }}</div>
            </div>
            <span class="wrong-count">错 {{ item.wrong_count }} 次</span>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script>
import apiClient from '../api/axios'

export default {
  name: 'LearningStats',
  data() {
    return {
      loading: true,
      today: { total: 0, known: 0, unknown: 0 },
      overall: {
        total_words: 0,
        learned_words: 0,
        known_words: 0,
        unknown_words: 0,
        wrong_words: 0,
        progress: 0,
        accuracy: 0,
        custom_words: 0
      },
      weekly: [],
      topWrong: [],
      streakDays: 0,
      maxWeeklyCount: 1
    }
  },
  computed: {
    todayAccuracy() {
      if (this.today.total === 0) return 0
      return Math.round((this.today.known / this.today.total) * 100)
    }
  },
  mounted() {
    this.fetchStats()
  },
  methods: {
    async fetchStats() {
      try {
        const response = await apiClient.get('/learn/stats')
        const data = response.data
        this.today = data.today
        this.overall = {
          ...data.overall,
          custom_words: data.custom_words
        }
        this.weekly = data.weekly
        this.topWrong = data.top_wrong
        this.streakDays = data.streak_days

        // 计算最大学习数量用于柱状图高度
        if (data.weekly && data.weekly.length > 0) {
          this.maxWeeklyCount = Math.max(...data.weekly.map(d => d.count), 1)
        }
      } catch (error) {
        console.error('获取学习统计失败:', error)
      } finally {
        this.loading = false
      }
    },
    barHeight(count) {
      const pct = (count / this.maxWeeklyCount) * 100
      return Math.max(pct, 4) + '%'
    },
    formatDayLabel(dateStr) {
      const days = ['日', '一', '二', '三', '四', '五', '六']
      const d = new Date(dateStr)
      return '周' + days[d.getDay()]
    },
    goBack() {
      this.$router.push('/learn')
    }
  }
}
</script>

<style scoped>
.stats-container {
  padding: 20px;
  min-height: 100vh;
  background-color: var(--bg-color);
  max-width: 800px;
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
  font-size: 1.8em;
  margin: 0;
}

.header-placeholder {
  width: 100px;
}

.back-button {
  background: none;
  border: 2px solid var(--border-color);
  padding: 8px 18px;
  border-radius: 20px;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 4px;
}

.back-button:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
  transform: translateY(-1px);
}

.back-arrow {
  font-size: 18px;
}

.loading {
  text-align: center;
  padding: 60px;
  font-size: 18px;
  color: #999;
}

.stats-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.06);
  border: 1px solid var(--border-color);
}

.section-title {
  font-size: 18px;
  color: #333;
  margin: 0 0 20px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-icon {
  font-size: 22px;
}

/* 今日统计 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.stat-card {
  text-align: center;
  padding: 20px 12px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-value {
  font-size: 2em;
  font-weight: 700;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.today-total {
  color: var(--primary-color);
}

.today-known {
  color: #4caf50;
}

.today-unknown {
  color: #ff9800;
}

.today-accuracy {
  color: #2196f3;
}

/* 累计统计 */
.overall-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.overall-item {
  text-align: center;
  padding: 16px 8px;
  background: #f8f9fa;
  border-radius: 10px;
}

.overall-value {
  font-size: 1.5em;
  font-weight: 700;
  color: #333;
  margin-bottom: 4px;
}

.overall-label {
  font-size: 13px;
  color: #888;
}

.progress-block {
  margin-bottom: 16px;
}

.progress-block:last-child {
  margin-bottom: 0;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
  color: #555;
}

.progress-text {
  font-weight: 600;
  color: var(--primary-color);
}

.progress-bar-bg {
  height: 16px;
  background: #e9ecef;
  border-radius: 8px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  border-radius: 8px;
  background: linear-gradient(90deg, var(--primary-color), #ffb74d);
  transition: width 0.8s ease;
}

.accuracy-bg {
  background: #e3f2fd;
}

.accuracy-fill {
  background: linear-gradient(90deg, #42a5f5, #2196f3);
}

/* 连续天数 */
.streak-display {
  text-align: center;
  padding: 20px;
}

.streak-number {
  font-size: 4em;
  font-weight: 800;
  color: #ff6d00;
  line-height: 1;
}

.streak-unit {
  font-size: 1.5em;
  font-weight: 600;
  color: #666;
  margin-left: 8px;
}

/* 柱状图 */
.chart-container {
  padding: 10px 0;
}

.bar-chart {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  height: 180px;
  gap: 8px;
}

.bar-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  justify-content: flex-end;
}

.bar-value {
  font-size: 13px;
  font-weight: 600;
  color: var(--primary-color);
  margin-bottom: 6px;
}

.bar {
  width: 100%;
  max-width: 48px;
  border-radius: 6px 6px 0 0;
  background: linear-gradient(180deg, var(--primary-color), #ffcc80);
  transition: height 0.6s ease;
  min-height: 4px;
}

.bar-today {
  background: linear-gradient(180deg, #ff6d00, #ffab40);
  box-shadow: 0 2px 8px rgba(255, 109, 0, 0.3);
}

.bar-label {
  font-size: 12px;
  color: #888;
  margin-top: 8px;
  font-weight: 500;
}

/* 常错单词 */
.wrong-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.wrong-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 10px;
  transition: transform 0.2s;
}

.wrong-item:hover {
  transform: translateX(4px);
}

.wrong-rank {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #fff3e0;
  color: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.wrong-info {
  flex: 1;
}

.wrong-word {
  font-weight: 700;
  font-size: 16px;
  color: #333;
}

.wrong-translation {
  font-size: 14px;
  color: #e65100;
  margin-top: 2px;
}

.wrong-count {
  font-size: 13px;
  color: #e53935;
  font-weight: 600;
  white-space: nowrap;
  padding: 4px 10px;
  background: #ffebee;
  border-radius: 12px;
}

@media (max-width: 600px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .overall-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .header h1 {
    font-size: 1.4em;
  }
}
</style>
