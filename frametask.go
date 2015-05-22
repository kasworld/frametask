// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// task manage by framestep
package frametask

import (
	"fmt"
	"sync"

	"github.com/kasworld/idgen"
)

type TaskI interface {
	GetID() idgen.IDInt
	SetFrame(f int64)
	GetFrame() int64
}

type TaskList []TaskI

type TaskQueue struct {
	frame2task map[int64]TaskList
	id2task    map[idgen.IDInt]TaskI
	mutex      sync.Mutex
}

func (ftq TaskQueue) String() string {
	return fmt.Sprintf("TaskQueue %v %v", len(ftq.frame2task), len(ftq.id2task))
}

func New() *TaskQueue {
	return &TaskQueue{
		frame2task: make(map[int64]TaskList),
		id2task:    make(map[idgen.IDInt]TaskI),
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

func (ftq *TaskQueue) GetTaskByID(id idgen.IDInt) TaskI {
	return ftq.id2task[id]
}

func (ftq *TaskQueue) CancelTaskByID(id idgen.IDInt) TaskI {
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
