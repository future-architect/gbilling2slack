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

type ProjectCost struct {
	ProjectID   string
	MonthlyCost int64
	DailyCost   int64
}

// SummaryReport is represent billing summary message. For example is below.
//
// < 06/01 - 06/22 >  Invoice YYYY/MM ( MM/DD 00:00-24:00 )
// dev-pj                    |          0 円 (     0 → )
// stg-pj                    |     15,044 円 (   347 ↑ )
// rcv-pj                    |      5,551 円 (   114 ↑ )
// ―――――――――――――――――――――――――――――
// Sum                       |     20,596 円 (   462 ↑ )
type SummaryReport struct {
	ProjectCostList  []ProjectCost
	MonthlyTotalCost int64
	DailyTotalCost   int64
}

func NewSummaryReport(monthBilling, dayBilling invoice.CostList) *SummaryReport {
	monthCosts := monthBilling.GroupBy()
	dailyCosts := dayBilling.GroupBy()

	var pcList []ProjectCost
	for projectID, costList := range monthCosts {
		pcList = append(pcList, ProjectCost{
			ProjectID:   projectID,
			MonthlyCost: costList.CalcTotalCost(),
			DailyCost:   dailyCosts[projectID].CalcTotalCost(),
		})
	}

	// sort by A-Z
	sort.Slice(pcList, func(i, j int) bool {
		return strings.ToLower(pcList[i].ProjectID) < strings.ToLower(pcList[j].ProjectID)
	})

	return &SummaryReport{
		ProjectCostList:  pcList,
		MonthlyTotalCost: monthBilling.CalcTotalCost(),
		DailyTotalCost:   dayBilling.CalcTotalCost(),
	}
}

func (r SummaryReport) GetMaxKeyLength(initVal int) int {
	length := initVal
	for _, pc := range r.ProjectCostList {
		if length < len(pc.ProjectID) {
			length = len(pc.ProjectID)
		}
	}
	return length
}
