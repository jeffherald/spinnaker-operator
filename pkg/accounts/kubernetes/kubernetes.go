package kubernetes

import (
	"errors"
	"os"

	"strings"

	"github.com/armory/spinnaker-operator/pkg/accounts/account"
	"github.com/armory/spinnaker-operator/pkg/apis/spinnaker/interfaces"
)

// Kubernetes accounts have a deeper integration than other accounts.
// When read from Spinnaker settings, they support `kubeconfigFile`, `kubeconfigContents`, or oauth via `oauth2l`.
// When read from the CRD, user can reference a Kubernetes secret, pass the kubeconfig file inlined,
// reference a secret (s3, gcs...), or pass provider options to make the kubeconfig on the fly.
const (
	KubeconfigFileSettings        = "kubeconfigFile"
	KubeconfigFileContentSettings = "kubeconfigContents"
	UseServiceAccount             = "serviceAccount"
)

var TypesFactory interfaces.TypesFactory

type AccountType struct{}

func (k *AccountType) GetType() interfaces.AccountType {
	return interfaces.KubernetesAccountType
}

func (k *AccountType) GetAccountsKey() string {
	return "kubernetes.accounts"
}

func (k *AccountType) GetConfigAccountsKey() string {
	return "providers.kubernetes.accounts"
}

func (k *AccountType) GetServices() []string {
	// TODO(maybe): this was formerly hardcoded and should probably be derived from the spinnaker service
	// configuration, but we will punt for now and use an environment variable due to the number of changes
	// that would be needed to support that.
	val := os.Getenv("CLOUDDRIVER_DEPLOYMENTS")
	if val == "" {
		return []string{"clouddriver"}
	}
	svcMap := map[string]bool{}
	for _, s := range strings.Split(val, ",") {
		svcMap[strings.TrimSpace(s)] = true
	}
	svcs := make([]string, len(svcMap))
	idx := 0
	for s := range svcMap {
		svcs[idx] = s
		idx += 1
	}
	return svcs
}

func (k *AccountType) GetPrimaryAccountsKey() string {
	return "providers.kubernetes.primaryAccount"
}

func (k *AccountType) newAccount() *Account {
	return &Account{
		Env: Env{},
	}
}

func (k *AccountType) GetValidationSettings(spinsvc interfaces.SpinnakerService) *interfaces.ValidationSetting {
	v := spinsvc.GetSpinnakerValidation()
	for n, s := range v.Providers {
		if strings.ToLower(n) == strings.ToLower(string(interfaces.KubernetesAccountType)) {
			return &s
		}
	}
	return v.GetValidationSettings()
}

type Env struct {
	Namespaces      []string                   `json:"namespaces,omitempty"`
	OmitNamespaces  []string                   `json:"omitNamespaces,omitempty"`
	Kinds           []string                   `json:"kinds,omitempty"`
	OmitKinds       []string                   `json:"omitKinds,omitempty"`
	CustomResources []CustomKubernetesResource `json:"customResources,omitempty"`
}

type CustomKubernetesResource struct {
	KubernetesKind string `json:"kubernetesKind,omitempty"`
	SpinnakerKind  string `json:"spinnakerKind,omitEmpty"`
	Versioned      bool   `json:"versioned,omitempty"`
}

type Account struct {
	*account.BaseAccount
	Name     string `json:"name,omitempty"`
	Auth     *interfaces.KubernetesAuth
	Env      Env                 `json:"env,omitempty"`
	Settings interfaces.FreeForm `json:"settings,omitempty"`
}

func (k *Account) GetType() interfaces.AccountType {
	return interfaces.KubernetesAccountType
}

func (k *Account) GetName() string {
	return k.Name
}

func (k *Account) GetSettings() *interfaces.FreeForm {
	return &k.Settings
}

func (k *Account) validateFormat() error {
	if k.Name == "" {
		return errors.New("Spinnaker account must have a name")
	}
	return nil
}

func (k *Account) NewValidator() account.AccountValidator {
	return &kubernetesAccountValidator{account: k}
}
