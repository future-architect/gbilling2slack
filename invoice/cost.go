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

type Cost struct {
	ProjectID   string  `bigquery:"project"`
	ServiceName string  `bigquery:"service"`
	Cost        float64 `bigquery:"cost"`
}

type CostList []Cost

func (l CostList) CalcTotalCost() int64 {
	var totalCost float64
	for _, sc := range l {
		totalCost += sc.Cost
	}
	return int64(totalCost)
}

func (l CostList) GroupBy() map[string]CostList {
	result := make(map[string]CostList)
	for _, sc := range l {
		result[sc.ProjectID] = append(result[sc.ProjectID], sc)
	}
	return result
}

func (l CostList) GetCost(projectID, serviceName string) int64 {
	for _, cost := range l {
		if cost.ProjectID == projectID && cost.ServiceName == serviceName {
			return int64(cost.Cost)
		}
	}
	return 0
}
