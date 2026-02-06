// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package permission

import (
	"log"

	"github.com/casbin/casbin/v2"
)

// Service wraps the Casbin enforcer
type Service struct {
	Enforcer *casbin.Enforcer
}

// NewService initializes the permission service.
// modelConf: path to rbac_model.conf or text.
// adapter: file path (string) or a persist.Adapter implementation.
func NewService(modelConf string, adapter interface{}) (*Service, error) {
	e, err := casbin.NewEnforcer(modelConf, adapter)
	if err != nil {
		return nil, err
	}

	// Load policies from adapter
	if err := e.LoadPolicy(); err != nil {
		log.Printf("Failed to load generic access policies: %v", err)
		// We define this as non-fatal during init unless strictly required? 
		// Usually fatal if policy cannot be loaded.
		return nil, err
	}

	return &Service{Enforcer: e}, nil
}

// Check verifies if the subject has permission
func (s *Service) Check(sub, obj, act string) (bool, error) {
	return s.Enforcer.Enforce(sub, obj, act)
}

// AddPolicy adds a specific permission rule
func (s *Service) AddPolicy(sub, obj, act string) (bool, error) {
	return s.Enforcer.AddPolicy(sub, obj, act)
}

// RemovePolicy removes a specific permission rule
func (s *Service) RemovePolicy(sub, obj, act string) (bool, error) {
	return s.Enforcer.RemovePolicy(sub, obj, act)
}

// AddRoleForUser assigns a role to a user
func (s *Service) AddRoleForUser(user, role string) (bool, error) {
	return s.Enforcer.AddGroupingPolicy(user, role)
}

// GetRolesForUser gets the roles for a user
func (s *Service) GetRolesForUser(user string) ([]string, error) {
	return s.Enforcer.GetRolesForUser(user)
}
