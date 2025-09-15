package task

import (
	"github.com/google/uuid"
)

// TaskID はタスクの識別子を表すValue Object
type TaskID struct {
	value uuid.UUID
}

// NewTaskID は新しいTaskIDを生成する
func NewTaskID() TaskID {
	return TaskID{value: uuid.New()}
}

// NewTaskIDFromUUID は既存のUUIDからTaskIDを作成する
func NewTaskIDFromUUID(id uuid.UUID) TaskID {
	return TaskID{value: id}
}

// NewTaskIDFromString は文字列からTaskIDを作成する
func NewTaskIDFromString(id string) (TaskID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return TaskID{}, err
	}
	return TaskID{value: parsedUUID}, nil
}

// Value はUUID値を返す
func (t TaskID) Value() uuid.UUID {
	return t.value
}

// String は文字列表現を返す
func (t TaskID) String() string {
	return t.value.String()
}

// IsNil は空のIDかを判定する
func (t TaskID) IsNil() bool {
	return t.value == uuid.Nil
}

// Equals は別のTaskIDと等価かを判定する
func (t TaskID) Equals(other TaskID) bool {
	return t.value == other.value
}