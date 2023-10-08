// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with  soikawallet. If not, see <http://www.gnu.org/licenses/>.

package settings

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/tview"
)

func (p *pageSettings) tabLabels() *tview.Flex {
	var (
		selectedOption     = uint8(1)
		selectedLabelIndex = uint32(0)
	)

	root := tview.NewTreeNode(`Labels`)

	treeLabels := tview.NewTreeView().SetRoot(root)
	treeLabels.SetBorder(true)

	inputFieldLabel := tview.NewInputField().
		SetLabel("Name").
		SetFieldWidth(20)

	updateLabels := func() {
		var labels map[uint32]string
		root.ClearChildren()

		// Hide

		if selectedOption == 1 {
			root.SetText("Account labels")
			labels = p.API().GetAccountLabels()
		} else if selectedOption == 2 {
			root.SetText("Address labels")
			labels = p.API().GetAddressLabels()
		}
		for index, title := range labels {
			root.AddChild(tview.NewTreeNode(title).SetReference(index))
		}
	}

	formControls := tview.NewForm().
		SetItemPadding(1).
		AddDropDown("Standard", []string{"Account", "Address"}, 0, func(option string, optionIndex int) {
			if p.API() != nil {
				selectedOption = uint8(optionIndex) + 1
				inputFieldLabel.SetText("")
				updateLabels()
			}
		}).
		AddFormItem(inputFieldLabel).
		AddButton(p.Tr().T("ui.button", "add"), func() {
			if p.API() != nil {
				_, err := p.API().AddLabel(&dto.AddLabelDTO{
					LabelType: selectedOption,
					Title:     inputFieldLabel.GetText(),
				})
				if err != nil {
					p.Emit(
						events.EventLogError,
						fmt.Sprintf("Cannot add label: %s", err),
					)
				} else {
					inputFieldLabel.SetText("")
				}
				updateLabels()
			}
		})
	formControls.SetBorder(true)

	formDetails := tview.NewForm().
		SetItemPadding(1).
		AddInputField("Label", "", 20, nil, nil).
		AddButton("Save", func() {}).
		AddButton("Remove", func() {
			if p.API() != nil && selectedLabelIndex > 0 {
				err := p.API().RemoveLabel(&dto.RemoveLabelDTO{
					LabelType: selectedOption,
					Index:     selectedLabelIndex,
				})
				if err != nil {
					p.Emit(
						events.EventLogError,
						fmt.Sprintf("Cannot remove label: %s", err),
					)
				}
				selectedLabelIndex = 0
				updateLabels()
			}
		})

	formDetails.SetBorder(true)

	treeLabels.SetSelectedFunc(func(node *tview.TreeNode) {
		// Hide
		formDetails.GetFormItem(0).(*tview.InputField).SetText("")

		reference := node.GetReference()
		if reference != nil {
			selectedLabelIndex = reference.(uint32)
			formDetails.GetFormItem(0).(*tview.InputField).SetText(node.GetText())
		} else {
			selectedLabelIndex = 0
		}
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(formControls, 0, 1, false).
		AddItem(treeLabels, 0, 1, false).
		AddItem(formDetails, 0, 1, false)

	return layout
}
