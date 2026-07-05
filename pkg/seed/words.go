package seed

// WordDetails represents all dictionary & metadata attributes for a word.
type WordDetails struct {
	Word           string              `json:"word"`
	Phonetic       string              `json:"phonetic"`
	PartOfSpeech   string              `json:"part_of_speech"`
	Definitions    []Definition        `json:"definitions"`
	CorpusCount    int64               `json:"corpus_count"`
	ZipfScore      float64             `json:"zipf_score"`
	RarityTier     string              `json:"rarity_tier"` // "Uncommon", "Elegant", "Obscure", "Rare"
	RelatedWords   map[string][]string `json:"related_words"`
	Examples       []Example           `json:"examples"`
	AttributionText string             `json:"attribution_text"`
}

type Definition struct {
	Text            string `json:"text"`
	PartOfSpeech    string `json:"part_of_speech"`
	AttributionText string `json:"attribution_text"`
}

type Example struct {
	Text   string `json:"text"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// GetSeedWords returns a curated dataset of uncommon but highly expressive English vocabulary.
func GetSeedWords() []WordDetails {
	return []WordDetails{
		{
			Word:         "perspicacious",
			Phonetic:     "/ˌpɜː.spɪˈkeɪ.ʃəs/",
			PartOfSpeech: "adjective",
			CorpusCount:  210,
			ZipfScore:    2.8,
			RarityTier:   "Uncommon",
			AttributionText: "from The Century Dictionary and Wordnik",
			Definitions: []Definition{
				{
					Text:            "Having keen mental perception and understanding; discerning; acute; insightful.",
					PartOfSpeech:    "adjective",
					AttributionText: "from The Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"shrewd", "astute", "perceptive", "discerning", "sagacious", "clear-sighted"},
				"antonym": {"obtuse", "foolish", "ignorant", "undiscerning"},
				"same-context": {"acumen", "discernment", "insight"},
			},
			Examples: []Example{
				{
					Text:   "Her perspicacious analysis of the financial market saved the firm from catastrophic loss.",
					Title:  "Financial Quarterly Review",
					Author: "A. Vance",
				},
				{
					Text:   "He was a perspicacious observer of human character.",
					Title:  "Biographical Essays",
					Author: "C. R. Hall",
				},
			},
		},
		{
			Word:         "synecdoche",
			Phonetic:     "/sɪˈnɛk.də.ki/",
			PartOfSpeech: "noun",
			CorpusCount:  145,
			ZipfScore:    2.5,
			RarityTier:   "Elegant",
			AttributionText: "from American Heritage Dictionary",
			Definitions: []Definition{
				{
					Text:            "A figure of speech in which a part is used to represent the whole (e.g. 'hired hands' for workers), or the whole for a part (e.g. 'the law' for police officers).",
					PartOfSpeech:    "noun",
					AttributionText: "from American Heritage Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"metonymy", "trope", "figure of speech", "allegory"},
				"antonym": {"literalism"},
				"same-context": {"representation", "symbolism"},
			},
			Examples: []Example{
				{
					Text:   "Using 'suits' to refer to corporate executives is a classic synecdoche.",
					Title:  "Linguistic Stylistics",
					Author: "E. Bennet",
				},
			},
		},
		{
			Word:         "verisimilitude",
			Phonetic:     "/ˌvɛr.ɪ.sɪˈmɪl.ɪ.tjuːd/",
			PartOfSpeech: "noun",
			CorpusCount:  380,
			ZipfScore:    3.1,
			RarityTier:   "Uncommon",
			AttributionText: "from Century Dictionary",
			Definitions: []Definition{
				{
					Text:            "The appearance or semblance of truth or reality; likelihood; probability.",
					PartOfSpeech:    "noun",
					AttributionText: "from Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"plausibility", "realism", "credibility", "authenticity", "likeness"},
				"antonym": {"implausibility", "falsity", "unreality"},
			},
			Examples: []Example{
				{
					Text:   "The novel achieved a remarkable degree of verisimilitude through detailed historical research.",
					Title:  "Literary Craft",
					Author: "M. Sterling",
				},
			},
		},
		{
			Word:         "peripatetic",
			Phonetic:     "/ˌpɛr.ɪ.pəˈtɛt.ɪk/",
			PartOfSpeech: "adjective",
			CorpusCount:  290,
			ZipfScore:    2.9,
			RarityTier:   "Uncommon",
			AttributionText: "from WordNet 3.0",
			Definitions: []Definition{
				{
					Text:            "Traveling from place to place, especially working or tutoring in various places.",
					PartOfSpeech:    "adjective",
					AttributionText: "from WordNet 3.0",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"itinerant", "nomadic", "wandering", "wayfaring", "roving"},
				"antonym": {"sedentary", "stationary", "settled"},
			},
			Examples: []Example{
				{
					Text:   "He led a peripatetic life as a freelance consultant across South Europe.",
					Title:  "Modern Wayfarers",
					Author: "L. Thorne",
				},
			},
		},
		{
			Word:         "mellifluous",
			Phonetic:     "/mɛˈlɪf.lʊ.əs/",
			PartOfSpeech: "adjective",
			CorpusCount:  175,
			ZipfScore:    2.6,
			RarityTier:   "Elegant",
			AttributionText: "from Century Dictionary",
			Definitions: []Definition{
				{
					Text:            "Flowing sweetly or smoothly; sweet-sounding; pleasing to the ear (as if honeyed).",
					PartOfSpeech:    "adjective",
					AttributionText: "from Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"dulcet", "euphonious", "honeyed", "harmonious", "soothing"},
				"antonym": {"cacophonous", "harsh", "strident", "grating"},
			},
			Examples: []Example{
				{
					Text:   "The cellist played with a mellifluous tone that captivated the entire auditorium.",
					Title:  "Symphony Notes",
					Author: "R. Sterling",
				},
			},
		},
		{
			Word:         "sesquipedalian",
			Phonetic:     "/ˌsɛs.kwɪ.pɪˈdeɪ.lɪ.ən/",
			PartOfSpeech: "adjective",
			CorpusCount:  42,
			ZipfScore:    1.9,
			RarityTier:   "Obscure",
			AttributionText: "from American Heritage Dictionary",
			Definitions: []Definition{
				{
					Text:            "Given to using long words; (of a word) containing many syllables or a foot and a half long.",
					PartOfSpeech:    "adjective",
					AttributionText: "from American Heritage Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"polysyllabic", "grandiloquent", "pedantic", "verbose"},
				"antonym": {"concise", "pithy", "succinct", "monosyllabic"},
			},
			Examples: []Example{
				{
					Text:   "His sesquipedalian prose alienated casual readers who preferred plain English.",
					Title:  "Critique of Rhetoric",
					Author: "P. Davies",
				},
			},
		},
		{
			Word:         "quixotic",
			Phonetic:     "/kwɪkˈsɒt.ɪk/",
			PartOfSpeech: "adjective",
			CorpusCount:  510,
			ZipfScore:    3.4,
			RarityTier:   "Uncommon",
			AttributionText: "from Century Dictionary",
			Definitions: []Definition{
				{
					Text:            "Exceedingly idealistic, unrealistic, or impractical, especially in pursuit of noble goals.",
					PartOfSpeech:    "adjective",
					AttributionText: "from Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"idealistic", "chivalrous", "romantic", "impractical", "starry-eyed"},
				"antonym": {"pragmatic", "realistic", "utilitarian", "cynical"},
			},
			Examples: []Example{
				{
					Text:   "His quixotic quest to eliminate bureaucracy single-handedly ended in frustration.",
					Title:  "Political Memoirs",
					Author: "H. Finch",
				},
			},
		},
		{
			Word:         "liminal",
			Phonetic:     "/ˈlɪm.ɪ.nəl/",
			PartOfSpeech: "adjective",
			CorpusCount:  420,
			ZipfScore:    3.2,
			RarityTier:   "Uncommon",
			AttributionText: "from Wiktionary",
			Definitions: []Definition{
				{
					Text:            "Relating to a transitional or initial stage of a process; occupying a position at, or on both sides of, a boundary or threshold.",
					PartOfSpeech:    "adjective",
					AttributionText: "from Wiktionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"transitional", "intermediate", "threshold", "in-between"},
				"antonym": {"established", "permanent", "fixed"},
			},
			Examples: []Example{
				{
					Text:   "Dusk is a liminal space where day gradually surrenders to night.",
					Title:  "Atmospheric Studies",
					Author: "K. Owens",
				},
			},
		},
		{
			Word:         "ephemeral",
			Phonetic:     "/ɪˈfɛm.ər.əl/",
			PartOfSpeech: "adjective",
			CorpusCount:  680,
			ZipfScore:    3.6,
			RarityTier:   "Uncommon",
			AttributionText: "from WordNet 3.0",
			Definitions: []Definition{
				{
					Text:            "Lasting for a very short time; transient; fleeting.",
					PartOfSpeech:    "adjective",
					AttributionText: "from WordNet 3.0",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"transient", "fleeting", "evanescent", "fugitive", "momentary"},
				"antonym": {"eternal", "enduring", "permanent", "perennial"},
			},
			Examples: []Example{
				{
					Text:   "Fame in the digital age can be notoriously ephemeral.",
					Title:  "Media Trends",
					Author: "J. Mercer",
				},
			},
		},
		{
			Word:         "cacophony",
			Phonetic:     "/kəˈkɒf.ə.ni/",
			PartOfSpeech: "noun",
			CorpusCount:  540,
			ZipfScore:    3.5,
			RarityTier:   "Uncommon",
			AttributionText: "from Century Dictionary",
			Definitions: []Definition{
				{
					Text:            "A harsh, discordant mixture of sounds.",
					PartOfSpeech:    "noun",
					AttributionText: "from Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"din", "racket", "discord", "dissonance", "clamor"},
				"antonym": {"euphony", "harmony", "symphony", "peace"},
			},
			Examples: []Example{
				{
					Text:   "A cacophony of car horns and construction drills filled the morning air.",
					Title:  "Urban Life",
					Author: "D. Ross",
				},
			},
		},
		{
			Word:         "pulchritudinous",
			Phonetic:     "/ˌpʌl.krɪˈtjuː.dɪ.nəs/",
			PartOfSpeech: "adjective",
			CorpusCount:  35,
			ZipfScore:    1.8,
			RarityTier:   "Rare",
			AttributionText: "from Century Dictionary",
			Definitions: []Definition{
				{
					Text:            "Characterized by great physical beauty and comeliness.",
					PartOfSpeech:    "adjective",
					AttributionText: "from Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"beautiful", "comely", "resplendent", "gorgeous"},
				"antonym": {"hideous", "unsightly", "plain"},
			},
			Examples: []Example{
				{
					Text:   "The portrait depicted a pulchritudinous maiden adorned in silk.",
					Title:  "Victorian Aesthetics",
					Author: "G. Bell",
				},
			},
		},
		{
			Word:         "defenestration",
			Phonetic:     "/diːˌfɛn.ɪˈstreɪ.ʃən/",
			PartOfSpeech: "noun",
			CorpusCount:  120,
			ZipfScore:    2.4,
			RarityTier:   "Elegant",
			AttributionText: "from Wiktionary",
			Definitions: []Definition{
				{
					Text:            "The act of throwing someone or something out of a window; (figuratively) the dismissal of a leader or political figure.",
					PartOfSpeech:    "noun",
					AttributionText: "from Wiktionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"dismissal", "ejection", "ouster", "expulsion"},
				"antonym": {"appointment", "installation"},
			},
			Examples: []Example{
				{
					Text:   "The Defenestration of Prague triggered a major European conflict in 1618.",
					Title:  "History of Europe",
					Author: "W. Miller",
				},
			},
		},
		{
			Word:         "apocryphal",
			Phonetic:     "/əˈpɒk.rɪ.fəl/",
			PartOfSpeech: "adjective",
			CorpusCount:  460,
			ZipfScore:    3.3,
			RarityTier:   "Uncommon",
			AttributionText: "from American Heritage Dictionary",
			Definitions: []Definition{
				{
					Text:            "Of doubtful authenticity, although widely circulated as being true.",
					PartOfSpeech:    "adjective",
					AttributionText: "from American Heritage Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"spurious", "fictitious", "dubious", "unauthenticated", "legendary"},
				"antonym": {"authentic", "genuine", "verifiable", "factual"},
			},
			Examples: []Example{
				{
					Text:   "The story about George Washington and the cherry tree is largely apocryphal.",
					Title:  "Myths of History",
					Author: "E. Wright",
				},
			},
		},
		{
			Word:         "anachronism",
			Phonetic:     "/əˈnæk.rə.nɪz.əm/",
			PartOfSpeech: "noun",
			CorpusCount:  620,
			ZipfScore:    3.5,
			RarityTier:   "Uncommon",
			AttributionText: "from Century Dictionary",
			Definitions: []Definition{
				{
					Text:            "A thing belonging or appropriate to a period other than that in which it exists, especially a conspicuous error in chronology.",
					PartOfSpeech:    "noun",
					AttributionText: "from Century Dictionary",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"misplacement", "chronological error", "relic", "archaism"},
				"antonym": {"synchronism"},
			},
			Examples: []Example{
				{
					Text:   "A wristwatch on a medieval knight in the film was a hilarious anachronism.",
					Title:  "Cinema Errors",
					Author: "T. Brooks",
				},
			},
		},
		{
			Word:         "penultimate",
			Phonetic:     "/pɪˈnʌl.tɪ.mət/",
			PartOfSpeech: "adjective",
			CorpusCount:  750,
			ZipfScore:    3.7,
			RarityTier:   "Uncommon",
			AttributionText: "from WordNet 3.0",
			Definitions: []Definition{
				{
					Text:            "Next to last; second to final in a series.",
					PartOfSpeech:    "adjective",
					AttributionText: "from WordNet 3.0",
				},
			},
			RelatedWords: map[string][]string{
				"synonym": {"second-last", "next-to-last"},
				"antonym": {"ultimate", "final", "first"},
			},
			Examples: []Example{
				{
					Text:   "The penultimate chapter delivers the dramatic twist before the resolution.",
					Title:  "Novel Structure",
					Author: "S. King",
				},
			},
		},
	}
}
