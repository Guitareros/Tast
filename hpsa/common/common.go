// Copyright 2023 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package common

import (
	"chromiumos/tast/local/apps"
	"chromiumos/tast/local/chrome"
	"chromiumos/tast/local/chrome/ash"
	"chromiumos/tast/local/chrome/browser"
	"chromiumos/tast/local/chrome/browser/browserfixt"
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/nodewith"
	"chromiumos/tast/local/chrome/uiauto/role"
	"chromiumos/tast/local/screenshot"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.chromium.org/tast/core/ctxutil"
	"go.chromium.org/tast/core/errors"
	"go.chromium.org/tast/core/testing"
)

const (
	//ExtensionDir is the path of HPSA
	ExtensionDir = "/var/chrome_extension_hpsa_itg/"
	//Proxy is using to test HP ITG environment
	Proxy = "--proxy-server=http://web-proxy.sgp.hp.com:8080"
	//Language is the ChromeOS language
	Language = "--lang=en-US"
	//AllLanguage is the list of the language

	//AppURLITG is the ITG URL for HPSA
	AppURLITG = "https://hpcs-appschr-itg.hpcloud.hp.com"
	//ScreenshotPath is the path to save the screenshot
	ScreenshotPath = "/var/hpsa_test_pictures/"
)

// AllLanguage is the language list of all test languages
var AllLanguage = [...]string{"ar-SA", "bg-BG", "cs-CZ", "da-DK", "de-DE", "el-GR", "en-US", "es-ES", "et-EE", "fi-FI", "fr-FR", "he-IL", "hr-HR", "hu-HU", "it-IT", "ja-JP", "ko-KR", "lt-LT", "lv-LV", "nb-NO", "nl-NL", "pl-PL", "pt-BR", "pt-PT", "ro-RO", "ru-RU", "sk-SK", "sl-SI", "sr-BA", "sv-SE", "th-TH", "tr-TR", "uk-UA", "zh-CN", "zh-HK", "zh-TW"}

// ManualInstallHPSA is a function to do below steps
// 1)Start Chrome
// 2)Launch browser and Navigate to the appURL
// 3)Install the Extension
func ManualInstallHPSA(ctx context.Context, tconn *chrome.TestConn, cr *chrome.Chrome, browserType browser.Type, appURL string) (string, error) {
	cleanupCtx := ctx
	ctx, cancel := ctxutil.Shorten(ctx, 5*time.Second)
	defer cancel()

	conn, br, closeBrowser, err := browserfixt.SetUpWithURL(ctx, cr, browserType, appURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to set up browser")
	}

	defer func(ctx context.Context) {
		closeBrowser(ctx)
		conn.Close()
	}(cleanupCtx)

	ui := uiauto.New(tconn).WithInterval(2 * time.Second)
	installIcon := nodewith.HasClass("PwaInstallView").Role(role.Button)
	installButton := nodewith.Name("Install").Role(role.Button)

	if err := testing.Poll(ctx, func(ctx context.Context) error {
		// Wait for longer time after second launch, since it can be delayed on low-end devices.
		if err := ui.WithTimeout(30 * time.Second).WaitUntilExists(installIcon)(ctx); err != nil {
			testing.ContextLog(ctx, "Install button is not shown initially. See b/230413572")
			testing.ContextLog(ctx, "Refresh page to enable installation")
			if reloadErr := br.ReloadActiveTab(ctx); reloadErr != nil {
				return testing.PollBreak(errors.Wrap(reloadErr, "failed to reload page"))
			}
			return err
		}
		return nil
	}, &testing.PollOptions{Timeout: 3 * time.Minute}); err != nil {
		return "", errors.Wrap(err, "failed to wait for Cursive to be installable")
	}

	if err := uiauto.Combine("",
		ui.LeftClick(installIcon),
		ui.LeftClick(installButton))(ctx); err != nil {
		return "", err
	}

	HPSAAppID, err := ash.WaitForChromeAppByNameInstalled(ctx, tconn, apps.HPSA.Name, 1*time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "failed to wait for installed app")
	}
	return HPSAAppID, nil
}

