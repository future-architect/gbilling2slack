/**
 * Copyright (c) 2019-present Future Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package invoice

import (
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

type invoice struct {
	projectID string
	tableName string
}

func NewInvoice(projectID, tableName string) *invoice {
	return &invoice{
		projectID: projectID,
		tableName: tableName,
	}
}

func (i *invoice) FetchBillingDay(ctx context.Context) (CostList, error) {
	stmt := `
		SELECT
			project.id as project,
			service.description as service,
			IFNULL(sum(cost), 0) as cost
		FROM	
			` + "`" + i.tableName + "`" + `
		WHERE
			DATE(_PARTITIONTIME) = ` + time.Now().Format("'2006-01-02'") + `
		AND
			project.id IS NOT NULL
		GROUP BY
			project,
			service
		ORDER BY
			project`

	return i.fetchBilling(ctx, stmt)
}

func (i *invoice) FetchBillingMonth(ctx context.Context) (CostList, error) {

	stmt := `
		SELECT
			project.id as project,
			service.description as service,
			IFNULL(sum(cost), 0) as cost
		FROM
			` + "`" + i.tableName + "`" + `
		WHERE
			invoice.month = ` + time.Now().Format("'200601'") + `
		AND
			project.id IS NOT NULL
		GROUP BY
			project,
			service
		ORDER BY
			project`

	return i.fetchBilling(ctx, stmt)
}

func (i *invoice) fetchBilling(ctx context.Context, stmt string) (CostList, error) {

	client, err := bigquery.NewClient(ctx, i.projectID)
	if err != nil {
		return nil, err
	}

	iter, err := client.Query(stmt).Read(ctx)
	if err != nil {
		return nil, err
	}

	billing := make([]Cost, 0)
	for {
		var sc Cost
		err := iter.Next(&sc)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		billing = append(billing, sc)
	}

	return billing, nil
}
