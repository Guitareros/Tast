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
		Func:         Hpsa09stresscpu,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json", "profile.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa09stresscpu(ctx context.Context, s *testing.State) {
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
	s.Log(ctx.Deadline())
	ctx, cancel := ctxutil.Shorten(ctx, 10*time.Second)
	// dl, ok := ctx.Deadline()
	// if !ok {
	// 	s.Fatal("Failed to set up deadline: ", err)
	// }
	// context.WithDeadline(ctx, dl.Add(5*time.Minute))
	s.Log(ctx.Deadline())
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
	var dashboardPath = s.DataPath("dashboard.json")
	common.CloseLastBrowser(ctx, topWindowName, s, bt, ui)
	// Do pretest after oobe
	common.PreTest(ctx, s, bt, ui, path)

	//Check CPU screenshot
	var checkCPUClass, _, _ = common.GetJSONDashboard(common.CheckCPU, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.CheckCPU, checkCPUClass); err != nil {
		s.Fatalf("Failed to click %v button : %v", common.CheckCPU, err)
	}
	common.TakeScreenshot(ctx, s, "Hpsa08screenshotfornotoption_checkCPU.png", common.ScreenshotPath)
	var runCPUCheckClass, _, _ = common.GetJSONDashboard(common.RunBatteryCheck, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.RunBatteryCheck, runCPUCheckClass); err != nil {
		s.Fatalf("Can not click %v : %v", common.RunBatteryCheck, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.CPUCheckCancel, bt)
	// var cpuCheckPassClass, _, _ = common.GetJSONDashboard(common.CPUCheckPassImage, dashboardPath)
	// if err := testing.Poll(ctx, func(ctx context.Context) error {
	// 	if err := uiauto.Combine(
	// 		fmt.Sprintf("Click the %v button in %v browser", common.CPUCheckPassImage, bt),
	// 		ui.WaitUntilExists(nodewith.HasClass(cpuCheckPassClass).Name("Passed").First()),
	// 		ui.LeftClick(nodewith.HasClass(cpuCheckPassClass).Name("Passed").First()),
	// 	)(ctx); err != nil {
	// 		s.Logf("Failed to find and click the %v button in %v: %v", common.CPUCheckPassImage, bt, err)
	// 		return err
	// 	}
	// 	return nil
	// }, &testing.PollOptions{Interval: 2 * time.Minute,
	// 	Timeout: 6 * time.Minute}); err != nil {
	// 	s.Log("Can not finish the action: ", err)
	// }
	//GoBigSleepLint wait for run CPU Check enable
	// testing.Sleep(ctx, 3*time.Minute)
	if err := testing.Poll(ctx, func(ctx context.Context) error {
		if err := uiauto.Combine(
			fmt.Sprintf("Click the %v button in %v browser", common.RunBatteryCheck, bt),
			ui.WaitUntilExists(nodewith.HasClass(runCPUCheckClass).Name("Run CPU check")),
			ui.LeftClick(nodewith.HasClass(runCPUCheckClass).Name("Run CPU check")),
		)(ctx); err != nil {
			s.Logf("Failed to find and click the %v button in %v: %v", common.RunBatteryCheck, bt, err)
			return err
		}
		return nil
	}, &testing.PollOptions{Interval: 10 * time.Second,
		Timeout: 3 * time.Minute}); err != nil {
		s.Fatal("Can not finish the action: ", err)
	}
	s.Log("All Complete")
	// testing.Sleep(ctx, 40*time.Second)
	// s.Fatal("Get ui dump")

}
