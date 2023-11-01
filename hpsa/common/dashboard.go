// Copyright 2023 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package common

import (
	"chromiumos/tast/local/chrome/browser"
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/nodewith"
	"chromiumos/tast/local/input"
	"context"
	"fmt"
	"time"

	"go.chromium.org/luci/common/logging"
	"go.chromium.org/tast/core/testing"
)

const (
	//ClosePinPopup is the name of ClosePinPopup
	ClosePinPopup = "Close pin popup"
	//Specifications is the name of Specifications
	Specifications = "Specifications"
	//CreateAccountOrSignIn is the name of CreateAccountOrSignIn
	CreateAccountOrSignIn = "createAccountOrSignIn"
	//UserName is the name of UserName
	UserName = "Username"
	//Profile is the name of Profile
	Profile = "Profile"
	//SignOut is the name of SignOut
	SignOut = "Sign out"
	//SignOutConfirm is the name of SignOutConfirm
	SignOutConfirm = "Yes in Sign out popup"
	//SpecificationsList is the name of SpecificationsList
	SpecificationsList = "SpecificationsList"
	//SpecificationsClose is the button to close the specification list
	SpecificationsClose = "SpecificationsClose"
	//WarrantyCard is the card on dashboard to show warranty
	WarrantyCard = "WarrantyCard"
	//WarrantyBack is the back button on warranty card
	WarrantyBack = "WarrantyBack"
	//AdditionalInformation is the link in the warranty card
	AdditionalInformation = "AdditionalInformation"
	//BatteryCheck is the diagnostic tool button
	BatteryCheck = "BatteryCheck"
	//BatteryCheckBack is the back button in battery check
	BatteryCheckBack = "BatteryCheckBack"
	//CheckCPU is the diagnostic tool button
	CheckCPU = "CheckCPU"
	//CheckCPUBack is the back button in cpu check
	CheckCPUBack = "CheckCPUBack"
	//CheckSystemMemory is diagonstic tool button
	CheckSystemMemory = "CheckSystemMemory"
	//CheckSystemMemoryBack is the back button in memory check
	CheckSystemMemoryBack = "CheckSystemMemoryBack"
	//CheckConnectivity is diagnostic tool button
	CheckConnectivity = "CheckConnectivity"
	//CheckConnectivityBack is the back button in connectivity check
	CheckConnectivityBack = "CheckConnectivityBack"
	//ComponentTest is diagnostic tool button
	ComponentTest = "ComponentTest"
	//ComponentTestBack is the back button in component test
	ComponentTestBack = "ComponentTestBack"
	//CheckStorage is diagnostic tool button
	CheckStorage = "CheckStorage"
	//CheckStorageBack is the back button in storage check
	CheckStorageBack = "CheckStorageBack"
	//Settings is the button on the menu bar
	Settings = "Settings"
	//AboutHPSA is the link in settings
	AboutHPSA = "AboutHPSA"
	//SeeAll is the button to naviage to support page
	SeeAll = "SeeAll"
	//Feedback is the button on menu bar
	Feedback = "Feedback"
	//OneStar is the star button on feedback
	OneStar = "OneStar"
	//TwoStars is the two stars button on feedback
	TwoStars = "TwoStars"
	//ThreeStars is the three stars button on feedback
	ThreeStars = "ThreeStars"
	//FourStars is the four stars button on feedback
	FourStars = "FourStars"
	//FiveStars is the five stars button on feedback
	FiveStars = "FiveStars"
	//FeedbackTextboxunselect is the text box for feedback popup
	FeedbackTextboxunselect = "FeedbackTextboxunselect"
	//FeedbackTextboxselect is the text bos for feedback popup with select status
	FeedbackTextboxselect = "FeedbackTextboxselect"
	//FeedbackLink is the HP privacy statement link
	FeedbackLink = "FeedbackLink"
	//FeedbackCancel is the cancel button on feedback
	FeedbackCancel = "FeedbackCancel"
	//Network is the name in specification list
	Network = "Network"
	//Audio is the name in specification list
	Audio = "Audio"
	//Battery is the name in specification list
	Battery = "Battery"
	//Video is the name in specification list
	Video = "Video"
	//DeviceName is the title of the device
	DeviceName = "DeviceName"
	//SerialNumber is the sn on dashboard
	SerialNumber = "SerialNumber"
	//ProductNumber is the pn on dashboard also sku
	ProductNumber = "ProductNumber"
	//RunBatteryCheck is the button to run battery check
	RunBatteryCheck = "RunBatteryCheck"
	//ExceptionBtn is the exception popup button
	ExceptionBtn = "ExceptionBtn"
	//RunBatteryCheckDisabled is the disabled status for run battery check
	RunBatteryCheckDisabled = "RunBatteryCheckDisabled"
	//CPUCheckCancel is cpu cancel button
	CPUCheckCancel = "CPUCheckCancel"
	//CPUCheckPassImage is the image for pass in cpu check
	CPUCheckPassImage = "CPUCheckPassImage"
	//LoggedIn is the element to verify log in status
	LoggedIn = "LoggedIn"
	//WarrantyCardGetDetail is the warranty card with out select option
	WarrantyCardGetDetail = "WarrantyCardGetDetail"
	//WarrantyCardGetDetailYES is the yes button in warranty popup
	WarrantyCardGetDetailYES = "WarrantyCardGetDetailYES"
	//VirtualAgent is the button of the VA pop
	VirtualAgent = "VirtualAgent"
	//VirtualAgentDown is the expend down button in va popup
	VirtualAgentDown = "VirtualAgentDown"
	//VirtualAgentUp is the expend up button in va popup
	VirtualAgentUp = "VirtualAgentUp"
	//VirtualAgentClose is the close button of va popup
	VirtualAgentClose = "VirtualAgentClose"
)

