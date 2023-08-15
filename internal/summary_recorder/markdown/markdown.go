package markdown

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// MarkdownSummaryRecorder represents a component that records a summary in markdown
type MarkdownSummaryRecorder struct {
	summaryLock sync.Mutex
	summaryRaw  *strings.Builder
	summary     string
	evaporation float32
	efficiency  float32
	timeline    string
}

// NewMarkdownSummaryRecorder creates a new MarkdownSummaryRecorder
func NewMarkdownSummaryRecorder() *MarkdownSummaryRecorder {
	var summaryRaw strings.Builder
	summaryRaw.WriteString("# Summary\n\n")
	summaryRaw.WriteString("The following summary was generated on " + time.Now().Format("2006-01-02 15:04:05") + "\n\n")
	return &MarkdownSummaryRecorder{
		summaryRaw: &summaryRaw,
	}
}

// addNewLine adds a new line to the summary
func (r *MarkdownSummaryRecorder) addNewLine(line string) {
	r.summaryLock.Lock()
	defer r.summaryLock.Unlock()
	r.summaryRaw.WriteString(line + "\n")
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (r *MarkdownSummaryRecorder) AddMashTemp(temp float64, notes string) {
	r.addNewLine("## Mash")
	line := fmt.Sprintf("- **Mashing temperature**: %.2f°C (%s)", temp, notes)
	r.addNewLine(line)
}

// AddRast adds a rast to the summary and notes related to it
func (r *MarkdownSummaryRecorder) AddRast(temp float64, duration float64, notes string) {
	line := fmt.Sprintf("- **Rast**: %.2f°C for %.2f minutes (%s)", temp, duration, notes)
	r.addNewLine(line)
}

// // AddLauternNotes adds lautern notes to the summary
func (r *MarkdownSummaryRecorder) AddLaunternNotes(notes string) {
	r.addNewLine("## Lautern")
	r.addNewLine(notes)
	r.addNewLine("")
	r.addNewLine("## Hopping")
}

// AddHopping adds a hopping to the summary and notes related to it
func (r *MarkdownSummaryRecorder) AddHopping(name string, amount float32, alpha float32, notes string) {
	line := fmt.Sprintf("- **%s**: %.2fg (%.2f%% alpha) (%s)", name, amount, alpha, notes)
	r.addNewLine(line)
}

// AddMeasuredVolume adds a measured volume to the summary
func (r *MarkdownSummaryRecorder) AddMeasuredVolume(name string, amount float32, notes string) {
	r.addNewLine("")
	line := fmt.Sprintf("- **Measured volume - %s**: %.2fL (%s)", name, amount, notes)
	r.addNewLine(line)
}

// AddEvaporation adds an evaporation to the summary
func (r *MarkdownSummaryRecorder) AddEvaporation(amount float32) {
	r.evaporation = amount
}

// AddCooling adds a cooling to the summary and notes related to it
func (r *MarkdownSummaryRecorder) AddCooling(finalTemp, coolingTime float32, notes string) {
	r.addNewLine("## Cooling")
	line := fmt.Sprintf("Reached %.2f°C in %.2f minutes (%s)", finalTemp, coolingTime, notes)
	r.addNewLine(line)
}

// AddSummaryPreFermentation adds a summary of the pre fermentation
func (r *MarkdownSummaryRecorder) AddSummaryPreFermentation(volume float32, sg float32, notes string) {
	r.addNewLine("## Pre fermentation")
	line := fmt.Sprintf("- **Volume**: %.2fL", volume)
	r.addNewLine(line)
	line = fmt.Sprintf("- **Specific gravity**: %.3f", sg)
	r.addNewLine(line)
	r.addNewLine(notes)
}

// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
func (r *MarkdownSummaryRecorder) AddEfficiency(efficiencyPercentage float32) {
	r.efficiency = efficiencyPercentage
}

// AddYeastStart adds the yeast start to the summary
func (r *MarkdownSummaryRecorder) AddYeastStart(temperature, notes string) {
	r.addNewLine("## Yeast start")
	line := fmt.Sprintf("- **Temperature**: %s°C", temperature)
	r.addNewLine(line)
	r.addNewLine(notes)
}

// Close closes the summary recorder
func (r *MarkdownSummaryRecorder) Close() {
	r.addNewLine("## Calculations")
	r.addNewLine(fmt.Sprintf("- **Evaporation**: %.2f%%/h", r.evaporation))
	r.addNewLine(fmt.Sprintf("- **Efficiency**: %.2f%%", r.efficiency))
	r.addNewLine("")
	r.addNewLine("## Timeline")
	r.addNewLine("Timestamp | Event")
	r.addNewLine("--- | ---")
	r.addNewLine(r.timeline)
	r.summary = r.summaryRaw.String()
}

// GetSummary returns the summary
func (r *MarkdownSummaryRecorder) GetSummary() string {
	return r.summary
}

// GetExtention returns the extention of the summary
func (r *MarkdownSummaryRecorder) GetExtention() string {
	return "md"
}

// AddTimeline adds a timeline to the summary
func (r *MarkdownSummaryRecorder) AddTimeline(timeline []string) {
	r.timeline = strings.Join(timeline, "\n")
}
