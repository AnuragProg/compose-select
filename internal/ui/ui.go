package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AnuragProg/compose-select/internal/parser"
)

const (
	FILENAME_SCREEN			= 1
	SERVICENAME_SCREEN		= 2
	OUTPUTFILENAME_SCREEN	= 3
	END_SCREEN					= 4
)


type UI struct {

	err error

	page uint // 0 - filename, 1 - service
	filenameInput textinput.Model
	serviceNameInput textinput.Model
	outputFilenameInput textinput.Model
	composeFile *parser.ComposeFile
}

func NewUI() UI{
	filenameInput := textinput.New()
	filenameInput.Placeholder = "Enter filename..."
	filenameInput.Focus()
	filenameInput.CharLimit =	150
	filenameInput.Width = 20

	serviceNameInput := textinput.New()
	serviceNameInput.Placeholder = "Enter service name..."
	serviceNameInput.Focus()
	serviceNameInput.CharLimit =	150
	serviceNameInput.Width = 20

	outputFilenameInput := textinput.New()
	outputFilenameInput.Placeholder = "Enter output file name..."
	outputFilenameInput.Focus()
	outputFilenameInput.CharLimit =	150
	outputFilenameInput.Width = 20
	return UI{
		page: 1,
		filenameInput: filenameInput,
		serviceNameInput: serviceNameInput,
		outputFilenameInput: outputFilenameInput,
	}
}

func (ui UI) Run() error {
	program := tea.NewProgram(ui)
	if _, err := program.Run(); err != nil{
		return err
	}
	return nil
}

func (ui UI) Init() tea.Cmd {
	return textinput.Blink
}

func (ui UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// updating filename input field
	var cmd tea.Cmd

	switch ui.page {
	case 1:
	ui.filenameInput, cmd = ui.filenameInput.Update(msg)
	case 2:
	ui.serviceNameInput, cmd = ui.serviceNameInput.Update(msg)
	case 3:
	ui.outputFilenameInput, cmd = ui.outputFilenameInput.Update(msg)
	}

	if msgKey, ok := msg.(tea.KeyMsg); ok {
		if ui.page == END_SCREEN {
			return ui, tea.Quit
		}

		// handle enter events
		if msgKey.Type == tea.KeyEnter {    
			switch ui.page{
			case FILENAME_SCREEN:
				// file input screen
				ui.composeFile, ui.err = parser.NewComposeFile(ui.filenameInput.Value())
				if ui.err == nil {
					ui.filenameInput.Blur()
					ui.page = SERVICENAME_SCREEN
				}

			case SERVICENAME_SCREEN:
				// service input screen
				if err := ui.composeFile.GetDependentServicesYAML(ui.serviceNameInput.Value()); err != nil{
					ui.err = err
				}else {
					ui.err = nil
					ui.serviceNameInput.Blur()
					ui.page = OUTPUTFILENAME_SCREEN
				}

			case OUTPUTFILENAME_SCREEN:
				// output file screen
				if err := ui.composeFile.WriteYAML(ui.outputFilenameInput.Value()); err != nil {
					ui.err = err
				}else {
					ui.err = nil
					ui.outputFilenameInput.Blur()
					ui.page = END_SCREEN
				}
			}
		}

		if msgKey.Type == tea.KeyEsc {
			return ui, tea.Quit
		}
	}

	return ui, cmd
}

func (ui UI) View() string {
	screen := ""
	if ui.err != nil {
		screen += "Error: " + ui.err.Error() + "\n\n"
	}
	switch ui.page {
		case FILENAME_SCREEN:
		screen += fmt.Sprintf("Enter filename:\n\n%s", ui.filenameInput.View())
		case SERVICENAME_SCREEN:
		screen += fmt.Sprintf("Enter service name:\n\n%s", ui.serviceNameInput.View())
		case OUTPUTFILENAME_SCREEN:
		screen += fmt.Sprintf("Enter output filename:\n\n%s", ui.outputFilenameInput.View())
		case END_SCREEN:
		screen += "File created successfully. (Press any button to exit)\n\n"
	}
	if ui.page != END_SCREEN{
		screen += "\n\n(esc to quit)\n\n"
	}
	return screen
}
