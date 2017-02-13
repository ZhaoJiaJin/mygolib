package shell

import(
	"time"
)

type Task struct{
	Host string
	User string
	Passwd string
	Path string
	Cmd string
	Tid int
}

type TaskStatus struct{
	Tid int
	Begin time.Time
	End time.Time
	Status int //101:正在运行, 200:正常结束 2xx:异常结束 301:结束并且已经检查
	Ret int //退出码
	Info string
	Err string
}

//传入一个task  和 task 和 taskstatus

func (task Task) Run_task(status *TaskStatus){
	if task.Host == "0.0.0.0" || task.Host == "127.0.0.1"{
		exec_shell(task,status)
	} else {
		ssh_exec(task,status)
	}
}


