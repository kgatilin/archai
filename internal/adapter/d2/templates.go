package d2

// legendTemplate returns the D2 legend block content.
// The legend is positioned at top-right and shows the DDD color coding.
func legendTemplate() string {
	return `legend: {
  label: "Color Legend (DDD)"
  style.stroke: "` + ColorLegendBorder + `"
  style.fill: "` + ColorLegendFill + `"
  near: top-right

  aggregate: {
    label: "Domain Model"
    shape: class
    style.fill: "` + ColorBlue + `"
    style.font-color: "#000"
  }
  service: {
    label: "Service"
    shape: class
    style.fill: "` + ColorPurple + `"
    style.font-color: "#000"
  }
  factory: {
    label: "Factory"
    shape: class
    style.fill: "` + ColorGreen + `"
    style.font-color: "#000"
  }
  options: {
    label: "Value Object"
    shape: class
    style.fill: "` + ColorGray + `"
    style.font-color: "#000"
  }
}`
}
