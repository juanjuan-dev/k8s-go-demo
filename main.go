package main

import (
	"flag"
	"fmt"
	"k8s-demo/controller/k8sClient"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	var config *rest.Config
	var kubeconfig *string

	if home := homeDir(); home != "" {
		fmt.Println("(可选)kubeconfig 文件的绝对路径")
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选)kubeconfig 文件的绝对路径")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件的绝对路径")
		fmt.Println("kubeconfig 文件的绝对路径")
	}
	flag.Parse()

	fmt.Println("kubeconfig = {} ", kubeconfig)
	// 首先使用 inCluster 模式(需要区配置对应的RBAC 权限,默认的sa是default-->是没有获取deployment的List权限)
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用Kubeonfig文件配置集群配置Config对象
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}
	// 已经获得了rest.Config对象
	// 创建Clientset对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("start")
	flag.Parse()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		// test
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/namespaces", func(c *gin.Context) {
		namespaces, err := k8sClient.GetAllNamespaces(clientset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, namespaces)
	})

	router.GET("/pods", func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace parameter is required"})
			return
		}

		pods, err := k8sClient.GetAllPodsInNamespace(clientset, namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pods)
	})

	router.Run(":30300")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") //windows
}
