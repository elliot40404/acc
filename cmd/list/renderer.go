package list

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elliot40404/acc/pkg/database"
)

var page int = 1
var QC *database.TransactionConfig
var DB database.TransactionRepository = database.NewTransactionRepository()
var TotalTx int

func newModel(qc database.TransactionConfig, totalPages int) model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 1
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(totalPages)
	return model{
		paginator:   p,
		interactive: true,
		qc:          qc,
		totalPages:  totalPages,
	}
}

type model struct {
	paginator   paginator.Model
	interactive bool
	qc          database.TransactionConfig
	totalPages  int
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("acc list")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "h", "left":
			m.paginator.PrevPage()
			page--
			if page < 1 {
				page = 1
			}
			QC.Page = page
			return m, nil
		case "l", "right":
			m.paginator.NextPage()
			page++
			if page > m.totalPages {
				page = m.totalPages
			}
			QC.Page = page
			return m, nil
		case "G":
			m.paginator.Page = m.totalPages
			page = m.totalPages
			QC.Page = page
			return m, nil
		case "g":
			m.paginator.Page = 1
			page = 1
			QC.Page = page
			return m, nil
		}
	}
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	transactions, err := DB.GetTransactionsWithConfig(*QC)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	table := tableWriter(transactions, TotalTx, *QC, true)
	b.WriteString(table)
	if m.totalPages <= 50 {
		b.WriteString("  " + m.paginator.View())
	}
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return b.String()
}

func InteractiveListRenderer(queryConfig database.TransactionConfig) {
	totalTx, err := DB.GetTransactionCountWithConfig(queryConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	if queryConfig.Dry {
		return
	}
	TotalTx = totalTx
	QC = &queryConfig
	totalPages := int(math.Ceil(float64(totalTx) / float64(queryConfig.Limit)))
	model := newModel(queryConfig, totalPages)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func NonInteractiveListRenderer(queryConfig database.TransactionConfig) {
	transactions, err := DB.GetTransactionsWithConfig(queryConfig)
	totalTx, _ := DB.GetTransactionCountWithConfig(queryConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	if queryConfig.Dry {
		return
	}
	switch queryConfig.Format {
	case "json":
		jsonWriter(transactions, queryConfig.IsPretty)
	case "csv":
		csvWriter(transactions)
	default:
		tableWriter(transactions, totalTx, queryConfig, false)
	}
}
