package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
    "os"
    "strconv"
    "fmt"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}


// worker请求master分配任务的reply标志，worker以此来判断需要进行的操作
type ReqTaskReplyFlag int

// 枚举worker请求任务的reply标志
const (
	TaskGetted    ReqTaskReplyFlag = iota // master给worker分配任务成功
	WaitPlz                               // 当前阶段暂时没有尚未分配的任务，worker的本次请求获取不到任务
	FinishAndExit                         // mapreduce工作已全部完成，worker准备退出
)


// RPC请求master分配任务时的传入参数——事实上什么都不用传入，因为只是worker获取一个任务。判断分配什么类型的任务由Master根据执行阶段决定
type TaskArgs struct {
}

// RPC请求master分配任务时的返回参数——即返回一个Master分配的任务
type TaskReply struct {
	Answer ReqTaskReplyFlag
	Task   Task
}


// Add your RPC definitions here.
// 向master请求RPC,获取任务
func RequestTask() TaskReply {
	args := TaskArgs{}
	reply := TaskReply{}

	// send the RPC request, wait for the reply.
	// 请求调用Master的AssignTask方法
	ok := call("Master.AssignTask", &args, &reply)

	if !ok { // 请求失败
		fmt.Println("Call failed!")
	}

	return reply
}

// worker完成工作后调用RPC通知master更新任务状态的传入参数——传入要更新的taskId即可
type FinArgs struct {
	TaskId int
}

// worker完成工作后调用RPC通知master更新任务状态的返回参数——事实上不需要传出什么，master自己修改TaskMap中的对应任务状态即可
type FinReply struct {
}


// worker完成任务后调用RPC告知master并请求修改任务状态
func FinishTaskAndReport(id int) {
	args := FinArgs{TaskId: id}
	reply := FinReply{}

	// 请求调用Master的UpdateTaskState方法
	ok := call("Master.UpdateTaskState", &args, &reply)
	if !ok { //请求失败
		fmt.Println("Call failed!")
	}
}




// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
