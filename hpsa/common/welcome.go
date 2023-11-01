// Copyright 2023 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package common

import (
	"chromiumos/tast/local/chrome/browser"
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/nodewith"
	"context"
	"fmt"
	"time"

	"go.chromium.org/tast/core/testing"
)

const (
	//Letsstart is the name of Letsstart
	Letsstart = "let's get start"
	//LaunchHPSupportAssistant is the name of LaunchHPSupportAssistant
	LaunchHPSupportAssistant = "Launch HP Support Assistant"
	//SelectRegion is the name of SelectRegion
	SelectRegion = "Select region"
	//DropMenu is the name of DropMenu
	DropMenu = "drop region menu"
	//SelectRegionUS is the name of SelectRegionUS
	SelectRegionUS = "selectregion"
	//ContinueBTN is name of the ContinueBTN
	ContinueBTN = "Continue"
	//DonotShowAgain is the name of DonotShowAgain
	DonotShowAgain = "Don't show again"
	//ContinueAsGuest is the name of ContinueAsGuest
	ContinueAsGuest = "Continue as Guest"
	//WarrantyOption is the name of WarrantyOption
	WarrantyOption = "warranty option"
	//UsageData is name of UsageData
	UsageData = "usage data"
	//ImproveMyExperience is the name of the ImproveMyExperience
	ImproveMyExperience = "improve my experience"
	//CreateAccount is the sign button
	CreateAccount = "CreateAccount"
	//Details is the link in data popup
	Details = "Details"
	//LetsShareLater is the link on data popup
	LetsShareLater = "LetsShareLater"
)

// ClickWelcomeBtns is using to click all element in welcome
func ClickWelcomeBtns(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, element, elementClass string) (string, error) {
	s.Logf("Asserting that mouse click works on the %v button in %v browser", element, bt)
	if err := testing.Poll(ctx, func(ctx context.Context) error {
		if err := uiauto.Combine(
			fmt.Sprintf("Click the %v button in %v browser", element, bt),
			ui.WaitUntilExists(nodewith.HasClass(elementClass).First()),
			ui.LeftClick(nodewith.HasClass(elementClass).First()),
		)(ctx); err != nil {
			s.Logf("Failed to find and click the %v button in %v: %v", element, bt, err)
			return err
		}
		return nil
	}, &testing.PollOptions{Interval: 1 * time.Minute,
		Timeout: time.Minute}); err != nil {
		s.Log("Can not finish the action: ", err)
	}
	return "Sucessfully clicked", nil
}

// ClickWelcomeBtnsNTH is using to click all element in welcome
func ClickWelcomeBtnsNTH(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, element, elementClass string, nth int) (string, error) {
	s.Logf("Asserting that mouse click works on the %v button in %v browser", element, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", element, bt),
		ui.WaitUntilExists(nodewith.HasClass(elementClass).Nth(nth)),
		ui.LeftClick(nodewith.HasClass(elementClass).Nth(nth)),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", element, bt, err)
		return "Failed to click let's get start", err
	}
	return "Sucessfully clicked", nil
}

