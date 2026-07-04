package ths

const AsIsClientLoggedIn = `osascript -e '
on isClientLoggedIn()
	tell application "System Events"
		tell application "System Events" to set isRunning to exists (processes where name is "同花顺")
		if isRunning then
			tell application "同花顺" to activate
			delay 4
			tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
				delay 0.25
				set val to get value of attribute "AXTitle" of button of tradeWindow
				if val contains "游客登录" then
					return false
				end if
				return true
			end tell
		end if
	end tell
end isClientLoggedIn
isClientLoggedIn()
'`

const AsIsBrokerLoggedIn = `osascript -e '
on isBrokerLoggedIn()
	tell application "同花顺" to activate
	delay 0.8
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			click button 6 of tradeWindow
			click button "模拟" of tradeWindow
			click button "A股" of tradeWindow
			try
				set info to get value of attribute "AXTitle" of button of UI element 2 of row 10 of table 1 of scroll area 1 of tradeWindow
				if info contains "退出" then
					return true
				end if
			on error
				return false
			end try
			return false
		end tell
	end tell
end isBrokerLoggedIn
isBrokerLoggedIn()
'`

const AsLoginClient = `osascript -e '
on loginClientHelp(userid, pwd)
	try
		tell application "同花顺" to quit
		delay 1
	end try
	tell application "同花顺" to activate
	delay 3
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				set value of text field 1 of tradeWindow to userid
				set value of text field 2 of tradeWindow to pwd
				click button "登 录" of tradeWindow
				delay 10
			end try
		end tell
	end tell
end loginClientHelp

on loginClient(userid, pwd)
	set isRunning to false
	repeat 3 times
		tell application "System Events" to set isRunning to exists (processes where name is "同花顺")
		if not isRunning then
			loginClientHelp(userid, pwd)
			tell application "System Events" to set isRunning to exists (processes where name is "同花顺")
		end if
		if isRunning then
			exit repeat
		end if
	end repeat
	if isRunning then
		return "successed"
	else
		return "failed"
	end if
end loginClient

on run {userid, pwd}
	loginClient(userid, pwd)
end run
'`

const AsLogoutClient = `osascript -e '
on logoutClient()
	try
		tell application "同花顺" to quit
		return "successed"
	on error
		return "failed"
	end try
end logoutClient
logoutClient()
'`

const AsLoginBroker = `osascript -e '
on loginBroker(broker_name, trade_account, trade_pwd)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			click button 6 of tradeWindow
			click button "A股" of tradeWindow
			try
				set info to get value of attribute "AXTitle" of button of UI element 2 of row 10 of table 1 of scroll area 1 of tradeWindow
				if info contains "退出" then
					return "successed"
				end if
			on error
				try
					click button "立即登录" of tradeWindow
					click button 1 of combo box 1 of sheet 1 of tradeWindow
					set brokerName to broker_name
					set historyBrokers to get value of text field of list 1 of scroll area 1 of combo box 1 of sheet 1 of tradeWindow
					if historyBrokers is {} then
						return "failed"
					end if
					if historyBrokers contains brokerName then
						repeat with rowNum from 1 to length of historyBrokers
							set theCurrentListItem to item rowNum of historyBrokers
							if theCurrentListItem contains brokerName then
								exit repeat
							end if
						end repeat
					else
						return "failed"
					end if
					select text field rowNum of list 1 of scroll area 1 of combo box 1 of sheet 1 of tradeWindow
					set po to get position of text field rowNum of list 1 of scroll area 1 of combo box 1 of sheet 1 of tradeWindow
					set po1 to get item 1 of po
					set po2 to get item 2 of po
					do shell script "if [ -x /opt/homebrew/bin/cliclick ]; then /opt/homebrew/bin/cliclick c:" & po1 & "," & po2 & "; else /usr/local/bin/cliclick c:" & po1 & "," & po2 & "; fi"
					delay 0.25
					set value of checkbox 1 of sheet 1 of tradeWindow to trade_account
					set value of text field 1 of sheet 1 of tradeWindow to trade_pwd
					set verificationCodeList to get value of static text of sheet 1 of tradeWindow
					set verificationCode to ""
					repeat with x from 8 to 12
						set verificationCode to verificationCode & item x of verificationCodeList
					end repeat
					set value of text field 2 of sheet 1 of tradeWindow to verificationCode
					click button "登录" of sheet 1 of tradeWindow
					try
						delay 2
						set info to get value of static text of sheet 1 of tradeWindow
						set warningFlag to item 1 of info
						if warningFlag is "连接委托主站失败！可能是以下原因：" then
							click button "确定" of sheet 1 of tradeWindow
						end if
					end try
					return "successed"
				on error
					return "failed"
				end try
			end try
		end tell
	end tell
end loginBroker

on run {broker_name, trade_account, trade_pwd}
	loginBroker(broker_name, trade_account, trade_pwd)
end run
'`

