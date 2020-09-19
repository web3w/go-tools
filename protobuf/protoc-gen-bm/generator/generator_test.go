package generator

import (
	"os"
	"os/exec"
	"testing"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func TestGenerateParseCommandLineParamsError(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		g := &bm{}
		g.Generate(&plugin.CodeGeneratorRequest{
			Parameter: proto.String("invalid"),
		})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGenerateParseCommandLineParamsError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

type Permission struct {
	Name        string
	Url         string
	Description string
}

var Perms = []Permission{
	Permission{"PermissionAsSiteRegist", "/regist", ""},
	Permission{"PermissionAsSiteUserActive", "/regist/active", ""},
	Permission{"PermissionAsSiteLogin", "/login", ""},
	Permission{"PermissionAsSiteLogout", "/logout", ""},
	Permission{"PermissionAsSiteUserInvite", "/user/invite", ""},
	Permission{"PermissionAsSiteApplyResetOtherPass", "/user/other/pass/reset", ""},
	Permission{"PermissionAsSiteApplyResetOther2FA", "/user/other/2fa/reset", ""},
	Permission{"PermissionAsSiteUpdateOtherUser", "/user/other/update", ""},
	Permission{"PermissionAsSiteResetMePass", "/user/me/pass/reset", ""},
	Permission{"PermissionAsSiteResetMe2FA", "/user/me/2fa/reset", ""},
	Permission{"PermissionAsSiteGetMe", "/user/me", ""},
	Permission{"PermissionAsSiteUpdateMe", "/user/me/update", ""},
	Permission{"PermissionAsSiteVerifyMePass", "/user/me/pass/verify", ""},
	Permission{"PermissionAsSiteListUser", "/user/list", ""},
	Permission{"PermissionAsSiteListRoleAndPermission", "/user/roleandpermission/list", ""},
	Permission{"PermissionAsSiteAddApiUser", "/user/apiuser/add", ""},
	Permission{"PermissionAsSiteListApiUser", "/user/apiuser/list", ""},
	Permission{"PermissionAsSiteUpdateApiUser", "/user/apiuser/update", ""},
	Permission{"PermissionAsSiteDeleteApiUser", "/user/apiuser/del", ""},
	Permission{"PermissionAsSiteAddWallet", "/wallet/add", ""},
	Permission{"PermissionAsSiteUpdateWallet", "/wallet/update", ""},
	Permission{"PermissionAsSiteRemoveWallet", "/wallet/delete", ""},
	Permission{"PermissionAsSiteListWallet", "/wallet/list", ""},
	Permission{"PermissionAsSiteListCoinInfo", "/coin/list", ""},
	Permission{"PermissionAsSiteAddWalletCoin", "/wallet/coin/add", ""},
	Permission{"PermissionAsSiteRemoveWalletCoin", "/wallet/coin/delete", ""},
	Permission{"PermissionAsSiteListWalletCoin", "/wallet/coin/list", ""},
	Permission{"PermissionAsSiteGetNewAddress", "/wallet/coin/address/add", ""},
	Permission{"PermissionAsSiteHideAddress", "/wallet/coin/address/hide", ""},
	Permission{"PermissionAsSiteListAddress", "/wallet/coin/address/list", ""},
	Permission{"PermissionAsSiteHasAddress", "/wallet/coin/address/exist", ""},
	Permission{"PermissionAsSiteCheckAddress", "/wallet/coin/address/check", ""},
	Permission{"PermissionAsSiteGetCoinFee", "/wallet/coin/fee", ""},
	Permission{"PermissionAsSiteListDepositTx", "/wallet/coin/deposit/tx/list", ""},
	Permission{"PermissionAsSiteListWithdrawTx", "/wallet/coin/withdraw/tx/list", ""},
	Permission{"PermissionAsSiteNewWithdraw", "/wallet/coin/withdraw/add", ""},
	Permission{"PermissionAsSiteGetWithdrawDetail", "/wallet/coin/withdraw/tx/get", ""},
	Permission{"PermissionAsSiteGetWalletAsset", "/wallet/asset/get", ""},
	Permission{"PermissionAsSiteListWalletAsset", "/wallet/asset/list", ""},
	Permission{"PermissionAsSiteAddWithdrawSetting", "/wallet/coin/withdraw/setting/add", ""},
	Permission{"PermissionAsSiteUpdateWithdrawSetting", "/wallet/coin/withdraw/setting/update", ""},
	Permission{"PermissionAsSiteRemoveWithdrawSetting", "/wallet/coin/withdraw/setting/del", ""},
	Permission{"PermissionAsSiteGetWithdrawSetting", "/wallet/coin/withdraw/setting/get", ""},
	Permission{"PermissionAsSiteAddWithdrawQuota", "/wallet/coin/withdraw/setting/quota/add", ""},
	Permission{"PermissionAsSiteRemoveWithdrawQuota", "/wallet/coin/withdraw/setting/quota/del", ""},
	Permission{"PermissionAsSiteListWithdrawQuota", "/wallet/coin/withdraw/setting/quota/list", ""},
	Permission{"PermissionAsSiteAddWithdrawWhitelist", "/wallet/coin/withdraw/setting/whitelist/add", ""},
	Permission{"PermissionAsSiteRemoveWithdrawWhitelist", "/wallet/coin/withdraw/setting/whitelist/del", ""},
	Permission{"PermissionAsSiteListWithdrawWhitelist", "/wallet/coin/withdraw/setting/whitelist/list", ""},
	Permission{"PermissionAsSiteListWithdrawPolicy", "/wallet/coin/withdraw/setting/policy/list", ""},
	Permission{"PermissionAsSiteAddCallbackHistory", "/wallet/callback/history/add", ""},
	Permission{"PermissionAsSiteUpdateCallbackHistory", "/wallet/callback/history/update", ""},
	Permission{"PermissionAsSiteListCallbackHistory", "/wallet/callback/history/list", ""},
}
