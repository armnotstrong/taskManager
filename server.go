package main

import (
    "html/template"
    "net/http"
    "fmt"
    "bufio"
    "strings"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "strconv"
    "time"
)

type Status struct{
    Ratio int
    Err   error
    Stdout []byte
    ErrString string
    StdoutString string
    TaskName string
    ShortDesc string
    JobId int
    StartTime string
    EndTime string
}

type Task struct {
    ShortDesc string
    Desc string
    Cmd string
    Args []string
    progressMark []string
    UploadFileNeeded bool
    uploadFileName string
    InputArgsNeeded bool
    InputArgs map[string]string
    InputArgsDesc []map[string]string
}

type Job struct {
    Cmd string
    JobId   int
    Args []string
    progressMark []string
    TaskName string
}

type TaskMap map[string] Task

var taskMap TaskMap

func getFileNames(dir string) []string{
    fileNames := make([]string,0)
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil
    }
    for _, f := range files {
        if f.IsDir() != true {
            fileNames = append(fileNames, f.Name())
        }
    }
    return fileNames
}

func getCurrentTime() string{
    t := time.Now().Local()
    return t.Format("2006-01-02 15:04:05")
}

func (status *Status) convertOutStr() error{
    status.StdoutString = string(status.Stdout)
    return nil
}

var statusArray []Status

func (job *Job) doJob() error {
    var args []string
    for _, arg := range job.Args{
        switch {
        case strings.Contains(arg,"$UPLOADFILE$"):
            args = append(args, strings.Replace(arg,"$UPLOADFILE$",
                                taskMap[job.TaskName].uploadFileName,-1))
        case strings.Contains(arg, "$$"):
            args = append(args, taskMap[job.TaskName].InputArgs[strings.Replace(arg,"$","",-1)])

        default:
           args = append(args, arg)
        }
    }
    cmd := exec.Command(job.Cmd, args...)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        statusArray[job.JobId].Err = err
        return err
    }
    if err := cmd.Start(); err != nil{
        statusArray[job.JobId].Err = err
        return err
    }

    in := bufio.NewScanner(stdout)

    for in.Scan(){
        s := in.Text()
        for key, value := range job.progressMark{
            if strings.Contains(s, value){
                statusArray[job.JobId].Ratio = key
            }
        }
        s = s + "\n"
        statusArray[job.JobId].Stdout = append(statusArray[job.JobId].Stdout, s...)
    }

    err = cmd.Wait()
    statusArray[job.JobId].EndTime = getCurrentTime()
    if err != nil {
        statusArray[job.JobId].Ratio = 100
        statusArray[job.JobId].Err = err
        statusArray[job.JobId].ErrString = err.Error()
    }else{
        statusArray[job.JobId].Ratio = 100

    }
    return err
}


func getprogressHandler(w http.ResponseWriter, r *http.Request){
    jobId := r.FormValue("jobId")
    jobIdInt,_ := strconv.Atoi(jobId)
    status := statusArray[jobIdInt]
    status.convertOutStr()
    jsonStr,_ := json.Marshal(status)
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonStr)
}

const MAX_MEMORY = 1 * 1024 * 1024

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    var taskName string
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

    for key, value := range r.MultipartForm.Value {
        switch key {
        case "taskName":
            taskName = value[0]
            fmt.Printf("%T", taskName)

        }
    }

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			path := fmt.Sprintf("upload/%s", fileHeader.Filename)
            task := taskMap[taskName]
            task.uploadFileName = fileHeader.Filename
            taskMap[taskName] = task
			buf, _ := ioutil.ReadAll(file)
			ioutil.WriteFile(path, buf, os.ModePerm)
		}
	}
    fmt.Fprintf(w,"%s","{}")
}

func parseTaskMap(config string) error{
    file, err := ioutil.ReadFile(config)
    if err != nil {
        log.Fatal(err)
    }
    err = json.Unmarshal(file, &taskMap)

    if err != nil {
        log.Fatal(err)
    }
    return err
}
type ModleView struct{
    TaskMap TaskMap
    StatusArray []Status
    DownLoadFiles []string
}
func delHandler(w http.ResponseWriter, r *http.Request){
    fileName := r.FormValue("filename")
    os.Remove("download/" + fileName)
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("{\"Result\":true}"))
}

func jobHandler(w http.ResponseWriter, r *http.Request){
    //jobId := r.FormValue("jobId")
    action := r.FormValue("action")
    //jobIdInt, err := strconv.Atoi(jobId)
    //if err != nil {
        //status := statusArray[jobIdInt]
        //_ = status
    //}else{
        //job := Job{Cmd: CmdStr, jobId:len(statusArray)}
        //_ = job
    //}
    switch action {
    case "addJob":
        taskName := r.FormValue("taskName")
        if(taskMap[taskName].InputArgsNeeded){
            argsMap := make(map[string]string)
            for _, value := range taskMap[taskName].InputArgsDesc{
                for subkey, _:= range value {
                    argsMap[subkey] = r.FormValue(subkey)
                }
            }
            task := taskMap[taskName]
            task.InputArgs = argsMap
            taskMap[taskName] = task
        }
        var job Job
        job.Cmd = taskMap[taskName].Cmd
        job.Args = taskMap[taskName].Args
        job.progressMark = taskMap[taskName].progressMark
        job.TaskName = taskName
        job.JobId = len(statusArray)
        statusArray = append(statusArray, Status{StartTime: getCurrentTime(),Ratio:0,Err:nil,Stdout:nil,JobId:job.JobId,TaskName:taskName,ShortDesc:taskMap[taskName].ShortDesc})
        go job.doJob()
        _ = statusArray
        //jsonStr, _ := json.Marshal(statusArray[len(statusArray) -1])
        //jsonStr, _ := json.Marshal(statusArray)
        jsonStr, _ := json.Marshal(job)
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonStr)
    }
}

func mainHandler(w http.ResponseWriter, r *http.Request){
    t, _ := template.ParseFiles("template/index.tmpl")
    modalView := ModleView{
        TaskMap: taskMap,
        StatusArray: statusArray,
        DownLoadFiles: getFileNames("download/"),
    }
    t.Execute(w,modalView)
}
const CONFIG = "./config/task.json"

func main(){
    err := parseTaskMap(CONFIG)
    if err != nil {
        log.Fatal(err)
        os.Exit(2)
    }
    http.HandleFunc("/", mainHandler)
    http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request){
        http.ServeFile(w,r,r.URL.Path[1:])
    })
    http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request){
        http.ServeFile(w,r,r.URL.Path[1:])
    })
    http.HandleFunc("/upload",  uploadHandler)
    http.HandleFunc("/api/del",delHandler)
    http.HandleFunc("/api/job", jobHandler)
    http.HandleFunc("/api/progress", getprogressHandler)
    http.ListenAndServe(":8085",nil)
}
