package worker

type Request struct {
	Type TaskType
	Key  string
	Val  []byte
}

type Response struct {
	Type TaskType
	Val  []byte
}

type TaskType int

const (
	TaskTypeInValid = iota
	TaskTypeGet
	TaskTypeSet
	TaskTypeDelete
)

type Task struct {
	tType TaskType
	key   string
	val   []byte
	done  chan struct{}
	err   error
}

func toTask(req *Request) *Task {
	return &Task{
		tType: req.Type,
		key:   req.Key,
		val:   req.Val,
	}
}

func toResponse(t *Task) *Response {
	return &Response{
		Type: t.tType,
		Val:  t.val,
	}
}
