package shell

import(
	"golang.org/x/crypto/ssh"
	"syscall"
	"bytes"
	"io/ioutil"
	"os/exec"
	"time"
	"strconv"
	"strings"
)


func exec_shell(one_task Task,one_task_sta *TaskStatus){
	//one_task_sta := TaskStatus{one_task.Tid,time.Now(),time.Now(),101,-999,"",""}
	//saveTaskSta(&one_task_sta)

	//defer exitsave(&one_task_sta)
	//init status
	one_task_sta.Tid = one_task.Tid
	one_task_sta.Begin = time.Now()

	task_arr := strings.Split(one_task.Cmd," ")
	Cmd := exec.Cmd{Path:task_arr[0],
				Args:task_arr,
				Dir:one_task.Path}
	stdout, err := Cmd.StdoutPipe()
	if err != nil {
		one_task_sta.Status = 201
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}
	defer stdout.Close()
	stderr, err := Cmd.StderrPipe()
	if err != nil {
		one_task_sta.Status = 202
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}
	defer stderr.Close()

	if err := Cmd.Start(); err != nil {
		one_task_sta.Status = 203
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		one_task_sta.Status = 204
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}
	err_output, err := ioutil.ReadAll(stderr)
	if err != nil {
		one_task_sta.Status = 205
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}
	if err := Cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if Status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				one_task_sta.Ret = Status.ExitStatus()
			}else{
				one_task_sta.Status = 207
				one_task_sta.End = time.Now()
				return
			}
		} else {
			one_task_sta.Status = 0
		}
	}else{
		one_task_sta.Ret = 0
	}
	one_task_sta.Status = 200
	one_task_sta.Info = string(output)
	one_task_sta.Err = string(err_output)
	one_task_sta.End = time.Now()
	return
}


func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func ssh_exec(one_task Task,one_task_sta *TaskStatus){
	//one_task_sta := TaskStatus{one_task.Tid,time.Now(),time.Now(),101,-999,"",""}
	//defer exitsave(&one_task_sta)
	//saveTaskSta(&one_task_sta)
	one_task_sta.Tid = one_task.Tid
	one_task_sta.Begin = time.Now()

	config := &ssh.ClientConfig{
		User: one_task.User,
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/root/.ssh/id_rsa"),
		},
	}


	if one_task.Passwd != "" {
		config = &ssh.ClientConfig{
		User: one_task.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(one_task.Passwd),
		},
		}
	}

	if ! strings.Contains(one_task.Host,":"){
		one_task.Host = one_task.Host+":22"
	}
	//return
	client, err := ssh.Dial("tcp", one_task.Host, config)
	try_count := 1
	for err != nil && try_count < 3  {
		time.Sleep(1000 * time.Millisecond)
		client, err = ssh.Dial("tcp", one_task.Host, config)
		try_count ++
		//one_task_sta.Status = 211
		//one_task_sta.Err = err.Error()
		//one_task_sta.End = time.Now()
		//return
	}
	if err != nil {
		one_task_sta.Status = 211
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}
	defer client.Close()
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		one_task_sta.Status = 212
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}
	defer session.Close()
	// Once a Session is created, you can execute a single command on
	var b bytes.Buffer
	var c bytes.Buffer
	session.Stdout = &b
	session.Stderr = &c
	if err := session.Run("cd "+one_task.Path + " && " + one_task.Cmd + "; echo $?"); err != nil {
		one_task_sta.Status = 213
		one_task_sta.Err = err.Error()
		one_task_sta.End = time.Now()
		return
	}


	one_task_sta.Status = 210
	b_out := b.String()
	b_arr := strings.Split(b_out,"\n")
	b_len := len(b_arr)
	ret_code, err := strconv.Atoi(b_arr[b_len-2])
	if err != nil {
		ret_code = -888
	}
	one_task_sta.Ret = ret_code
	one_task_sta.Info = strings.Join(b_arr[0:b_len-2],"\n")
	one_task_sta.Err = c.String()
	one_task_sta.End = time.Now()
	return

}

