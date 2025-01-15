package Constants

import (
	"fmt"
	"time"
)

var CURRENT_YEAR = time.Now().Year()
var EVENT_BASE_URL = "https://www.vlr.gg/event/matches/%d/"        // %d should be our Event ID (e.g; 2274 would be for Champions Tour 2025: Americas Kickoff)
var VCT_BASE_URL = "https://www.vlr.gg/vct-%d"                     // This is where all the Events per region live so we can go to their pages
var VCT_BASE_URL_CURRENT = fmt.Sprintf(VCT_BASE_URL, CURRENT_YEAR) // Should just fetch the current VCT circuit for this year
