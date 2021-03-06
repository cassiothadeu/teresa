package cloudprovider

import (
	"github.com/luizalabs/teresa/pkg/server/service"
	"github.com/pkg/errors"
)

type Operations interface {
	CreateOrUpdateSSL(appName, cert string, port int) error
	SSLInfo(appName string) (*service.SSLInfo, error)
}

type K8sOperations interface {
	CloudProviderName() (string, error)
	SetServiceAnnotations(namespace, service string, annotations map[string]string) error
	ServiceAnnotations(namespace, service string) (map[string]string, error)
	IsNotFound(err error) bool
}

func NewOperations(k8s K8sOperations) (Operations, error) {
	name, err := k8s.CloudProviderName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cloud provider name")
	}
	switch name {
	case "aws":
		return &awsOperations{k8s: k8s}, nil
	case "gce":
		return &gceOperations{k8s: k8s}, nil
	default:
		return &fallbackOperations{}, nil
	}
}
