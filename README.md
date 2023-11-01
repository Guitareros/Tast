# HPSA_Chromebook_AutoTest_TAST
ChromeOS auto test code using the new TAST framework and GO Language

It's a project for HPSANEO test cases which use chrome tast.

How to use:
1. Copy imports.go to {$CHROMIUMOS}/src/platform/tast-tests/src/chromiumos/tast/local/bundles/cros
2. Copy apps.go to {$CHROMIUMOS}/src/platform/tast-tests/src/chromiumos/tast/local/apps
3. Copy hpsa folder to {$CHROMIUMOS}/src/platform/tast-tests/src/chromiumos/tast/local/bundles/cros
4. Copy common to {$CHROMIUMOS}/src/platform/tast-tests/src/chromiumos/tast/
5. Use terminal to chroot environment
6. Use the commond tast -verbose run <ip> <testfolder>.<Testcase> e.g. tast -verbose run <ip> hpsa.Walkthrough01
