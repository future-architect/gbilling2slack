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
package notice

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/future-architect/gbilling2slack/report"
	"github.com/nlopes/slack"
	"time"
)

type slackNotifier struct {
	slackAPIToken string
	slackChannel  string
}

func NewSlackNotifier(slackAPIToken, slackChannel string) *slackNotifier {
	return &slackNotifier{
		slackAPIToken: slackAPIToken,
		slackChannel:  slackChannel,
	}
}

func getArrow(cost int64) string {
	if cost > 0 {
		return "↑"
	}
	return "→"
}

func getDate() (int, int, int, string) {
	year, month, day := time.Now().Date()
	period := fmt.Sprintf("< %v/01 - %v/%v > ", int(month), int(month), day)
	return year, int(month), day, period
}

func createHead() string {
	year, month, day, period := getDate()

	text := period
	text += fmt.Sprintf(" Invoice %v/%v ", year, month)
	text += fmt.Sprintf("( %v/%v 00:00-24:00 )\n", month, day)
	return text
}

func insertHeaderPerProject(projectID string, pad int) string {
	_, _, _, period := getDate()

	text := period
	text += fmt.Sprintf("%v\n", projectID)
	text += fmt.Sprintf("%*v | %13v %7v\n", pad, "service name", "month cost", "( day cost )")
	text += fmt.Sprintf("------------------------------------------------------\n")
	return text
}

// Return post message's timestamp to post in the thread
func (n *slackNotifier) postInline(text string) (string, error) {
	_, ts, err := slack.New(n.slackAPIToken).PostMessage(
		n.slackChannel,
		slack.MsgOptionText("```"+text+"```", false),
	)
	return ts, err
}

// ts is parent message's timestamp to post in the thread
func (n *slackNotifier) postThreadInline(text, ts string) error {
	_, _, err := slack.New(n.slackAPIToken).PostMessage(
		n.slackChannel,
		slack.MsgOptionText("```"+text+"```", false),
		slack.MsgOptionTS(ts),
	)
	return err
}

func (n *slackNotifier) PostBilling(summaryReport *report.SummaryReport) (string, error) {

	// padding degree
	padDegree := summaryReport.GetMaxKeyLength(25) * (-1)

	text := createHead()

	// this loop create cost list per project
	for _, cost := range summaryReport.ProjectCostList {
		projectID := cost.ProjectID
		text += fmt.Sprintf("%*v | %10v 円 ( %5v %v )\n",
			padDegree,
			projectID, humanize.Comma(cost.MonthlyCost),
			humanize.Comma(cost.DailyCost),
			getArrow(cost.DailyCost))
	}

	text += fmt.Sprintf("―――――――――――――――――――――――――――――――――――――――――――――――――――――\n")
	text += fmt.Sprintf("%*v | %10v 円 ( %5v %v )\n",
		padDegree,
		"Sum",
		humanize.Comma(summaryReport.MonthlyTotalCost),
		humanize.Comma(summaryReport.DailyTotalCost),
		getArrow(summaryReport.DailyTotalCost))

	// get parent message's timestamp to create thead
	return n.postInline(text)
}

func (n *slackNotifier) PostBillingThread(parentTS string, detailReport *report.DetailReport) error {
	padDegree := detailReport.GetMaxKeyLength(25) * (-1)

	text := insertHeaderPerProject(detailReport.ProjectID, padDegree)

	for _, c := range detailReport.ServiceCostList {
		if c.MonthlyServiceCost == 0 {
			continue
		}

		text += fmt.Sprintf("%*v | %10v 円 ( %5v %v )\n",
			padDegree,
			c.ServiceName,
			humanize.Comma(c.MonthlyServiceCost),
			humanize.Comma(c.DailyServiceCost),
			getArrow(c.DailyServiceCost))
	}
	return n.postThreadInline(text, parentTS)
}
