package rbac

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	core_util "kmodules.xyz/client-go/core/v1"
	rbac_util "kmodules.xyz/client-go/rbac/v1"
	api_v1beta1 "stash.appscode.dev/stash/apis/stash/v1beta1"
)

const (
	StashCronJob = "stash-cron-job"
)

func EnsureCronJobRBAC(kubeClient kubernetes.Interface, resource *core.ObjectReference, sa string, psps []string, labels map[string]string) error {
	// ensure CronJob cluster role
	err := ensureCronJobClusterRole(kubeClient, psps, labels)
	if err != nil {
		return err
	}

	// ensure RoleBinding
	err = ensureCronJobRoleBinding(kubeClient, resource, sa, labels)
	return nil
}

func ensureCronJobClusterRole(kubeClient kubernetes.Interface, psps []string, labels map[string]string) error {
	meta := metav1.ObjectMeta{
		Name:   StashCronJob,
		Labels: labels,
	}
	_, _, err := rbac_util.CreateOrPatchClusterRole(kubeClient, meta, func(in *rbac.ClusterRole) *rbac.ClusterRole {
		in.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{api_v1beta1.SchemeGroupVersion.Group},
				Resources: []string{api_v1beta1.ResourcePluralBackupSession},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{api_v1beta1.SchemeGroupVersion.Group},
				Resources: []string{api_v1beta1.ResourcePluralBackupConfiguration},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"events"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups:     []string{policy.GroupName},
				Resources:     []string{"podsecuritypolicies"},
				Verbs:         []string{"use"},
				ResourceNames: psps,
			},
			{
				APIGroups: []string{apps.GroupName},
				Resources: []string{"deployments", "statefulsets", "replicasets", "daemonsets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"replicationcontrollers", "persistentvolumeclaims"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"apps.openshift.io"},
				Resources: []string{"deploymentconfigs"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"appcatalog.appscode.com"},
				Resources: []string{"*"},
				Verbs:     []string{"get"},
			},
		}
		return in

	})
	return err
}

func ensureCronJobRoleBinding(kubeClient kubernetes.Interface, resource *core.ObjectReference, sa string, labels map[string]string) error {
	meta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", StashCronJob, resource.Name),
		Namespace: resource.Namespace,
		Labels:    labels,
	}

	// ensure role binding
	_, _, err := rbac_util.CreateOrPatchRoleBinding(kubeClient, meta, func(in *rbac.RoleBinding) *rbac.RoleBinding {
		core_util.EnsureOwnerReference(&in.ObjectMeta, resource)

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     KindClusterRole,
			Name:     StashCronJob,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind,
				Name:      sa,
				Namespace: resource.Namespace,
			},
		}
		return in
	})
	if err != nil {
		return err
	}
	return nil
}
