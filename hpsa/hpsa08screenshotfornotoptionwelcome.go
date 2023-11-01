// Copyright 2021 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package hpsa

import (

	// Standard library packages
	"context"
	"fmt"
	"path/filepath"
	"time"

	//chromiumos/ packages
	"chromiumos/tast/local/bundles/cros/hpsa/common"
	"chromiumos/tast/local/chrome"
	"chromiumos/tast/local/chrome/ash"
	"chromiumos/tast/local/chrome/browser"
	"chromiumos/tast/local/chrome/browser/browserfixt"
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/faillog"
	"chromiumos/tast/local/chrome/uiauto/nodewith"

	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/testing"
)

func init() {
	testing.AddTest(&testing.Test{
		Func:         Hpsa08screenshotfornotoptionwelcome,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json", "profile.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa08screenshotfornotoptionwelcome(ctx context.Context, s *testing.State) {
	// for _, language := range common.AllLanguage {
	// hpsaSteps(ctx, s, language)
	//Need copy the file to the path
	extDir := filepath.Dir(common.ExtensionDir)
	extID, err := chrome.ComputeExtensionID(extDir)
	if err != nil {
		s.Fatalf("Failed to compute extension ID for %v: %v", extDir, err)
	}
	s.Log("Extension ID is ", extID)
	// for _, language := range common.AllLanguage {
	// var languageSet = "--lang" + language
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
	br, closeBrowser, err := browserfixt.SetUp(ctx, cr, browser.TypeAsh)

	// br, closeBrowser, err := browserfixt.SetUp(ctx, cr, browser.TypeAsh)
	if err != nil {
		s.Fatal("Failed to set up browser: ", err)
	}
	defer closeBrowser(cleanupCtx)
	tconn, err := cr.TestAPIConn(ctx)
	if err != nil {
		s.Fatal("Failed to create Test API connection: ", err)
	}
	ui := uiauto.New(tconn)
	common.SetUpBrowser(ctx, ui, br, s, common.Language)
	var topWindowName string
	switch bt {
	case browser.TypeAsh:
		topWindowName = "BrowserFrame"
	case browser.TypeLacros:
		topWindowName = "ExoShellSurface"
	default:
		s.Fatal("Unrecognized browser type: ", bt)
	}
	const tabletMode = false
	cleanup, err := ash.EnsureTabletModeEnabled(ctx, tconn, tabletMode)
	if err != nil {
		s.Fatalf("Failed to ensure the tablet mode is set to %v: %v", tabletMode, err)
	}
	defer cleanup(cleanupCtx)
	_, err = common.ManualInstallHPSA(ctx, tconn, cr, bt, common.AppURLITG)
	if err != nil {
		s.Fatal("Failed to manually install HPSA: ", err)
	}
	defer faillog.DumpUITreeOnError(ctx, s.OutDir(), s.HasError, tconn)
	var path = s.DataPath("hpsa.json")
	// var dashboardPath = s.DataPath("dashboard.json")
	common.CloseLastBrowser(ctx, topWindowName, s, bt, ui)
	// Do pretest after oobe
	// common.PreTestWithNoOPT(ctx, s, bt, ui, path)
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_LetsStart.png", common.ScreenshotPath)
	letsstartclass, _, err := common.GetJSON(common.Letsstart, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.Letsstart, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.Letsstart, letsstartclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_LaunchHPSA.png", common.ScreenshotPath)
	launchHPSupportAssistantclass, _, err := common.GetJSON(common.LaunchHPSupportAssistant, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.LaunchHPSupportAssistant, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_Welcome.png", common.ScreenshotPath)
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.LaunchHPSupportAssistant, launchHPSupportAssistantclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_SelectRegion.png", common.ScreenshotPath)
	selectRegionclass, _, err := common.GetJSON(common.SelectRegion, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.SelectRegion, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.SelectRegion, selectRegionclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_RegionDrop.png", common.ScreenshotPath)
	var dropMenuclass, _, _ = common.GetJSON(common.DropMenu, path)
	var selectRegionUSclass, _, _ = common.GetJSON(common.SelectRegionUS, path)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.SelectRegionUS, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.SelectRegionUS, bt),
		ui.WaitUntilExists(nodewith.HasClass(dropMenuclass).First()),
		ui.FocusAndWait(nodewith.HasClass(dropMenuclass).First()),
		ui.LeftClick(nodewith.HasClass(selectRegionUSclass).First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.SelectRegionUS, bt, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_SelectUS.png", common.ScreenshotPath)
	continueBTNclass, _, err := common.GetJSON(common.ContinueBTN, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.ContinueBTN, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.ContinueBTN, continueBTNclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_Continue.png", common.ScreenshotPath)
	donotShowAgainclass, _, err := common.GetJSON(common.DonotShowAgain, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.DonotShowAgain, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.DonotShowAgain, donotShowAgainclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_DonotShowagain.png", common.ScreenshotPath)
	continueAsGuestclass, _, err := common.GetJSON(common.ContinueAsGuest, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.ContinueAsGuest, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.ContinueAsGuest, continueAsGuestclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_ContinueAsGuest.png", common.ScreenshotPath)
	detailClass, _, err := common.GetJSON(common.Details, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.Details, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.Details, detailClass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_Detail.png", common.ScreenshotPath)
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.Details, detailClass); err != nil {
		s.Fatal(tips, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.Details, detailClass); err != nil {
		s.Fatal(tips, err)
	}

	letsharelaterClass, _, err := common.GetJSON(common.LetsShareLater, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.LetsShareLater, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_LetShareLater.png", common.ScreenshotPath)
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.LetsShareLater, letsharelaterClass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_Pinpopup.png", common.ScreenshotPath)
	closePinPopupclass, _, err := common.GetJSON(common.ClosePinPopup, path)
	if err != nil {
		s.Fatal("Can not get the json data for "+common.ClosePinPopup, err)
	}
	if tips, err := common.ClickWelcomeBtns(ctx, s, bt, ui, common.ClosePinPopup, closePinPopupclass); err != nil {
		s.Fatal(tips, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_Dashboard.png", common.ScreenshotPath)
}
