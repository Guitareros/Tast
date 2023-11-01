// Copyright 2021 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package hpsa

import (

	// Standard library packages
	"context"
	"fmt"
	"path/filepath"
	"strings"
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
	"chromiumos/tast/local/chrome/uiauto/role"

	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/testing"
)

func init() {
	testing.AddTest(&testing.Test{
		Func:         Hpsa03checksnpn,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa03checksnpn(ctx context.Context, s *testing.State) {
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
	defer faillog.DumpUITreeOnError(ctx, s.OutDir(), s.HasError, tconn)
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
	// var dashboardPath = s.DataPath("dashboard.json")

	//Do pretest after oobe
	common.PreTest(ctx, s, bt, ui, path)
	// var screenshotName string = "Tast_Test_Screenshot.png"
	// common.TakeScreenshot(ctx, s, screenshotName, common.ScreenshotPath)
	modelName, err := common.ReadFromVpd("model_name")
	if err != nil {
		s.Fatalf("Failed to get %v : %v", modelName, err)
	}
	// var deviceNameClass, _, _ = common.GetJSONDashboard(common.DeviceName, dashboardPath)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.DeviceName, bt)
	// name = nodewith.HasClass(deviceNameClass).Nth(0).Name
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", modelName, bt),
		ui.WaitUntilExists(nodewith.NameContaining(modelName).Role(role.StaticText)),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.DeviceName, bt, err)
	}
	serialNumber, err := common.ReadFromVpd("serial_number")
	if err != nil {
		s.Fatalf("Failed to get %v : %v", serialNumber, err)
	}
	// var deviceNameClass, _, _ = common.GetJSONDashboard(common.DeviceName, dashboardPath)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", serialNumber, bt)
	// name = nodewith.HasClass(deviceNameClass).Nth(0).Name
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", serialNumber, bt),
		ui.WaitUntilExists(nodewith.NameContaining(serialNumber).Role(role.StaticText)),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.SerialNumber, bt, err)
	}
	skuNumber, err := common.ReadFromVpd("sku_number")
	skuNumber = strings.Replace(skuNumber, "-", "#", 1)
	if err != nil {
		s.Fatalf("Failed to get %v : %v", skuNumber, err)
	}
	// var deviceNameClass, _, _ = common.GetJSONDashboard(common.DeviceName, dashboardPath)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", skuNumber, bt)
	// name = nodewith.HasClass(deviceNameClass).Nth(0).Name
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", skuNumber, bt),
		ui.WaitUntilExists(nodewith.NameContaining(skuNumber).Role(role.StaticText)),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.ProductNumber, bt, err)
	}
	// var dashboardPath = s.DataPath("dashboard.json")

}