const AsLogoutBroker = `osascript -e '
on logoutBroker()
	tell application "同花顺" to activate
	delay 0.25
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				delay 0.25
				click button "A股" of tradeWindow
				click button "退出" of UI element 2 of row 10 of table 1 of scroll area 1 of tradeWindow
				return "successed"
			on error
				return "failed"
			end try
		end tell
	end tell
end logoutBroker
logoutBroker()
'`

const AsTransferBank2Broker = `osascript -e '
on transfer(transferType, amount, bank_pwd, trade_pwd)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 6 of tradeWindow
				click button "模拟" of tradeWindow
				click button "A股" of tradeWindow
				delay 0.25
				click button "转账" of UI element 2 of row 9 of table 1 of scroll area 1 of tradeWindow
				delay 0.25
				set value of text field 1 of window "银证转账" of application process "同花顺" of application "System Events" to amount
				set value of text field 2 of window "银证转账" of application process "同花顺" of application "System Events" to bank_pwd
				set value of text field 3 of window "银证转账" of application process "同花顺" of application "System Events" to trade_pwd
				delay 0.1
				click button "确定转入券商" of window "银证转账" of application process "同花顺" of application "System Events"
				click button "确认" of sheet 1 of window "银证转账" of application process "同花顺" of application "System Events"
				delay 0.2
				try
					set info to get value of static text of sheet 1 of window "银证转账" of application process "同花顺" of application "System Events"
					if info contains "警告" then
						click button "确认" of sheet 1 of window "银证转账" of application process "同花顺" of application "System Events"
						click button 6 of window "银证转账" of application process "同花顺" of application "System Events"
						return {"failed", "警告:外部机构[5200]不支持7*24银证业务"}
					end if
				end try
				try
					click button 1 of window "银证转账" of application process "同花顺" of application "System Events"
				end try
				return "successed"
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end transfer

on run {transferType, amount, bank_pwd, trade_pwd}
	transfer(transferType, amount, bank_pwd, trade_pwd)
end run
'`

const AsTransferBroker2Bank = `osascript -e '
on transfer(transferType, amount, bank_pwd, trade_pwd)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 6 of tradeWindow
				click button "模拟" of tradeWindow
				click button "A股" of tradeWindow
				delay 0.25
				click button "转账" of UI element 2 of row 9 of table 1 of scroll area 1 of tradeWindow
				delay 0.25
				click button 2 of window "银证转账" of application process "同花顺" of application "System Events"
				set value of text field 1 of window "银证转账" of application process "同花顺" of application "System Events" to amount
				set value of text field 2 of window "银证转账" of application process "同花顺" of application "System Events" to bank_pwd
				set value of text field 3 of window "银证转账" of application process "同花顺" of application "System Events" to trade_pwd
				delay 0.1
				click button "确定转入银行" of window "银证转账" of application process "同花顺" of application "System Events"
				click button "确认" of sheet 1 of window "银证转账" of application process "同花顺" of application "System Events"
				delay 0.2
				try
					set info to get value of static text of sheet 1 of window "银证转账" of application process "同花顺" of application "System Events"
					if info contains "警告" then
						click button "确认" of sheet 1 of window "银证转账" of application process "同花顺" of application "System Events"
						click button 6 of window "银证转账" of application process "同花顺" of application "System Events"
						return {"failed", "警告:外部机构[5200]不支持7*24银证业务"}
					end if
				end try
				try
					click button 1 of window "银证转账" of application process "同花顺" of application "System Events"
				end try
				return "successed"
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end transfer

on run {transferType, amount, bank_pwd, trade_pwd}
	transfer(transferType, amount, bank_pwd, trade_pwd)
end run
'`

const AsRevokeAllEntrust = `osascript -e '
on revokeAllEntrust()
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "股票" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "全撤" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try
				click button "科创板盘后" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "全撤" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try
				click button "创业板盘后" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "全撤" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try
				return {"successed"}
			on error
				return {"failed"}
			end try
		end tell
	end tell
end revokeAllEntrust
revokeAllEntrust()
'`