// ClickDashboardBtns is using to click all element in welcome
func ClickDashboardBtns(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, element, elementClass string) (string, error) {
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
	}, &testing.PollOptions{Timeout: 3 * time.Minute}); err != nil {
		s.Logf("Failed to find and click the %v button in 3 mins : %v", element, err)
	}
	return "Sucessfully clicked", nil
}

// ClickDashboardBtnsNTH is using to click all element in welcome
func ClickDashboardBtnsNTH(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, element, elementClass string, nth int) (string, error) {
	s.Logf("Asserting that mouse click works on the %v button in %v browser", element, bt)
	if err := testing.Poll(ctx, func(ctx context.Context) error {
		if err := uiauto.Combine(
			fmt.Sprintf("Click the %v button in %v browser", element, bt),
			ui.WaitUntilExists(nodewith.HasClass(elementClass).Nth(nth)),
			ui.LeftClick(nodewith.HasClass(elementClass).Nth(nth)),
		)(ctx); err != nil {
			s.Fatalf("Failed to find and click the %v button in %v: %v", element, bt, err)
			return err
		}
		return nil
	}, &testing.PollOptions{Timeout: 3 * time.Minute}); err != nil {
		s.Logf("Failed to find and click the %v button in 3 mins : %v", element, err)
	}
	return "Sucessfully clicked", nil
}

// InputDashboardText is using to click all element in welcome
func InputDashboardText(ctx context.Context, s *testing.State, bt browser.Type, ui *uiauto.Context, element, elementClass, inputContext string) (string, error) {
	s.Logf("Asserting that mouse click works on the %v button in %v browser", element, bt)
	kb, _ := input.Keyboard(ctx)
	// Poll for a minute to make sure DUT connection is ready.
	if err := testing.Poll(ctx, func(ctx context.Context) error {

		if err := uiauto.Combine(
			fmt.Sprintf("Click the %v button in %v browser", element, bt),
			ui.WaitUntilExists(nodewith.HasClass(elementClass).First()),
			ui.LeftClick(nodewith.HasClass(elementClass).First()),
			kb.TypeAction(inputContext),
		)(ctx); err != nil {
			return err
		}
		return nil
	}, &testing.PollOptions{Interval: 30 * time.Second,
		Timeout: time.Minute}); err != nil {
		logging.Infof(ctx, "Can not finish the action %v", err)
	}
	return "Sucessfully clicked", nil
}
