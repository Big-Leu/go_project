package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)
func main(){
home,_ := os.UserHomeDir()
 kubeConfigPath := filepath.Join(home,".kube/config")
 
  config, err:= clientcmd.BuildConfigFromFlags("",kubeConfigPath)

  if err != nil {
	panic(err.Error())
  }

  client:= kubernetes.NewForConfigOrDie(config)

  pods,err :=client.CoreV1().Pods("default").List(context.Background(),metav1.ListOptions{})

  if err != nil {
	panic(err.Error())
  }

  for i, pod := range pods.Items{
      fmt.Printf("Name of %dth pods:%s\n",i,pod.Name)
  }
}