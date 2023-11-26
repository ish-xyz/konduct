package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

func NewKubeClient(clientset *kubernetes.Clientset, dynclient dynamic.Interface, config *rest.Config) *KubeClient {
	return &KubeClient{
		Client:    clientset,
		DynClient: dynclient,
		Config:    config,
	}
}

func getRESTMapper(restconfig *rest.Config) (meta.RESTMapper, error) {
	// Init discovery client and mapper
	dc, err := discovery.NewDiscoveryClientForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// Get GVR
	groupResources, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	return mapper, nil
}

func getNamespace(obj *unstructured.Unstructured, mapping *meta.RESTMapping) string {
	// Default to "default" namespace if not specified
	namespace := obj.GetNamespace()
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace && namespace == "" {
		namespace = DEFAULT_NAMESPACE
	}
	return namespace
}

// Implementation of Kubernetes server side apply
func (k *KubeClient) Apply(ctx context.Context, obj *unstructured.Unstructured) error {

	var dr dynamic.ResourceInterface

	mapper, err := getRESTMapper(k.Config)
	if err != nil {
		return err
	}

	mapping, err := mapper.RESTMapping(schema.ParseGroupKind(obj.GroupVersionKind().GroupKind().String()))
	if err != nil {
		return err
	}

	namespace := getNamespace(obj, mapping)

	dr = k.DynClient.Resource(mapping.Resource)
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = k.DynClient.Resource(mapping.Resource).Namespace(namespace)
	}

	// Check if namespace is empty and if resource is namespaced or not
	data, _ := json.Marshal(obj)
	_, err = dr.Patch(
		ctx,
		obj.GetName(),
		types.ApplyPatchType,
		data,
		metav1.PatchOptions{
			FieldManager: "go-kubetest",
		},
	)

	return err

}

// Kubernetes Get/List operation with dynamic client for custom and core types.
func (k *KubeClient) Get(ctx context.Context, apiVersion, kind, namespace, name, labelSelector string) (*unstructured.UnstructuredList, error) {

	var dr dynamic.ResourceInterface

	// Use empty group name if core v1
	group := strings.Split(apiVersion, "/")[0]
	if group == apiVersion {
		group = ""
	}

	mapper, err := getRESTMapper(k.Config)
	if err != nil {
		return nil, err
	}

	mapping, err := mapper.RESTMapping(schema.GroupKind{Kind: kind, Group: group})
	if err != nil {
		return nil, err
	}

	if mapping.Scope.Name() == meta.RESTScopeNameNamespace && namespace == "" {
		namespace = DEFAULT_NAMESPACE
	}

	// Init dynamic client
	dr = k.DynClient.Resource(mapping.Resource)
	if namespace != "" {
		dr = k.DynClient.Resource(mapping.Resource).Namespace(namespace)
	}

	var outerr error
	objects := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{},
	}
	if name != "" {
		item, err := dr.Get(ctx, name, metav1.GetOptions{})
		outerr = err
		if item != nil {
			objects.Items = append(objects.Items, *item)
		}
	} else {
		list, err := dr.List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		outerr = err
		if list != nil {
			objects = list
		}
	}

	return objects, outerr
}

// Kubernetes delete operation with dynamic client
func (k *KubeClient) Delete(ctx context.Context, obj *unstructured.Unstructured) error {

	var dr dynamic.ResourceInterface

	mapper, err := getRESTMapper(k.Config)
	if err != nil {
		return err
	}

	mapping, err := mapper.RESTMapping(schema.ParseGroupKind(obj.GroupVersionKind().GroupKind().String()))
	if err != nil {
		return err
	}

	namespace := getNamespace(obj, mapping)

	dr = k.DynClient.Resource(mapping.Resource)
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = k.DynClient.Resource(mapping.Resource).Namespace(namespace)
	}

	// Exec rest request to API
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err = dr.Delete(ctx, obj.GetName(), deleteOptions)

	return err
}

// Pod exec method to run commands and fetch outputs (stdout/stderr)
func (k *KubeClient) Exec(ctx context.Context, name string, namespace string, cmd []string) (string, error) {

	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	stdout := bufio.NewWriter(&stdoutBuff)
	stderr := bufio.NewWriter(&stderrBuff)

	req := k.Client.CoreV1().RESTClient().Post().Resource("pods").Name(name).Namespace(namespace).SubResource("exec")
	opts := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(
		opts,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(k.Config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return "", err
	}

	fmt.Println(stdoutBuff.String())
	fmt.Println(stderrBuff.String())

	return "", nil
}
