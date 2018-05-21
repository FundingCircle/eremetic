package mesos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/mesos/mesos-go/api/v0/auth"
	"github.com/mesos/mesos-go/api/v0/mesosproto"
	mesossched "github.com/mesos/mesos-go/api/v0/scheduler"
	"golang.org/x/net/context"
)

func getFrameworkID(scheduler *Scheduler) *mesosproto.FrameworkID {
	if scheduler.frameworkID != "" {
		return &mesosproto.FrameworkID{
			Value: proto.String(scheduler.frameworkID),
		}
	}
	return nil
}

func getPrincipalID(credential *mesosproto.Credential) *string {
	if credential != nil {
		return credential.Principal
	}
	return nil
}

func getCredential(settings *Settings) (*mesosproto.Credential, error) {
	if settings.CredentialFile != "" {
		content, err := ioutil.ReadFile(settings.CredentialFile)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"credential_file": settings.CredentialFile,
			}).Error("Unable to read credential_file")
			return nil, err
		}
		fields := strings.Fields(string(content))

		if len(fields) != 2 {
			err := errors.New("Unable to parse credentials")
			logrus.WithError(err).WithFields(logrus.Fields{
				"credential_file": settings.CredentialFile,
			}).Error("Should only contain a key and a secret separated by whitespace")
			return nil, err
		}

		logrus.WithField("principal", fields[0]).Info("Successfully loaded principal")
		return &mesosproto.Credential{
			Principal: proto.String(fields[0]),
			Secret:    proto.String(fields[1]),
		}, nil
	}
	logrus.Debug("No credentials specified in configuration")
	return nil, nil
}

func getAuthContext(ctx context.Context) context.Context {
	return auth.WithLoginProvider(ctx, "SASL")
}

func getIPAddressFromConsul(hostname string) (net.IP, error) {
	client, err := api.NewClient(&api.Config{
		Address:    "consul.service.consul:8500",
		Scheme:     "http",
		HttpClient: http.DefaultClient,
	})

	if err != nil {
		fmt.Println("Error with consul client")
		fmt.Println(err)
		return nil, err
	}
	catalog := client.Catalog()
	nodes, meta, err := catalog.Nodes(nil)
	if err != nil {
		fmt.Println("Error with consul nodes")
		fmt.Println(meta)
		fmt.Println(err)
		return nil, err
	}
	for _, n := range nodes {
		fmt.Println("Node")
		fmt.Println(n.Node)
		fmt.Println("Address")
		fmt.Println(n.Address)
		if n.Node == hostname {
			return net.ParseIP(n.Address), nil
		}
	}
	err = errors.New("couldn't find in consul")

	return nil, err
}

func createDriver(scheduler *Scheduler, settings *Settings) (*mesossched.MesosSchedulerDriver, error) {
	publishedAddr := net.ParseIP(settings.MessengerAddress)
	bindingPort := settings.MessengerPort
	credential, err := getCredential(settings)

	consulIP, err := getIPAddressFromConsul(settings.MessengerAddress)
	fmt.Println("Setting up mesos shite")
	fmt.Println("Setting: ", settings.MessengerAddress)
	fmt.Println("Setting: ", publishedAddr)
	fmt.Println("Consul:", consulIP)

	if err != nil {
		return nil, err
	}

	return mesossched.NewMesosSchedulerDriver(mesossched.DriverConfig{
		Master: settings.Master,
		Framework: &mesosproto.FrameworkInfo{
			Id:              getFrameworkID(scheduler),
			Name:            proto.String(settings.Name),
			User:            proto.String(settings.User),
			Checkpoint:      proto.Bool(settings.Checkpoint),
			FailoverTimeout: proto.Float64(settings.FailoverTimeout),
			Hostname:        proto.String(settings.MessengerAddress),
			Principal:       getPrincipalID(credential),
		},
		Scheduler:        scheduler,
		BindingAddress:   net.ParseIP("0.0.0.0"),
		PublishedAddress: consulIP,
		BindingPort:      bindingPort,
		HostnameOverride: settings.MessengerAddress,
		Credential:       credential,
		WithAuthContext:  getAuthContext,
	})
}
