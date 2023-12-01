package client

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const DEFAULT_NAMESPACE = "default"

type Client interface {
	Apply(ctx context.Context, obj []*unstructured.Unstructured) *Response
	Get(ctx context.Context, apiVersion, kind, namespace, name, labelSelector string) *Response
	Delete(ctx context.Context, obj *unstructured.Unstructured) *Response
	Exec(ctx context.Context, name string, namespace string, cmd []string) (string, error)
}

type KubeClient struct {
	Client    *kubernetes.Clientset
	DynClient dynamic.Interface
	Config    *rest.Config
}

type Response struct {
	Error   string
	Output  string
	Objects []map[string]interface{}
}