// PreTest is a function to navigate to dashboard for HPSA
func PreTest(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, path string) (string, error) {
	letsstartclass, _, err := GetJSON(Letsstart, path)
	if err != nil {
		return "Can not get the json data for " + Letsstart, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, Letsstart, letsstartclass); err != nil {
		return tips, err
	}
	launchHPSupportAssistantclass, _, err := GetJSON(LaunchHPSupportAssistant, path)
	if err != nil {
		return "Can not get the json data for " + LaunchHPSupportAssistant, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, LaunchHPSupportAssistant, launchHPSupportAssistantclass); err != nil {
		return tips, err
	}
	selectRegionclass, _, err := GetJSON(SelectRegion, path)
	if err != nil {
		return "Can not get the json data for " + SelectRegion, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, SelectRegion, selectRegionclass); err != nil {
		return tips, err
	}
	if tips, err := selectRegion(ctx, s, bt, ui, path); err != nil {
		return tips, err
	}
	continueBTNclass, _, err := GetJSON(ContinueBTN, path)
	if err != nil {
		return "Can not get the json data for " + ContinueBTN, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ContinueBTN, continueBTNclass); err != nil {
		return tips, err
	}
	donotShowAgainclass, _, err := GetJSON(DonotShowAgain, path)
	if err != nil {
		return "Can not get the json data for " + DonotShowAgain, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, DonotShowAgain, donotShowAgainclass); err != nil {
		return tips, err
	}
	continueAsGuestclass, _, err := GetJSON(ContinueAsGuest, path)
	if err != nil {
		return "Can not get the json data for " + ContinueAsGuest, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ContinueAsGuest, continueAsGuestclass); err != nil {
		return tips, err
	}
	warrantyOptionclass, warrantyOptionnth, err := GetJSON(WarrantyOption, path)
	if err != nil {
		return "Can not get the json data for " + WarrantyOption, err
	}
	if tips, err := ClickWelcomeBtnsNTH(ctx, s, bt, ui, WarrantyOption, warrantyOptionclass, warrantyOptionnth); err != nil {
		return tips, err
	}
	usageDataclass, usageDatanth, err := GetJSON(UsageData, path)
	if err != nil {
		return "Can not get the json data for " + UsageData, err
	}
	if tips, err := ClickWelcomeBtnsNTH(ctx, s, bt, ui, UsageData, usageDataclass, usageDatanth); err != nil {
		return tips, err
	}
	improveMyExperienceclass, _, err := GetJSON(ImproveMyExperience, path)
	if err != nil {
		return "Can not get the json data for " + ImproveMyExperience, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ImproveMyExperience, improveMyExperienceclass); err != nil {
		return tips, err
	}

	closePinPopupclass, _, err := GetJSON(ClosePinPopup, path)
	if err != nil {
		return "Can not get the json data for " + ClosePinPopup, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ClosePinPopup, closePinPopupclass); err != nil {
		return tips, err
	}
	return "Successful navigate to dashboard", nil
}

// PreTestToSignin is a function to navigate to dashboard for HPSA
func PreTestToSignin(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, path string) (string, error) {
	letsstartclass, _, err := GetJSON(Letsstart, path)
	if err != nil {
		return "Can not get the json data for " + Letsstart, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, Letsstart, letsstartclass); err != nil {
		return tips, err
	}
	launchHPSupportAssistantclass, _, err := GetJSON(LaunchHPSupportAssistant, path)
	if err != nil {
		return "Can not get the json data for " + LaunchHPSupportAssistant, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, LaunchHPSupportAssistant, launchHPSupportAssistantclass); err != nil {
		return tips, err
	}
	selectRegionclass, _, err := GetJSON(SelectRegion, path)
	if err != nil {
		return "Can not get the json data for " + SelectRegion, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, SelectRegion, selectRegionclass); err != nil {
		return tips, err
	}
	if tips, err := selectRegion(ctx, s, bt, ui, path); err != nil {
		return tips, err
	}
	continueBTNclass, _, err := GetJSON(ContinueBTN, path)
	if err != nil {
		return "Can not get the json data for " + ContinueBTN, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ContinueBTN, continueBTNclass); err != nil {
		return tips, err
	}
	return "Successful navigate to welcome", nil
}

// PretestAfterSignin is the next test steps after sign in on welcome page
func PretestAfterSignin(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, path string) (string, error) {
	warrantyOptionclass, warrantyOptionnth, err := GetJSON(WarrantyOption, path)
	if err != nil {
		return "Can not get the json data for " + WarrantyOption, err
	}
	if tips, err := ClickWelcomeBtnsNTH(ctx, s, bt, ui, WarrantyOption, warrantyOptionclass, warrantyOptionnth); err != nil {
		return tips, err
	}
	usageDataclass, usageDatanth, err := GetJSON(UsageData, path)
	if err != nil {
		return "Can not get the json data for " + UsageData, err
	}
	if tips, err := ClickWelcomeBtnsNTH(ctx, s, bt, ui, UsageData, usageDataclass, usageDatanth); err != nil {
		return tips, err
	}
	improveMyExperienceclass, _, err := GetJSON(ImproveMyExperience, path)
	if err != nil {
		return "Can not get the json data for " + ImproveMyExperience, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ImproveMyExperience, improveMyExperienceclass); err != nil {
		return tips, err
	}

	closePinPopupclass, _, err := GetJSON(ClosePinPopup, path)
	if err != nil {
		return "Can not get the json data for " + ClosePinPopup, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ClosePinPopup, closePinPopupclass); err != nil {
		return tips, err
	}
	return "Successful navigate to dashboard", nil
}

