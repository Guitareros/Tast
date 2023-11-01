// Copyright 2021 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package hpsa

import (

	// Standard library packages
	"context"
	"path/filepath"
	"time"

	//chromiumos/ packages
	"chromiumos/tast/local/bundles/cros/hpsa/common"
	"chromiumos/tast/local/chrome"
	"chromiumos/tast/local/chrome/ash"
	"chromiumos/tast/local/chrome/browser"
	"chromiumos/tast/local/chrome/browser/browserfixt"
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/nodewith"

	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/testing"
)

func init() {
	testing.AddTest(&testing.Test{
		Func:         Hpsa01walkthrough,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa01walkthrough(ctx context.Context, s *testing.State) {
	//Need copy the file to the path
	extDir := filepath.Dir(common.ExtensionDir)
	extID, err := chrome.ComputeExtensionID(extDir)
	if err != nil {
		s.Fatalf("Failed to compute extension ID for %v: %v", extDir, err)
	}
	s.Log("Extension ID is ", extID)
	//Create the chrome with the extra arguments
	cr, err := chrome.New(ctx, chrome.UnpackedExtension(extDir),
		chrome.ExtraArgs(common.Proxy),
		chrome.ExtraArgs(common.Language),
	)
	if err != nil {
		s.Fatal("Chrome login failed: ", err)
	}
	defer cr.Close(ctx)

	bt := browser.TypeAsh
	// Reserve ten seconds for cleanup.
	cleanupCtx := ctx
	ctx, cancel := ctxutil.Shorten(ctx, 10*time.Second)
	defer cancel()
	_, closeBrowser, err := browserfixt.SetUp(ctx, cr, browser.TypeAsh)
	// br, closeBrowser, err := browserfixt.SetUp(ctx, cr, browser.TypeAsh)
	if err != nil {
		s.Fatal("Failed to set up browser: ", err)
	}
	defer closeBrowser(cleanupCtx)
	tconn, err := cr.TestAPIConn(ctx)
	if err != nil {
		s.Fatal("Failed to create Test API connection: ", err)
	}
	const tabletMode = false
	cleanup, err := ash.EnsureTabletModeEnabled(ctx, tconn, tabletMode)
	if err != nil {
		s.Fatalf("Failed to ensure the tablet mode is set to %v: %v", tabletMode, err)
	}
	defer cleanup(cleanupCtx)
	ui := uiauto.New(tconn)
	_, err = common.ManualInstallHPSA(ctx, tconn, cr, bt, common.AppURLITG)
	if err != nil {
		s.Fatal("Failed to manually install HPSA: ", err)
	}
	var path = s.DataPath("hpsa.json")
	//Do pretest after oobe
	common.PreTest(ctx, s, bt, ui, path)
	// var screenshotName string = "Tast_Test_Screenshot.png"
	// common.TakeScreenshot(ctx, s, screenshotName, common.ScreenshotPath)
	var dashboardPath = s.DataPath("dashboard.json")
	var warrantyCardClass, _, _ = common.GetJSONDashboard(common.WarrantyCard, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.WarrantyCard, warrantyCardClass); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.WarrantyCard, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_warrantyCard.png", common.ScreenshotPath)
	var additionalInformationClass, _, _ = common.GetJSONDashboard(common.AdditionalInformation, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.AdditionalInformation, additionalInformationClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.AdditionalInformation, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_additionalInformation.png", common.ScreenshotPath)
	var warrantyBackClass, _, _ = common.GetJSONDashboard(common.WarrantyBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.WarrantyBack, warrantyBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.WarrantyBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_closeWarrantyCard.png", common.ScreenshotPath)
	//System memory check screenshot
	var checkSystemMemoryClass, checkSystemMemorynth, _ = common.GetJSONDashboard(common.CheckSystemMemory, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.CheckSystemMemory, checkSystemMemoryClass, checkSystemMemorynth); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckSystemMemory, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkSystemMemory.png", common.ScreenshotPath)
	var checkSystemMemoryBackClass, _, _ = common.GetJSONDashboard(common.CheckSystemMemoryBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.CheckSystemMemoryBack, checkSystemMemoryBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckSystemMemoryBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkSystemMemoryClose.png", common.ScreenshotPath)
	//Battery check screenshot
	var batteryCheckClass, _, _ = common.GetJSONDashboard(common.BatteryCheck, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.BatteryCheck, batteryCheckClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.BatteryCheck, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_batteryCheck.png", common.ScreenshotPath)
	var batteryCheckBackClass, _, _ = common.GetJSONDashboard(common.BatteryCheckBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.BatteryCheckBack, batteryCheckBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.BatteryCheckBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_batteryCheckClose.png", common.ScreenshotPath)
	//Check component screenshot
	var componentTestClass, componentTestnth, _ = common.GetJSONDashboard(common.ComponentTest, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.ComponentTest, componentTestClass, componentTestnth); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.ComponentTest, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_component.png", common.ScreenshotPath)
	var componentTestBackClass, _, _ = common.GetJSONDashboard(common.ComponentTestBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.ComponentTestBack, componentTestBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.ComponentTestBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_componentClose.png", common.ScreenshotPath)

	//Storage check screenshot
	var checkStorageClass, checkStoragenth, _ = common.GetJSONDashboard(common.CheckStorage, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.CheckStorage, checkStorageClass, checkStoragenth); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckStorage, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkStorage.png", common.ScreenshotPath)
	var checkStorageBackClass, _, _ = common.GetJSONDashboard(common.CheckStorageBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.CheckStorageBack, checkStorageBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckStorageBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkStorageClose.png", common.ScreenshotPath)
	//Check CPU screenshot
	var checkCPUClass, _, _ = common.GetJSONDashboard(common.CheckCPU, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.CheckCPU, checkCPUClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckCPU, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkCPU.png", common.ScreenshotPath)
	var checkCPUBackClass, _, _ = common.GetJSONDashboard(common.CheckCPUBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.CheckCPUBack, checkCPUBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckCPUBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkCPUClose.png", common.ScreenshotPath)
	//Check connectivity screenshot
	var checkConnectivityClass, checkConnectivitynth, _ = common.GetJSONDashboard(common.CheckConnectivity, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.CheckConnectivity, checkConnectivityClass, checkConnectivitynth); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckConnectivity, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkConnectivity.png", common.ScreenshotPath)
	var checkConnectivityBackClass, _, _ = common.GetJSONDashboard(common.CheckConnectivityBack, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.CheckConnectivityBack, checkConnectivityBackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckConnectivityBack, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_checkConnectivityClose.png", common.ScreenshotPath)

	var settingsClass, _, _ = common.GetJSONDashboard(common.Settings, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.Settings, settingsClass); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.Settings, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_settings.png", common.ScreenshotPath)
	var aboutHPSAClass, aboutHPSAnth, _ = common.GetJSONDashboard(common.AboutHPSA, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.AboutHPSA, aboutHPSAClass, aboutHPSAnth); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.AboutHPSA, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_about.png", common.ScreenshotPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.Settings, settingsClass); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.Settings, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_settingsClose.png", common.ScreenshotPath)
	var seeAllClass, _, _ = common.GetJSONDashboard(common.SeeAll, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.SeeAll, seeAllClass); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.SeeAll, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_support.png", common.ScreenshotPath)
	var feedbackClass, _, _ = common.GetJSONDashboard(common.Feedback, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.Feedback, feedbackClass); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.Feedback, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedback.png", common.ScreenshotPath)
	var oneStarClass, oneStarNTH, _ = common.GetJSONDashboard(common.OneStar, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.OneStar, oneStarClass, oneStarNTH); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.OneStar, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackOneStar.png", common.ScreenshotPath)
	var twoStarsClass, twoStarsNTH, _ = common.GetJSONDashboard(common.TwoStars, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.TwoStars, twoStarsClass, twoStarsNTH); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.TwoStars, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackTwoStar.png", common.ScreenshotPath)
	var threeStarsClass, threeStarsNTH, _ = common.GetJSONDashboard(common.ThreeStars, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.ThreeStars, threeStarsClass, threeStarsNTH); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.ThreeStars, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackThreeStar.png", common.ScreenshotPath)
	var fourStarsClass, fourStarsNTH, _ = common.GetJSONDashboard(common.FourStars, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.FourStars, fourStarsClass, fourStarsNTH); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.FourStars, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackFourStar.png", common.ScreenshotPath)
	var fiveStarsClass, fiveStarsNTH, _ = common.GetJSONDashboard(common.FiveStars, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.FiveStars, fiveStarsClass, fiveStarsNTH); err != nil {
		s.Fatalf("Failed to click %v button : %v ", common.FiveStars, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackFiveStar.png", common.ScreenshotPath)
	//Remove input test
	// var feedbackTextboxUnSelectClass, _, _ = common.GetJSONDashboard(common.FeedbackTextboxunselect, dashboardPath)
	// var feedbackTextboxSelectClass, _, _ = common.GetJSONDashboard(common.FeedbackTextboxselect, dashboardPath)
	// var inputContext string = "When you play the game of thrones, you win or you die. There is no middle ground.When you play the game of thrones, you win or you die. There is no middle ground.When you play the game of thrones, you win or you die. There is no middle ground.When you play the game of thrones, you win or you die. There is no middle ground."
	// if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.FeedbackTextboxunselect, feedbackTextboxUnSelectClass); err != nil {
	// 	s.Fatalf("Failed to click %v button : %v ", common.FeedbackTextboxunselect, err)
	// }
	// if _, err := common.InputDashboardText(ctx, s, bt, ui, common.FeedbackTextboxunselect, feedbackTextboxUnSelectClass, inputContext); err != nil {
	// 	s.Fatalf("Failed to input context to  %v  : %v ", common.FeedbackTextboxselect, err)
	// }
	// common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackInput.png", common.ScreenshotPath)
	var feedbackCancelClass, feedbackCancelNTH, _ = common.GetJSONDashboard(common.FeedbackCancel, dashboardPath)
	if _, err := common.ClickDashboardBtnsNTH(ctx, s, bt, ui, common.FeedbackCancel, feedbackCancelClass, feedbackCancelNTH); err != nil {
		s.Fatalf("Failed to click  %v button : %v ", common.FeedbackCancel, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_feedbackClose.png", common.ScreenshotPath)
	var specificationsclass, _, _ = common.GetJSON(common.Specifications, path)
	var networkclass, networknth, _ = common.GetJSON(common.Network, path)
	var specificationsListclass, _, _ = common.GetJSON(common.SpecificationsList, path)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.Specifications, specificationsclass); err != nil {
		s.Fatalf("Failed to click  %v button : %v ", common.FeedbackCancel, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_specifications.png", common.ScreenshotPath)
	networtElement := nodewith.HasClass(networkclass).Nth(networknth)
	specificationsList := nodewith.HasClass(specificationsListclass).First()
	if err := common.ScrollToElement(ctx, s, ui, specificationsList, networtElement); err != nil {
		s.Fatalf("Failed to scroll to element  %v : %v ", common.Audio, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa01walkthrough_scrollToNetWork.png", common.ScreenshotPath)
}
