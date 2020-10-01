package value

import "time"

// Metadata basic metadata for each entity
type Metadata struct {
	createTime time.Time
	updateTime time.Time
	active     bool
}

func NewMetadata() *Metadata {
	return &Metadata{
		createTime: time.Now(),
		updateTime: time.Now(),
		active:     true,
	}
}

func (m Metadata) GetCreateTime() time.Time {
	return m.createTime
}

func (m Metadata) GetUpdateTime() time.Time {
	return m.updateTime
}

func (m Metadata) GetState() bool {
	return m.active
}

func (m *Metadata) SetCreateTime(t time.Time) error {
	m.createTime = t
	return nil
}

func (m *Metadata) SetUpdateTime(t time.Time) error {
	m.updateTime = t
	return nil
}

func (m *Metadata) SetState(s bool) error {
	m.active = s
	return nil
}

func (m *Metadata) ToggleSate() {
	if m.active {
		m.active = false
		return
	}

	m.active = true
}
