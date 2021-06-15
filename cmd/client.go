package cmd

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// newClient() shall return a dynamic client
func newClient() *client {

	dynamicClient, err := dynamic.NewForConfig(newConfig())
	if err != nil {
		panic(err.Error())
	}

	return &client{dynamicClient}
}

// newConfig() shall return a config
func newConfig() *rest.Config {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}