const AsOneKeyIPO = `osascript -e '
on oneKeyIPO()
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "新股申购" of tradeWindow
				delay 1
				click button "一键申购" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "警告" then
						click button "确认" of sheet 1 of tradeWindow
						return {"failed", "请输入正确的委托数量"}
					end if
				end try
				return "successed"
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end oneKeyIPO
oneKeyIPO()
'`

const AsGetAccountInfo = `osascript -e '
on getAccountInfo()
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				delay 0.5
				tell table 1 of scroll area 1 of tradeWindow
					set accountInfo to get value of every static text of every UI element of every row of table 1 of scroll area 1 of tradeWindow
					return {"successed", accountInfo}
				end tell
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getAccountInfo
getAccountInfo()
'`

const AsGetHoldingSharesStock = `osascript -e '
on getHoldingShares(assetType)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "股票" of tradeWindow
				click button "持仓" of tradeWindow
				delay 0.1
				try
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 4 of tradeWindow
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 5 of tradeWindow
				end try
				return {"successed", comments, holdingShares}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getHoldingShares

on run {assetType}
	getHoldingShares(assetType)
end run
'`

const AsGetHoldingSharesSciTech = `osascript -e '
on getHoldingShares(assetType)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "科创板盘后" of tradeWindow
				click button "持仓" of tradeWindow
				delay 0.1
				try
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 4 of tradeWindow
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 5 of tradeWindow
				end try
				return {"successed", comments, holdingShares}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getHoldingShares

on run {assetType}
	getHoldingShares(assetType)
end run
'`

const AsGetHoldingSharesGem = `osascript -e '
on getHoldingShares(assetType)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "创业板盘后" of tradeWindow
				click button "持仓" of tradeWindow
				delay 0.1
				try
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 4 of tradeWindow
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 5 of tradeWindow
				end try
				return {"successed", comments, holdingShares}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getHoldingShares

on run {assetType}
	getHoldingShares(assetType)
end run
'`

