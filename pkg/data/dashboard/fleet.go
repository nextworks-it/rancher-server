package dashboard

import (
	"reflect"

	"github.com/rancher/rancher/pkg/features"
	"github.com/rancher/rancher/pkg/wrangler"
	rbacv1 "github.com/rancher/wrangler/pkg/generated/controllers/rbac/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddFleetRoles(wrangler *wrangler.Context) error {
	f, err := wrangler.Mgmt.Feature().Get("fleet", metav1.GetOptions{})
	if err != nil {
		return err
	}

	if !features.IsEnabled(f) {
		return nil
	}

	return ensureFleetRoles(wrangler.RBAC)
}

func ensureFleetRoles(rbac rbacv1.Interface) error {
	fleetWorkspaceAdminRole := v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fleetworkspace-admin",
		},
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{
					"fleet.cattle.io",
				},
				Resources: []string{
					"clusterregistrationtokens",
					"gitreporestrictions",
					"clusterregistrations",
					"clusters",
					"gitrepos",
					"bundles",
					"clustergroups",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"rbac.authorization.k8s.io",
				},
				Resources: []string{
					"rolebindings",
				},
				Verbs: []string{
					"*",
				},
			},
		},
	}

	fleetWorkspaceMemberRole := v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fleetworkspace-member",
		},
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{
					"fleet.cattle.io",
				},
				Resources: []string{
					"gitrepos",
					"bundles",
				},
				Verbs: []string{
					"*",
				},
			},
			{
				APIGroups: []string{
					"fleet.cattle.io",
				},
				Resources: []string{
					"clusterregistrationtokens",
					"gitreporestrictions",
					"clusterregistrations",
					"clusters",
					"clustergroups",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
		},
	}

	fleetWorkspaceReadonlyRole := v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fleetworkspace-readonly",
		},
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{
					"fleet.cattle.io",
				},
				Resources: []string{
					"clusterregistrationtokens",
					"gitreporestrictions",
					"clusterregistrations",
					"clusters",
					"gitrepos",
					"bundles",
					"clustergroups",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
			},
		},
	}

	clusterRoles := []v1.ClusterRole{
		fleetWorkspaceAdminRole,
		fleetWorkspaceMemberRole,
		fleetWorkspaceReadonlyRole,
	}

	for _, role := range clusterRoles {
		existing, err := rbac.ClusterRole().Get(role.Name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return err
		} else if errors.IsNotFound(err) {
			if _, err := rbac.ClusterRole().Create(&role); err != nil {
				return err
			}
		} else {
			if !reflect.DeepEqual(existing.Rules, role.Rules) {
				if _, err := rbac.ClusterRole().Update(&role); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
