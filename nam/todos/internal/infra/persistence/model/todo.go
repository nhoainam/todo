package model

import "time"

// Todo is the GORM model for the todos table.
// It maps to the database schema defined in database/todos.schema.
type Todo struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	ListID    int64      `gorm:"column:todo_list_id;not null"`
	Title     string     `gorm:"column:title;not null"`
	Content   string     `gorm:"column:content"`
	Status    int        `gorm:"column:status;not null;default:0"`
	Priority  int        `gorm:"column:priority;not null;default:0"`
	DueDate   *time.Time `gorm:"column:due_date"`
	CreatorID int64      `gorm:"column:creator_id;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Todo) TableName() string { return "todos" }
