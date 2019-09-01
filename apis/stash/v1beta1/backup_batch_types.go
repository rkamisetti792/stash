package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindBackupBatch     = "BackupBatch"
	ResourceSingularBackupBatch = "backupbatch"
	ResourcePluralBackupBatch   = "backupbatchs"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackupBatch struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupBatchSpec `json:"spec,omitempty"`
}

type BackupBatchSpec struct {
	Schedule string `json:"schedule,omitempty"`

	// Members specify the backup configurations that are part of this batch
	// +optional
	Members []BackupConfiguration `json:"members,omitempty"`

	// Indicates that the BackupBatch is paused from taking backup. Default value is 'false'
	// +optional
	Paused bool `json:"paused,omitempty"`

	// Actions that Stash should take in response to backup sessions.
	// Cannot be updated.
	// +optional
	Hooks *Hooks `json:"hooks,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BackupBatchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupBatch `json:"items,omitempty"`
}
