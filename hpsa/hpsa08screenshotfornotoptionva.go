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
	"chromiumos/tast/local/chrome/uiauto/faillog"

	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/testing"
)

func init() {
	testing.AddTest(&testing.Test{
		Func:         Hpsa08screenshotfornotoptionva,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json", "profile.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa08screenshotfornotoptionva(ctx context.Context, s *testing.State) {
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
	var dashboardPath = s.DataPath("dashboard.json")
	common.CloseLastBrowser(ctx, topWindowName, s, bt, ui)
	// Do pretest after oobe
	common.PreTestWithNoOPT(ctx, s, bt, ui, path)

	// VirtualAgent test
	var vaClass, _, _ = common.GetJSONDashboard(common.VirtualAgent, dashboardPath)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.VirtualAgent, vaClass); err != nil {
		s.Fatalf("Failed to click to element  %v : %v ", common.VirtualAgent, err)
	}
	//GoBigSleepLint for va loading
	testing.Sleep(ctx, time.Minute)
	common.TakeScreenshot(ctx, s, "Hpsa08screenshot_vapopup.png", common.ScreenshotPath)
	// s.Fatal("Get ui dump")

}
