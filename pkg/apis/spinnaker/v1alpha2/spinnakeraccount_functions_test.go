package v1alpha2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSpinnakerAccountGetItems(t *testing.T) {
	RegisterTypes()
	accounts := SpinnakerAccountList{}
	accounts.Items = []SpinnakerAccount{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "account1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "account2",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "account3",
			},
		},
	}
	u := map[string]bool{}
	for _, a := range accounts.GetItems() {
		_, ok := u[a.GetName()]
		assert.Equalf(t, false, ok, "duplicate account %q", a.GetName())
		u[a.GetName()] = true
	}
	assert.Len(t, u, 3)
}