// PreTestWithNoOPT is a function to navigate to dashboard for HPSA without select option in welcome
func PreTestWithNoOPT(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, path string) (string, error) {
	letsstartclass, _, err := GetJSON(Letsstart, path)
	if err != nil {
		return "Can not get the json data for " + Letsstart, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, Letsstart, letsstartclass); err != nil {
		return tips, err
	}
	launchHPSupportAssistantclass, _, err := GetJSON(LaunchHPSupportAssistant, path)
	if err != nil {
		return "Can not get the json data for " + LaunchHPSupportAssistant, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, LaunchHPSupportAssistant, launchHPSupportAssistantclass); err != nil {
		return tips, err
	}
	selectRegionclass, _, err := GetJSON(SelectRegion, path)
	if err != nil {
		return "Can not get the json data for " + SelectRegion, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, SelectRegion, selectRegionclass); err != nil {
		return tips, err
	}
	if tips, err := selectRegion(ctx, s, bt, ui, path); err != nil {
		return tips, err
	}
	continueBTNclass, _, err := GetJSON(ContinueBTN, path)
	if err != nil {
		return "Can not get the json data for " + ContinueBTN, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ContinueBTN, continueBTNclass); err != nil {
		return tips, err
	}
	donotShowAgainclass, _, err := GetJSON(DonotShowAgain, path)
	if err != nil {
		return "Can not get the json data for " + DonotShowAgain, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, DonotShowAgain, donotShowAgainclass); err != nil {
		return tips, err
	}
	continueAsGuestclass, _, err := GetJSON(ContinueAsGuest, path)
	if err != nil {
		return "Can not get the json data for " + ContinueAsGuest, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ContinueAsGuest, continueAsGuestclass); err != nil {
		return tips, err
	}

	detailClass, _, err := GetJSON(Details, path)
	if err != nil {
		return "Can not get the json data for " + Details, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, Details, detailClass); err != nil {
		return tips, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, Details, detailClass); err != nil {
		return tips, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, Details, detailClass); err != nil {
		return tips, err
	}
	letsharelaterClass, _, err := GetJSON(LetsShareLater, path)
	if err != nil {
		return "Can not get the json data for " + LetsShareLater, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, LetsShareLater, letsharelaterClass); err != nil {
		return tips, err
	}

	closePinPopupclass, _, err := GetJSON(ClosePinPopup, path)
	if err != nil {
		return "Can not get the json data for " + ClosePinPopup, err
	}
	if tips, err := ClickWelcomeBtns(ctx, s, bt, ui, ClosePinPopup, closePinPopupclass); err != nil {
		return tips, err
	}
	return "Successful navigate to dashboard", nil
}

func selectRegion(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, path string) (string, error) {
	var dropMenuclass, _, _ = GetJSON(DropMenu, path)
	var selectRegionUSclass, _, _ = GetJSON(SelectRegionUS, path)
	s.Logf("Asserting that mouse click works on the %v button in %v browser", SelectRegionUS, bt)
	if err := uiauto.Combine(
		fmt.Sprintf("Click the %v button in %v browser", SelectRegionUS, bt),
		ui.WaitUntilExists(nodewith.HasClass(dropMenuclass).First()),
		ui.FocusAndWait(nodewith.HasClass(dropMenuclass).First()),
		ui.LeftClick(nodewith.HasClass(selectRegionUSclass).First()),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the %v button in %v: %v", SelectRegionUS, bt, err)
		return "Failed to click select_region", err
	}
	return "", nil
}
