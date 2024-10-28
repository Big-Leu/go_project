package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Command struct {
    PodName string `json:"podname"`
    Namespace string `json:"namespace"`
    ContainerName string `json:"containername"`
    EndPoint string `json:"endpoint"`
    FunctionName string `json:"functionname"`
    Route string `json:"route"`
    FunctionBody string `json:"functionbody"`
}

type PodInfo struct {
    Name       string        `json:"name"`
    Containers []ContainerInfo `json:"containers"`
    Status     string        `json:"status"`
}

type ContainerInfo struct {
    Name  string `json:"name"`
    Image string `json:"image"`
}


type PodCreateSchema struct {
    PodName  string `json:"podname"`
    NameSpace string `json:"namespace"`
    ContainerName   string    `json:"containername"`
	Image string `json:"image"`
}


type LogStreamer struct{
    b bytes.Buffer
}

func (l *LogStreamer) String() string {
    return l.b.String()
}

func (l *LogStreamer) Write(p []byte) (n int, err error) {
    a := strings.TrimSpace(string(p))
    l.b.WriteString(a)
    return len(p), nil
}


func creatClient() (*kubernetes.Clientset,*rest.Config){
  home,_ := os.UserHomeDir()
  kubeConfigPath := filepath.Join(home,".kube/config")
 
  config, err:= clientcmd.BuildConfigFromFlags("",kubeConfigPath)

  if err != nil {
	panic(err.Error())
  }

  client:= kubernetes.NewForConfigOrDie(config)

  return client,config
}

func createPods(c *gin.Context){
	client,_ := creatClient()
     
     var pod PodCreateSchema
    if err := c.ShouldBindJSON(&pod); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	fmt.Print(pod)

	podDefintion := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pod.PodName,
			Namespace: pod.NameSpace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
				Name: pod.ContainerName,
				Image: pod.Image,
				},
			},
		},
	}

	newPods,err := client.CoreV1().Pods(pod.NameSpace).Create(context.Background(),podDefintion, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}
	c.IndentedJSON(http.StatusAccepted, newPods)

}



func execCommandInPod(c *gin.Context) {
    client,config := creatClient()
    var pod Command
    if err := c.ShouldBindJSON(&pod); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    pythonCommand := fmt.Sprintf(`python script_to.py --endpoint_type %s --function_name %s --route %s --function_file '%s'`, pod.EndPoint,
    pod.FunctionName, pod.Route, pod.FunctionBody)

    fmt.Print(pythonCommand)
    command := []string{"/bin/sh", "-c",pythonCommand}

    // Create the request
    req := client.CoreV1().RESTClient().
        Post().
        Resource("pods").
        Name(pod.PodName).
        Namespace(pod.Namespace).
        SubResource("exec").
        Param("container",pod.ContainerName).
        Param("stdout", "true").
        Param("stderr", "true").
        Param("tty", "false")

    // Add command parameters correctly
    for _, cmd := range command {
        req.Param("command", cmd)
    }
    l := &LogStreamer{}
    Executor,err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())

    Executor.StreamWithContext(context.Background(),remotecommand.StreamOptions{
        Stdin:  os.Stdin,
        Stdout: l,
        Stderr: nil,
        Tty:    true,
    })

    if err != nil {
        c.String(http.StatusInternalServerError, "Error executing command: %s", err.Error())
        return
    }

    c.String(http.StatusOK, "Command executed successfully")
}




func getPods(c *gin.Context){
    client,_:= creatClient()
	pods,err :=client.CoreV1().Pods("default").List(context.Background(),metav1.ListOptions{})

	if err != nil {
	  panic(err.Error())
	}
  
    var podDetails []PodInfo
    for _, pod := range pods.Items {
        var containers []ContainerInfo

        for _, container := range pod.Spec.Containers {
            containers = append(containers, ContainerInfo{
                Name:  container.Name,
                Image: container.Image,
            })
        }

        podDetails = append(podDetails, PodInfo{
            Name:       pod.Name,
            Containers: containers,
            Status:     string(pod.Status.Phase),
        })
    }

    if err!= nil {
        panic(err.Error())
    }


	c.IndentedJSON(http.StatusAccepted,podDetails)
}

func main() {
  router := gin.Default()
  router.GET("/api/v1/getpods",getPods)
  router.POST("/api/v1/createpods",createPods)
  router.POST("/api/v1/add",execCommandInPod)
  router.Run("localhost:8080")
}