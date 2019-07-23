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
	"encoding/json"
	"fmt"
	"github.com/future-architect/gbilling2slack/invoice"
	"testing"
)

var monthlyInput = `[
		{"ProjectID": "dev-pj", "ServiceName": "Compute Engine", "Cost": 14995.11},
		{"ProjectID": "dev-pj", "ServiceName": "App Engine", "Cost": 3.14},
		{"ProjectID": "stg-dx", "ServiceName": "Compute Engine", "Cost": 0.1},
		{"ProjectID": "btn-pj", "ServiceName": "Cloud Scheduler", "Cost": 10.1},
		{"ProjectID": "btn-pj", "ServiceName": "Cloud SQL", "Cost": 741.9},
		{"ProjectID": "btn-pj", "ServiceName": "Compute Engine", "Cost": 4780.1}
	]`

var dailyInput = `[
		{"ProjectID": "dev-pj", "ServiceName": "Compute Engine", "Cost": 298.1},
		{"ProjectID": "dev-pj", "ServiceName": "App Engine", "Cost": 0.2},
		{"ProjectID": "stg-dx", "ServiceName": "Compute Engine", "Cost": 0.01},
		{"ProjectID": "btn-pj", "ServiceName": "Cloud Scheduler", "Cost": 2.1},
		{"ProjectID": "btn-pj", "ServiceName": "Cloud SQL", "Cost": 0.4},
		{"ProjectID": "btn-pj", "ServiceName": "Compute Engine", "Cost": 96.12}
	]`

func TestSummary(t *testing.T) {

	var monthCosts []invoice.Cost
	if err := json.Unmarshal([]byte(monthlyInput), &monthCosts); err != nil {
		t.Fatal("test data parse is failed", err)
	}

	var dayCosts []invoice.Cost
	if err := json.Unmarshal([]byte(dailyInput), &dayCosts); err != nil {
		t.Fatal("test data parse is failed", err)
	}

	report := NewSummaryReport(monthCosts, dayCosts)

	if report.MonthlyTotalCost != 20530 {
		t.Error("monthlyTotalCost is expected 20531", report.MonthlyTotalCost)
	}

	if report.DailyTotalCost != 396 {
		t.Error("dailyTotalCost is expected 397", report.DailyTotalCost)
	}

	// sort check
	if report.ProjectCostList[0].ProjectID != "btn-pj" &&
		report.ProjectCostList[1].ProjectID != "dev-dx" &&
		report.ProjectCostList[2].ProjectID != "stg-dx" {
		t.Error("sort error", report.ProjectCostList[0], report.ProjectCostList[1], report.ProjectCostList[2])
	}

	fmt.Printf("%+v", report)
}

func TestDetail(t *testing.T) {

	var monthCosts []invoice.Cost
	if err := json.Unmarshal([]byte(monthlyInput), &monthCosts); err != nil {
		t.Fatal("test data parse is failed", err)
	}

	var dayCosts []invoice.Cost
	if err := json.Unmarshal([]byte(dailyInput), &dayCosts); err != nil {
		t.Fatal("test data parse is failed", err)
	}

	report := NewDetailReportList(monthCosts, dayCosts)

	if len(report) != 3 {
		t.Fatal("report size is invalid", len(report))
	}

	// check sort
	if report[0].ProjectID != "btn-pj" &&
		report[1].ProjectID != "stg-dx" &&
		report[2].ProjectID != "dev-pj" {
		t.Error("sort error", report[0].ProjectID)
	}

	// check detail sort check
	if report[0].ServiceCostList[0].ServiceName != "Cloud Scheduler" &&
		report[0].ServiceCostList[1].ServiceName != "Cloud SQL" &&
		report[0].ServiceCostList[2].ServiceName != "Compute Engine" {
		t.Error("sort detail list error", report[0].ProjectID)
	}

	// check monthly cost
	if report[0].MonthlyTotalCost != 5532 {
		t.Error("MonthlyTotalCost is expected 5532", report[0].MonthlyTotalCost)
	}

	// check daily cost
	if report[0].DailyTotalCost != 98 {
		t.Error("DailyTotalCost is expected 5532", report[0].DailyTotalCost)
	}

	// check service cost
	if report[0].ServiceCostList[0].ServiceName != "Cloud Scheduler" {
		t.Error("service name is illegal", report[0].ServiceCostList[0].ServiceName)
	}
	if report[0].ServiceCostList[0].MonthlyServiceCost != 10 {
		t.Error("scheduler monthly cost is invalid", report[0].ServiceCostList[0].MonthlyServiceCost)
	}
	if report[0].ServiceCostList[0].ServiceName != "Cloud Scheduler" {
		t.Error("service name is illegal", report[0].ServiceCostList[0].ServiceName)
	}

	fmt.Printf("%+v", report[0])
}