const AsIssuingEntrustBuyStock = `osascript -e '
on issuingEntrust(tradingAction, assetType, stockCode, price, amount)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "股票" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 2 of tradeWindow)
					if price is "- -" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 3 of tradeWindow)
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				click button "确定买入" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "提示信息" then
						delay 0.01
						click button "确认" of sheet 1 of tradeWindow
					end if
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"买入委托"} or info contains {"卖出委托"} then
					delay 0.01
					click button "确认" of sheet 1 of tradeWindow
				end if
				set flag to 0
				set info to ""
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"警告"} then
					set flag to -1
				end if
				delay 0.1
				if flag is not 0 then
					click button "确认" of sheet 1 of tradeWindow
					return {"failed", flag, info}
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				delay 0.6
				try
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set contractNoList to {}
				set contractNoList1 to {}
				set contractNoList2 to {}
				repeat with x from 1 to length of revocableEntrustment1
					set end of contractNoList1 to item 11 of item x of revocableEntrustment1
				end repeat
				repeat with x from 1 to length of revocableEntrustment2
					set end of contractNoList2 to item 11 of item x of revocableEntrustment2
				end repeat
				repeat with x from 1 to length of contractNoList2
					set curitem to item x of contractNoList2
					if contractNoList1 does not contain curitem then
						set end of contractNoList to curitem
					end if
				end repeat
				if contractNoList is {} then
					set info to "委托失败"
					return {"failed", "委托失败"}
				else
					return {"successed", contractNoList}
				end if
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end issuingEntrust

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrust(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsIssuingEntrustSellStock = `osascript -e '
on issuingEntrust(tradingAction, assetType, stockCode, price, amount)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "股票" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				click button "卖出" of tradeWindow
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 3 of tradeWindow)
					if price is "- -" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 2 of tradeWindow)
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				click button "确定卖出" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "提示信息" then
						delay 0.01
						click button "确认" of sheet 1 of tradeWindow
					end if
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"买入委托"} or info contains {"卖出委托"} then
					delay 0.01
					click button "确认" of sheet 1 of tradeWindow
				end if
				set flag to 0
				set info to ""
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"警告"} then
					set flag to -1
				end if
				delay 0.1
				if flag is not 0 then
					click button "确认" of sheet 1 of tradeWindow
					return {"failed", flag, info}
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				delay 0.6
				try
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set contractNoList to {}
				set contractNoList1 to {}
				set contractNoList2 to {}
				repeat with x from 1 to length of revocableEntrustment1
					set end of contractNoList1 to item 11 of item x of revocableEntrustment1
				end repeat
				repeat with x from 1 to length of revocableEntrustment2
					set end of contractNoList2 to item 11 of item x of revocableEntrustment2
				end repeat
				repeat with x from 1 to length of contractNoList2
					set curitem to item x of contractNoList2
					if contractNoList1 does not contain curitem then
						set end of contractNoList to curitem
					end if
				end repeat
				if contractNoList is {} then
					set info to "委托失败"
					return {"failed", "委托失败"}
				else
					return {"successed", contractNoList}
				end if
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end issuingEntrust

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrust(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsIssuingEntrustBuySciTech = `osascript -e '
on issuingEntrust(tradingAction, assetType, stockCode, price, amount)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "科创板盘后" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 2 of tradeWindow)
					if price is "- -" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 3 of tradeWindow)
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				click button "确定买入" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "提示信息" then
						delay 0.01
						click button "确认" of sheet 1 of tradeWindow
					end if
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"买入委托"} or info contains {"卖出委托"} then
					delay 0.01
					click button "确认" of sheet 1 of tradeWindow
				end if
				set flag to 0
				set info to ""
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"警告"} then
					set flag to -1
				end if
				delay 0.1
				if flag is not 0 then
					click button "确认" of sheet 1 of tradeWindow
					return {"failed", flag, info}
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				delay 0.6
				try
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set contractNoList to {}
				set contractNoList1 to {}
				set contractNoList2 to {}
				repeat with x from 1 to length of revocableEntrustment1
					set end of contractNoList1 to item 11 of item x of revocableEntrustment1
				end repeat
				repeat with x from 1 to length of revocableEntrustment2
					set end of contractNoList2 to item 11 of item x of revocableEntrustment2
				end repeat
				repeat with x from 1 to length of contractNoList2
					set curitem to item x of contractNoList2
					if contractNoList1 does not contain curitem then
						set end of contractNoList to curitem
					end if
				end repeat
				if contractNoList is {} then
					set info to "委托失败"
					return {"failed", "委托失败"}
				else
					return {"successed", contractNoList}
				end if
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end issuingEntrust

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrust(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsIssuingEntrustSellSciTech = `osascript -e '
on issuingEntrust(tradingAction, assetType, stockCode, price, amount)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "科创板盘后" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				click button "卖出" of tradeWindow
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 3 of tradeWindow)
					if price is "- -" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 2 of tradeWindow)
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				click button "确定卖出" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "提示信息" then
						delay 0.01
						click button "确认" of sheet 1 of tradeWindow
					end if
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"买入委托"} or info contains {"卖出委托"} then
					delay 0.01
					click button "确认" of sheet 1 of tradeWindow
				end if
				set flag to 0
				set info to ""
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"警告"} then
					set flag to -1
				end if
				delay 0.1
				if flag is not 0 then
					click button "确认" of sheet 1 of tradeWindow
					return {"failed", flag, info}
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				delay 0.6
				try
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set contractNoList to {}
				set contractNoList1 to {}
				set contractNoList2 to {}
				repeat with x from 1 to length of revocableEntrustment1
					set end of contractNoList1 to item 11 of item x of revocableEntrustment1
				end repeat
				repeat with x from 1 to length of revocableEntrustment2
					set end of contractNoList2 to item 11 of item x of revocableEntrustment2
				end repeat
				repeat with x from 1 to length of contractNoList2
					set curitem to item x of contractNoList2
					if contractNoList1 does not contain curitem then
						set end of contractNoList to curitem
					end if
				end repeat
				if contractNoList is {} then
					set info to "委托失败"
					return {"failed", "委托失败"}
				else
					return {"successed", contractNoList}
				end if
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end issuingEntrust

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrust(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsIssuingEntrustBuyGem = `osascript -e '
on issuingEntrust(tradingAction, assetType, stockCode, price, amount)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "创业板盘后" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 2 of tradeWindow)
					if price is "- -" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 3 of tradeWindow)
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				click button "确定买入" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "提示信息" then
						delay 0.01
						click button "确认" of sheet 1 of tradeWindow
					end if
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"买入委托"} or info contains {"卖出委托"} then
					delay 0.01
					click button "确认" of sheet 1 of tradeWindow
				end if
				set flag to 0
				set info to ""
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"警告"} then
					set flag to -1
				end if
				delay 0.1
				if flag is not 0 then
					click button "确认" of sheet 1 of tradeWindow
					return {"failed", flag, info}
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				delay 0.6
				try
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set contractNoList to {}
				set contractNoList1 to {}
				set contractNoList2 to {}
				repeat with x from 1 to length of revocableEntrustment1
					set end of contractNoList1 to item 11 of item x of revocableEntrustment1
				end repeat
				repeat with x from 1 to length of revocableEntrustment2
					set end of contractNoList2 to item 11 of item x of revocableEntrustment2
				end repeat
				repeat with x from 1 to length of contractNoList2
					set curitem to item x of contractNoList2
					if contractNoList1 does not contain curitem then
						set end of contractNoList to curitem
					end if
				end repeat
				if contractNoList is {} then
					set info to "委托失败"
					return {"failed", "委托失败"}
				else
					return {"successed", contractNoList}
				end if
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end issuingEntrust

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrust(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsIssuingEntrustSellGem = `osascript -e '
on issuingEntrust(tradingAction, assetType, stockCode, price, amount)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "创业板盘后" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				click button "卖出" of tradeWindow
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 3 of tradeWindow)
					if price is "- -" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 2 of tradeWindow)
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				click button "确定卖出" of tradeWindow
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "提示信息" then
						delay 0.01
						click button "确认" of sheet 1 of tradeWindow
					end if
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"买入委托"} or info contains {"卖出委托"} then
					delay 0.01
					click button "确认" of sheet 1 of tradeWindow
				end if
				set flag to 0
				set info to ""
				try
					set info to get value of static text of sheet 1 of tradeWindow
				end try
				if info contains {"警告"} then
					set flag to -1
				end if
				delay 0.1
				if flag is not 0 then
					click button "确认" of sheet 1 of tradeWindow
					return {"failed", flag, info}
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if checkboxStatus is true then click theCheckbox
				end tell
				delay 0.6
				try
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment2 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments2 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set contractNoList to {}
				set contractNoList1 to {}
				set contractNoList2 to {}
				repeat with x from 1 to length of revocableEntrustment1
					set end of contractNoList1 to item 11 of item x of revocableEntrustment1
				end repeat
				repeat with x from 1 to length of revocableEntrustment2
					set end of contractNoList2 to item 11 of item x of revocableEntrustment2
				end repeat
				repeat with x from 1 to length of contractNoList2
					set curitem to item x of contractNoList2
					if contractNoList1 does not contain curitem then
						set end of contractNoList to curitem
					end if
				end repeat
				if contractNoList is {} then
					set info to "委托失败"
					return {"failed", "委托失败"}
				else
					return {"successed", contractNoList}
				end if
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end issuingEntrust

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrust(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsGetEntrustToday = `osascript -e '
on getEntrust(assetType, dateRange, isRevocable)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				if assetType is "stock" then
					click button "股票" of tradeWindow
				else if assetType is "sciTech" then
					click button "科创板盘后" of tradeWindow
				else if assetType is "gem" then
					click button "创业板盘后" of tradeWindow
				end if
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.01
				click button "今天" of pop over 1 of tradeWindow
				delay 0.1
				set theCheckbox to checkbox 1 of tradeWindow
				tell theCheckbox
					set checkboxStatus to value of theCheckbox as boolean
					if isRevocable is "true" then
						if checkboxStatus is false then click theCheckbox
					else
						if checkboxStatus is true then click theCheckbox
					end if
				end tell
				try
					set revocableEntrustment to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "警告" then
						click button "确认" of sheet 1 of tradeWindow
						return {"failed", {"警告"}, {"业务提示: 超过 93 天"}}
					end if
					return {"successed", comments, revocableEntrustment}
				on error
					return {"successed", comments, revocableEntrustment}
				end try
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getEntrust

on run {assetType, dateRange, isRevocable}
	getEntrust(assetType, dateRange, isRevocable)
end run
'`

const AsGetClosedDealsToday = `osascript -e '
on getClosedDeals(assetType, dateRange)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				if assetType is "stock" then
					click button "股票" of tradeWindow
				else if assetType is "sciTech" then
					click button "科创板盘后" of tradeWindow
				else if assetType is "gem" then
					click button "创业板盘后" of tradeWindow
				end if
				click button "成交" of tradeWindow
				click button "今天" of tradeWindow
				click button "今天" of pop over 1 of tradeWindow
				delay 0.45
				try
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					set closedDeals to get value of every static text of every row of table 1 of scroll area 4 of tradeWindow
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set closedDeals to get value of every static text of every row of table 1 of scroll area 5 of tradeWindow
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "警告" then
						click button "确认" of sheet 1 of tradeWindow
						return {"failed", "警告, 业务提示: 超过 93 天"}
					end if
					return {"successed", comments, closedDeals}
				on error
					return {"successed", comments, closedDeals}
				end try
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getClosedDeals

on run {assetType, dateRange}
	getClosedDeals(assetType, dateRange)
end run
'`

const AsGetBidsStock = `osascript -e '
on getBids(assetType, stockCode)
	tell application "同花顺" to activate
	delay 0.25
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			click button 6 of tradeWindow
			click button "模拟" of tradeWindow
			click button "A股" of tradeWindow
			click button "股票" of tradeWindow
			try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				set value of text field 2 of tradeWindow to stockCode
				delay 0.1
				set bidsSPrice to get value of attribute "AXTitle" of every button of every UI element of every row of table 1 of scroll area 2 of tradeWindow
				set bidsBPrice to get value of attribute "AXTitle" of every button of every UI element of every row of table 1 of scroll area 3 of tradeWindow
				set bidsS_vol to get value of every static text of every UI element of every row of table 1 of scroll area 2 of tradeWindow
				set bidsB_vol to get value of every static text of every UI element of every row of table 1 of scroll area 3 of tradeWindow
				return {"successed", bidsSPrice, bidsBPrice, bidsS_vol, bidsB_vol}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getBids

on run {assetType, stockCode}
	getBids(assetType, stockCode)
end run
'`

const AsGetBidsSciTech = `osascript -e '
on getBids(assetType, stockCode)
	tell application "同花顺" to activate
	delay 0.25
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			click button 6 of tradeWindow
			click button "模拟" of tradeWindow
			click button "A股" of tradeWindow
			click button "科创板盘后" of tradeWindow
			try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				set value of text field 2 of tradeWindow to stockCode
				delay 0.1
				set bidsSPrice to get value of attribute "AXTitle" of every button of every UI element of every row of table 1 of scroll area 2 of tradeWindow
				set bidsBPrice to get value of attribute "AXTitle" of every button of every UI element of every row of table 1 of scroll area 3 of tradeWindow
				set bidsS_vol to get value of every static text of every UI element of every row of table 1 of scroll area 2 of tradeWindow
				set bidsB_vol to get value of every static text of every UI element of every row of table 1 of scroll area 3 of tradeWindow
				return {"successed", bidsSPrice, bidsBPrice, bidsS_vol, bidsB_vol}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getBids

on run {assetType, stockCode}
	getBids(assetType, stockCode)
end run
'`

const AsGetBidsGem = `osascript -e '
on getBids(assetType, stockCode)
	tell application "同花顺" to activate
	delay 0.25
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			click button 6 of tradeWindow
			click button "模拟" of tradeWindow
			click button "A股" of tradeWindow
			click button "创业板盘后" of tradeWindow
			try
				set value of text field 2 of tradeWindow to stockCode
				click button "卖出" of tradeWindow
				click button "买入" of tradeWindow
				set value of text field 2 of tradeWindow to stockCode
				delay 0.1
				set bidsSPrice to get value of attribute "AXTitle" of every button of every UI element of every row of table 1 of scroll area 2 of tradeWindow
				set bidsBPrice to get value of attribute "AXTitle" of every button of every UI element of every row of table 1 of scroll area 3 of tradeWindow
				set bidsS_vol to get value of every static text of every UI element of every row of table 1 of scroll area 2 of tradeWindow
				set bidsB_vol to get value of every static text of every UI element of every row of table 1 of scroll area 3 of tradeWindow
				return {"successed", bidsSPrice, bidsBPrice, bidsS_vol, bidsB_vol}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getBids

on run {assetType, stockCode}
	getBids(assetType, stockCode)
end run
'`

const AsRevokeEntrust = `osascript -e '
on revokeEntrust(revokeType, assetType, contractNo)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				if assetType is "stock" then
					click button "股票" of tradeWindow
				else if assetType is "sciTech" then
					click button "科创板盘后" of tradeWindow
				else if assetType is "gem" then
					click button "创业板盘后" of tradeWindow
				else
					return {"failed", "wrong option: " & assetType}
				end if
				click button "委托" of tradeWindow
				delay 0.2
				if revokeType is "allBuyAndSell" then
					click button "全撤" of tradeWindow
				else if revokeType is "contractNo" then
					set idarea to 4
					try
						set EntrustmentList to get value of static text of row of table 1 of scroll area 4 of tradeWindow
					on error
						set idarea to 5
						set EntrustmentList to get value of static text of row of table 1 of scroll area 5 of tradeWindow
					end try
					if EntrustmentList is {} then
						return {"successed", "nothing to revoke"}
					end if
					repeat with rowNum from 1 to length of EntrustmentList
						set theCurrentListItem to item rowNum of EntrustmentList
						if theCurrentListItem contains contractNo then
							exit repeat
						end if
					end repeat
					set len to length of EntrustmentList
					if rowNum is len then
						if theCurrentListItem does not contain contractNo then
							return {"successed", "contract No. " & contractNo & " was not found"}
						end if
					end if
					if idarea is 4 then
						set po to get position of (get item 11 of (get static text of row rowNum of table 1 of scroll area 4 of tradeWindow))
					else if idarea is 5 then
						set po to get position of (get item 11 of (get static text of row rowNum of table 1 of scroll area 5 of tradeWindow))
					end if
					set po1 to get item 1 of po
					set po2 to get item 2 of po
					do shell script "if [ -x /opt/homebrew/bin/cliclick ]; then /opt/homebrew/bin/cliclick dc:" & po1 & "," & po2 & "; else /usr/local/bin/cliclick dc:" & po1 & "," & po2 & "; fi"
				else
					return {"failed", "wrong revoke type: " & revokeType}
				end if
				try
					click button "确认" of sheet 1 of tradeWindow
					return {"successed", "revoke " & revokeType & " " & assetType & " is successed"}
				on error
					return {"successed", "nothing to revoke"}
				end try
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end revokeEntrust

on run {revokeType, assetType, contractNo}
	revokeEntrust(revokeType, assetType, contractNo)
end run
'`

const AsGetTodayIPO = `osascript -e '
on getTodayIPO()
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "新股申购" of tradeWindow
				delay 1.2
				set comment to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 2 of tradeWindow
				set todayipo to get value of static text of UI element of row of table 1 of scroll area 2 of tradeWindow
				set nums to get value of text field 1 of UI element 4 of row of table 1 of scroll area 2 of tradeWindow
				repeat with idx from 1 to length of nums
					set curItem to item idx of nums
					set item 4 of item idx of todayipo to curItem
				end repeat
				return {"successed", comment, todayipo}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getTodayIPO
getTodayIPO()
'`

const AsGetTransferRecords = `osascript -e '
on getTransferRecords(dateRange)
	tell application "同花顺" to activate
	delay 0.4
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 6 of tradeWindow
				click button "A股" of tradeWindow
				click button "流水" of UI element 2 of row 9 of table 1 of scroll area 1 of tradeWindow
				click button "今天" of window "查询流水" of application process "同花顺" of application "System Events"
				delay 0.2
				if dateRange is "today" then
					click button "今天" of pop over 1 of window "查询流水" of application process "同花顺" of application "System Events"
				else if dateRange is "thisWeek" then
					click button "本周" of pop over 1 of window "查询流水" of application process "同花顺" of application "System Events"
				else if dateRange is "thisMonth" then
					click button "本月" of pop over 1 of window "查询流水" of application process "同花顺" of application "System Events"
				else if dateRange is "thisSeason" then
					click button "本季" of pop over 1 of window "查询流水" of application process "同花顺" of application "System Events"
				else if dateRange is "thisYear" then
					click button "本年" of pop over 1 of window "查询流水" of application process "同花顺" of application "System Events"
				end if
				delay 0.01
				set info to get value of static text of every row of table 1 of scroll area 1 of window "查询流水" of application process "同花顺" of application "System Events"
				set comment to get value of attribute "AXTitle" of every button of group 1 of table 1 of scroll area 1 of window "查询流水" of application process "同花顺" of application "System Events"
				click button 1 of window "查询流水" of application process "同花顺" of application "System Events"
				return {"successed", comment, info}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getTransferRecords

on run {dateRange}
	getTransferRecords(dateRange)
end run
'`

const AsRevokeAllBuyEntrust = `osascript -e '
on revokeAllBuyEntrust()
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow

				click button "股票" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "撤买" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try

				click button "科创板盘后" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "撤买" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try

				click button "创业板盘后" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "撤买" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try
				return "successed"
			on error
				return "failed"
			end try
		end tell
	end tell
end revokeAllBuyEntrust

revokeAllBuyEntrust()
'`

const AsRevokeAllSellEntrust = `osascript -e '
on revokeAllSellEntrust()
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow

				click button "股票" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "撤卖" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try

				click button "科创板盘后" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "撤卖" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try

				click button "创业板盘后" of tradeWindow
				click button "委托" of tradeWindow
				delay 0.08
				click button "撤卖" of tradeWindow
				try
					click button "确认" of sheet 1 of tradeWindow
				end try
				return "successed"
			on error
				return "failed"
			end try
		end tell
	end tell
end revokeAllSellEntrust

revokeAllSellEntrust()
'`

const AsGetCapitalDetails = `osascript -e '
on getCapitalDetails(assetType, dateRange)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow

				if assetType is "stock" then
					click button "股票" of tradeWindow
				else if assetType is "sciTech" then
					click button "科创板盘后" of tradeWindow
				else if assetType is "gem" then
					click button "创业板盘后" of tradeWindow
				end if

				click button "资金明细" of tradeWindow
				click button "今天" of tradeWindow

				if dateRange is "today" then
					click button "今天" of pop over 1 of tradeWindow
				else if dateRange is "thisWeek" then
					click button "本周" of pop over 1 of tradeWindow
				else if dateRange is "thisMonth" then
					click button "本月" of pop over 1 of tradeWindow
				else if dateRange is "thisSeason" then
					click button "本季" of pop over 1 of tradeWindow
				else if dateRange is "thisYear" then
					click button "本年" of pop over 1 of tradeWindow
				end if
				delay 0.1

				try
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					set closedDeals to get value of every text field of every row of table 1 of scroll area 4 of tradeWindow
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set closedDeals to get value of every text field of every row of table 1 of scroll area 5 of tradeWindow
				end try

				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "警告" then
						click button "确认" of sheet 1 of tradeWindow
						return {"failed", {"警告"}, {"业务提示: 超过 93 天"}}
					end if
					return {"successed", comments, closedDeals}
				on error
					return {"successed", comments, closedDeals}
				end try
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getCapitalDetails

on run {assetType, dateRange}
	getCapitalDetails(assetType, dateRange)
end run
'`

const AsGetIPO = `osascript -e '
on getIPO(queryType, dateRange)
	tell application "同花顺" to activate
	delay 0.5
	tell application "System Events"
		tell process "同花顺"
			set tradeWindow to window 1
			repeat with candidateWindow in windows
				try
					if (count of buttons of candidateWindow) > 0 or (count of scroll areas of candidateWindow) > 0 then
						set tradeWindow to candidateWindow
						exit repeat
					end if
				end try
			end repeat
			try
				click button 1 of tradeWindow
				click button 6 of tradeWindow
				click button "A股" of tradeWindow

				click button "新股申购" of tradeWindow

				if queryType is "entrust" then
					click button "申购委托" of tradeWindow
				else if queryType is "allotmentNo" then
					click button "配号查询" of tradeWindow
				else if queryType is "winningLots" then
					click button "中签查询" of tradeWindow
				else
					return {"failed", "have no " & queryType & "queryType"}
				end if

				if queryType is not "entrust" then
					click button "今天" of tradeWindow
					if dateRange is "today" then
						click button "今天" of pop over 1 of tradeWindow
					else if dateRange is "thisWeek" then
						click button "本周" of pop over 1 of tradeWindow
					else if dateRange is "thisMonth" then
						click button "本月" of pop over 1 of tradeWindow
					else if dateRange is "thisSeason" then
						click button "本季" of pop over 1 of tradeWindow
					else if dateRange is "thisYear" then
						click button "本季" of pop over 1 of tradeWindow
					end if
				else
					if dateRange is not "today" then
						return {"failed", "<queryType> entrust only supports <dateRange> today"}
					end if
				end if

				set comment to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				set res to get value of static text of every row of table 1 of scroll area 4 of tradeWindow

				return {"successed", comment, res}
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getIPO

on run {queryType, dateRange}
	getIPO(queryType, dateRange)
end run
'`
