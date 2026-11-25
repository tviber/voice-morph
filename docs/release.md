# 发布 macOS 安装包

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
xcrun stapler validate voicemorph-0.2.0.dmg
hdiutil verify voicemorph-0.2.0.dmg
shasum -a 256 voicemorph-0.2.0.dmg > SHA256SUMS
```

确认 `Authority` 为预期的 Developer ID、Gatekeeper 返回 `accepted`，并进行 BlackHole 2ch 试录。在 GitHub 草稿中上传公证后的 DMG 和校验文件，移除候选测试包，核对标签目标与构建提交，再发布正式 Release。
