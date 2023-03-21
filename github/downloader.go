package github

// primaryQueue

// github.HasRequestSlot()

// func (g *GitHub) HasRequestSlot() bool {
// 	var reservedSlots = 5
// 	if g.RateLimitRemaning == 0 && g.RateLimit == 0 {
// 		fmt.Println("Can not get rate limit, github just created")
// 	}

// 	if g.RateLimitRemaning < g.RateLimit-reservedSlots {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func (g *GitHub) download() {

// 	if g.HasRequestSlot() {
// 		g.Download()
// 	}
// }
