package cmd

import (
	"context"
	"encoding/json"
	"io"
	"sort"
	"text/tabwriter"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

type druidIoWriter struct {
	out io.Writer
	w   tabwriter.Writer
}

// GVK for Druid CR
var GVK = schema.GroupVersionResource{
	Group:    "druid.apache.org",
	Version:  "v1alpha1",
	Resource: "druids",
}

//  patchValue specifies a patch operation.
type patchValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// constructor for patchValue{}
func NewPatchValue(op, path string, value interface{}) []byte {
	patchPayload := make([]patchValue, 1)

	patchPayload[0].Op = op
	patchPayload[0].Path = path
	patchPayload[0].Value = value

	bytes, _ := json.Marshal(patchPayload)
	return bytes
}

type dynamicInterface interface {
	writers
	readers
	patcher
}

// readers interface
type readers interface {
	listDruidCR(namespace string) (map[string][]string, error)
	getDruidNodeNames(namespaces, cr string) (map[string][]string, error)
}

// writers interface
type writers interface {
	writerDruidNodeSpecReplicas(nodeName, namespace, cr string, replica int64) (bool, error)
	writerDruidNodeImages(nodeName, namespace, cr, image string) (bool, error)
}

// patchers interface
type patcher interface {
	patcherDruidDeleteOrphanPvc(namespace, cr string, value bool) (bool, error)
	patcherDruidRollingDeploy(namespace, cr string, value bool) (bool, error)
}

// client struct holds the dynamic client
type client struct {
	dynamic.Interface
}

// initalize dynamicInterface
var di dynamicInterface = client{newClient()}

// getDruidNodeNames gets all the druid nodes in a namespace
// response map[namespace][]nameNames
func (c client) getDruidNodeNames(namespaces, cr string) (map[string][]string, error) {

	var err error

	crd, err := c.Resource(GVK).Namespace(namespaces).Get(context.TODO(), cr, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var response = make(map[string][]string, 0)

	nameLists, _, _ := unstructured.NestedMap(crd.Object, "spec", "nodes")
	for nameList := range nameLists {
		response[crd.GetNamespace()] = append(response[crd.GetNamespace()], nameList)
	}

	sort.Strings(response[crd.GetNamespace()])

	return response, nil

}

// listDruidCR lists all the druid CR in a namespace or all namespaces
// response map[namespace][]nameCR
func (c client) listDruidCR(namespace string) (map[string][]string, error) {

	var err error

	crd, err := c.Resource(GVK).Namespace(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var response = make(map[string][]string, 0)
	for _, d := range crd.Items {
		response[d.GetNamespace()] = append(response[d.GetNamespace()], d.GetName())
	}

	return response, nil

}

// writerNodeSpecReplicas writer nodespec replica
func (c client) writerDruidNodeSpecReplicas(nodeName, namespace, cr string, replica int64) (bool, error) {
	var err error
	dcr, err := c.Resource(GVK).Namespace(namespace).Get(context.TODO(), cr, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	if err := unstructured.SetNestedField(dcr.Object, int64(replica), "spec", "nodes", nodeName, "replicas"); err != nil {
		return false, err
	}

	_, err = c.Resource(GVK).Namespace(namespace).Update(context.TODO(), dcr, v1.UpdateOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// writerDruidNodeImages writers updates nodes images
func (c client) writerDruidNodeImages(nodeName, namespace, cr, image string) (bool, error) {
	var err error

	dcr, err := c.Resource(GVK).Namespace(namespace).Get(context.TODO(), cr, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	if err := unstructured.SetNestedField(dcr.Object, image, "spec", "nodes", nodeName, "image"); err != nil {
		return false, err
	}

	_, err = c.Resource(GVK).Namespace(namespace).Update(context.TODO(), dcr, v1.UpdateOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// patcherDruidDeleteOrphanPvc patches DeleteOrphanPvc flag
func (c client) patcherDruidDeleteOrphanPvc(namespace, cr string, value bool) (bool, error) {
	var err error

	patchBytes := NewPatchValue("replace", "/spec/deleteOrphanPvc", value)

	_, err = c.Resource(GVK).Namespace(namespace).Patch(context.TODO(), cr, types.JSONPatchType, patchBytes, v1.PatchOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// patcherDruidRollingDeploy patches rollingDeploy flag
func (c client) patcherDruidRollingDeploy(namespace, cr string, value bool) (bool, error) {
	var err error

	patchBytes := NewPatchValue("replace", "/spec/rollingDeploy", value)

	_, err = c.Resource(GVK).Namespace(namespace).Patch(context.TODO(), cr, types.JSONPatchType, patchBytes, v1.PatchOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}
