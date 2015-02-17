// task manage by framestep
package frametask

import (
	"fmt"
	"sync"
)

type TaskI interface {
	GetID() int64
	SetFrame(f int64)
	GetFrame() int64
}

type TaskList []TaskI

type TaskQueue struct {
	frame2task map[int64]TaskList
	id2task    map[int64]TaskI
	mutex      sync.Mutex
}

func (ftq TaskQueue) String() string {
	return fmt.Sprintf("TaskQueue %v %v", len(ftq.frame2task), len(ftq.id2task))
}

func New() *TaskQueue {
	return &TaskQueue{
		frame2task: make(map[int64]TaskList),
		id2task:    make(map[int64]TaskI),
	}
}

func (ftq *TaskQueue) AddTaskToFrame(task TaskI, step int64) {
	ftq.mutex.Lock()
	defer ftq.mutex.Unlock()
	task.SetFrame(step)
	ftq.frame2task[step] = append(ftq.frame2task[step], task)
	ftq.id2task[task.GetID()] = task
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
		delete(ftq.id2task, v.GetID())
	}
	delete(ftq.frame2task, step)
}

func (ftq *TaskQueue) GetTaskByID(id int64) TaskI {
	return ftq.id2task[id]
}

func (ftq *TaskQueue) CancelTaskByID(id int64) TaskI {
	ftq.mutex.Lock()
	defer ftq.mutex.Unlock()
	task := ftq.GetTaskByID(id)
	if task == nil {
		return task
	}
	tasks := ftq.GetTaskListByFrame(task.GetFrame())
	for i, v := range tasks {
		if v == nil {
			continue
		}
		if v.GetID() == id {
			delete(ftq.id2task, v.GetID())
			tasks[i] = nil
			return task
		}
	}
	return task
}
