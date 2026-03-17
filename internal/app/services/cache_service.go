package services

import (
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// CacheService 缓存服务
type CacheService struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

// NewCacheService 创建缓存服务
func NewCacheService() *CacheService {
	return &CacheService{
		items: make(map[string]*CacheItem),
	}
}

// Set 设置缓存项
func (c *CacheService) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expirationTime int64
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration).UnixNano()
	}

	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: expirationTime,
	}
}

// Get 获取缓存项
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

// Delete 删除缓存项
func (c *CacheService) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear 清空缓存
func (c *CacheService) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}

// Count 获取缓存项数量
func (c *CacheService) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Cleanup 清理过期缓存
func (c *CacheService) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.items, key)
		}
	}
}

// StartCleanupRoutine 启动定期清理协程
func (c *CacheService) StartCleanupRoutine(interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.Cleanup()
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}



// ConnectionPool 连接池
type ConnectionPool struct {
	connections map[string]interface{}
	maxSize     int
	mu          sync.RWMutex
}

// NewConnectionPool 创建连接池
func NewConnectionPool(maxSize int) *ConnectionPool {
	return &ConnectionPool{
		connections: make(map[string]interface{}),
		maxSize:     maxSize,
	}
}

// Add 添加连接
func (p *ConnectionPool) Add(id string, conn interface{}) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.connections) >= p.maxSize {
		return false
	}

	p.connections[id] = conn
	return true
}

// Get 获取连接
func (p *ConnectionPool) Get(id string) (interface{}, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	conn, exists := p.connections[id]
	return conn, exists
}

// Remove 移除连接
func (p *ConnectionPool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.connections, id)
}

// Count 获取连接数
func (p *ConnectionPool) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.connections)
}

// Clear 清空连接池
func (p *ConnectionPool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connections = make(map[string]interface{})
}

// ForEach 遍历连接池
func (p *ConnectionPool) ForEach(fn func(id string, conn interface{})) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for id, conn := range p.connections {
		fn(id, conn)
	}
}

// PerformanceMonitor 性能监控
type PerformanceMonitor struct {
	stats map[string]*StatItem
	mu    sync.RWMutex
}

// StatItem 统计项
type StatItem struct {
	Count     int64
	Min       int64
	Max       int64
	Total     int64
	LastValue int64
}

// NewPerformanceMonitor 创建性能监控
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		stats: make(map[string]*StatItem),
	}
}

// Record 记录性能数据
func (p *PerformanceMonitor) Record(key string, value int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, exists := p.stats[key]
	if !exists {
		item = &StatItem{}
		p.stats[key] = item
	}

	item.Count++
	item.Total += value
	item.LastValue = value

	if value < item.Min || item.Min == 0 {
		item.Min = value
	}
	if value > item.Max {
		item.Max = value
	}
}

// GetStats 获取统计数据
func (p *PerformanceMonitor) GetStats(key string) (*StatItem, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	item, exists := p.stats[key]
	if !exists {
		return nil, false
	}

	// 返回副本
	return &StatItem{
		Count:     item.Count,
		Min:       item.Min,
		Max:       item.Max,
		Total:     item.Total,
		LastValue: item.LastValue,
	}, true
}

// GetAverage 获取平均值
func (p *PerformanceMonitor) GetAverage(key string) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	item, exists := p.stats[key]
	if !exists || item.Count == 0 {
		return 0
	}

	return float64(item.Total) / float64(item.Count)
}

// Reset 重置统计数据
func (p *PerformanceMonitor) Reset(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.stats, key)
}

// GetAllStats 获取所有统计数据
func (p *PerformanceMonitor) GetAllStats() map[string]*StatItem {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]*StatItem)
	for key, item := range p.stats {
		result[key] = &StatItem{
			Count:     item.Count,
			Min:       item.Min,
			Max:       item.Max,
			Total:     item.Total,
			LastValue: item.LastValue,
		}
	}
	return result
}
