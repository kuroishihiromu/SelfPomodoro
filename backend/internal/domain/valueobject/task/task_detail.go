package task

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// TaskDetail はタスクの詳細を表すValue Object
type TaskDetail struct {
	value string
}

const (
	MinDetailLength = 1
	MaxDetailLength = 500
)

// NewTaskDetail は新しいTaskDetailを作成する
func NewTaskDetail(detail string) (TaskDetail, error) {
	// 前後の空白を除去
	trimmed := strings.TrimSpace(detail)
	
	// 空文字列チェック
	if trimmed == "" {
		return TaskDetail{}, errors.New("タスク詳細は空にできません")
	}
	
	// 文字数チェック（UTF-8文字数）
	charCount := utf8.RuneCountInString(trimmed)
	if charCount < MinDetailLength {
		return TaskDetail{}, errors.New("タスク詳細は1文字以上入力してください")
	}
	if charCount > MaxDetailLength {
		return TaskDetail{}, errors.New("タスク詳細は500文字以内で入力してください")
	}
	
	// 改行文字の正規化（CRLFをLFに統一）
	normalized := strings.ReplaceAll(trimmed, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	
	// 連続する改行を制限（最大2個まで）
	for strings.Contains(normalized, "\n\n\n") {
		normalized = strings.ReplaceAll(normalized, "\n\n\n", "\n\n")
	}
	
	return TaskDetail{value: normalized}, nil
}

// Value は詳細文字列を返す
func (t TaskDetail) Value() string {
	return t.value
}

// Length は文字数を返す
func (t TaskDetail) Length() int {
	return utf8.RuneCountInString(t.value)
}

// IsEmpty は空かを判定する
func (t TaskDetail) IsEmpty() bool {
	return strings.TrimSpace(t.value) == ""
}

// IsShort は短い詳細（10文字以下）かを判定する
func (t TaskDetail) IsShort() bool {
	return t.Length() <= 10
}

// IsLong は長い詳細（100文字以上）かを判定する
func (t TaskDetail) IsLong() bool {
	return t.Length() >= 100
}

// HasMultipleLines は複数行かを判定する
func (t TaskDetail) HasMultipleLines() bool {
	return strings.Contains(t.value, "\n")
}

// GetLineCount は行数を返す
func (t TaskDetail) GetLineCount() int {
	if t.value == "" {
		return 0
	}
	return strings.Count(t.value, "\n") + 1
}

// GetPreview はプレビュー用の短縮文字列を返す
func (t TaskDetail) GetPreview(maxLength int) string {
	if maxLength <= 0 {
		maxLength = 50
	}
	
	// 改行があれば最初の行のみ
	firstLine := strings.Split(t.value, "\n")[0]
	
	// 指定文字数で切り詰め
	if utf8.RuneCountInString(firstLine) <= maxLength {
		return firstLine
	}
	
	runes := []rune(firstLine)
	if len(runes) > maxLength-3 {
		return string(runes[:maxLength-3]) + "..."
	}
	return string(runes[:maxLength])
}

// GetFirstLine は最初の行を返す
func (t TaskDetail) GetFirstLine() string {
	lines := strings.Split(t.value, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
}

// ContainsKeywords はキーワードが含まれているかを判定する（大小文字を区別しない）
func (t TaskDetail) ContainsKeywords(keywords ...string) bool {
	lowerValue := strings.ToLower(t.value)
	for _, keyword := range keywords {
		if strings.Contains(lowerValue, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// Update は詳細を更新する
func (t TaskDetail) Update(newDetail string) (TaskDetail, error) {
	return NewTaskDetail(newDetail)
}

// Append は既存の詳細に文字列を追加する
func (t TaskDetail) Append(additional string) (TaskDetail, error) {
	if additional == "" {
		return t, nil
	}
	
	newValue := t.value
	if !strings.HasSuffix(newValue, "\n") && additional != "" {
		newValue += "\n"
	}
	newValue += additional
	
	return NewTaskDetail(newValue)
}

// Prepend は既存の詳細の前に文字列を追加する
func (t TaskDetail) Prepend(prefix string) (TaskDetail, error) {
	if prefix == "" {
		return t, nil
	}
	
	newValue := prefix
	if !strings.HasSuffix(prefix, "\n") && t.value != "" {
		newValue += "\n"
	}
	newValue += t.value
	
	return NewTaskDetail(newValue)
}

// String は文字列表現を返す
func (t TaskDetail) String() string {
	return t.value
}

// Equals は別のTaskDetailと等価かを判定する
func (t TaskDetail) Equals(other TaskDetail) bool {
	return t.value == other.value
}