package cmd

import (
	"context"
	"encoding/json"
	"sort"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

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

// dynamicInterface holds writers,reader and patcher interfaces
type dynamicInterface interface {
	writers
	readers
	patcher
}

// readers interface
type readers interface {
	listDruidCR(namespaces string) ([]string, error)
	getDruidNodeNames(namespaces, CR string) ([]string, error)
}

// writers interface
type writers interface {
	writerDruidNodeSpecReplicas(nodeName, namespace, CR string, replica int64) (bool, error)
	writerDruidNodeImages(nodeName, namespace, CR, image string) (bool, error)
}

// patchers interface
type patcher interface {
	patcherDruidDeleteOrphanPvc(namespace, CR string, value bool) (bool, error)
	patcherDruidRollingDeploy(namespace, CR string, value bool) (bool, error)
}

// client struct holds the dynamic client
type client struct {
	dynamic.Interface
}

// initalize dynamicInterface
var di dynamicInterface = client{newClient()}

// getDruidNodeNames gets all the druid nodes in a namespace
func (c client) getDruidNodeNames(namespaces, CR string) ([]string, error) {

	var err error

	druidNodeName, err := c.Resource(GVK).Namespace(namespaces).Get(context.TODO(), CR, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var names []string

	nameLists, _, _ := unstructured.NestedMap(druidNodeName.Object, "spec", "nodes")
	for nameList := range nameLists {
		names = append(names, nameList)
		sort.Strings(names)
	}

	return names, nil

}

// listDruidCR lists all the druid CR in a namespace or all namespaces
func (c client) listDruidCR(namespaces string) ([]string, error) {

	var err error

	druidList, err := c.Resource(GVK).Namespace(namespaces).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, d := range druidList.Items {
		names = append(names, d.GetName())
		sort.Strings(names)
	}

	return names, nil

}

// writerNodeSpecReplicas writer nodespec replica
func (c client) writerDruidNodeSpecReplicas(nodeName, namespace, CR string, replica int64) (bool, error) {
	var err error
	cr, err := c.Resource(GVK).Namespace(namespace).Get(context.TODO(), CR, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	if err := unstructured.SetNestedField(cr.Object, int64(replica), "spec", "nodes", nodeName, "replicas"); err != nil {
		return false, err
	}

	_, err = c.Resource(GVK).Namespace(namespace).Update(context.TODO(), cr, v1.UpdateOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// writerDruidNodeImages writers updates nodes images
func (c client) writerDruidNodeImages(nodeName, namespace, CR, image string) (bool, error) {
	var err error

	cr, err := c.Resource(GVK).Namespace(namespace).Get(context.TODO(), CR, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	if err := unstructured.SetNestedField(cr.Object, image, "spec", "nodes", nodeName, "image"); err != nil {
		return false, err
	}

	_, err = c.Resource(GVK).Namespace(namespace).Update(context.TODO(), cr, v1.UpdateOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// patcherDruidDeleteOrphanPvc patches DeleteOrphanPvc flag
func (c client) patcherDruidDeleteOrphanPvc(namespace, CR string, value bool) (bool, error) {
	var err error

	patchBytes := NewPatchValue("replace", "/spec/deleteOrphanPvc", value)

	_, err = c.Resource(GVK).Namespace(namespace).Patch(context.TODO(), CR, types.JSONPatchType, patchBytes, v1.PatchOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// patcherDruidRollingDeploy patches rollingDeploy flag
func (c client) patcherDruidRollingDeploy(namespace, CR string, value bool) (bool, error) {
	var err error

	patchBytes := NewPatchValue("replace", "/spec/rollingDeploy", value)

	_, err = c.Resource(GVK).Namespace(namespace).Patch(context.TODO(), CR, types.JSONPatchType, patchBytes, v1.PatchOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}
