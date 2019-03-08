package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
)

// obsoleteReasons are the set of reasons that mean a CSV should no longer be processed as active
var obsoleteReasons = map[ConditionReason]struct{}{
	CSVReasonReplaced:      {},
	CSVReasonBeingReplaced: {},
}

func (c *ClusterServiceVersion) SetPhaseWithEvent(phase ClusterServiceVersionPhase, reason ConditionReason, message string, now metav1.Time, recorder record.EventRecorder) {
	var eventtype string
	if phase == CSVPhaseFailed {
		eventtype = v1.EventTypeWarning
	} else {
		eventtype = v1.EventTypeNormal
	}
	go recorder.Event(c, eventtype, string(reason), message)
	c.SetPhase(phase, reason, message, now)
}

// SetPhase sets the current phase and adds a condition if necessary
func (c *ClusterServiceVersion) SetPhase(phase ClusterServiceVersionPhase, reason ConditionReason, message string, now metav1.Time) {
	c.Status.LastUpdateTime = now
	if c.Status.Phase != phase {
		c.Status.Phase = phase
		c.Status.LastTransitionTime = now
	}
	c.Status.Message = message
	c.Status.Reason = reason
	if len(c.Status.Conditions) == 0 {
		c.Status.Conditions = append(c.Status.Conditions, ClusterServiceVersionCondition{
			Phase:              c.Status.Phase,
			LastTransitionTime: c.Status.LastTransitionTime,
			LastUpdateTime:     c.Status.LastUpdateTime,
			Message:            message,
			Reason:             reason,
		})
	}
	previousCondition := c.Status.Conditions[len(c.Status.Conditions)-1]
	if previousCondition.Phase != c.Status.Phase || previousCondition.Reason != c.Status.Reason {
		c.Status.Conditions = append(c.Status.Conditions, ClusterServiceVersionCondition{
			Phase:              c.Status.Phase,
			LastTransitionTime: c.Status.LastTransitionTime,
			LastUpdateTime:     c.Status.LastUpdateTime,
			Message:            message,
			Reason:             reason,
		})
	}
}

// SetRequirementStatus adds the status of all requirements to the CSV status
func (c *ClusterServiceVersion) SetRequirementStatus(statuses []RequirementStatus) {
	c.Status.RequirementStatus = statuses
}

// IsObsolete returns if this CSV is being replaced or is marked for deletion
func (c *ClusterServiceVersion) IsObsolete() bool {
	for _, condition := range c.Status.Conditions {
		_, ok := obsoleteReasons[condition.Reason]
		if ok {
			return true
		}
	}
	return false
}
