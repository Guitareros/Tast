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
	"chromiumos/tast/local/bundles/cros/hpsa/sign"
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
		Func:         Smokeextension,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "profile.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Smokeextension(ctx context.Context, s *testing.State) {
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
	br, closeBrowser, err := browserfixt.SetUp(ctx, cr, browser.TypeAsh)
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
	var screenshotName string = "Tast_Test_Screenshot.png"
	common.TakeScreenshot(ctx, s, screenshotName, common.ScreenshotPath)
	var profilePath = s.DataPath(("profile.json"))
	s.Log("Get the profile json path : ", profilePath)
	username, password, err := common.GetProfileJSON("1", profilePath)
	if err != nil {
		s.Fatal("Failed to find json: ", err)
	}
	// s.Logf("Find the profile %v, %v", username, password)
	sign.Signin(ctx, s, bt, ui, tconn, br, path, username, password)
	var specificationsclass, _, _ = common.GetJSON(common.Specifications, path)
	var networkclass, networknth, _ = common.GetJSON(common.Network, path)
	var specificationsListclass, _, _ = common.GetJSON(common.SpecificationsList, path)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.Specifications, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.Specifications, bt),
		ui.WaitUntilExists(nodewith.HasClass(specificationsclass).First()),
		ui.LeftClick(nodewith.HasClass(specificationsclass).First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.Specifications, bt, err)
	}

	networtElement := nodewith.HasClass(networkclass).Nth(networknth)
	specificationsList := nodewith.HasClass(specificationsListclass).First()
	common.ScrollToElement(ctx, s, ui, specificationsList, networtElement)

	var specificationsCloseclass, _, _ = common.GetJSON(common.SpecificationsClose, path)
	if _, err := common.ClickDashboardBtns(ctx, s, bt, ui, common.SpecificationsClose, specificationsCloseclass); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.SpecificationsClose, bt, err)
	}
	sign.Signout(ctx, s, bt, ui, tconn, br, path)

}
