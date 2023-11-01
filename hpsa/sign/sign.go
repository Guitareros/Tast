// Copyright 2023 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sign

import (
	"chromiumos/tast/local/bundles/cros/hpsa/common"
	"chromiumos/tast/local/chrome"
	"chromiumos/tast/local/chrome/browser"
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/nodewith"
	"chromiumos/tast/local/input"
	"context"
	"fmt"
	"time"

	"go.chromium.org/tast/core/testing"
)

// Signin is a function to send username and password for HPID
func Signin(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, tconn *chrome.TestConn, br *browser.Browser, path, username, password string) (string, error) {
	var createAccountOrSignInclass, _, _ = common.GetJSON(common.CreateAccountOrSignIn, path)

	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.CreateAccountOrSignIn, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.CreateAccountOrSignIn, bt),
		ui.WaitUntilExists(nodewith.HasClass(createAccountOrSignInclass).First()),
		ui.LeftClick(nodewith.HasClass(createAccountOrSignInclass).First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.CreateAccountOrSignIn, bt, err)
		return "Failed to click sign in button", err
	}
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.UserName, bt)
	kb, _ := input.Keyboard(ctx)
	// GoBigSleepLint Wait for load to sign in page
	testing.Sleep(ctx, 20*time.Second)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.UserName, bt),
		kb.TypeAction(username),
		kb.AccelAction("Enter"),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.UserName, bt, err)
		return "Failed to click warranty option", err
	}
	// GoBigSleepLint Wait for navigate to password page
	testing.Sleep(ctx, 5*time.Second)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.UserName, bt),
		kb.TypeAction(password),
		kb.AccelAction("Enter"),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", common.UserName, bt, err)
		return "Failed to click warranty option", err
	}
	// GoBigSleepLint Wait for finish sign in
	testing.Sleep(ctx, 10*time.Second)
	defer kb.Close(ctx)
	return "Successful sign in account ", nil
}

// Signout is the function for sign out in HPSA
func Signout(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, tconn *chrome.TestConn, br *browser.Browser, path string) error {
	var profileclass, _, _ = common.GetJSON(common.Profile, path)
	var profileElement = nodewith.HasClass(profileclass).First()
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.Profile, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.Profile, bt),
		ui.WaitUntilExists(profileElement),
		ui.LeftClick(profileElement),
	)(ctx); err != nil {
		return err
	}
	var signOutclass, signOutNTH, _ = common.GetJSON(common.SignOut, path)

	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.SignOut, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.SignOut, bt),
		ui.WaitUntilExists(nodewith.HasClass(signOutclass).Nth(signOutNTH)),
		ui.LeftClick(nodewith.HasClass(signOutclass).Nth(signOutNTH)),
	)(ctx); err != nil {
		return err
	}
	var signOutConfirmclass, _, _ = common.GetJSON(common.SignOutConfirm, path)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", common.SignOutConfirm, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", common.SignOutConfirm, bt),
		ui.WaitUntilExists(nodewith.HasClass(signOutConfirmclass).First()),
		ui.LeftClick(nodewith.HasClass(signOutConfirmclass).First()),
	)(ctx); err != nil {
		return err
	}
	return nil
}
