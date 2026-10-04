# Plots are archived; plantings are deleted

A plot accumulates history, so a mis-click must be undoable: plots are archived and can be restored. A planting carries nothing beyond a mistaken entry (no yield, no outcome), so it is deleted for good after a confirmation screen, and `plantings.archived_at` stays unused. A failed crop is never deleted: it is closed with `ended_at` and a note, because the rotation needs to know the plot was occupied. If harvests or incidents are later attached to plantings, revisit this.
