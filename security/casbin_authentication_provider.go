package security

import (
	"github.com/casbin/casbin"
)

var (
	casbinAuthenticationProvider *TCasbinAuthenticationProvider
)

// TCasbinAuthenticationProvider an implementation that retrieves user details from a UserDetailService
type TCasbinAuthenticationProvider struct {
	enforcer *casbin.Enforcer
}

type TCasbinPolicyDetails struct {
	Sub     string
	Domain  string
	Obj     string
	Act     string
	Service string
	Eft     string
}

func newCasbinAuthenticationProvider() IAuthenticationProvider {
	if casbinAuthenticationProvider == nil {
		return &TCasbinAuthenticationProvider{
			enforcer: GetCasbinEnforcer(),
		}
	}
	return casbinAuthenticationProvider
}

func (provider *TCasbinAuthenticationProvider) Authenticate(authentication IAuthentication) IAuthentication {
	if requestAuthenticationToken, ok := authentication.(*TRequestAuthenticationToken); ok {
		details := authentication.GetDetails()
		if policy, ok := details.(TCasbinPolicyDetails); ok {
			var param []interface{}

			if policy.Sub != "" {
				param = append(param, policy.Sub)
			}
			if policy.Domain != "" {
				param = append(param, policy.Domain)
			}
			if policy.Obj != "" {
				param = append(param, policy.Obj)
			}
			if policy.Act != "" {
				param = append(param, policy.Act)
			}
			if policy.Service != "" {
				param = append(param, policy.Service)
			}
			if policy.Eft != "" {
				param = append(param, policy.Eft)
			}
			if provider.enforcer.Enforce(param...) {
				requestAuthenticationToken.SetAuthenticated(true)
			} else {
				Log.WithField("param:", param).Debug("not promisson")
			}
		}
		return authentication
	}
	return nil
}
