package middleware

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitState 熔断器状态
type CircuitState int

const (
	StateClosed   CircuitState = iota // 正常，请求通过
	StateOpen                         // 熔断，请求快速失败
	StateHalfOpen                     // 半开，允许部分请求试探
)

// CircuitBreaker 熔断器
// 用于保护数据库和Redis，防止级联故障
// 当错误率达到阈值时打开熔断，快速失败；
// 经过恢复时间后进入半开状态，允许试探请求
type CircuitBreaker struct {
	mu               sync.RWMutex
	state            CircuitState
	failureCount     int64
	successCount     int64
	lastFailureTime  time.Time
	halfOpenSuccess  int64

	name            string
	threshold       int64        // 触发熔断的错误数
	recoveryTimeout time.Duration // 从Open恢复到Half-Open的等待时间
	halfOpenMaxReq  int64        // 半开状态允许的最大请求数
	resetInterval   time.Duration // 关闭状态下重置计数的时间间隔
}

var (
	// DBCircuitBreaker 数据库熔断器
	DBCircuitBreaker = &CircuitBreaker{
		name:            "database",
		threshold:       10,
		recoveryTimeout: 30 * time.Second,
		halfOpenMaxReq:  3,
		resetInterval:   1 * time.Minute,
	}

	// RedisCircuitBreaker Redis熔断器
	RedisCircuitBreaker = &CircuitBreaker{
		name:            "redis",
		threshold:       10,
		recoveryTimeout: 15 * time.Second,
		halfOpenMaxReq:  3,
		resetInterval:   30 * time.Second,
	}
)

// Allow 判断请求是否允许通过
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	state := cb.state
	cb.mu.RUnlock()

	switch state {
	case StateClosed:
		return true
	case StateOpen:
		// 检查是否达到恢复时间
		if time.Since(cb.lastFailureTime) > cb.recoveryTimeout {
			cb.mu.Lock()
			// 双重检查
			if cb.state == StateOpen {
				cb.state = StateHalfOpen
				cb.halfOpenSuccess = 0
				log.Printf("[CircuitBreaker] %s 进入半开状态", cb.name)
			}
			cb.mu.Unlock()
			return true
		}
		return false
	case StateHalfOpen:
		cb.mu.RLock()
		count := atomic.LoadInt64(&cb.halfOpenSuccess)
		cb.mu.RUnlock()
		return count < cb.halfOpenMaxReq
	default:
		return true
	}
}

// Success 记录成功请求
func (cb *CircuitBreaker) Success() {
	atomic.AddInt64(&cb.successCount, 1)
	atomic.StoreInt64(&cb.failureCount, 0)

	cb.mu.Lock()
	if cb.state == StateHalfOpen {
		cb.halfOpenSuccess++
		if cb.halfOpenSuccess >= cb.halfOpenMaxReq {
			cb.state = StateClosed
			cb.halfOpenSuccess = 0
			log.Printf("[CircuitBreaker] %s 恢复正常（关闭状态）", cb.name)
		}
	}
	cb.mu.Unlock()
}

// Failure 记录失败请求
func (cb *CircuitBreaker) Failure() {
	atomic.AddInt64(&cb.failureCount, 1)
	cb.lastFailureTime = time.Now()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.state = StateOpen
		cb.halfOpenSuccess = 0
		log.Printf("[CircuitBreaker] %s 半开状态探测失败，回到断开状态", cb.name)
		return
	}

	if cb.state == StateClosed && atomic.LoadInt64(&cb.failureCount) >= cb.threshold {
		cb.state = StateOpen
		log.Printf("[CircuitBreaker] %s 熔断已打开（%d次连续失败）", cb.name, cb.threshold)
	}
}

// Reset 手动重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
	cb.halfOpenSuccess = 0
	log.Printf("[CircuitBreaker] %s 已手动重置", cb.name)
}

// State 获取当前状态
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// DBExecute 通过熔断器执行数据库操作
// fn: 实际数据库操作函数
// fallback: 降级函数（可选），在熔断时调用
func DBExecute(fn func() error, fallback func() error) error {
	if !DBCircuitBreaker.Allow() {
		// 熔断：执行降级逻辑
		if fallback != nil {
			return fallback()
		}
		return nil
	}

	err := fn()
	if err != nil {
		DBCircuitBreaker.Failure()
		if fallback != nil {
			return fallback()
		}
		return err
	}

	DBCircuitBreaker.Success()
	return nil
}

// RedisExecute 通过熔断器执行Redis操作
func RedisExecute(fn func() error, fallback func() error) error {
	if !RedisCircuitBreaker.Allow() {
		// 熔断：执行降级逻辑
		if fallback != nil {
			return fallback()
		}
		return nil
	}

	err := fn()
	if err != nil {
		RedisCircuitBreaker.Failure()
		// Redis熔断时不执行fallback，数据继续走数据库
		return err
	}

	RedisCircuitBreaker.Success()
	return nil
}
