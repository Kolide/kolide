package inmem

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	kolide_errors "github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (orm *Inmem) NewLabel(label *kolide.Label) (*kolide.Label, error) {
	newLabel := *label

	orm.mtx.Lock()
	for _, l := range orm.labels {
		if l.Name == label.Name {
			return nil, kolide_errors.ErrExists
		}
	}

	newLabel.ID = orm.nextID(label)
	orm.labels[newLabel.ID] = &newLabel
	orm.mtx.Unlock()

	return &newLabel, nil
}

func (orm *Inmem) ListLabelsForHost(hid uint) ([]kolide.Label, error) {
	// First get IDs of label executions for the host
	resLabels := []kolide.Label{}

	orm.mtx.Lock()
	for _, lqe := range orm.labelQueryExecutions {
		if lqe.HostID == hid && lqe.Matches {
			if label := orm.labels[lqe.LabelID]; label != nil {
				resLabels = append(resLabels, *label)
			}
		}
	}
	orm.mtx.Unlock()

	return resLabels, nil
}

func (orm *Inmem) LabelQueriesForHost(host *kolide.Host, cutoff time.Time) (map[string]string, error) {
	// Get post-cutoff executions for host
	execedIDs := map[uint]bool{}

	orm.mtx.Lock()
	for _, lqe := range orm.labelQueryExecutions {
		if lqe.HostID == host.ID && (lqe.UpdatedAt == cutoff || lqe.UpdatedAt.After(cutoff)) {
			execedIDs[lqe.LabelID] = true
		}
	}

	queries := map[string]string{}
	for _, label := range orm.labels {
		if (label.Platform == "" || strings.Contains(label.Platform, host.Platform)) && !execedIDs[label.ID] {
			queries[strconv.Itoa(int(label.ID))] = label.Query
		}
	}
	orm.mtx.Unlock()

	return queries, nil
}

func (orm *Inmem) getLabelByIDString(id string) (*kolide.Label, error) {
	labelID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("non-int label ID")
	}

	orm.mtx.Lock()
	label, ok := orm.labels[uint(labelID)]
	orm.mtx.Unlock()

	if !ok {
		return nil, errors.New("label ID not found: " + string(labelID))
	}

	return label, nil
}

func (orm *Inmem) RecordLabelQueryExecutions(host *kolide.Host, results map[string]bool, t time.Time) error {
	// Record executions
	for strLabelID, matches := range results {
		label, err := orm.getLabelByIDString(strLabelID)
		if err != nil {
			return err
		}

		updated := false
		orm.mtx.Lock()
		for _, lqe := range orm.labelQueryExecutions {
			if lqe.LabelID == label.ID && lqe.HostID == host.ID {
				// Update existing execution values
				lqe.UpdatedAt = t
				lqe.Matches = matches
				updated = true
				break
			}
		}

		if !updated {
			// Create new execution
			lqe := kolide.LabelQueryExecution{
				HostID:    host.ID,
				LabelID:   label.ID,
				UpdatedAt: t,
				Matches:   matches,
			}
			lqe.ID = orm.nextID(lqe)
			orm.labelQueryExecutions[lqe.ID] = &lqe
		}
		orm.mtx.Unlock()
	}

	return nil
}

func (orm *Inmem) DeleteLabel(lid uint) error {
	orm.mtx.Lock()
	delete(orm.labels, lid)
	orm.mtx.Unlock()

	return nil
}

func (orm *Inmem) Label(lid uint) (*kolide.Label, error) {
	orm.mtx.Lock()
	label, ok := orm.labels[lid]
	orm.mtx.Unlock()

	if !ok {
		return nil, errors.New("Label not found")
	}
	return label, nil
}

func (orm *Inmem) ListLabels(opt kolide.ListOptions) ([]*kolide.Label, error) {
	// We need to sort by keys to provide reliable ordering
	keys := []int{}

	orm.mtx.Lock()
	for k, _ := range orm.labels {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	labels := []*kolide.Label{}
	for _, k := range keys {
		labels = append(labels, orm.labels[uint(k)])
	}
	orm.mtx.Unlock()

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":         "ID",
			"created_at": "CreatedAt",
			"updated_at": "UpdatedAt",
			"name":       "Name",
		}
		if err := sortResults(labels, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := orm.getLimitOffsetSliceBounds(opt, len(labels))
	labels = labels[low:high]

	return labels, nil
}

func (orm *Inmem) SearchLabels(query string, omit ...uint) ([]kolide.Label, error) {
	omitLookup := map[uint]bool{}
	for _, o := range omit {
		omitLookup[o] = true
	}

	var results []kolide.Label

	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, l := range orm.labels {
		if len(results) == 10 {
			break
		}

		if strings.Contains(l.Name, query) && !omitLookup[l.ID] {
			results = append(results, *l)
			continue
		}
	}

	return results, nil
}

func (orm *Inmem) ListHostsInLabel(lid uint) ([]kolide.Host, error) {
	var hosts []kolide.Host

	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, lqe := range orm.labelQueryExecutions {
		if lqe.LabelID == lid && lqe.Matches {
			hosts = append(hosts, *orm.hosts[lqe.HostID])
		}
	}

	return hosts, nil
}

func (orm *Inmem) ListUniqueHostsInLabels(labels []uint) ([]kolide.Host, error) {
	var hosts []kolide.Host

	labelSet := map[uint]bool{}
	hostSet := map[uint]bool{}

	for _, label := range labels {
		labelSet[label] = true
	}

	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	for _, lqe := range orm.labelQueryExecutions {
		if labelSet[lqe.LabelID] && lqe.Matches {
			if !hostSet[lqe.HostID] {
				hosts = append(hosts, *orm.hosts[lqe.HostID])
				hostSet[lqe.HostID] = true
			}
		}
	}

	return hosts, nil
}
