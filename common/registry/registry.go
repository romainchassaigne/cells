/*
 * Copyright (c) 2018. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

// Package registry provides the main glue between services
//
// It wraps micro registry (running services declared to the discovery server) into a more generic registry where all
// actual plugins are self-declared.
package registry

import (
	"sync"

	"github.com/gyuho/goraph"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/pydio/cells/common"
	"github.com/pydio/cells/common/service/defaults"
	"github.com/spf13/pflag"
)

var (
	// Default registry
	Default = NewRegistry()
)

// Registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option)
	Register(Service, ...RegisterOption) error
	Deregister(Service) error
	GetService(string) ([]Service, error)
	GetServicesByName(string) []Service
	ListServices(withExcluded ...bool) ([]Service, error)
	ListRunningServices() ([]Service, error)
	ListServicesWithMicroMeta(string, ...string) ([]Service, error)
	SetServiceStopped(string) error
	Filter(func(Service) bool) error
	Watch() (Watcher, error)
	String() string
	Options() Options

	BeforeInit() error
	AfterInit() error
}

type pydioregistry struct {
	*sync.Mutex
	graph goraph.Graph

	flags    pflag.FlagSet
	register map[string]Service

	runningmutex *sync.Mutex
	running      []*registry.Service

	opts Options
}

// Init the default registry
func Init(opts ...Option) {
	Default.Init(opts...)
}

// ListServices returns the list of services that are started in the system
func ListServices(withExcluded ...bool) ([]Service, error) {
	return Default.ListServices(withExcluded...)
}

// ListRunningServices returns the list of services that are started in the system
func ListRunningServices() ([]Service, error) {
	return Default.ListRunningServices()
}

// Watch triggers a watch of the default registry
func Watch() (Watcher, error) {
	return Default.Watch()
}

// NewRegistry provides a new registry object
func NewRegistry(opts ...Option) Registry {
	r := &pydioregistry{
		Mutex:        &sync.Mutex{},
		graph:        goraph.NewGraph(),
		register:     make(map[string]Service),
		runningmutex: &sync.Mutex{},
		opts:         newOptions(opts...),
	}

	return r
}

// Init the registry with the options in arguments
func (c *pydioregistry) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&c.opts)
	}

	if c.opts.Name != "" {
		services, err := c.ListServices()
		if err != nil {
			return
		}

		for _, s := range services {
			s.AddDependency(c.opts.Name)
		}

		// Retrieving the node that got the main registry
		registryNode := goraph.NewNode(c.opts.Name)
		if node, err := c.graph.GetNode(registryNode); err == nil && node != nil {
			registryNode = node
		} else {
			c.graph.AddNode(registryNode)
		}

		for _, s := range services {
			if c.opts.Name == s.Name() {
				continue
			}
			serviceNode := goraph.NewNode(s.Name())
			if node, err := c.graph.GetNode(serviceNode); err == nil && node != nil {
				serviceNode = node
			} else {
				c.graph.AddNode(serviceNode)
			}

			c.graph.AddEdge(registryNode.ID(), serviceNode.ID(), 1)
		}
	}
}

// Deregister sets a service as excluded in the registry
func (c *pydioregistry) Deregister(s Service) error {
	// delete our hash of the service
	c.Lock()
	c.register[s.Name()].SetExcluded(true)
	c.Unlock()

	return nil
}

// Register a new service. Manages dependencies
func (c *pydioregistry) Register(s Service, opts ...RegisterOption) error {

	var options RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	c.Lock()
	defer c.Unlock()

	id := s.Name()

	c.register[id] = s

	mainNode := goraph.NewNode(id)
	if node, err := c.graph.GetNode(mainNode); err == nil && node != nil {
		mainNode = node
	} else {
		c.graph.AddNode(mainNode)
	}

	// Handling dependencies
	for _, dependency := range options.Dependencies {
		dependencyNode := goraph.NewNode(dependency)
		if node, err := c.graph.GetNode(dependencyNode); err == nil && node != nil {
			dependencyNode = node
		} else {
			c.graph.AddNode(dependencyNode)
		}

		c.graph.AddEdge(dependencyNode.ID(), mainNode.ID(), 1)
	}

	for _, flag := range options.Flags {
		c.flags.AddFlag(flag)
	}

	return nil
}

// GetService returns the service by name
func (c *pydioregistry) GetService(string) ([]Service, error) {
	return nil, nil
}

// GetClient returns the default client for the service name
func GetClient(name string) (string, client.Client) {
	return common.SERVICE_GRPC_NAMESPACE_ + name, defaults.NewClient()
}

// Filter the service out of the registry
func (c *pydioregistry) Filter(f func(s Service) bool) error {
	services, err := c.ListServices()
	if err != nil {
		return err
	}

	for _, s := range services {
		if f(s) {
			c.Deregister(s)
		}
	}

	return nil
}

// Watch the registry for new entries
func (c *pydioregistry) Watch() (Watcher, error) {
	return newWatcher(), nil
}

// String representation of the registry
func (c *pydioregistry) String() string {
	return "pydio"
}

// Options returns the list of options set to the registry
func (c *pydioregistry) Options() Options {
	return c.opts
}

// BeforeInit calls the before init function for each service in the registry
func (c *pydioregistry) BeforeInit() error {
	services, err := c.ListServices()
	if err != nil {
		return err
	}

	for _, s := range services {
		s.BeforeInit()
	}

	return nil
}

// AfterInit calls the after init function for each service in the registry
func (c *pydioregistry) AfterInit() error {
	services, err := c.ListServices()
	if err != nil {
		return err
	}

	for _, s := range services {
		s.AfterInit()
	}

	return nil
}
