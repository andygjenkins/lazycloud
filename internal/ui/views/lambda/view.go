package lambda

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	
	lambdaService "lazycloud/internal/aws/lambda"
)

type View struct {
	*tview.Flex
	
	functionList   *tview.List
	functionDetail *tview.TextView
	statusBar      *tview.TextView
	
	service    *lambdaService.Service
	functions  []*lambdaService.Function
	loading    bool
}

func NewView(service *lambdaService.Service) *View {
	v := &View{
		service: service,
	}
	
	v.setupUI()
	v.setupKeybindings()
	
	return v
}

func (v *View) setupUI() {
	// Create function list
	v.functionList = tview.NewList().ShowSecondaryText(true)
	v.functionList.SetBorder(true).SetTitle(" Lambda Functions ").SetTitleAlign(tview.AlignLeft)
	v.functionList.SetHighlightFullLine(true)
	v.functionList.SetSelectedFunc(v.onFunctionSelected)
	
	// Create function detail view
	v.functionDetail = tview.NewTextView()
	v.functionDetail.SetBorder(true).SetTitle(" Function Details ").SetTitleAlign(tview.AlignLeft)
	v.functionDetail.SetWordWrap(true)
	v.functionDetail.SetDynamicColors(true)
	
	// Create status bar
	v.statusBar = tview.NewTextView()
	v.statusBar.SetText("Press 'r' to refresh, 'q' to quit")
	v.statusBar.SetTextAlign(tview.AlignLeft)
	
	// Create main layout
	mainFlex := tview.NewFlex().
		AddItem(v.functionList, 0, 1, true).
		AddItem(v.functionDetail, 0, 2, false)
	
	v.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(v.statusBar, 1, 0, false)
		
	// Initial load
	go v.loadFunctions()
}

func (v *View) setupKeybindings() {
	v.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			go v.loadFunctions()
			return nil
		case 'q':
			// This will be handled by the main app
			return event
		}
		return event
	})
}

func (v *View) loadFunctions() {
	v.loading = true
	v.updateStatus("Loading Lambda functions...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	functions, err := v.service.ListFunctions(ctx)
	if err != nil {
		v.updateStatus(fmt.Sprintf("Error: %v", err))
		v.loading = false
		return
	}
	
	v.functions = functions
	v.updateFunctionList()
	v.updateStatus(fmt.Sprintf("Loaded %d functions", len(functions)))
	v.loading = false
}

func (v *View) updateFunctionList() {
	v.functionList.Clear()
	
	if len(v.functions) == 0 {
		v.functionList.AddItem("No Lambda functions found", "", 0, nil)
		v.functionDetail.SetText("No functions available")
		return
	}
	
	for i, fn := range v.functions {
		primaryText := fn.Name
		secondaryText := fmt.Sprintf("%s | %dMB | %ds timeout", 
			fn.Runtime, fn.Memory, fn.Timeout)
		
		// Add status indicator
		statusColor := "green"
		if fn.Status != "Active" {
			statusColor = "yellow"
		}
		
		primaryText = fmt.Sprintf("[%s]●[white] %s", statusColor, fn.Name)
		
		v.functionList.AddItem(primaryText, secondaryText, rune('1'+i), nil)
	}
	
	// Select first function if available
	if len(v.functions) > 0 {
		v.functionList.SetCurrentItem(0)
		v.showFunctionDetails(0)
	}
}

func (v *View) onFunctionSelected(index int, primaryText, secondaryText string, shortcut rune) {
	v.showFunctionDetails(index)
}

func (v *View) showFunctionDetails(index int) {
	if index < 0 || index >= len(v.functions) {
		return
	}
	
	fn := v.functions[index]
	
	details := strings.Builder{}
	details.WriteString(fmt.Sprintf("[yellow]Function Name:[white] %s\n", fn.Name))
	details.WriteString(fmt.Sprintf("[yellow]Runtime:[white] %s\n", fn.Runtime))
	details.WriteString(fmt.Sprintf("[yellow]Handler:[white] %s\n", fn.Handler))
	details.WriteString(fmt.Sprintf("[yellow]Memory:[white] %d MB\n", fn.Memory))
	details.WriteString(fmt.Sprintf("[yellow]Timeout:[white] %d seconds\n", fn.Timeout))
	details.WriteString(fmt.Sprintf("[yellow]Status:[white] %s\n", fn.Status))
	
	if fn.Description != "" {
		details.WriteString(fmt.Sprintf("[yellow]Description:[white] %s\n", fn.Description))
	}
	
	if !fn.LastModified.IsZero() {
		details.WriteString(fmt.Sprintf("[yellow]Last Modified:[white] %s\n", 
			fn.LastModified.Format("2006-01-02 15:04:05")))
	}
	
	// Environment variables
	if len(fn.Environment) > 0 {
		details.WriteString("\n[yellow]Environment Variables:[white]\n")
		for k, v := range fn.Environment {
			details.WriteString(fmt.Sprintf("  %s = %s\n", k, v))
		}
	}
	
	// Add some sample actions
	details.WriteString("\n[blue]Available Actions:[white]\n")
	details.WriteString("  [green]Enter[white] - View logs\n")
	details.WriteString("  [green]i[white] - Invoke function\n")
	details.WriteString("  [green]r[white] - Refresh list\n")
	
	v.functionDetail.SetText(details.String())
}

func (v *View) updateStatus(message string) {
	// Update status in the main thread
	go func() {
		v.statusBar.SetText(message)
	}()
}

func (v *View) GetFunctionList() *tview.List {
	return v.functionList
}