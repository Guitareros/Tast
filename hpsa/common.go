// Copyright 2021 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package hpsa

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

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
		Func:         Common,
		LacrosStatus: testing.LacrosVariantExists,
		Desc:         "POC for HPSA Tast",
		Contacts:     []string{"xinyang.li@hp.com"},
		BugComponent: "",
		Attr:         []string{"group:mainline"},
		SoftwareDeps: []string{"chrome"},
	})
}

func Common(ctx context.Context, s *testing.State) {

	extDir := filepath.Dir("/var/chrome_extension_hpsa_itg/")

	extID, err := chrome.ComputeExtensionID(extDir)
	if err != nil {
		s.Fatalf("Failed to compute extension ID for %v: %v", extDir, err)
	}
	s.Log("Extension ID is ", extID)

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
	// TODO(crbug.com/1240344): Ensure the tablet mode is turned off until it is supported on Lacros.
	const tabletMode = false
	cleanup, err := ash.EnsureTabletModeEnabled(ctx, tconn, tabletMode)
	if err != nil {
		s.Fatalf("Failed to ensure the tablet mode is set to %v: %v", tabletMode, err)
	}
	defer cleanup(cleanupCtx)
	var topWindowName string
	switch bt {
	case browser.TypeAsh:
		topWindowName = "BrowserFrame"
	case browser.TypeLacros:
		topWindowName = "ExoShellSurface"
	default:
		s.Fatal("Unrecognized browser type: ", bt)
	}
	topLevelWindow := nodewith.Role(role.Window).HasClass(topWindowName)
	s.Logf("Opening a new tab in %v browser", bt)
	ui := uiauto.New(tconn)
	defer faillog.DumpUITreeOnError(cleanupCtx, s.OutDir(), s.HasError, tconn)
	//set up browser
	common.SetUpBrowser(ctx, ui, br, s, common.Language)
	s.Logf("Asserting that UI elements on browser window frame are accessible in %v browser", bt)
	for _, e := range []struct {
		name   string
		finder *nodewith.Finder
	}{
		{"Browser: New Tab", nodewith.HasClass("NewTabButton").Role(role.Button).Ancestor(topLevelWindow).First()},
		{"Browser: Tab Close", nodewith.HasClass("TabCloseButton").Role(role.Button).Ancestor(topLevelWindow).First()},
		{"Browser: Let's get Start", nodewith.HasClass("btn flex-row-center hp-button-primary").Role(role.Button)},
	} {
		if err = ui.WaitUntilExists(e.finder)(ctx); err != nil {
			s.Fatalf("Failed to find the UI element (%v) in %v: %v", e.name, bt, err)
		}
	}

	s.Logf("Asserting that the a11y node (rootWebArea) on the webview are accessible inside %v browser", bt)
	rootWebArea := nodewith.Role("rootWebArea").Ancestor(topLevelWindow).First()
	if err := ui.WaitUntilExists(rootWebArea)(ctx); err != nil {
		s.Fatalf("Failed to find the rootWebArea inside %v browser: %v", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "let's get start", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "let's get start", bt),
		ui.WaitUntilExists(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
		ui.LeftClick(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "let's get start", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "Launch HP Support Assistant", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "Launch HP Support Assistant", bt),
		ui.WaitUntilExists(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
		ui.LeftClick(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "Launch HP Support Assistant", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "Select region", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "Select region", bt),
		ui.WaitUntilExists(nodewith.HasClass("hp-dropdown-c").First()),
		ui.LeftClick(nodewith.HasClass("hp-dropdown-c").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "Select region", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "select_region", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "select_region", bt),
		ui.WaitUntilExists(nodewith.HasClass("menu dd-menu").First()),
		ui.FocusAndWait(nodewith.HasClass("menu dd-menu").First()),
		ui.LeftClick(nodewith.HasClass("text-line-1").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "select_region", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "Continue", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "Continue", bt),
		ui.WaitUntilExists(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
		ui.LeftClick(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "Continue", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "Don't show again", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "Don't show again", bt),
		ui.WaitUntilExists(nodewith.HasClass("cb hp-checkbox").First()),
		ui.LeftClick(nodewith.HasClass("cb hp-checkbox").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "Don't show again", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "Continue as Guest", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "Continue as Guest", bt),
		ui.WaitUntilExists(nodewith.HasClass("hp-link-underline ng-star-inserted").First()),
		ui.LeftClick(nodewith.HasClass("hp-link-underline ng-star-inserted").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "Continue as Guest", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "warranty option", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "warranty option", bt),
		ui.WaitUntilExists(nodewith.HasClass("cb hp-checkbox").Nth(1)),
		ui.LeftClick(nodewith.HasClass("cb hp-checkbox").Nth(1)),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "warranty option", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "usage data", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "usage data", bt),
		ui.WaitUntilExists(nodewith.HasClass("cb hp-checkbox").Nth(2)),
		ui.LeftClick(nodewith.HasClass("cb hp-checkbox").Nth(2)),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "usage data", bt, err)
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", "improve my experience", bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", "improve my experience", bt),
		ui.WaitUntilExists(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
		ui.LeftClick(nodewith.HasClass("btn flex-row-center hp-button-primary").First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", "improve my experience", bt, err)
	}
	s.Log("Sleep to wait for pin popup display")

}
