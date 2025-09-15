package task

import (
	"errors"
	"fmt"
)

// TaskStatusType はタスクステータスの型を定義
type TaskStatusType int

const (
	// StatusIncomplete は未完了状態
	StatusIncomplete TaskStatusType = iota
	// StatusCompleted は完了状態
	StatusCompleted
)

// TaskStatus はタスクの状態を表すValue Object
type TaskStatus struct {
	status TaskStatusType
}

// NewTaskStatus は新しいTaskStatusを作成する
func NewTaskStatus(isCompleted bool) TaskStatus {
	if isCompleted {
		return TaskStatus{status: StatusCompleted}
	}
	return TaskStatus{status: StatusIncomplete}
}

// NewIncompleteTaskStatus は未完了のTaskStatusを作成する
func NewIncompleteTaskStatus() TaskStatus {
	return TaskStatus{status: StatusIncomplete}
}

// NewCompletedTaskStatus は完了のTaskStatusを作成する
func NewCompletedTaskStatus() TaskStatus {
	return TaskStatus{status: StatusCompleted}
}

// NewTaskStatusFromString は文字列からTaskStatusを作成する
func NewTaskStatusFromString(status string) (TaskStatus, error) {
	switch status {
	case "incomplete", "false", "0":
		return NewIncompleteTaskStatus(), nil
	case "completed", "true", "1":
		return NewCompletedTaskStatus(), nil
	default:
		return TaskStatus{}, errors.New("無効なタスクステータスです")
	}
}

// IsCompleted は完了状態かを判定する
func (t TaskStatus) IsCompleted() bool {
	return t.status == StatusCompleted
}

// IsIncomplete は未完了状態かを判定する
func (t TaskStatus) IsIncomplete() bool {
	return t.status == StatusIncomplete
}

// ToBool はbool値を返す
func (t TaskStatus) ToBool() bool {
	return t.IsCompleted()
}

// ToInt は整数値を返す（0: 未完了, 1: 完了）
func (t TaskStatus) ToInt() int {
	if t.IsCompleted() {
		return 1
	}
	return 0
}

// Toggle は状態を切り替える
func (t TaskStatus) Toggle() TaskStatus {
	if t.IsCompleted() {
		return NewIncompleteTaskStatus()
	}
	return NewCompletedTaskStatus()
}

// Complete は完了状態にする
func (t TaskStatus) Complete() TaskStatus {
	return NewCompletedTaskStatus()
}

// Incomplete は未完了状態にする
func (t TaskStatus) Incomplete() TaskStatus {
	return NewIncompleteTaskStatus()
}

// GetStatusType はステータスタイプを返す
func (t TaskStatus) GetStatusType() TaskStatusType {
	return t.status
}

// ToDisplayString は表示用の文字列を返す
func (t TaskStatus) ToDisplayString() string {
	if t.IsCompleted() {
		return "完了"
	}
	return "未完了"
}

// ToDisplayIcon は表示用のアイコンを返す
func (t TaskStatus) ToDisplayIcon() string {
	if t.IsCompleted() {
		return "✅"
	}
	return "⭕"
}

// ToDisplayEmoji は表示用の絵文字を返す
func (t TaskStatus) ToDisplayEmoji() string {
	if t.IsCompleted() {
		return "🎉"
	}
	return "📝"
}

// ToAPIString はAPI用の文字列を返す
func (t TaskStatus) ToAPIString() string {
	if t.IsCompleted() {
		return "completed"
	}
	return "incomplete"
}

// GetProgressMessage は進捗メッセージを返す
func (t TaskStatus) GetProgressMessage() string {
	if t.IsCompleted() {
		return "このタスクは完了しています"
	}
	return "このタスクは作業中です"
}

// CanTransitionTo は指定した状態に遷移可能かを判定する
func (t TaskStatus) CanTransitionTo(newStatus TaskStatus) bool {
	// 任意の状態から任意の状態への遷移を許可
	return true
}

// GetValidTransitions は有効な遷移先を返す
func (t TaskStatus) GetValidTransitions() []TaskStatus {
	if t.IsCompleted() {
		return []TaskStatus{NewIncompleteTaskStatus()}
	}
	return []TaskStatus{NewCompletedTaskStatus()}
}

// String は文字列表現を返す
func (t TaskStatus) String() string {
	return t.ToDisplayString()
}

// GoString はGo言語の文字列表現を返す
func (t TaskStatus) GoString() string {
	return fmt.Sprintf("TaskStatus{status: %v}", t.status)
}

// Equals は別のTaskStatusと等価かを判定する
func (t TaskStatus) Equals(other TaskStatus) bool {
	return t.status == other.status
}

// Compare は他のTaskStatusとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
// 未完了 < 完了 の順序
func (t TaskStatus) Compare(other TaskStatus) int {
	if t.status < other.status {
		return -1
	}
	if t.status > other.status {
		return 1
	}
	return 0
}