// TakeScreenshot is the function to take screenshot to test unit
func TakeScreenshot(ctx context.Context, s *testing.State, screenshotName, screenshotPath string) (string, error) {
	if find := strings.HasSuffix(screenshotName, ".png"); find {
		s.Log("Get screenshot name is ", screenshotName)
	} else {
		screenshotName += ".png"
		s.Log("Get screenshot name is ", screenshotName)
	}
	screenshotFile := filepath.Join(screenshotPath, screenshotName)
	s.Log("Save screenshot to ", screenshotFile)
	if err := screenshot.Capture(ctx, screenshotFile); err != nil {
		testing.ContextLog(ctx, "Failed to take screenshot: ", err)
		return "Take screenshot Fail", errors.Wrap(err, "failed to take screenshot")
	}
	return "", nil
}

// ReadFromVpd is the func to use command to get vpd
func ReadFromVpd(vpdField string) (string, error) {
	cmd := exec.Command("vpd", "-i", "RO_VPD", "-g", vpdField)
	VpdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(VpdOutput), nil
}

// FindException is the function to get the exception popup
func FindException(ctx context.Context, ui *uiauto.Context, s *testing.State, name string) error {
	var exceptionClass, _, _ = GetJSONDashboard(ExceptionBtn, s.DataPath("dashboard.json"))
	if err := uiauto.Combine("",
		ui.WaitUntilExists(nodewith.HasClass(exceptionClass).Role(role.Button)),
	)(ctx); err == nil {
		TakeScreenshot(ctx, s, name, ScreenshotPath)
		return errors.Wrap(err, "throw exception in app, please check the screenshot")
	}
	return nil
}

// CheckExceptionFailed is checking the exception file exist or not
func CheckExceptionFailed(filename string) bool {
	_, err := os.Lstat(filepath.Join(ScreenshotPath, filename))
	return !os.IsNotExist(err)
}

// SetUpBrowser is a function to add localstorage
func SetUpBrowser(ctx context.Context, ui *uiauto.Context, br *browser.Browser, s *testing.State, lang string) error {
	// Visit the page and create a history entry.
	conn, err := br.NewConn(ctx, AppURLITG)
	if err != nil {
		return errors.Wrap(err, "failed to open page")
	}
	defer conn.Close()
	// Set localStorage on the page.
	if err := conn.Call(ctx, nil, `(key, value) => localStorage.setItem(key, value)`, "HP_ENV", "pro"); err != nil {
		return errors.Wrap(err, "failed to set localStorage value HP_ENV")
	}
	if err := conn.Call(ctx, nil, `(key, value) => localStorage.setItem(key, value)`, "test_lang", lang); err != nil {
		return errors.Wrap(err, "failed to set localStorage value test_lang")
	}
	if err := conn.Call(ctx, nil, `(key, value) => localStorage.setItem(key, value)`, "HP_Disable_Firebase", "true"); err != nil {
		return errors.Wrap(err, "failed to set localStorage value HP_Disable_Firebase")
	}
	if err := conn.Call(ctx, nil, `(key, value) => localStorage.setItem(key, value)`, "isFullDebug", "true"); err != nil {
		return errors.Wrap(err, "failed to set localStorage value isFullDebug")
	}
	if err := conn.Call(ctx, nil, `(key, value) => localStorage.setItem(key, value)`, "HP_Survey", "false"); err != nil {
		return errors.Wrap(err, "failed to set localStorage value HP_Survey")
	}
	if err := conn.Call(ctx, nil, `(key, value) => localStorage.setItem(key, value)`, "HP_Survey_Delay", "5000"); err != nil {
		return errors.Wrap(err, "failed to set localStorage value HP_Survey_Delay")
	}
	return nil
}

// CloseLastBrowser is the func to close the browser
func CloseLastBrowser(ctx context.Context, topWindowName string, s *testing.State, bt browser.Type, ui *uiauto.Context) (string, error) {
	topLevelWindow := nodewith.Role(role.Window).HasClass(topWindowName).First()
	closeButton := nodewith.HasClass("FrameCaptionButton").Name("Close").Role(role.Button).Ancestor(topLevelWindow).First()
	if err := uiauto.Combine(
		fmt.Sprintf("Click the close button in %v browser", bt),
		ui.WaitUntilExists(closeButton),
		ui.LeftClick(closeButton),
	)(ctx); err != nil {
		s.Fatalf("Failed to find and click the close button in %v: %v", bt, err)
		return "Close browser Fail", errors.Wrap(err, "failed to close browser")
	}
	return "", nil
	//Check close the browser
	//Launch HPSA extension
	// func LaunchHPSAExtension(ctx context.Context, tconn *chrome.TestConn, cr *chrome.Chrome, appID, appURL string)

}
