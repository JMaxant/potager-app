# "Sown in place" is recorded, not inferred

A planting dated only by its sowing may be in the plot (sown in place) or in a tray waiting to be planted out, and the dates cannot tell which;
the species' `maturity_from` gives the usual practice, not what the gardener did. The planting therefore stores `sown_in_place`, 
defaulted from the species in the form. 
Occupancy starts at the sowing when it is set and at the transplant otherwise, so a seedling in its tray occupies nothing 
yet. The *start date* stays the earlier of the two dates: it dates the crop, occupancy dates the plot. 
The scaffold's `COALESCE(transplanted_at, sown_at)` is not used, because it puts a tray seedling in the plot from its sowing.
