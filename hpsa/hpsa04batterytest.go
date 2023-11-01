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

	"go.chromium.org/luci/common/logging"
	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/testing"
)

func init() {
	testing.AddTest(&testing.Test{
		Func:         Hpsa04batterytest,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa04batterytest(ctx context.Context, s *testing.State) {
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
	defer faillog.DumpUITreeOnError(ctx, s.OutDir(), s.HasError, tconn)
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
	//Battery check screenshot
	var batteryCheckClass, _, _ = common.GetJSONDashboard(common.BatteryCheck, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.BatteryCheck, batteryCheckClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.BatteryCheck, err)
	}
	common.TakeScreenshot(ctx, s, "HPSA_hpsa04batterytest_batteryCheck.png", common.ScreenshotPath)
	var runBatteryCheck, _, _ = common.GetJSONDashboard(common.RunBatteryCheck, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.RunBatteryCheck, runBatteryCheck); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.RunBatteryCheck, err)
	}
	// Poll for a minute to make sure DUT connection is ready.
	if err := testing.Poll(ctx, func(ctx context.Context) error {
		if err := uiauto.Combine(
			fmt.Sprintf("Click the %v button in %v browser", nodewith.HasClass(runBatteryCheck), bt),
			ui.WaitUntilExists(nodewith.HasClass(runBatteryCheck).First()),
		)(ctx); err != nil {
			return err
		}
		return nil
	}, &testing.PollOptions{Interval: 3 * time.Minute,
		Timeout: time.Minute}); err != nil {
		logging.Infof(ctx, "Can not finish the action %v", err)
	}
	if err := testing.Poll(ctx, func(ctx context.Context) error {
		if err := common.FindException(ctx, ui, s, "hpsa04batterytest_Exception.png"); err != nil {
			return err
		}
		return nil
	}, &testing.PollOptions{Interval: 1 * time.Minute,
		Timeout: time.Minute}); err != nil {
		// s.Fatal("Can not finish the action: ", err)
	}
	if common.CheckExceptionFailed("hpsa04batterytest_Exception.png") {
		s.Fatal("Test failed, find the exception popup")
	}
	// s.Fatal("Geting the ui dump")
}
