// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package upgrades

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/security/username"
	"github.com/cockroachdb/cockroach/pkg/sql/isql"
	_ "github.com/cockroachdb/cockroach/pkg/sql/tablemetadatacache" // Ensure job implementation is linked.
	"github.com/cockroachdb/cockroach/pkg/upgrade"
)

func createUpdateTableMetadataCacheJob(
	ctx context.Context, _ clusterversion.ClusterVersion, d upgrade.TenantDeps,
) error {
	if d.TestingKnobs != nil && d.TestingKnobs.SkipUpdateTableMetadataCacheBootstrap {
		return nil
	}

	return d.DB.Txn(ctx, func(ctx context.Context, txn isql.Txn) error {
		jr := jobs.Record{
			JobID:         jobs.UpdateTableMetadataCacheJobID,
			Description:   jobspb.TypeUpdateTableMetadataCache.String(),
			Details:       jobspb.UpdateTableMetadataCacheDetails{},
			Progress:      jobspb.UpdateTableMetadataCacheProgress{},
			CreatedBy:     &jobs.CreatedByInfo{Name: username.NodeUser, ID: username.NodeUserID},
			Username:      username.NodeUserName(),
			NonCancelable: true,
		}
		return d.JobRegistry.CreateIfNotExistAdoptableJobWithTxn(ctx, jr, txn)
	})
}
