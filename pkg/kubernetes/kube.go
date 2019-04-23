package kube

import (
	"log"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClientSet does something.
func ClientSet() (kubernetes.Interface, error) {
	log.Print("[info][kubernetes] connecting to kubernetes API...")

	c, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the kubernetes API")
	}

	log.Print("[info][kubernetes] connected to kubernetes API")

	return kubernetes.NewForConfig(c)
}
