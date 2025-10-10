package readerCharV3

type Lorebook struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	ScanDepth         int                    `json:"scan_depth"`
	TokenBudget       int                    `json:"token_budget"`
	RecursiveScanning bool                   `json:"recursive_scanning"`
	Extensions        map[string]interface{} `json:"extensions"`
	Entries           []struct {
		Keys           []string               `json:"keys"`
		Content        string                 `json:"content"`
		Extensions     map[string]interface{} `json:"extensions"`
		Enabled        bool                   `json:"enabled"`
		InsertionOrder int                    `json:"insertion_order"`
		CaseSensitive  *bool                  `json:"case_sensitive"`
		UseRegex       bool                   `json:"use_regex"`
		Constant       *bool                  `json:"constant"`
		Name           *string                `json:"name"`
		Priority       *int                   `json:"priority"`
		Id             *string                `json:"id"`
		Comment        *string                `json:"comment"`
		Selective      *bool                  `json:"selective"`
		SecondaryKeys  *[]string              `json:"secondary_keys"`
		Position       *string                `json:"position"`
	} `json:"entries"`
}

type CharacterCardV3 struct {
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	Tags                    []string               `json:"tags"`
	Creator                 string                 `json:"creator"`
	CharacterVersion        string                 `json:"character_version"`
	MesExample              string                 `json:"mes_example"`
	Extensions              map[string]interface{} `json:"extensions"`
	SystemPrompt            string                 `json:"system_prompt"`
	PostHistoryInstructions string                 `json:"post_history_instructions"`
	FirstMes                string                 `json:"first_mes"`
	AlternateGreetings      []string               `json:"alternate_greetings"`
	Personality             string                 `json:"personality"`
	Scenario                string                 `json:"scenario"`
	CreatorNotes            string                 `json:"creator_notes"`
	CharacterBook           *Lorebook              `json:"character_book"`
	Assets                  []struct {
		Type string `json:"type"`
		Uri  string `json:"uri"`
		Name string `json:"name"`
		Ext  string `json:"ext"`
	} `json:"assets"`
	Nickname                 string            `json:"nickname"`
	CreatorNotesMultilingual map[string]string `json:"creator_notes_multilingual"`
	Source                   []string          `json:"source"`
	GroupOnlyGreetings       []string          `json:"group_only_greetings"`
	CreationDate             *int              `json:"creation_date"`
	ModificationDate         *int              `json:"modification_date"`
}

func ReadPNG(file string) (string, error) {
	return "", nil
}
