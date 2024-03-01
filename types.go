package localdb

import "time"

type FindResponse[T Model] struct {
	StopOnFirst bool
	Query       bool
}

type AddResponse[T Model] struct {
	Query bool
	Value T
}

type UpdateResponse[T Model] struct {
	StopOnFirst bool
	Query       bool
	Value       T
}

type DeleteResponse[T Model] FindResponse[T]

type DBFindCallback[T Model] func(model T) *FindResponse[T]
type DBAddCallback[T Model] func(model T) *AddResponse[T]
type DBDeleteCallback[T Model] func(model T) *DeleteResponse[T]
type DBUpdateCallback[T Model] func(model T) *UpdateResponse[T]

type Model interface {
	GetID() string
	SetID(id string)

	GetCreatedAt() time.Time
	SetCreatedAt(c time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(u time.Time)
}

type Base struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

func (b *Base) GetID() string {
	return b.ID
}
func (b *Base) GetCreatedAt() time.Time {
	return b.CreatedAt
}
func (b *Base) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}
func (b *Base) SetID(id string) {
	b.ID = id
}
func (b *Base) SetCreatedAt(c time.Time) {
	b.CreatedAt = c
}
func (b *Base) SetUpdatedAt(u time.Time) {
	b.UpdatedAt = u
}
func NewBase() *Base {
	return &Base{}
}
