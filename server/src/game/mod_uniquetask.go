package game

import (
	"fmt"
)

type TaskInfo struct {
	TaskId int
	State  int
}

type ModUniqueTask struct {
	MyTaskInfo map[int]*TaskInfo
	// Locker     *sync.RWMutex
}

func (mu *ModUniqueTask) IsTaskFinish(taskId int) bool {

	//if taskId == 10001 || taskId == 10002 {
	//	return true
	//}

	task, ok := mu.MyTaskInfo[taskId]
	fmt.Println(mu.MyTaskInfo)
	if !ok {
		return false
	}
	return task.State == TASK_STATE_FINISH
}
