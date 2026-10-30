# Potager

A kitchen-garden manager: what is growing in which plot, since when, until when,
and whether the succession makes agronomic sense.

## Language

**Plot** (*parcelle*):
A bounded growing area of the garden: open ground, raised bed, greenhouse or coldframe.
_Avoid_: bed, field, garden

**Catalogue** (*catalogue*):
The botanical reference the garden draws on: families, species and varieties. Seeded by migration and read-only in step 1; the gardener records plots and plantings, never the catalogue. Not a type in code.
_Avoid_: listing, library, reference data

**Planting** (*culture*):
One act of sowing or transplanting in one plot, with one intent and one set of dates. A cover crop mix is a single planting.
_Avoid_: crop, sowing, culture (in code)

**Purpose**:
The gardener's intent for a planting: food or cover. It belongs to the planting, never to the species.

**Food crop** (*culture vivrière*):
A planting meant to be harvested and eaten.

**Cover crop** (*engrais vert*):
A planting grown for the soil, never for harvest. It is terminated before it fruits.
_Avoid_: green manure (in the interface), catch crop

**Termination** (*terminaison*):
The end of a cover crop, by any means: mowing, rolling, crimping, frost or tillage. The garden is no-till, so it is never assumed to be buried.
_Avoid_: incorporation, burying (*enfouissement*), destruction

**Harvest** (*récolte*):
The end of a food crop.

**Expected end** (*fin prévue*; *récolte prévue* or *terminaison prévue* in the interface):
When a planting should leave the ground. Either entered by hand, which always wins, or derived for a food crop: the date the species counts maturity from (sowing or transplant), plus the variety's days to maturity. Never derived from the other date when that one is missing, and never derived for a cover crop.
_Avoid_: due date, harvest date (that is when it actually ended)

**Soil effect** (*effet sur le sol*):
What a planting leaves for whatever follows it: enriching, neutral, slightly depleting or depleting. Always derived, never stored.

**Season** (*saison*):
The year a planting leaves the ground, not the year it went in.

**Occupancy**:
The span of dates during which a planting holds its plot.
