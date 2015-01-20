# Task Manager

## SYNOPSIS

### Brief
This is a system:

* execute command using the http web ui

### What does it look like
![screen_shot] (https://cloud.githubusercontent.com/assets/3839691/5801983/8eda7b96-a02b-11e4-90f1-66bf28670bc2.png)
![screen_shot] (https://cloud.githubusercontent.com/assets/3839691/5801990/997296ce-a02b-11e4-8405-a33e389356ad.png)


### How to run

####For those who already have a `go` ENV
* just clone this project. Tweak the listen port if you like(default is `8085`), you'll find it very easy to specify another port in the `main func` at the buttom of `server.go`.
* change dir to the project root.
* `go run server.go`
* visit http://localhost:8085 and check it out.

####For those who dont hava a `go` ENV
* I have already build a `./server` binary for you, but you will not be able to tweak the listen port which is `8085` **yet**.
* This binary is build on `Linux AMD64` other architecture will not run this binary functionly. And this binary **may not** keep updated.

------

## DIRS

    +.
    ├── config
    │   └── task.json
    ├── download
    ├── README.md
    ├── resource
    │   ├── ksider.keystore
    │   └── sign.sh
    ├── server.go
    ├── static
    │   ├── css
    │   │   ├── bootstrap.min.css
    │   │   └── fileinput.css
    │   ├── fonts
    │   │   └── glyphicons-halflings-regular.woff
    │   ├── img
    │   │   └── loading.gif
    │   └── js
    │       ├── bootstrap.min.js
    │       ├── fileinput.js
    │       └── jquery-1.10.2.js
    ├── template
    │   ├── index.tmpl
    │   └── process.tmpl
    └── upload  

------

## CONTENT OF DIRS

   * **config**:
      - stored task.json which defined the task using json.
   * **download**:
      - place to store the files for user to download.
   * **resource**:
      - place for programmers to store the script of config stuff needed by task.
   * **upload**
      - user upload files will be stored here.
   * server.go:
      - the main server program.
   * static:
      - mostly some js/css pulgins.
   * template:
      - the templates for web server.

------

## CONFIG

  - `config/task.json` used for config task, only configed task could be saw and ran.
  - a sample of `config/task.json`:

```javascript
{
    "findJob":{
        // short description of the task
        "ShortDesc":"find files in /tmp/", 
        // more detail description of the task
        "Desc": "this is a sample of finding files under the /tmp/ dir", 
        // system command or script or executable program the task should run
        "Cmd": "find", 
        // args list, could using the placeHolder of `$UPLOADFILE$` or `$$Args$$` which will be speicifeid by the next chapter.
        "Args": ["/tmp/"], 
        // progress bar marks generated by program to Stdout.
        "progressMark":["google","unix","Job","One"], 
        // will this program need a uploaded files?
        "UploadFileNeeded":false, 
        //indicate whether user input args are needed.
        "InputArgsNeeded":false 
    },
    "ls":{
        "ShortDesc":"list files the use specified",
        "Desc": "this is a sample of list files in the $dir user specified",
        "Cmd": "ls",
        "Args": ["$$dir$$","$$additionOption$$"],
        "progressMark":["google","unix","Job","One"],
        "UploadFileNeeded":false,
        "InputArgsNeeded":true,
        "InputArgsDesc":[
            {"dir":"which dir to list"},
            {"additionOption":"some args like -al"}
        ] 
    },
    "cat":{
        "ShortDesc":"cat file",
        "Desc": "just cat the file uploaded",
        "Cmd": "./resource/myCat.sh",
        "Args": ["$UPLOADFILE$"],
        "progressMark":[""],
        "UploadFileNeeded":true
    },
    "appendDate":{
        "ShortDesc":"appendDate to txt file",
        "Desc":"this is a sample of modify file uploaded and prepare for user to download",
        "Cmd": "./resource/appendDate.sh",
        "Args": ["$UPLOADFILE$","$UPLOADFILE$.modified"],
        "progressMark":[""],
        "UploadFileNeeded":true
    }

}
```
------

## NOTICE

* User upload file will be stored in `upload` dir.
* Files to be download and generated by the **Task Program** should be in the `download` dir, if you wan't the user to see and download.
* File which user uploaded could be referred in `task.json` Using the PlaceHolder of `$UPLOADFILE$`(there is **one** dollar sign before and after) but be sure that the **Task Program** using the right relative path to accessing the uploaded file.
* If you want user to input args, you should use the PlaceHolder `$$args$$`(there are **two** dollar sign before and after) as args and define a map `InputArgsDesc` to descripe the args you want the user to input.
* Uploaded file is limited to just **one** per task yet.
* User input args could be multi.
* PlaceHolder `$UPLOADFILE$` could be a replacer in the arg list.
   for example, if a script named `modify_upload_file.sh` behave that when there is  a file `a.txt` and invoke the script

```bash
   ./modify_upload_file.sh a.txt a.txt.modified 

```
there will be a file named `a.txt.modified` generated under `download/`.

In that case `Args` in `task.json` could be writen like:
```
   ...

   "Cmd":"./resouce/modify_upload_file.sh",
   "UploadFileNeeded":true,
   "Args": ["$UPLOADFILE$","$UPLOADFILE$.modified"],
   
   ...
```
## PS
------
* Some samples are given in the `task.json` by default.
* Concurrently execution of the **same** task is not supported yet.
