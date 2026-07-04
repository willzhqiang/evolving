package ths

const AsGetAccountInfoSim = `osascript -e '
on getAccountInfoSim()
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
				click button "模拟" of tradeWindow
				delay 0.6
				tell table 1 of scroll area 1 of tradeWindow
					set simulationAccountInfo to get value of every static text of every UI element of every row of table 1 of scroll area 1 of tradeWindow
					return {"successed", simulationAccountInfo}
				end tell
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getAccountInfoSim
getAccountInfoSim()
'`

const AsIssuingEntrustSim = `osascript -e '
on issuingEntrustSim(tradingAction, assetType, stockCode, price, amount)
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
				click button "模拟" of tradeWindow
				if assetType is "stock" then
					click button "股票" of tradeWindow
				else
					return {"failed", "wrong option: " & assetType}
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
				try
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
				on error
					set revocableEntrustment1 to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments1 to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
				end try
				set value of text field 2 of tradeWindow to stockCode
				if tradingAction is "buy" then
					click button "卖出" of tradeWindow
					click button "买入" of tradeWindow
				else if tradingAction is "sell" then
					click button "卖出" of tradeWindow
					click button "买入" of tradeWindow
					click button "卖出" of tradeWindow
				end if
				set value of attribute "AXFocused" of text field 2 of tradeWindow to true
				set value of text field 2 of tradeWindow to stockCode
				if price is "None" then
					delay 0.05
					if tradingAction is "buy" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 2 of tradeWindow)
						if price is "- -" then
							set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 1 of table 1 of scroll area 3 of tradeWindow)
						end if
					else if tradingAction is "sell" then
						set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 3 of tradeWindow)
						if price is "- -" then
							set price to item 1 of item 1 of (get value of attribute "AXTitle" of every button of every UI element of row 5 of table 1 of scroll area 2 of tradeWindow)
						end if
					end if
				end if
				delay 0.25
				set value of text field 1 of tradeWindow to price
				set value of text field 3 of tradeWindow to amount
				if tradingAction is "buy" then
					click button "确定买入" of tradeWindow
				else if tradingAction is "sell" then
					click button "确定卖出" of tradeWindow
				end if
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
				delay 0.25
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
end issuingEntrustSim

on run {tradingAction, assetType, stockCode, price, amount}
	issuingEntrustSim(tradingAction, assetType, stockCode, price, amount)
end run
'`

const AsGetHoldingSharesSim = `osascript -e '
on getHoldingSharesSim(assetType)
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
				click button "模拟" of tradeWindow
				if assetType is "stock" then
					click button "股票" of tradeWindow
				else
					return {"failed", "wrong option: " & assetType}
				end if
				click button "持仓" of tradeWindow
				try
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 4 of tradeWindow
					return {"successed", comments, holdingShares}
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set holdingShares to get value of every static text of every row of table 1 of scroll area 5 of tradeWindow
					return {"successed", comments, holdingShares}
				end try
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getHoldingSharesSim

on run {assetType}
	getHoldingSharesSim(assetType)
end run
'`

const AsGetEntrustSim = `osascript -e '
on getEntrustSim(assetType, dateRange, isRevocable)
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
				click button "模拟" of tradeWindow
				click button "股票" of tradeWindow
				click button "持仓" of tradeWindow
				click button "委托" of tradeWindow
				click button "今天" of tradeWindow
				delay 0.1
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
					set info to get value of static text of sheet 1 of tradeWindow
					if info is {"警告", "不支持历史委托查询"} then
						click button "确认" of sheet 1 of tradeWindow
						return {"failed", info}
					end if
				end try
				try
					set revocableEntrustment to get value of static text of every row of table 1 of scroll area 4 of tradeWindow
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 4 of tradeWindow
					return {"successed", comments, revocableEntrustment}
				on error
					set revocableEntrustment to get value of static text of every row of table 1 of scroll area 5 of tradeWindow
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					return {"successed", comments, revocableEntrustment}
				end try
			on error
				return {"failed", "unknown err"}
			end try
		end tell
	end tell
end getEntrustSim

on run {assetType, dateRange, isRevocable}
	getEntrustSim(assetType, dateRange, isRevocable)
end run
'`

const AsRevokeEntrustSim = `osascript -e '
on revokeEntrustSim(revokeType, assetType, contractNo)
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
				click button "模拟" of tradeWindow
				if assetType is "stock" then
					click button "股票" of tradeWindow
				else
					return {"failed", "wrong option"}
				end if
				click button "委托" of tradeWindow
				delay 0.2
				if revokeType is "allBuyAndSell" then
					click button "全撤" of tradeWindow
				else if revokeType is "allBuy" then
					click button "撤买" of tradeWindow
				else if revokeType is "allSell" then
					click button "撤卖" of tradeWindow
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
end revokeEntrustSim

on run {revokeType, assetType, contractNo}
	revokeEntrustSim(revokeType, assetType, contractNo)
end run
'`

const AsGetClosedDealsSim = `osascript -e '
on getClosedDealsSim(assetType, dateRange)
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
				click button "模拟" of tradeWindow
				click button "股票" of tradeWindow
				click button "成交" of tradeWindow
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
					set closedDeals to get value of every static text of every row of table 1 of scroll area 4 of tradeWindow
				on error
					set comments to get value of attribute "AXTitle" of button of group 1 of table 1 of scroll area 5 of tradeWindow
					set closedDeals to get value of every static text of every row of table 1 of scroll area 5 of tradeWindow
				end try
				try
					set info to get value of static text of sheet 1 of tradeWindow
					if info contains "警告" then
						click button "确认" of sheet 1 of tradeWindow
						return {"failed", {"警告, 业务提示: 查询时间区间必须在30天以内"}}
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
end getClosedDealsSim

on run {assetType, dateRange}
	getClosedDealsSim(assetType, dateRange)
end run
'`

const AsGetCapitalDetailsSim = `osascript -e '
on getCapitalDetailsSim(assetType, dateRange)
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
				click button "模拟" of tradeWindow
				delay 0.1
				click button "股票" of tradeWindow
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
						return {"failed", "警告"}
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
end getCapitalDetailsSim

on run {assetType, dateRange}
	getCapitalDetailsSim(assetType, dateRange)
end run
'`
