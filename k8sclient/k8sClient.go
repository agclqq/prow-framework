package k8sclient

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	pubErr "github.com/agclqq/prow-framework/error"
	"github.com/agclqq/prow-framework/logger"
)

// InitK8sClientConn 初始化k8s ClientSet
func InitK8sClientConn(ctx context.Context, KubeConfPath string) (*kubernetes.Clientset, error) {
	kubeConf, err := clientcmd.BuildConfigFromFlags("", KubeConfPath)
	if err != nil {
		logger.ErrorfWithTrace(ctx, "初始化k8s ClientSet err %v", err)
		return nil, err
	}
	// 实例化ClientSet对象
	clientSet, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		logger.ErrorfWithTrace(ctx, "%v 初始化k8s ClientSet err %v", pubErr.ErrLabelK8S, err)
		return nil, err
	}
	logger.InfoWithTrace(ctx, "初始化k8s ClientSet Success ...")
	return clientSet, nil
}
