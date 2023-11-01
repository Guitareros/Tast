// Copyright 2023 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package common

import (
	"chromiumos/tast/local/chrome/uiauto"
	"chromiumos/tast/local/chrome/uiauto/nodewith"
	"chromiumos/tast/local/input"
	"context"
	"time"

	"go.chromium.org/tast/core/errors"
	"go.chromium.org/tast/core/testing"
)

// ScrollToElement scroll to the Element element
func ScrollToElement(ctx context.Context, s *testing.State, ui *uiauto.Context, scrollbarElement, targetElement *nodewith.Finder) error {
	//Display the network part
	mew, err := input.Mouse(ctx)
	if err != nil {
		s.Fatal("Failed to setup the mouse: ", err)
	}
	defer mew.Close(ctx)
	// move mouse to the collections container so that we can scroll the mouse.
	if err := uiauto.Combine("wait for collections and move mouse to collections area",
		ui.WaitUntilExists(scrollbarElement),
		ui.MouseMoveTo(targetElement, 10*time.Millisecond),
	)(ctx); err != nil {
		s.Fatal("Failed to load or move to collections: ", err)
	}
	return ScrollDownUntilSucceeds(ctx, SelectCollectionNode(ui, targetElement), mew)

}

// SelectCollectionNode is collection node
func SelectCollectionNode(ui *uiauto.Context, collectionNode *nodewith.Finder) uiauto.Action {
	return uiauto.Combine("select collection node",
		ui.WaitUntilExists(collectionNode),
		ui.MakeVisible(collectionNode),
		ui.DoDefault(collectionNode),
	)
}

// ScrollDownUntilSucceeds scrolls the mouse down until an action is achieved.
func ScrollDownUntilSucceeds(ctx context.Context, action uiauto.Action, mew *input.MouseEventWriter) error {
	const (
		maxNumSelectRetries = 4
		numScrolls          = 100
	)
	var actionErr error
	for i := 0; i < maxNumSelectRetries; i++ {
		if actionErr = action(ctx); actionErr == nil {
			return nil
		}
		for j := 0; j < numScrolls; j++ {
			if err := mew.ScrollDown(); err != nil {
				return errors.Wrap(err, "failed to scroll down")
			}
		}
	}
	return actionErr
}
