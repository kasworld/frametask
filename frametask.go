// task manage by framestep
package frametask

import (
	"fmt"
	"sync"
)

type Task struct {
	ID    int64
	Obj   interface{}
	frame int64
}

func (t Task) String() string {
	return fmt.Sprintf("Task ID:%v, Obj:%v, Frame:%v",
		t.ID, t.Obj, t.frame)
}

func (t *Task) GetFrame() int64 {
	return t.frame
}

type TaskList []*Task

type TaskQueue struct {
	frame2task map[int64]TaskList
	id2task    map[int64]*Task
	mutex      sync.Mutex
}

func (ftq TaskQueue) String() string {
	return fmt.Sprintf("TaskQueue %v %v", len(ftq.frame2task), len(ftq.id2task))
}

func New() *TaskQueue {
	return &TaskQueue{
		frame2task: make(map[int64]TaskList),
		id2task:    make(map[int64]*Task),
	}
}

func (ftq *TaskQueue) AddTaskToFrame(task *Task, step int64) {
	ftq.mutex.Lock()
	defer ftq.mutex.Unlock()
	task.frame = step
	ftq.frame2task[step] = append(ftq.frame2task[step], task)
	ftq.id2task[task.ID] = task
}

func (ftq *TaskQueue) GetTaskListByFrame(step int64) TaskList {
	return ftq.frame2task[step]
}

func (ftq *TaskQueue) ClearFrame(step int64) {
	ftq.mutex.Lock()
	defer ftq.mutex.Unlock()
	tasks := ftq.GetTaskListByFrame(step)
	for _, v := range tasks {
		if v == nil {
			continue
		}
		delete(ftq.id2task, v.ID)
	}
	delete(ftq.frame2task, step)
}

func (ftq *TaskQueue) GetTaskByID(id int64) *Task {
	return ftq.id2task[id]
}

func (ftq *TaskQueue) CancelTaskByID(id int64) *Task {
	ftq.mutex.Lock()
	defer ftq.mutex.Unlock()
	task := ftq.GetTaskByID(id)
	if task == nil {
		return task
	}
	tasks := ftq.GetTaskListByFrame(task.frame)
	for i, v := range tasks {
		if v == nil {
			continue
		}
		if v.ID == id {
			delete(ftq.id2task, v.ID)
			tasks[i] = nil
			return task
		}
	}
	return task
}
