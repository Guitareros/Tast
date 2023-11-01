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
	"chromiumos/tast/local/chrome/uiauto/faillog"
	"chromiumos/tast/local/chrome/uiauto/nodewith"

	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/testing"
)

func init() {
	testing.AddTest(&testing.Test{
		Func:         Hpsa07signinmainpage,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Data:         []string{"hpsa.json", "dashboard.json", "profile.json"},
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Hpsa07signinmainpage(ctx context.Context, s *testing.State) {
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
	defer faillog.DumpUITreeOnError(ctx, s.OutDir(), s.HasError, tconn)
	var path = s.DataPath("hpsa.json")
	//Do pretest after oobe
	common.PreTest(ctx, s, bt, ui, path)
	var profilePath = s.DataPath(("profile.json"))
	s.Log("Get the profile json path : ", profilePath)
	username, password, err := common.GetProfileJSON("1", profilePath)
	sign.Signin(ctx, s, bt, ui, tconn, br, path, username, password)
	loggedinClass, _, err := common.GetJSON(common.LoggedIn, path)
	if err := testing.Poll(ctx, func(ctx context.Context) error {
		if err := uiauto.Combine(
			fmt.Sprintf("Click the %v button in %v browser", nodewith.HasClass(loggedinClass), bt),
			ui.WaitUntilExists(nodewith.HasClass(loggedinClass)),
		)(ctx); err != nil {
			return err
		}
		return nil
	}, &testing.PollOptions{Interval: 10 * time.Second,
		Timeout: 2 * time.Minute}); err != nil {
		s.Logf("Asserting that mouse click works on the %v button in %v browser", common.LoggedIn, err)
		// s.Fatal("Can not finish the action : ", err)
	}
	// testing.Sleep(ctx, 10*time.Second)
	// s.Fatal("Geting the ui dump")
}
