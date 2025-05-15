package tabmanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/chromedp/cdproto/target"
)

type TabSession struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

type TabManager struct {
	mu     sync.RWMutex
	nowTab target.ID
	tabs   map[target.ID]*TabSession
	order  []target.ID // 탭 추가 순서 보존
}

// NewTabManager initializes a TabManager with the given root tab.
func NewTabManager(tid target.ID, ctx context.Context, cancel context.CancelFunc) *TabManager {
	return &TabManager{
		tabs:   map[target.ID]*TabSession{tid: {Ctx: ctx, Cancel: cancel}},
		order:  []target.ID{tid},
		nowTab: tid,
	}
}

// Now returns the currently active tab session.
func (m *TabManager) Now() (*TabSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getTabNoLock(m.nowTab)
}

// AddNewTab adds a new tab session and makes it the active tab.
func (m *TabManager) AddNewTab(id target.ID, ctx context.Context, cancel context.CancelFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tabs[id] = &TabSession{Ctx: ctx, Cancel: cancel}
	m.order = append(m.order, id)
	m.nowTab = id
}

// SwitchTo changes the active tab to the specified ID.
func (m *TabManager) SwitchTo(id target.ID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tabs[id]; !ok {
		return fmt.Errorf("tab %s not found", id)
	}
	m.nowTab = id
	return nil
}

// GetTab retrieves the TabSession for the given ID.
func (m *TabManager) GetTab(id target.ID) (*TabSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getTabNoLock(id)
}

// CloseTab cancels and removes the tab session for the given ID.
// If the closed tab was the active one, it switches to the most recently added remaining tab.
func (m *TabManager) CloseTab(id target.ID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess, ok := m.tabs[id]
	if !ok {
		return
	}
	sess.Cancel()
	delete(m.tabs, id)

	// order에서 제거
	for i, tid := range m.order {
		if tid == id {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}

	// nowTab이 닫힌 경우 처리
	if m.nowTab == id {
		if len(m.order) > 0 {
			m.nowTab = m.order[len(m.order)-1] // 최근 추가된 탭으로 설정
		} else {
			m.nowTab = "" // 모두 닫힌 상태
		}
	}
}

// ListIDs returns a slice of all stored tab IDs in insertion order.
func (m *TabManager) ListIDs() []target.ID {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]target.ID(nil), m.order...) // 복사본 반환
}

// CloseAll cancels and removes all tab sessions.
func (m *TabManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, sess := range m.tabs {
		sess.Cancel()
	}
	m.tabs = make(map[target.ID]*TabSession)
	m.order = nil
	m.nowTab = ""
}

// getTabNoLock is an internal helper that retrieves a tab without acquiring locks.
func (m *TabManager) getTabNoLock(id target.ID) (*TabSession, bool) {
	sess, ok := m.tabs[id]
	return sess, ok
}
