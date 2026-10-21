package style

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	Primary   = lipgloss.Color("#7C3AED")
	Secondary = lipgloss.Color("#0D9488")
	Text      = lipgloss.Color("#1F2937")
	Muted     = lipgloss.Color("#6B7280")
	Error     = lipgloss.Color("#DC2626")
	Success   = lipgloss.Color("#15803D")
	Warning   = lipgloss.Color("#B45309")
)

var (
	Title          = lipgloss.NewStyle().Foreground(Primary).Bold(true)
	Normal         = lipgloss.NewStyle().Foreground(Text)
	Description    = lipgloss.NewStyle().Foreground(Muted)
	ErrorMessage   = lipgloss.NewStyle().Foreground(Error)
	SuccessMessage = lipgloss.NewStyle().Foreground(Success)
	WarningMessage = lipgloss.NewStyle().Foreground(Warning)
)

var HuhTheme huh.Theme = huh.ThemeFunc(func(isDark bool) *huh.Styles {
	styles := huh.ThemeBase(isDark)
	text := lipgloss.Color("#E5E7EB")
	muted := lipgloss.Color("#A1A1AA")
	primary := Primary
	secondary := Secondary
	err := Error

	if isDark {
		primary = lipgloss.Color("#A78BFA")
		secondary = lipgloss.Color("#2DD4BF")
		err = lipgloss.Color("#F87171")
	}

	styles.Form.Base = styles.Form.Base.Foreground(text)
	styles.Group.Title = styles.Group.Title.Foreground(primary).Bold(true)
	styles.Group.Description = styles.Group.Description.Foreground(muted)

	styles.Focused.Base = styles.Focused.Base.Foreground(text).BorderForeground(secondary)
	styles.Focused.Card = styles.Focused.Base
	styles.Focused.Title = styles.Focused.Title.Foreground(primary).Bold(true)
	styles.Focused.NoteTitle = styles.Focused.NoteTitle.Foreground(primary).Bold(true)
	styles.Focused.Description = styles.Focused.Description.Foreground(muted)
	styles.Focused.ErrorIndicator = styles.Focused.ErrorIndicator.Foreground(err)
	styles.Focused.ErrorMessage = styles.Focused.ErrorMessage.Foreground(err)
	styles.Focused.SelectSelector = styles.Focused.SelectSelector.Foreground(secondary)
	styles.Focused.NextIndicator = styles.Focused.NextIndicator.Foreground(secondary)
	styles.Focused.PrevIndicator = styles.Focused.PrevIndicator.Foreground(secondary)
	styles.Focused.Option = styles.Focused.Option.Foreground(text)
	styles.Focused.SelectedOption = styles.Focused.SelectedOption.Foreground(secondary)
	styles.Focused.UnselectedOption = styles.Focused.UnselectedOption.Foreground(text)
	styles.Focused.SelectedPrefix = styles.Focused.SelectedPrefix.Foreground(secondary)
	styles.Focused.UnselectedPrefix = styles.Focused.UnselectedPrefix.Foreground(muted)
	styles.Focused.FocusedButton = styles.Focused.FocusedButton.Foreground(text).Background(secondary)
	styles.Focused.BlurredButton = styles.Focused.BlurredButton.Foreground(text).Background(muted)
	styles.Focused.Next = styles.Focused.Next.Foreground(text).Background(secondary)
	styles.Focused.TextInput.Cursor = styles.Focused.TextInput.Cursor.Foreground(secondary)
	styles.Focused.TextInput.CursorText = styles.Focused.TextInput.CursorText.Foreground(text)
	styles.Focused.TextInput.Text = styles.Focused.TextInput.Text.Foreground(text)
	styles.Focused.TextInput.Prompt = styles.Focused.TextInput.Prompt.Foreground(secondary)
	styles.Focused.TextInput.Placeholder = styles.Focused.TextInput.Placeholder.Foreground(muted)

	styles.Blurred.Base = styles.Blurred.Base.Foreground(text)
	styles.Blurred.Card = styles.Blurred.Base
	styles.Blurred.Title = styles.Blurred.Title.Foreground(primary).Bold(true)
	styles.Blurred.NoteTitle = styles.Blurred.NoteTitle.Foreground(primary).Bold(true)
	styles.Blurred.Description = styles.Blurred.Description.Foreground(muted)
	styles.Blurred.ErrorIndicator = styles.Blurred.ErrorIndicator.Foreground(err)
	styles.Blurred.ErrorMessage = styles.Blurred.ErrorMessage.Foreground(err)
	styles.Blurred.SelectSelector = styles.Blurred.SelectSelector.Foreground(muted)
	styles.Blurred.NextIndicator = styles.Blurred.NextIndicator.Foreground(muted)
	styles.Blurred.PrevIndicator = styles.Blurred.PrevIndicator.Foreground(muted)
	styles.Blurred.Option = styles.Blurred.Option.Foreground(muted)
	styles.Blurred.SelectedOption = styles.Blurred.SelectedOption.Foreground(secondary)
	styles.Blurred.UnselectedOption = styles.Blurred.UnselectedOption.Foreground(muted)
	styles.Blurred.SelectedPrefix = styles.Blurred.SelectedPrefix.Foreground(secondary)
	styles.Blurred.UnselectedPrefix = styles.Blurred.UnselectedPrefix.Foreground(muted)
	styles.Blurred.FocusedButton = styles.Blurred.FocusedButton.Foreground(text).Background(secondary)
	styles.Blurred.BlurredButton = styles.Blurred.BlurredButton.Foreground(text).Background(muted)
	styles.Blurred.Next = styles.Blurred.Next.Foreground(text).Background(secondary)
	styles.Blurred.TextInput.Cursor = styles.Blurred.TextInput.Cursor.Foreground(muted)
	styles.Blurred.TextInput.CursorText = styles.Blurred.TextInput.CursorText.Foreground(text)
	styles.Blurred.TextInput.Text = styles.Blurred.TextInput.Text.Foreground(text)
	styles.Blurred.TextInput.Prompt = styles.Blurred.TextInput.Prompt.Foreground(muted)
	styles.Blurred.TextInput.Placeholder = styles.Blurred.TextInput.Placeholder.Foreground(muted)

	styles.Group.Title = styles.Focused.Title
	styles.Group.Description = styles.Focused.Description
	styles.Help.ShortKey = styles.Help.ShortKey.Foreground(primary)
	styles.Help.ShortDesc = styles.Help.ShortDesc.Foreground(muted)
	styles.Help.ShortSeparator = styles.Help.ShortSeparator.Foreground(muted)
	styles.Help.FullKey = styles.Help.FullKey.Foreground(primary)
	styles.Help.FullDesc = styles.Help.FullDesc.Foreground(muted)
	styles.Help.FullSeparator = styles.Help.FullSeparator.Foreground(muted)

	return styles
})
