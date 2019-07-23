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
package gbilling2slack

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/future-architect/gbilling2slack/invoice"
	"github.com/future-architect/gbilling2slack/notice"
	"github.com/future-architect/gbilling2slack/report"
	"log"
	"os"
)

func NotifyBilling(ctx context.Context, msg *pubsub.Message) error {

	var (
		projectID     = os.Getenv("GCP_PROJECT")
		tableName     = os.Getenv("TABLE_NAME")
		slackAPIToken = os.Getenv("SLACK_API_TOKEN")
		slackChannel  = os.Getenv("SLACK_CHANNEL")
	)

	if projectID == "" || tableName == "" || slackAPIToken == "" || slackChannel == "" {
		return fmt.Errorf("missing env")
	}

	inv := invoice.NewInvoice(projectID, tableName)

	monthBilling, err := inv.FetchBillingMonth(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	dayBilling, err := inv.FetchBillingDay(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	notifier := notice.NewSlackNotifier(slackAPIToken, slackChannel)

	// send summary report
	summaryReport := report.NewSummaryReport(monthBilling, dayBilling)
	parentTS, err := notifier.PostBilling(summaryReport)
	if err != nil {
		log.Println(err)
		return err
	}

	// send detail report
	detailReportList := report.NewDetailReportList(monthBilling, dayBilling)
	for _, detailReport := range detailReportList {
		if detailReport.MonthlyTotalCost == 0 {
			continue
		}
		if err := notifier.PostBillingThread(parentTS, detailReport); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
