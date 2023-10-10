package repository

import "time"

type FindResponse[T Model] struct {
	StopOnFirst bool
	Query       bool
}

type AddResponse[T Model] struct {
	Query bool
	Value T
}

type DeleteResponse[T Model] FindResponse[T]
type UpdateResponse[T Model] AddResponse[T]

type RepositoryFindCallback[T Model] func(model T) *FindResponse[T]
type RepositoryAddCallback[T Model] func(model T) *AddResponse[T]
type RepositoryDeleteCallback[T Model] func(model T) *DeleteResponse[T]
type RepositoryUpdateCallback[T Model] func(model T) *UpdateResponse[T]

type Model interface {
	GetID() string
	SetID(id string)

	GetCreatedAt() time.Time
	SetCreatedAt(c time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(u time.Time)
}

type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
