# 发布桌面安装包

将 `mygo.config.ts` 的版本改为目标版本，并同步 macOS 的版本字段。合入 `main` 后，在该提交上推送同名的 `v` 标签，例如 `v0.1.0-beta.1`。GitHub Actions 会构建 macOS、Windows 和 Linux 安装包，全部成功后上传到对应 Release、生成 SHA-256 校验清单并公开。带预发布后缀的标签会标记为 Pre-release。维护者可提前准备同标签的草稿及发布说明。失败后可在 Actions 重跑失败的作业。

CI 当前分发 ad-hoc 签名的 macOS 包和未签名的 Windows 包。若要分发经过 Apple 公证的版本，按下文构建，并在发布流程中使用公证后的文件。

## 在 macOS 构建机签名与公证

在 macOS 构建机上安装 Go、Bun 和 Xcode 命令行工具。将 Developer ID Application 证书及私钥导入钥匙串，用 `xcrun notarytool store-credentials voicemorph` 保存 Apple 公证凭据。私钥和凭据留在钥匙串中。

从要发布的提交运行以下步骤，核对 `mygo.config.ts` 中的版本及对应的 GitHub Release 草稿。签名身份使用 `security find-identity -v -p codesigning` 显示的完整身份名。

```sh
export VOICEMORPH_SIGN_IDENTITY='Developer ID Application: YOUR NAME (TEAMID)'
export VOICEMORPH_NOTARY_PROFILE=voicemorph
bun install --frozen-lockfile
CGO_ENABLED=0 go test ./...
bun test
bun run build -platform darwin/universal
cd build/darwin-universal
codesign --verify --deep --strict VoiceMorph.app
codesign -dv --verbose=4 VoiceMorph.app
spctl --assess --type execute --verbose=2 VoiceMorph.app
xcrun stapler validate VoiceMorph.app
xcrun stapler validate voicemorph-0.1.0-beta.1.dmg
hdiutil verify voicemorph-0.1.0-beta.1.dmg
shasum -a 256 voicemorph-0.1.0-beta.1.dmg > SHA256SUMS
```

确认 `Authority` 为预期的 Developer ID、Gatekeeper 返回 `accepted`，并进行 BlackHole 2ch 试录。在 GitHub 草稿中上传公证后的 DMG 和校验文件，移除候选测试包。发布 Beta 时使用 `v0.1.0-beta.1` 格式的标签，勾选 **Set as a pre-release**，并核对标签目标与构建提交。
