package engine

// Preset is a named voice (MorphVOX calls them voice aliases). Pitch and
// timbre values follow the examples in the MorphVOX Pro docs.
type Preset struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	NameZh string `json:"nameZh"`
	Hint   string `json:"hint"` // one-line description shown in the UI
	// Experimental marks a voice not yet validated by listening tests:
	// it stays renderable (voice-eval baselines, direct params) but the
	// UI badges it instead of presenting it as a finished product.
	Experimental bool   `json:"experimental,omitempty"`
	Params       Params `json:"params"`
}

// Presets: (pitch, timbre) pairs mirror the documented MorphVOX examples —
// (0,0)=normal, (-0.5,0)=low, (-0.5,-0.3)=giant, (+0.8,0)=high,
// (+0.8,+0.4)=child — extended to a full set of aliases.
var Presets = []Preset{
	{ID: "normal", Name: "Normal", NameZh: "原声", Hint: "不处理",
		Params: Params{Pitch: 0, Timbre: 0, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "warm-lecture", Name: "Warm Lecture", NameZh: "温暖讲堂", Hint: "稍低、柔和，适合长课",
		Params: Params{Speech: true, Pitch: -0.08, PitchRange: 0.95, Timbre: -0.09, Strength: 0.8, Gain: 0.95, EqPost: [10]float64{-3, -2, 0, 1, 0, 0, 0.5}}},
	{ID: "clear-demo", Name: "Clear Demo", NameZh: "清晰演示", Hint: "稍亮，突出操作讲解中的字音",
		Params: Params{Speech: true, Pitch: 0.055, PitchRange: 1, Timbre: 0.075, Strength: 0.85, Gain: 0.94, EqPost: [10]float64{-4, -2, -1, -0.5, 0, 0, 1, 1.5}}},
	{ID: "studio-narration", Name: "Studio Narration", NameZh: "沉稳旁白", Hint: "略厚、语调平稳，适合视频旁白",
		Params: Params{Speech: true, Pitch: -0.12, PitchRange: 0.9, Timbre: -0.12, Strength: 0.85, Gain: 0.94, EqPost: [10]float64{-4, -2, 0.5, 0, -1, 0, 1.5, 0.5}}},
	{ID: "casual-talk", Name: "Casual Talk", NameZh: "轻松对谈", Hint: "轻盈，保留问答的语调起伏",
		Params: Params{Speech: true, Pitch: 0.09, PitchRange: 1.04, Timbre: -0.045, Strength: 0.75, Gain: 0.97, EqPost: [10]float64{-2, -1, 0, 0, -0.5, 1, 0.5}}},
	{ID: "low", Name: "Low", NameZh: "低沉", Hint: "低半档音高",
		Params: Params{Pitch: -0.5, Timbre: 0, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "giant", Name: "Giant", NameZh: "巨人", Hint: "又低又宽",
		Params: Params{Pitch: -0.5, Timbre: -0.3, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "demon", Name: "Demon", NameZh: "恶魔", Hint: "极低共振峰",
		Params: Params{Pitch: -0.8, Timbre: -0.5, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "woman", Name: "Woman", NameZh: "女声", Hint: "偏高偏细",
		Params: Params{Pitch: 0.45, Timbre: 0.28, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "high", Name: "High", NameZh: "高音", Hint: "纯音高上移",
		Params: Params{Pitch: 0.8, Timbre: 0, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "child", Name: "Child", NameZh: "小孩", Hint: "高音高共振峰",
		Params: Params{Pitch: 0.8, Timbre: 0.4, Strength: 1, Gain: 1, PitchRange: 1}},
	{ID: "chipmunk", Name: "Chipmunk", NameZh: "花栗鼠", Hint: "极高玩梗",
		Params: Params{Pitch: 1, Timbre: 0.55, Strength: 1, Gain: 1, PitchRange: 1}},
	// Speech-path voices: an absolute F0 center instead of a pitch ratio,
	// so any input speaker lands near the same target register.
	{ID: "natural-woman", Name: "Natural Woman", NameZh: "自然女声", Hint: "实验 · 柔和", Experimental: true,
		Params: Params{Speech: true, PitchRange: 1.1, TargetF0: 220, Timbre: 0.20, FormantShape: 0.35, Strength: 1, Gain: 1}},
	{ID: "natural-child", Name: "Natural Child", NameZh: "自然童声", Hint: "实验 · 明亮", Experimental: true,
		Params: Params{Speech: true, PitchRange: 1.15, TargetF0: 280, Timbre: 0.35, FormantShape: 0.45, Strength: 1, Gain: 1}},
}
