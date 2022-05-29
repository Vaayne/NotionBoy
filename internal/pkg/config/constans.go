package config

const (
	MSG_BIND_SUCCESS = `<html>
<body>
<h1>恭喜 🎉 成功绑定 Notion</h1>
<p>
请关闭网页，返回微信使用 NotionBoy 发送消息
</p>
</body>
</html>`

	MSG_BINDING = `
欢迎使用 NotionBoy， 了解更多 NotionBoy 的功能，请参考：https://www.theboys.tech/notion-boy


您还未绑定，请访问下面的链接，然后选择允许访问的页面，并点击 "允许"，即可绑定 Notion。
(升级后的 NotionBoy 支持下面的链接进行绑定，不需要再手动配置 Bot，也不需要手动输入绑定信息)

注意，只选择一个页面，如果选择了多个页面，也只会使用选择的第一个页面。

如果微信打开有问题，请复制 URL 到浏览器打开。

`

	MSG_UNBIND_SUCCESS = "成功解除 Notion 绑定！"

	CMD_BIND       = "绑定"
	CMD_UNBIND     = "解绑"
	MSG_UNSUPPOERT = "不支持的消息类型!"
	MSG_PROCESSING = "正在处理，请稍后去 Notion 查看"
)
