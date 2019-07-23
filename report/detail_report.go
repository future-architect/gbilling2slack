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
package report

import (
	"github.com/future-architect/gbilling2slack/invoice"
	"sort"
	"strings"
)

type ServiceCost struct {
	ServiceName        string
	MonthlyServiceCost int64
	DailyServiceCost   int64
}

// DetailReport represent is detail cost per projectID. For example is below.
//
// < 06/01 - 06/22 > YOUR-PROJECT-ID
// service name              |    month cost ( day cost )
// ------------------------------------------------------
// Cloud Scheduler           |          1 円 (     0 → )
// Compute Engine            |     15,042 円 (   347 ↑ )
type DetailReport struct {
	ProjectID        string
	ServiceCostList  []ServiceCost
	MonthlyTotalCost int64
	DailyTotalCost   int64
}

func NewDetailReportList(monthBilling, dayBilling invoice.CostList) []*DetailReport {
	monthCosts := monthBilling.GroupBy()
	dayCosts := dayBilling.GroupBy()

	var result []*DetailReport
	for projectID, costList := range monthCosts {

		var scList []ServiceCost
		for _, cost := range costList {
			scList = append(scList, ServiceCost{
				ServiceName:        cost.ServiceName,
				MonthlyServiceCost: int64(cost.Cost),
				DailyServiceCost:   dayCosts[projectID].GetCost(projectID, cost.ServiceName),
			})
		}

		// sort by A-Z
		sort.Slice(scList, func(i, j int) bool {
			return strings.ToLower(scList[i].ServiceName) < strings.ToLower(scList[j].ServiceName)
		})

		result = append(result, &DetailReport{
			ProjectID:        projectID,
			ServiceCostList:  scList,
			MonthlyTotalCost: costList.CalcTotalCost(),
			DailyTotalCost:   dayCosts[projectID].CalcTotalCost(),
		})
	}

	// sort by A-Z
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].ProjectID) < strings.ToLower(result[j].ProjectID)
	})

	return result
}

func (r DetailReport) GetMaxKeyLength(initVal int) int {
	if initVal < len(r.ProjectID) {
		return len(r.ProjectID)
	}
	return initVal
}
