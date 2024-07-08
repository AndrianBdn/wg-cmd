package main

import (
	"strings"

	"github.com/andrianbdn/wg-cmd/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DynamicTableList struct {
	rows        [][]string
	headers     []string
	offset      int
	selected    int
	colRatio    []int
	colMinWidth []int
	tableSize   tea.WindowSizeMsg
}

func NewMainTable(headers []string, rows [][]string, colRatio []int, colMinWidth []int) DynamicTableList {
	return DynamicTableList{
		headers:     headers,
		rows:        rows,
		colRatio:    colRatio,
		colMinWidth: colMinWidth,
	}
}

func (m *DynamicTableList) Up() {
	if m.selected > 0 {
		m.selected--
	}
	h := m.tableSize.Height - 1

	if m.selected < m.offset {
		m.offset -= h / 2
		if m.offset < 0 {
			m.offset = 0
		}
	}
}

func (m *DynamicTableList) Down() {
	if m.selected < len(m.rows)-1 {
		m.selected++
	}
	h := m.tableSize.Height - 1

	if m.selected >= m.offset+h {
		m.offset += h / 2
		if m.offset > len(m.rows)-h/2 {
			m.offset = len(m.rows) - h/2
		}
	}
}

func (m *DynamicTableList) PageUp() {
	for i := 0; i < m.tableSize.Height-1; i++ {
		m.Up()
	}
}

func (m *DynamicTableList) PageDown() {
	for i := 0; i < m.tableSize.Height-1; i++ {
		m.Down()
	}
}

func (m *DynamicTableList) DeleteSelectedRow() {
	if m.selected >= len(m.rows) {
		return
	}
	m.rows = append(m.rows[:m.selected], m.rows[m.selected+1:]...)
	if m.selected >= len(m.rows) {
		m.selected = len(m.rows) - 1
	}
}

func (m *DynamicTableList) CalcWidth(w int) []int {
	result := make([]int, len(m.colRatio))
	rsum := 0
	for _, r := range m.colRatio {
		rsum += r
	}
	total := 0
	for i, r := range m.colRatio {
		result[i] = w * r / rsum
		if result[i] < m.colMinWidth[i] {
			result[i] = m.colMinWidth[i]
		}
		total += result[i]
	}
	if total < w {
		result[len(result)-1] += w - total
	}
	return result
}

func (m *DynamicTableList) GetSelectedIndex() int {
	return m.selected
}

func (m *DynamicTableList) SetTableSize(msg tea.WindowSizeMsg, offsetW, offsetH int) {
	m.tableSize = tea.WindowSizeMsg{Width: msg.Width + offsetW, Height: msg.Height + offsetH}
}

func (m *DynamicTableList) CopyTableState(table *DynamicTableList) {
	m.tableSize = table.tableSize
	m.offset = table.offset
	m.selected = table.selected
}

func (m *DynamicTableList) GetSelected() []string {
	return m.rows[m.selected]
}

func (m *DynamicTableList) RenderRow(style *lipgloss.Style, row []string, width []int) string {
	result := ""
	for i, c := range row {
		result += style.Width(width[i]).Render(c)
	}
	return result
}

func (m *DynamicTableList) Render() string {
	if m.tableSize.Height <= 0 {
		return ""
	}
	result := ""
	width := m.CalcWidth(m.tableSize.Width)
	result += m.RenderRow(&theme.Current.MainTableHeader, m.headers, width)
	result += "\n"
	h := m.tableSize.Height - 1

	rowsRendered := 0

	OffsetMax := m.offset + h
	if OffsetMax > len(m.rows) {
		OffsetMax = len(m.rows)
	}

	for i, r := range m.rows[m.offset:OffsetMax] {

		style := &theme.Current.MainTableBody
		if m.offset+i == 0 {
			style = &theme.Current.MainTableFirst
		}
		if m.offset+i == m.selected {
			style = &theme.Current.MainTableSelected
			if m.offset+i == 0 {
				style = &theme.Current.MainTableSelectedFirst
			}
		}

		result += m.RenderRow(style, r, width)
		result += "\n"

		rowsRendered++
	}

	if rowsRendered < h {
		emptyCount := h - rowsRendered
		for i := 0; i < emptyCount; i++ {
			result += m.RenderRow(&theme.Current.MainTableBody, make([]string, len(m.headers)), width)
			result += "\n"
		}
	}

	return strings.TrimRight(result, "\n")
}
