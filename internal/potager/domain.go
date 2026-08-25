// Package potager holds the domain types and the agronomic rules. It must not
// import net/http or database/sql: everything here is a plain value or a pure
// function, which is what makes the rules cheap to test.
package potager

import "time"

// NitrogenDemand is what a plant consumes while it grows. It says nothing about
// what the plant leaves behind — see SoilEffect for that.
type NitrogenDemand string

const (
	DemandHeavy    NitrogenDemand = "heavy"
	DemandModerate NitrogenDemand = "moderate"
	DemandLight    NitrogenDemand = "light"
)

// MaturityFrom names the date that days-to-maturity is counted from. Seed
// packets count solanaceae and brassicas from transplant, roots and salads from
// sowing.
type MaturityFrom string

const (
	MaturityFromSowing     MaturityFrom = "sowing"
	MaturityFromTransplant MaturityFrom = "transplant"
)

type PlotType string

const (
	PlotOpenGround PlotType = "open_ground"
	PlotRaisedBed  PlotType = "raised_bed"
	PlotGreenhouse PlotType = "greenhouse"
	PlotColdframe  PlotType = "coldframe"
)

// Purpose is the gardener's intent for a planting. It lives here rather than on
// the species because féverole, moutarde and sarrasin are all dual-use: the same
// species is a food crop or a green manure depending on what you mean to do
// with it.
type Purpose string

const (
	PurposeFood  Purpose = "food"
	PurposeCover Purpose = "cover"
)

type Family struct {
	ID                    int64
	Name                  string
	LatinName             string
	RotationIntervalYears int
	NitrogenFixing        bool
}

type Species struct {
	ID                  int64
	Family              Family
	Name                string
	NitrogenDemand      NitrogenDemand
	MaturityFrom        MaturityFrom
	SpacingCm           *int
	SuitableAsCover     bool
	TerminationLeadDays *int
}

type Variety struct {
	ID             int64
	Species        Species
	Name           string
	DaysToMaturity *int
}

type Plot struct {
	ID         int64
	Name       string
	Type       PlotType
	LengthCm   *int
	WidthCm    *int
	Position   int
	ArchivedAt *time.Time
}

// PlantingVariety is one component of a planting. A cover crop mix — trèfle +
// seigle + radis fourrager — is one Planting with three of these, not three
// plantings.
type PlantingVariety struct {
	Variety     Variety
	Quantity    *int
	SeedRateGM2 *float64
}

type Planting struct {
	ID         int64
	PlotID     int64
	Purpose    Purpose
	Varieties  []PlantingVariety
	Notes      string
	ArchivedAt *time.Time

	SownAt         *Date
	TransplantedAt *Date
	EndedAt        *Date

	// ExpectedEndAt is a stored override. When set it always wins; when nil the
	// expected end is computed for food plantings and left open for cover crops.
	ExpectedEndAt *Date

	MaturityDaysOverride *int
}

// Date is a calendar date with no time or zone. Everything in a garden happens
// on a day, never at an instant.
type Date struct {
	time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (d Date) AddDays(n int) Date {
	return Date{d.Time.AddDate(0, 0, n)}
}

func (d Date) Year() int { return d.Time.Year() }
