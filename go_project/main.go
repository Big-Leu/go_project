package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

func creatClient() *kubernetes.Clientset{
  home,_ := os.UserHomeDir()
  kubeConfigPath := filepath.Join(home,".kube/config")
 
  config, err:= clientcmd.BuildConfigFromFlags("",kubeConfigPath)

  if err != nil {
	panic(err.Error())
  }

  client:= kubernetes.NewForConfigOrDie(config)

  return client
}

func createPods(c *gin.Context){
	client := creatClient()
     
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

func getPods(c *gin.Context){
    client := creatClient()
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
  router.Run("localhost:8080")
}